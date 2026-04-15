package physics

import (
	"github.com/shaolei/cubism-go/internal/core"
	"github.com/shaolei/cubism-go/internal/id"
	"github.com/shaolei/cubism-go/internal/model"
)

// Physics constants matching the official SDK
const (
	AirResistance    float32 = 5.0
	MaximumWeight    float32 = 100.0
	MovementThreshold float32 = 0.001
	MaxDeltaTime     float32 = 5.0
)

// PhysicsManager manages the physics simulation for a model.
// Matches the official SDK's CubismPhysics class.
//
// The physics engine uses a pendulum (spring-damper) model:
//   - Each SubRig contains a strand of Particles
//   - Inputs map model parameters → normalized translation/angle
//   - Particles simulate a chain of connected masses with gravity and wind
//   - Outputs map particle positions → model parameter values
//   - Frame interpolation prevents jitter between physics timesteps
type PhysicsManager struct {
	rig                *Rig
	options            Options
	core               core.Core
	modelPtr           uintptr
	idManager          *id.CubismIdManager
	currentRigOutputs  []rigOutput
	previousRigOutputs []rigOutput
	currentRemainTime  float32
	parameterCaches    []float32
	parameterInputCaches []float32
}

// NewPhysicsManager creates a new physics manager from parsed physics JSON data.
func NewPhysicsManager(c core.Core, modelPtr uintptr, physicsJson model.PhysicsJson) *PhysicsManager {
	pm := &PhysicsManager{
		core:     c,
		modelPtr: modelPtr,
		options: Options{
			Gravity: Vector2{X: 0, Y: -1.0},
			Wind:    Vector2{},
		},
		currentRemainTime: 0,
	}

	// Parse the physics JSON into the internal rig representation
	pm.rig = ParsePhysicsJson(physicsJson)
	if pm.rig == nil {
		return nil
	}

	// Initialize output caches
	pm.currentRigOutputs = make([]rigOutput, pm.rig.SubRigCount)
	pm.previousRigOutputs = make([]rigOutput, pm.rig.SubRigCount)
	for i := 0; i < pm.rig.SubRigCount; i++ {
		outputCount := pm.rig.Settings[i].OutputCount
		pm.currentRigOutputs[i] = rigOutput{outputs: make([]float32, outputCount)}
		pm.previousRigOutputs[i] = rigOutput{outputs: make([]float32, outputCount)}
	}

	pm.initialize()

	return pm
}

// SetIdManager sets the CubismIdManager for fast parameter access by index.
func (pm *PhysicsManager) SetIdManager(idMgr *id.CubismIdManager) {
	pm.idManager = idMgr
}

// SetOptions sets the external gravity and wind options.
func (pm *PhysicsManager) SetOptions(options Options) {
	pm.options = options
}

// GetOptions returns the current physics options.
func (pm *PhysicsManager) GetOptions() Options {
	return pm.options
}

// Reset resets all physics state to initial values.
// Matches CubismPhysics::Reset.
func (pm *PhysicsManager) Reset() {
	pm.options = Options{
		Gravity: Vector2{X: 0, Y: -1.0},
		Wind:    Vector2{},
	}
	pm.rig.Gravity = Vector2{}
	pm.rig.Wind = Vector2{}
	pm.currentRemainTime = 0
	pm.initialize()
}

// initialize sets up initial particle positions and states.
// Matches CubismPhysics::Initialize.
func (pm *PhysicsManager) initialize() {
	for settingIndex := 0; settingIndex < pm.rig.SubRigCount; settingIndex++ {
		setting := &pm.rig.Settings[settingIndex]
		strand := pm.rig.Particles[setting.BaseParticleIndex:]

		// Initialize the top particle (anchor)
		strand[0].InitialPosition = Vector2{}
		strand[0].LastPosition = strand[0].InitialPosition
		strand[0].LastGravity = Vector2{X: 0, Y: 1.0} // Y flipped: -1 * -1 = 1
		strand[0].Velocity = Vector2{}
		strand[0].Force = Vector2{}

		// Initialize remaining particles in a chain
		for i := 1; i < setting.ParticleCount; i++ {
			radius := Vector2{X: 0, Y: strand[i].Radius}
			strand[i].InitialPosition = strand[i-1].InitialPosition.Add(radius)
			strand[i].Position = strand[i].InitialPosition
			strand[i].LastPosition = strand[i].InitialPosition
			strand[i].LastGravity = Vector2{X: 0, Y: 1.0} // Y flipped
			strand[i].Velocity = Vector2{}
			strand[i].Force = Vector2{}
		}
	}
}

