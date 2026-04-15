package look

import (
	"github.com/shaolei/cubism-go/internal/core"
	"github.com/shaolei/cubism-go/internal/id"
)

// LookParameterData defines the look/follow parameters for a single model parameter.
// The parameter value is calculated as:
//
//	value = FactorX * dragX + FactorY * dragY + FactorXY * (dragX * dragY)
//
// The value is then applied via AddParameterValue (with weight=1.0).
type LookParameterData struct {
	ParameterId string  // ID of the parameter to attach
	FactorX     float32 // Coefficient for drag input along the X-axis
	FactorY     float32 // Coefficient for drag input along the Y-axis
	FactorXY    float32 // Coefficient for the combined X-Y (cross) drag input
}

// LookManager manages the look/eye-tracking effect.
// Matches the official SDK's CubismLook design.
type LookManager struct {
	parameters []LookParameterData
	core       core.Core
	modelPtr   uintptr
	idManager  *id.CubismIdManager
	// Current drag position for interpolation
	dragX float32
	dragY float32
}

// NewLookManager creates a new look/eye-tracking manager.
func NewLookManager(c core.Core, modelPtr uintptr) *LookManager {
	return &LookManager{
		core:     c,
		modelPtr: modelPtr,
	}
}

// SetParameters sets the look parameter data collection.
func (l *LookManager) SetParameters(parameters []LookParameterData) {
	l.parameters = parameters
}

// GetParameters returns the current look parameter data collection.
func (l *LookManager) GetParameters() []LookParameterData {
	return l.parameters
}

// SetIdManager sets the CubismIdManager for fast parameter access by index
func (l *LookManager) SetIdManager(idMgr *id.CubismIdManager) {
	l.idManager = idMgr
}

// SetTarget sets the current target point for look tracking.
// dragX and dragY are in the range [-1.0, 1.0], where (0, 0) is center.
func (l *LookManager) SetTarget(dragX, dragY float32) {
	l.dragX = dragX
	l.dragY = dragY
}

// GetTarget returns the current target point.
func (l *LookManager) GetTarget() (float32, float32) {
	return l.dragX, l.dragY
}

// Update applies the look/eye-tracking effect to the model parameters.
// Matches the official SDK's CubismLook::UpdateParameters:
//
//	dragXY = dragX * dragY
//	for each parameter:
//	    AddParameterValue(id, FactorX * dragX + FactorY * dragY + FactorXY * dragXY)
func (l *LookManager) Update(deltaTimeSeconds float64) {
	_ = deltaTimeSeconds // not used in official SDK's look — direct application

	dragXY := l.dragX * l.dragY

	for i := range l.parameters {
		data := &l.parameters[i]
		value := data.FactorX*l.dragX + data.FactorY*l.dragY + data.FactorXY*dragXY

		if l.idManager != nil {
			handle := l.idManager.GetParameterId(data.ParameterId)
			if !handle.IsValid() {
				continue
			}
			index := int(handle)
			current := l.core.GetParameterValueByIndex(l.modelPtr, index)
			l.core.SetParameterValueByIndex(l.modelPtr, index, current+value)
		} else {
			current := l.core.GetParameterValue(l.modelPtr, data.ParameterId)
			l.core.SetParameterValue(l.modelPtr, data.ParameterId, current+value)
		}
	}
}
