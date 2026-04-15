package expression

import (
	"github.com/shaolei/cubism-go/internal/core"
	"github.com/shaolei/cubism-go/internal/id"
	"github.com/shaolei/cubism-go/internal/model"
	"github.com/shaolei/cubism-go/internal/motion"
)

// BlendMode defines how expression parameters are mixed with current values
type BlendMode int

const (
	BlendAdd       BlendMode = 0 // Add to current value
	BlendMultiply  BlendMode = 1 // Multiply with current value
	BlendOverwrite BlendMode = 2 // Directly set the value (overwrite)
)

// ExpressionParameter represents a single parameter override in an expression
type ExpressionParameter struct {
	Id    string
	Value float32
	Blend BlendMode
}

// ExpressionParameterValue tracks the three blend values for a single parameter
// across multiple active expressions. This matches the official SDK's design
// where each expression contributes to AdditiveValue, MultiplyValue, or OverwriteValue.
type ExpressionParameterValue struct {
	ParameterId    string
	AdditiveValue  float32 // Accumulated Add contributions
	MultiplyValue  float32 // Accumulated Multiply contributions
	OverwriteValue float32 // Latest Overwrite contribution
}

// CubismExpressionMotion represents an expression as a motion-like object
// that can be managed by the queue system. It extends the base Motion concept
// with expression-specific three-value tracking.
type CubismExpressionMotion struct {
	Motion     motion.Motion
	Parameters []ExpressionParameter
	Name       string
}

// NewCubismExpressionMotion creates a CubismExpressionMotion from parsed expression JSON
func NewCubismExpressionMotion(exp model.ExpJson) *CubismExpressionMotion {
	m := &CubismExpressionMotion{
		Name: exp.Name,
		Motion: motion.Motion{
			FadeInTime:  exp.FadeInTime,
			FadeOutTime: exp.FadeOutTime,
			// Expressions have infinite duration — they persist until replaced
			Meta: motion.Meta{
				Duration: -1.0,
			},
		},
		Parameters: make([]ExpressionParameter, len(exp.Parameters)),
	}

	// Default fade times if not specified in the file
	if m.Motion.FadeInTime == 0 {
		m.Motion.FadeInTime = 1.0
	}
	if m.Motion.FadeOutTime == 0 {
		m.Motion.FadeOutTime = 1.0
	}

	for i, p := range exp.Parameters {
		m.Parameters[i] = ExpressionParameter{
			Id:    p.Id,
			Value: float32(p.Value),
			Blend: blendModeFromString(p.Blend),
		}
	}

	return m
}

// blendModeFromString converts a JSON blend string to BlendMode enum
func blendModeFromString(s string) BlendMode {
	switch s {
	case "Add":
		return BlendAdd
	case "Multiply":
		return BlendMultiply
	case "Overwrite":
		return BlendOverwrite
	default:
		return BlendAdd // Default per official SDK
	}
}

// expressionQueueEntry tracks a single expression in the playback queue
type expressionQueueEntry struct {
	expression *CubismExpressionMotion
	// Lifecycle
	available bool
	finished  bool
	started   bool
	// Timing (global seconds)
	startTimeSeconds       float64
	fadeInStartTimeSeconds float64
	// Fade-out
	isTriggeredFadeOut bool
	fadeOutSeconds     float64
	// Computed
	fadeWeight float64
}

func newExpressionQueueEntry(exp *CubismExpressionMotion) *expressionQueueEntry {
	return &expressionQueueEntry{
		expression: exp,
		available:  true,
	}
}

// Start initializes the entry timing
func (e *expressionQueueEntry) Start(userTimeSeconds float64) {
	e.started = true
	e.startTimeSeconds = userTimeSeconds
	e.fadeInStartTimeSeconds = userTimeSeconds
}

// StartFadeout triggers an explicit fade-out
func (e *expressionQueueEntry) StartFadeout(fadeOutSeconds float64) {
	e.isTriggeredFadeOut = true
	e.fadeOutSeconds = fadeOutSeconds
}

// UpdateFadeWeight calculates the current fade weight
func (e *expressionQueueEntry) UpdateFadeWeight(userTimeSeconds float64) {
	var fadeInWeight float64 = 1.0
	var fadeOutWeight float64 = 1.0

	// Calculate fade-in weight
	if e.expression.Motion.FadeInTime > 0.0 {
		fadeInWeight = motionGetEasingSine((userTimeSeconds - e.fadeInStartTimeSeconds) / e.expression.Motion.FadeInTime)
	}

	// Calculate fade-out weight
	if e.isTriggeredFadeOut && e.fadeOutSeconds > 0.0 {
		elapsed := userTimeSeconds - e.startTimeSeconds - e.expression.Motion.Meta.Duration
		fadeOutWeight = motionGetEasingSine(1.0 - elapsed/e.fadeOutSeconds)
	}

	e.fadeWeight = fadeInWeight * fadeOutWeight
}

// CubismExpressionMotionManager manages expression playback with three-value tracking,
// matching the official SDK's CubismExpressionMotionManager design.
type CubismExpressionMotionManager struct {
	core       core.Core
	modelPtr   uintptr
	idManager  *id.CubismIdManager
	expressions map[string]*CubismExpressionMotion
	entries    []*expressionQueueEntry
	paramValues []ExpressionParameterValue
	userTime   float64
}

// NewCubismExpressionMotionManager creates a new expression manager
func NewCubismExpressionMotionManager(exps []model.ExpJson) *CubismExpressionMotionManager {
	em := &CubismExpressionMotionManager{
		expressions: make(map[string]*CubismExpressionMotion),
		entries:     make([]*expressionQueueEntry, 0),
		paramValues: make([]ExpressionParameterValue, 0),
	}

	for _, exp := range exps {
		em.expressions[exp.Name] = NewCubismExpressionMotion(exp)
	}

	return em
}