// Stabilization runs the physics simulation until it reaches a stable state
// using the current parameter values. Matches CubismPhysics::Stabilization.
func (pm *PhysicsManager) Stabilization() {
	parameterValues := pm.getParameterValues()
	parameterMinimumValues := pm.getParameterMinimumValues()
	parameterMaximumValues := pm.getParameterMaximumValues()
	parameterDefaultValues := pm.getParameterDefaultValues()
	parameterCount := pm.getParameterCount()

	if len(pm.parameterCaches) < parameterCount {
		pm.parameterCaches = make([]float32, parameterCount)
	}
	if len(pm.parameterInputCaches) < parameterCount {
		pm.parameterInputCaches = make([]float32, parameterCount)
	}

	for j := 0; j < parameterCount; j++ {
		pm.parameterCaches[j] = parameterValues[j]
		pm.parameterInputCaches[j] = parameterValues[j]
	}

	for settingIndex := 0; settingIndex < pm.rig.SubRigCount; settingIndex++ {
		totalAngle := float32(0.0)
		totalTranslation := Vector2{}
		setting := &pm.rig.Settings[settingIndex]
		inputs := pm.rig.Inputs[setting.BaseInputIndex:]
		outputs := pm.rig.Outputs[setting.BaseOutputIndex:]
		particles := pm.rig.Particles[setting.BaseParticleIndex:]

		// Load input parameters
		for i := 0; i < setting.InputCount; i++ {
			weight := inputs[i].Weight / MaximumWeight

			if inputs[i].SourceParameterIndex == -1 {
				inputs[i].SourceParameterIndex = pm.getParameterIndex(inputs[i].SourceId)
			}

			if inputs[i].SourceParameterIndex < 0 || inputs[i].SourceParameterIndex >= parameterCount {
				continue
			}

			pm.getInputNormalizedParameterValue(
				&totalTranslation, &totalAngle,
				parameterValues[inputs[i].SourceParameterIndex],
				parameterMinimumValues[inputs[i].SourceParameterIndex],
				parameterMaximumValues[inputs[i].SourceParameterIndex],
				parameterDefaultValues[inputs[i].SourceParameterIndex],
				&setting.NormalizationPosition,
				&setting.NormalizationAngle,
				inputs[i],
				weight,
			)

			pm.parameterCaches[inputs[i].SourceParameterIndex] = parameterValues[inputs[i].SourceParameterIndex]
		}

		radAngle := degreesToRadian(-totalAngle)
		totalTranslation = Vector2{
			X: totalTranslation.X*cosFloat32(radAngle) - totalTranslation.Y*sinFloat32(radAngle),
			Y: totalTranslation.X*sinFloat32(radAngle) + totalTranslation.Y*cosFloat32(radAngle),
		}

		// Calculate particles position for stabilization
		pm.updateParticlesForStabilization(
			particles, setting.ParticleCount,
			totalTranslation, totalAngle,
			pm.options.Wind,
			MovementThreshold*setting.NormalizationPosition.Maximum,
		)

		// Update output parameters
		for i := 0; i < setting.OutputCount; i++ {
			particleIndex := outputs[i].VertexIndex

			if outputs[i].DestinationParameterIndex == -1 {
				outputs[i].DestinationParameterIndex = pm.getParameterIndex(outputs[i].DestinationId)
			}

			if outputs[i].DestinationParameterIndex < 0 || outputs[i].DestinationParameterIndex >= parameterCount {
				continue
			}

			if particleIndex < 1 || particleIndex >= setting.ParticleCount {
				continue
			}

			translation := particles[particleIndex].Position.Sub(particles[particleIndex-1].Position)
			outputValue := pm.getOutputValue(
				translation, particles, particleIndex,
				outputs[i], pm.options.Gravity,
			)

			pm.currentRigOutputs[settingIndex].outputs[i] = outputValue
			pm.previousRigOutputs[settingIndex].outputs[i] = outputValue

			pm.updateOutputParameterValue(
				&parameterValues[outputs[i].DestinationParameterIndex],
				parameterMinimumValues[outputs[i].DestinationParameterIndex],
				parameterMaximumValues[outputs[i].DestinationParameterIndex],
				outputValue, &outputs[i],
			)

			pm.parameterCaches[outputs[i].DestinationParameterIndex] = parameterValues[outputs[i].DestinationParameterIndex]
		}
	}
}

