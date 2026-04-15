package physics

import "github.com/shaolei/cubism-go/internal/model"

// ParsePhysicsJson converts the existing PhysicsJson structure into the physics engine's
// internal Rig representation. This matches the official SDK's CubismPhysics::Parse logic.
func ParsePhysicsJson(json model.PhysicsJson) *Rig {
	rig := &Rig{
		SubRigCount: json.Meta.PhysicsSettingCount,
		Settings:    make([]SubRig, json.Meta.PhysicsSettingCount),
		Inputs:      make([]Input, json.Meta.TotalInputCount),
		Outputs:     make([]Output, json.Meta.TotalOutputCount),
		Particles:   make([]Particle, json.Meta.VertexCount),
		Gravity: Vector2{
			X: float32(json.Meta.EffectiveForces.Gravity.X),
			Y: float32(json.Meta.EffectiveForces.Gravity.Y),
		},
		Wind: Vector2{
			X: float32(json.Meta.EffectiveForces.Wind.X),
			Y: float32(json.Meta.EffectiveForces.Wind.Y),
		},
		Fps: float32(json.Meta.Fps),
	}

	inputIndex := 0
	outputIndex := 0
	particleIndex := 0

	for i := 0; i < rig.SubRigCount; i++ {
		setting := &json.PhysicsSettings[i]

		// Normalization
		rig.Settings[i].NormalizationPosition = Normalization{
			Minimum: float32(setting.Normalization.Position.Minimum),
			Maximum: float32(setting.Normalization.Position.Maximum),
			Default: float32(setting.Normalization.Position.Default),
		}
		rig.Settings[i].NormalizationAngle = Normalization{
			Minimum: float32(setting.Normalization.Angle.Minimum),
			Maximum: float32(setting.Normalization.Angle.Maximum),
			Default: float32(setting.Normalization.Angle.Default),
		}

		// Inputs
		rig.Settings[i].InputCount = len(setting.Input)
		rig.Settings[i].BaseInputIndex = inputIndex
		for j := 0; j < len(setting.Input); j++ {
			input := &setting.Input[j]
			rig.Inputs[inputIndex+j].SourceId = input.Source.Id
			rig.Inputs[inputIndex+j].SourceParameterIndex = -1 // resolved later
			rig.Inputs[inputIndex+j].Weight = float32(input.Weight)
			rig.Inputs[inputIndex+j].Reflect = input.Reflect

			switch input.Type {
			case "X":
				rig.Inputs[inputIndex+j].Type = PhysicsSourceX
			case "Y":
				rig.Inputs[inputIndex+j].Type = PhysicsSourceY
			case "Angle":
				rig.Inputs[inputIndex+j].Type = PhysicsSourceAngle
			}
		}
		inputIndex += len(setting.Input)

		// Outputs
		rig.Settings[i].OutputCount = len(setting.Output)
		rig.Settings[i].BaseOutputIndex = outputIndex
		for j := 0; j < len(setting.Output); j++ {
			output := &setting.Output[j]
			rig.Outputs[outputIndex+j].DestinationId = output.Destination.Id
			rig.Outputs[outputIndex+j].DestinationParameterIndex = -1 // resolved later
			rig.Outputs[outputIndex+j].VertexIndex = output.VertexIndex
			rig.Outputs[outputIndex+j].AngleScale = float32(output.Scale)
			rig.Outputs[outputIndex+j].Weight = float32(output.Weight)
			rig.Outputs[outputIndex+j].Reflect = output.Reflect

			switch output.Type {
			case "X":
				rig.Outputs[outputIndex+j].Type = PhysicsSourceX
				rig.Outputs[outputIndex+j].TranslationScale = Vector2{X: float32(output.Scale), Y: 0}
			case "Y":
				rig.Outputs[outputIndex+j].Type = PhysicsSourceY
				rig.Outputs[outputIndex+j].TranslationScale = Vector2{X: 0, Y: float32(output.Scale)}
			case "Angle":
				rig.Outputs[outputIndex+j].Type = PhysicsSourceAngle
				rig.Outputs[outputIndex+j].TranslationScale = Vector2{}
			}
		}
		outputIndex += len(setting.Output)

		// Particles
		rig.Settings[i].ParticleCount = len(setting.Vertices)
		rig.Settings[i].BaseParticleIndex = particleIndex
		for j := 0; j < len(setting.Vertices); j++ {
			vertex := &setting.Vertices[j]
			rig.Particles[particleIndex+j].Mobility = float32(vertex.Mobility)
			rig.Particles[particleIndex+j].Delay = float32(vertex.Delay)
			rig.Particles[particleIndex+j].Acceleration = float32(vertex.Acceleration)
			rig.Particles[particleIndex+j].Radius = float32(vertex.Radius)
			rig.Particles[particleIndex+j].Position = Vector2{
				X: float32(vertex.Position.X),
				Y: float32(vertex.Position.Y),
			}
		}
		particleIndex += len(setting.Vertices)
	}

	return rig
}
