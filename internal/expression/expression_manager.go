package expression

import (
	"github.com/shaolei/cubism-go/internal/core"
	"github.com/shaolei/cubism-go/internal/model"
)

// BlendMode defines how expression parameters are mixed with current values
type BlendMode string

const (
	BlendOverwrite BlendMode = "Overwrite" // Directly set the value
	BlendAdd       BlendMode = "Add"       // Add to current value
	BlendMultiply  BlendMode = "Multiply"  // Multiply with current value
)

// Expression holds the parsed expression data
type Expression struct {
	Name        string
	FadeInTime  float64
	FadeOutTime float64
	Parameters  []ExpressionParameter
}

// ExpressionParameter represents a single parameter override in an expression
type ExpressionParameter struct {
	Id    string
	Value float32
	Blend BlendMode
}

// ExpressionManager manages expression playback and blending
type ExpressionManager struct {
	expressions map[string]Expression
	currentName string
	fadeWeight  float64
	fadeState   int // 0: none, 1: fading in, 2: fading out
	currentTime float64
}

const (
	fadeStateNone = iota
	fadeStateIn
	fadeStateOut
)

// NewExpressionManager creates a new ExpressionManager from loaded expression data
func NewExpressionManager(exps []model.ExpJson) *ExpressionManager {
	em := &ExpressionManager{
		expressions: make(map[string]Expression),
	}

	for _, exp := range exps {
		e := Expression{
			Name:        exp.Name,
			FadeInTime:  exp.FadeInTime,
			FadeOutTime: exp.FadeOutTime,
			Parameters:  make([]ExpressionParameter, len(exp.Parameters)),
		}
		// Default fade times if not specified in the file
		if e.FadeInTime == 0 {
			e.FadeInTime = 1.0
		}
		if e.FadeOutTime == 0 {
			e.FadeOutTime = 1.0
		}
		for i, p := range exp.Parameters {
			e.Parameters[i] = ExpressionParameter{
				Id:    p.Id,
				Value: float32(p.Value),
				Blend: BlendMode(p.Blend),
			}
		}
		em.expressions[exp.Name] = e
	}

	return em
}

// PlayExpression starts playing an expression by name
func (em *ExpressionManager) PlayExpression(name string) {
	if _, ok := em.expressions[name]; !ok {
		return
	}
	em.currentName = name
	em.fadeState = fadeStateIn
	em.fadeWeight = 0.0
	em.currentTime = 0.0
}

// StopExpression stops the current expression with fade out
func (em *ExpressionManager) StopExpression() {
	if em.currentName == "" {
		return
	}
	em.fadeState = fadeStateOut
	em.fadeWeight = 1.0
	em.currentTime = 0.0
}

// GetCurrentExpression returns the name of the currently playing expression, or empty string
func (em *ExpressionManager) GetCurrentExpression() string {
	return em.currentName
}

// GetExpressionNames returns all available expression names
func (em *ExpressionManager) GetExpressionNames() []string {
	names := make([]string, 0, len(em.expressions))
	for name := range em.expressions {
		names = append(names, name)
	}
	return names
}

// Update applies the current expression parameters to the model
// Should be called after motion update but before core.Update()
func (em *ExpressionManager) Update(core core.Core, modelPtr uintptr, deltaTime float64) {
	if em.currentName == "" {
		return
	}

	exp, ok := em.expressions[em.currentName]
	if !ok {
		em.currentName = ""
		em.fadeState = fadeStateNone
		return
	}

	em.currentTime += deltaTime

	// Update fade weight based on state
	switch em.fadeState {
	case fadeStateIn:
		if exp.FadeInTime > 0 {
			em.fadeWeight = em.currentTime / exp.FadeInTime
			if em.fadeWeight >= 1.0 {
				em.fadeWeight = 1.0
				em.fadeState = fadeStateNone
			}
		} else {
			em.fadeWeight = 1.0
			em.fadeState = fadeStateNone
		}
	case fadeStateOut:
		if exp.FadeOutTime > 0 {
			em.fadeWeight = 1.0 - em.currentTime/exp.FadeOutTime
			if em.fadeWeight <= 0.0 {
				em.fadeWeight = 0.0
				em.currentName = ""
				em.fadeState = fadeStateNone
				return
			}
		} else {
			em.currentName = ""
			em.fadeState = fadeStateNone
			return
		}
	}

	// Apply expression parameters with blend mode and fade weight
	for _, param := range exp.Parameters {
		currentValue := core.GetParameterValue(modelPtr, param.Id)
		var newValue float32

		switch param.Blend {
		case BlendOverwrite:
			newValue = currentValue + (param.Value-currentValue)*float32(em.fadeWeight)
		case BlendAdd:
			newValue = currentValue + param.Value*float32(em.fadeWeight)
		case BlendMultiply:
			newValue = currentValue * (1.0 + (param.Value-1.0)*float32(em.fadeWeight))
		default:
			// Default to Overwrite
			newValue = currentValue + (param.Value-currentValue)*float32(em.fadeWeight)
		}

		core.SetParameterValue(modelPtr, param.Id, newValue)
	}
}