// Evaluate performs the physics simulation for one frame.
// Matches CubismPhysics::Evaluate.
func (pm *PhysicsManager) Evaluate(deltaTimeSeconds float64) {
	if deltaTimeSeconds <= 0 {
		return
	}

	parameterValues := pm.getParameterValues()
	parameterMinimumValues := pm.getParameterMinimumValues()
	parameterMaximumValues := pm.getParameterMaximumValues()
	parameterDefaultValues := pm.getParameterDefaultValues()
	parameterCount := pm.getParameterCount()

	pm.currentRemainTime += float32(deltaTimeSeconds)
	if pm.currentRemainTime > MaxDeltaTime {
		pm.currentRemainTime = 0.0
	}

	if len(pm.parameterCaches) < parameterCount {
		pm.parameterCaches = make([]float32, parameterCount)
	}
	if len(pm.parameterInputCaches) < parameterCount {
		pm.parameterInputCaches = make([]float32, parameterCount)
		for j := 0; j < parameterCount; j++ {
			pm.parameterInputCaches[j] = parameterValues[j]
		}
	}

	var physicsDeltaTime float32
	if pm.rig.Fps > 0 {
		physicsDeltaTime = 1.0 / pm.rig.Fps
	} else {
		physicsDeltaTime = float32(deltaTimeSeconds)
	}

	for pm.currentRemainTime >= physicsDeltaTime {
		// Copy current rig outputs to previous
		for settingIndex := 0; settingIndex < pm.rig.SubRigCount; settingIndex++ {
			setting := &pm.rig.Settings[settingIndex]
			for i := 0; i < setting.OutputCount; i++ {
				pm.previousRigOutputs[settingIndex].outputs[i] = pm.currentRigOutputs[settingIndex].outputs[i]
			}
		}

		// Linear interpolation of input between cached and current values
		inputWeight := physicsDeltaTime / pm.currentRemainTime
		for j := 0; j < parameterCount; j++ {
			pm.parameterCaches[j] = pm.parameterInputCaches[j]*(1.0-inputWeight) + parameterValues[j]*inputWeight
			pm.parameterInputCaches[j] = pm.parameterCaches[j]
		}

		for settingIndex := 0; settingIndex < pm.rig.SubRigCount; settingIndex++ {
			totalAngle := float32(0.0)
			totalTranslation := Vector2{}
			setting := &pm.rig.Settings[settingIndex]
			inputs := pm.rig.Inputs[setting.BaseInputIndex:]
			outputs := pm.rig.Outputs[setting.BaseOutputIndex:]
			particles := pm.rig.Particles[setting.BaseParticleIndex:]

			// Load input parameters
			for i := 0; i < setting.InputCount; i++ {
				weight := inputs[i].Weight / MaximumWeight

				if inputs[i].SourceParameterIndex == -1 {
					inputs[i].SourceParameterIndex = pm.getParameterIndex(inputs[i].SourceId)
				}

				if inputs[i].SourceParameterIndex < 0 || inputs[i].SourceParameterIndex >= parameterCount {
					continue
				}

				pm.getInputNormalizedParameterValue(
					&totalTranslation, &totalAngle,
					pm.parameterCaches[inputs[i].SourceParameterIndex],
					parameterMinimumValues[inputs[i].SourceParameterIndex],
					parameterMaximumValues[inputs[i].SourceParameterIndex],
					parameterDefaultValues[inputs[i].SourceParameterIndex],
					&setting.NormalizationPosition,
					&setting.NormalizationAngle,
					inputs[i],
					weight,
				)
			}

			radAngle := degreesToRadian(-totalAngle)
			totalTranslation = Vector2{
				X: totalTranslation.X*cosFloat32(radAngle) - totalTranslation.Y*sinFloat32(radAngle),
				Y: totalTranslation.X*sinFloat32(radAngle) + totalTranslation.Y*cosFloat32(radAngle),
			}

			// Calculate particles position
			pm.updateParticles(
				particles, setting.ParticleCount,
				totalTranslation, totalAngle,
				pm.options.Wind,
				MovementThreshold*setting.NormalizationPosition.Maximum,
				physicsDeltaTime,
				AirResistance,
			)

			// Update output parameters
			for i := 0; i < setting.OutputCount; i++ {
				particleIndex := outputs[i].VertexIndex

				if outputs[i].DestinationParameterIndex == -1 {
					outputs[i].DestinationParameterIndex = pm.getParameterIndex(outputs[i].DestinationId)
				}

				if outputs[i].DestinationParameterIndex < 0 || outputs[i].DestinationParameterIndex >= parameterCount {
					continue
				}

				if particleIndex < 1 || particleIndex >= setting.ParticleCount {
					continue
				}

				translation := particles[particleIndex].Position.Sub(particles[particleIndex-1].Position)
				outputValue := pm.getOutputValue(
					translation, particles, particleIndex,
					outputs[i], pm.options.Gravity,
				)

				pm.currentRigOutputs[settingIndex].outputs[i] = outputValue

				pm.updateOutputParameterValue(
					&pm.parameterCaches[outputs[i].DestinationParameterIndex],
					parameterMinimumValues[outputs[i].DestinationParameterIndex],
					parameterMaximumValues[outputs[i].DestinationParameterIndex],
					outputValue, &outputs[i],
				)
			}
		}

		pm.currentRemainTime -= physicsDeltaTime
	}

	alpha := pm.currentRemainTime / physicsDeltaTime
	pm.interpolate(parameterValues, parameterMinimumValues, parameterMaximumValues, parameterCount, alpha)
}