// InitWithCore initializes the manager with core access (lazy init)
func (em *CubismExpressionMotionManager) InitWithCore(c core.Core, modelPtr uintptr) {
	em.core = c
	em.modelPtr = modelPtr
}

// SetIdManager sets the CubismIdManager for fast parameter access by index
func (em *CubismExpressionMotionManager) SetIdManager(idMgr *id.CubismIdManager) {
	em.idManager = idMgr
}

// PlayExpression starts playing an expression by name
func (em *CubismExpressionMotionManager) PlayExpression(name string) {
	exp, ok := em.expressions[name]
	if !ok {
		return
	}

	// Trigger fade-out on all existing entries
	for _, entry := range em.entries {
		if entry.available && !entry.finished {
			fadeOutSeconds := exp.Motion.FadeOutTime
			if fadeOutSeconds <= 0.0 {
				fadeOutSeconds = 0.0
			}
			entry.StartFadeout(fadeOutSeconds)
		}
	}

	// Add new expression entry
	em.entries = append(em.entries, newExpressionQueueEntry(exp))
}

// StopExpression stops all expressions with fade-out
func (em *CubismExpressionMotionManager) StopExpression() {
	for _, entry := range em.entries {
		if entry.available && !entry.finished {
			fadeOutSeconds := entry.expression.Motion.FadeOutTime
			if fadeOutSeconds <= 0.0 {
				fadeOutSeconds = 0.0
			}
			entry.StartFadeout(fadeOutSeconds)
		}
	}
}

// GetCurrentExpression returns the name of the currently playing expression
func (em *CubismExpressionMotionManager) GetCurrentExpression() string {
	// Find the latest active expression entry
	for i := len(em.entries) - 1; i >= 0; i-- {
		entry := em.entries[i]
		if entry.available && !entry.finished {
			return entry.expression.Name
		}
	}
	return ""
}

// GetExpressionNames returns all available expression names
func (em *CubismExpressionMotionManager) GetExpressionNames() []string {
	names := make([]string, 0, len(em.expressions))
	for name := range em.expressions {
		names = append(names, name)
	}
	return names
}

// Update applies expression parameters using three-value tracking.
// The final parameter value is calculated as:
//
//	(OverwriteValue + AdditiveValue) * MultiplyValue
//
// where each contribution is scaled by the expression's fade weight.
func (em *CubismExpressionMotionManager) Update(c core.Core, modelPtr uintptr, deltaTime float64) {
	if em.core == nil {
		em.core = c
		em.modelPtr = modelPtr
	}

	em.userTime += deltaTime

	if len(em.entries) == 0 {
		return
	}

	// Process each expression entry
	activeEntries := make([]*expressionQueueEntry, 0, len(em.entries))
	for _, entry := range em.entries {
		if !entry.available {
			continue
		}

		// Start entry on first update
		if !entry.started {
			entry.Start(em.userTime)
		}

		// Update fade weight
		entry.UpdateFadeWeight(em.userTime)

		// Check if fade-out completed
		if entry.isTriggeredFadeOut && entry.fadeWeight <= 0.0 {
			entry.finished = true
			continue
		}

		if !entry.finished {
			activeEntries = append(activeEntries, entry)
		}
	}
	em.entries = activeEntries

	if len(em.entries) == 0 {
		return
	}

	// Reset parameter value tracking
	em.resetParamValues()

	// Process each active expression
	for _, entry := range em.entries {
		fadeWeight := entry.fadeWeight
		exp := entry.expression

		for _, param := range exp.Parameters {
			pv := em.getOrCreateParamValue(param.Id)
			switch param.Blend {
			case BlendAdd:
				pv.AdditiveValue += param.Value * float32(fadeWeight)
			case BlendMultiply:
				pv.MultiplyValue *= 1.0 + (param.Value-1.0)*float32(fadeWeight)
			case BlendOverwrite:
				pv.OverwriteValue = pv.OverwriteValue + (param.Value-pv.OverwriteValue)*float32(fadeWeight)
			}
		}
	}

	// Apply the final blended values to the model
	for i := range em.paramValues {
		pv := &em.paramValues[i]
		finalValue := (pv.OverwriteValue + pv.AdditiveValue) * pv.MultiplyValue
		if em.idManager != nil {
			handle := em.idManager.GetParameterId(pv.ParameterId)
			if handle.IsValid() {
				c.SetParameterValueByIndex(modelPtr, int(handle), finalValue)
			}
		} else {
			c.SetParameterValue(modelPtr, pv.ParameterId, finalValue)
		}
	}
}

// resetParamValues clears all parameter value tracking for this frame
func (em *CubismExpressionMotionManager) resetParamValues() {
	for i := range em.paramValues {
		pv := &em.paramValues[i]
		pv.AdditiveValue = 0.0
		pv.MultiplyValue = 1.0
		pv.OverwriteValue = 0.0
	}
}

// getOrCreateParamValue finds or creates an ExpressionParameterValue for the given parameter ID
func (em *CubismExpressionMotionManager) getOrCreateParamValue(id string) *ExpressionParameterValue {
	for i := range em.paramValues {
		if em.paramValues[i].ParameterId == id {
			return &em.paramValues[i]
		}
	}
	em.paramValues = append(em.paramValues, ExpressionParameterValue{
		ParameterId:    id,
		AdditiveValue:  0.0,
		MultiplyValue:  1.0,
		OverwriteValue: 0.0,
	})
	return &em.paramValues[len(em.paramValues)-1]
}

// motionGetEasingSine delegates to the motion package's easing function
func motionGetEasingSine(value float64) float64 {
	return motion.GetEasingSine(value)
}
