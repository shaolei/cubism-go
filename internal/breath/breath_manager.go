package breath

import (
	"math"

	"github.com/shaolei/cubism-go/internal/core"
	"github.com/shaolei/cubism-go/internal/id"
)

// BreathParameterData defines the breathing parameters for a single model parameter.
// The cyclic motion of breathing is set entirely by sine waves:
//
//	value = Offset + Peak * sin(currentTime * 2π / Cycle)
//
// The value is then applied to the parameter via AddParameterValue with the given Weight.
type BreathParameterData struct {
	ParameterId string  // ID of the parameter to attach
	Offset      float32 // Offset of the sine wave
	Peak        float32 // Peak value of the sine wave
	Cycle       float32 // Cycle of the sine wave (in seconds)
	Weight      float32 // Weight of the parameter application
}

// DefaultBreathParameters returns the standard breathing parameter set
// matching the official SDK's default configuration.
func DefaultBreathParameters() []BreathParameterData {
	return []BreathParameterData{
		{ParameterId: "ParamAngleX", Offset: 0.0, Peak: 15.0, Cycle: 6.5345, Weight: 1.0},
		{ParameterId: "ParamAngleY", Offset: 0.0, Peak: 8.0, Cycle: 3.5345, Weight: 1.0},
		{ParameterId: "ParamAngleZ", Offset: 0.0, Peak: 10.0, Cycle: 5.5345, Weight: 1.0},
		{ParameterId: "ParamBodyAngleX", Offset: 0.0, Peak: 4.0, Cycle: 15.5345, Weight: 1.0},
		{ParameterId: "ParamBreath", Offset: 0.0, Peak: 0.5, Cycle: 3.2345, Weight: 1.0},
	}
}

// BreathManager manages the breathing effect.
// Matches the official SDK's CubismBreath design.
type BreathManager struct {
	parameters []BreathParameterData
	currentTime float64
	core       core.Core
	modelPtr   uintptr
	idManager  *id.CubismIdManager
}

// NewBreathManager creates a new breathing effect manager.
func NewBreathManager(c core.Core, modelPtr uintptr) *BreathManager {
	return &BreathManager{
		parameters: DefaultBreathParameters(),
		core:       c,
		modelPtr:   modelPtr,
	}
}

// SetParameters sets the breathing parameter data collection.
func (b *BreathManager) SetParameters(parameters []BreathParameterData) {
	b.parameters = parameters
}

// GetParameters returns the current breathing parameter data collection.
func (b *BreathManager) GetParameters() []BreathParameterData {
	return b.parameters
}

// SetIdManager sets the CubismIdManager for fast parameter access by index
func (b *BreathManager) SetIdManager(idMgr *id.CubismIdManager) {
	b.idManager = idMgr
}

// Update applies the breathing effect to the model parameters.
// Matches the official SDK's CubismBreath::UpdateParameters:
//
//	_currentTime += deltaTimeSeconds;
//	t = _currentTime * 2.0f * Pi;
//	for each parameter:
//	    AddParameterValue(id, Offset + Peak * sin(t / Cycle), Weight)
func (b *BreathManager) Update(deltaTimeSeconds float64) {
	b.currentTime += deltaTimeSeconds

	t := b.currentTime * 2.0 * math.Pi

	for i := range b.parameters {
		data := &b.parameters[i]
		value := data.Offset + data.Peak*float32(math.Sin(t/float64(data.Cycle)))

		if b.idManager != nil {
			handle := b.idManager.GetParameterId(data.ParameterId)
			if !handle.IsValid() {
				continue
			}
			index := int(handle)
			current := b.core.GetParameterValueByIndex(b.modelPtr, index)
			b.core.SetParameterValueByIndex(b.modelPtr, index, current+value*data.Weight)
		} else {
			current := b.core.GetParameterValue(b.modelPtr, data.ParameterId)
			b.core.SetParameterValue(b.modelPtr, data.ParameterId, current+value*data.Weight)
		}
	}
}