// interpolate blends between previous and current rig outputs.
// Matches CubismPhysics::Interpolate.
func (pm *PhysicsManager) interpolate(parameterValues, parameterMinimumValues, parameterMaximumValues []float32, parameterCount int, weight float32) {
	for settingIndex := 0; settingIndex < pm.rig.SubRigCount; settingIndex++ {
		setting := &pm.rig.Settings[settingIndex]
		outputs := pm.rig.Outputs[setting.BaseOutputIndex:]

		for i := 0; i < setting.OutputCount; i++ {
			if outputs[i].DestinationParameterIndex < 0 || outputs[i].DestinationParameterIndex >= parameterCount {
				continue
			}

			interpolatedValue := pm.previousRigOutputs[settingIndex].outputs[i]*(1-weight) + pm.currentRigOutputs[settingIndex].outputs[i]*weight

			pm.updateOutputParameterValue(
				&parameterValues[outputs[i].DestinationParameterIndex],
				parameterMinimumValues[outputs[i].DestinationParameterIndex],
				parameterMaximumValues[outputs[i].DestinationParameterIndex],
				interpolatedValue, &outputs[i],
			)
		}
	}
}

// updateParticles performs the pendulum physics simulation for a strand of particles.
// Matches the UpdateParticles function in CubismPhysics.cpp.
func (pm *PhysicsManager) updateParticles(
	strand []Particle, strandCount int,
	totalTranslation Vector2, totalAngle float32,
	windDirection Vector2, thresholdValue float32,
	deltaTimeSeconds float32, airResistance float32,
) {
	strand[0].Position = totalTranslation

	totalRadian := degreesToRadian(totalAngle)
	currentGravity := radianToDirection(totalRadian).Normalize()

	for i := 1; i < strandCount; i++ {
		strand[i].Force = currentGravity.Scale(strand[i].Acceleration).Add(windDirection)
		strand[i].LastPosition = strand[i].Position

		delay := strand[i].Delay * deltaTimeSeconds * 30.0

		direction := strand[i].Position.Sub(strand[i-1].Position)
		radian := directionToRadian(strand[i].LastGravity, currentGravity) / airResistance

		// Apply rotation (note: this uses the non-standard rotation from the official SDK
		// where X is computed first and the new X is used in computing Y)
		newDirX := cosFloat32(radian)*direction.X - direction.Y*sinFloat32(radian)
		newDirY := sinFloat32(radian)*newDirX + direction.Y*cosFloat32(radian)
		direction = Vector2{X: newDirX, Y: newDirY}

		strand[i].Position = strand[i-1].Position.Add(direction)

		velocity := strand[i].Velocity.Scale(delay)
		force := strand[i].Force.Scale(delay * delay)
		strand[i].Position = strand[i].Position.Add(velocity).Add(force)

		newDirection := strand[i].Position.Sub(strand[i-1].Position).Normalize()
		strand[i].Position = strand[i-1].Position.Add(newDirection.Scale(strand[i].Radius))

		if absFloat32(strand[i].Position.X) < thresholdValue {
			strand[i].Position.X = 0.0
		}

		if delay != 0.0 {
			strand[i].Velocity = strand[i].Position.Sub(strand[i].LastPosition).Scale(1.0 / delay).Scale(strand[i].Mobility)
		}

		strand[i].Force = Vector2{}
		strand[i].LastGravity = currentGravity
	}
}

// updateParticlesForStabilization positions particles in their stable equilibrium.
// Matches the UpdateParticlesForStabilization function in CubismPhysics.cpp.
func (pm *PhysicsManager) updateParticlesForStabilization(
	strand []Particle, strandCount int,
	totalTranslation Vector2, totalAngle float32,
	windDirection Vector2, thresholdValue float32,
) {
	strand[0].Position = totalTranslation

	totalRadian := degreesToRadian(totalAngle)
	currentGravity := radianToDirection(totalRadian).Normalize()

	for i := 1; i < strandCount; i++ {
		strand[i].Force = currentGravity.Scale(strand[i].Acceleration).Add(windDirection)
		strand[i].LastPosition = strand[i].Position
		strand[i].Velocity = Vector2{}

		force := strand[i].Force.Normalize().Scale(strand[i].Radius)
		strand[i].Position = strand[i-1].Position.Add(force)

		if absFloat32(strand[i].Position.X) < thresholdValue {
			strand[i].Position.X = 0.0
		}

		strand[i].Force = Vector2{}
		strand[i].LastGravity = currentGravity
	}
}

// normalizeParameterValue maps a parameter value from its native range to the normalized range.
// Matches NormalizeParameterValue in CubismPhysics.cpp.
func normalizeParameterValue(
	value float32,
	parameterMinimum, parameterMaximum, parameterDefault float32,
	normalizedMinimum, normalizedMaximum, normalizedDefault float32,
	isInverted bool,
) float32 {
	maxValue := maxFloat32(parameterMaximum, parameterMinimum)
	if maxValue < value {
		value = maxValue
	}
	minValue := minFloat32(parameterMaximum, parameterMinimum)
	if minValue > value {
		value = minValue
	}

	minNormValue := minFloat32(normalizedMinimum, normalizedMaximum)
	maxNormValue := maxFloat32(normalizedMinimum, normalizedMaximum)
	middleNormValue := normalizedDefault
	middleValue := getDefaultValue(minValue, maxValue)
	paramValue := value - middleValue

	var result float32

	sign := signFloat32(paramValue)
	switch sign {
	case 1:
		nLength := maxNormValue - middleNormValue
		pLength := maxValue - middleValue
		if pLength != 0 {
			result = paramValue*(nLength/pLength) + middleNormValue
		}
	case -1:
		nLength := minNormValue - middleNormValue
		pLength := minValue - middleValue
		if pLength != 0 {
			result = paramValue*(nLength/pLength) + middleNormValue
		}
	case 0:
		result = middleNormValue
	}

	if isInverted {
		return result
	}
	return result * -1.0
}

// getInputNormalizedParameterValue applies input normalization based on the input type.
// Matches the function pointer dispatch in the official SDK.
func (pm *PhysicsManager) getInputNormalizedParameterValue(
	targetTranslation *Vector2, targetAngle *float32,
	value, parameterMinimum, parameterMaximum, parameterDefault float32,
	normalizationPosition, normalizationAngle *Normalization,
	input Input, weight float32,
) {
	var normalizedValue float32
	switch input.Type {
	case PhysicsSourceX:
		normalizedValue = normalizeParameterValue(
			value, parameterMinimum, parameterMaximum, parameterDefault,
			normalizationPosition.Minimum, normalizationPosition.Maximum, normalizationPosition.Default,
			input.Reflect,
		)
		targetTranslation.X += normalizedValue * weight
	case PhysicsSourceY:
		normalizedValue = normalizeParameterValue(
			value, parameterMinimum, parameterMaximum, parameterDefault,
			normalizationPosition.Minimum, normalizationPosition.Maximum, normalizationPosition.Default,
			input.Reflect,
		)
		targetTranslation.Y += normalizedValue * weight
	case PhysicsSourceAngle:
		normalizedValue = normalizeParameterValue(
			value, parameterMinimum, parameterMaximum, parameterDefault,
			normalizationAngle.Minimum, normalizationAngle.Maximum, normalizationAngle.Default,
			input.Reflect,
		)
		*targetAngle += normalizedValue * weight
	}
}

// getOutputValue extracts the output value from particle positions based on the output type.
// Matches the function pointer dispatch (GetValue) in the official SDK.
func (pm *PhysicsManager) getOutputValue(
	translation Vector2, particles []Particle, particleIndex int,
	output Output, parentGravity Vector2,
) float32 {
	var outputValue float32

	switch output.Type {
	case PhysicsSourceX:
		outputValue = translation.X
		if output.Reflect {
			outputValue *= -1.0
		}
	case PhysicsSourceY:
		outputValue = translation.Y
		if output.Reflect {
			outputValue *= -1.0
		}
	case PhysicsSourceAngle:
		if particleIndex >= 2 {
			parentGravity = particles[particleIndex-1].Position.Sub(particles[particleIndex-2].Position)
		} else {
			parentGravity = parentGravity.Scale(-1.0)
		}
		outputValue = directionToRadian(parentGravity, translation)
		if output.Reflect {
			outputValue *= -1.0
		}
	}

	return outputValue
}

// updateOutputParameterValue applies the physics output value to a model parameter.
// Matches UpdateOutputParameterValue in CubismPhysics.cpp.
func (pm *PhysicsManager) updateOutputParameterValue(
	parameterValue *float32,
	parameterValueMinimum, parameterValueMaximum float32,
	translation float32,
	output *Output,
) {
	var outputScale float32
	switch output.Type {
	case PhysicsSourceX:
		outputScale = output.TranslationScale.X
	case PhysicsSourceY:
		outputScale = output.TranslationScale.Y
	case PhysicsSourceAngle:
		outputScale = output.AngleScale
	}

	value := translation * outputScale

	if value < parameterValueMinimum {
		if value < output.ValueBelowMinimum {
			output.ValueBelowMinimum = value
		}
		value = parameterValueMinimum
	} else if value > parameterValueMaximum {
		if value > output.ValueExceededMaximum {
			output.ValueExceededMaximum = value
		}
		value = parameterValueMaximum
	}

	weight := output.Weight / MaximumWeight

	if weight >= 1.0 {
		*parameterValue = value
	} else {
		*parameterValue = *parameterValue*(1.0-weight) + value*weight
	}
}

// --- Parameter access helpers ---

func (pm *PhysicsManager) getParameterCount() int {
	return pm.core.GetParameterCount(pm.modelPtr)
}

func (pm *PhysicsManager) getParameterValues() []float32 {
	return pm.core.GetParameterValues(pm.modelPtr)
}

func (pm *PhysicsManager) getParameterMinimumValues() []float32 {
	return pm.core.GetParameterMinimumValues(pm.modelPtr)
}

func (pm *PhysicsManager) getParameterMaximumValues() []float32 {
	return pm.core.GetParameterMaximumValues(pm.modelPtr)
}

func (pm *PhysicsManager) getParameterDefaultValues() []float32 {
	return pm.core.GetParameterDefaultValues(pm.modelPtr)
}

func (pm *PhysicsManager) getParameterIndex(id string) int {
	if pm.idManager != nil {
		handle := pm.idManager.GetParameterId(id)
		if handle.IsValid() {
			return int(handle)
		}
		return -1
	}
	// Fallback: linear search via Core
	parameters := pm.core.GetParameters(pm.modelPtr)
	for i, p := range parameters {
		if p.Id == id {
			return i
		}
	}
	return -1
}

// --- Utility functions ---

func getRangeValue(min, max float32) float32 {
	return absFloat32(maxFloat32(min, max) - minFloat32(min, max))
}

func getDefaultValue(min, max float32) float32 {
	return minFloat32(min, max) + getRangeValue(min, max)/2.0
}

func signFloat32(value float32) int {
	if value > 0 {
		return 1
	} else if value < 0 {
		return -1
	}
	return 0
}
