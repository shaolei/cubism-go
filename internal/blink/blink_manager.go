package blink

import (
	"math/rand"

	"github.com/shaolei/cubism-go/internal/core"
	"github.com/shaolei/cubism-go/internal/id"
)

const (
	EyeStateFirst    = iota ///< Initial state
	EyeStateInterval        ///< State where the eyes are not blinking
	EyeStateClosing         ///< State where the eyelids are closing
	EyeStateClosed          ///< State where the eyelids are closed
	EyeStateOpening         ///< State where the eyelids are opening
)

type BlinkManager struct {
	core             core.Core
	modelPtr         uintptr
	idManager        *id.CubismIdManager
	ids              []string
	idHandles        []id.CubismIdHandle
	state            int
	interval         float64
	closing          float64
	opening          float64
	currentTime      float64
	stateStartTime   float64
	nextBlinkingTime float64
}

func NewBlinkManager(core core.Core, modelPtr uintptr, ids []string) *BlinkManager {
	return &BlinkManager{
		core:             core,
		modelPtr:         modelPtr,
		ids:              ids,
		state:            EyeStateFirst,
		interval:         4.0,
		closing:          0.1,
		opening:          0.15,
		currentTime:      0,
		stateStartTime:   0,
		nextBlinkingTime: 0,
	}
}

// SetIdManager sets the CubismIdManager for fast parameter access by index
func (b *BlinkManager) SetIdManager(idMgr *id.CubismIdManager) {
	b.idManager = idMgr
	// Pre-resolve all blink parameter IDs to handles
	b.idHandles = make([]id.CubismIdHandle, len(b.ids))
	for i, idStr := range b.ids {
		if idMgr != nil {
			b.idHandles[i] = idMgr.GetParameterId(idStr)
		} else {
			b.idHandles[i] = id.InvalidHandle
		}
	}
}

func (b *BlinkManager) DetermineNextBlinkingTiming() float64 {
	r := rand.Float64()
	return b.currentTime + (r * (2.0*b.interval - 1.0))
}

func (b *BlinkManager) Update(delta float64) {
	b.currentTime += delta

	var value float32

	switch b.state {
	case EyeStateFirst:
		b.state = EyeStateInterval
		b.nextBlinkingTime = b.DetermineNextBlinkingTiming()
		value = 1.0
	case EyeStateInterval:
		if b.currentTime >= b.nextBlinkingTime {
			b.state = EyeStateClosing
			b.stateStartTime = b.currentTime
		}
		value = 1.0
	case EyeStateClosing:
		t := (b.currentTime - b.stateStartTime) / b.closing
		if t >= 1 {
			b.state = EyeStateClosed
			b.stateStartTime = b.currentTime
		}
		value = 1.0 - float32(t)
	case EyeStateClosed:
		t := (b.currentTime - b.stateStartTime) / b.closing
		if t >= 1 {
			b.state = EyeStateOpening
			b.stateStartTime = b.currentTime
		}
		value = 0.0
	case EyeStateOpening:
		t := (b.currentTime - b.stateStartTime) / b.opening
		if t >= 1 {
			t = 1
			b.state = EyeStateInterval
			b.nextBlinkingTime = b.DetermineNextBlinkingTiming()
		}
		value = float32(t)
	}

	for i, id := range b.ids {
		if b.idManager != nil && b.idHandles[i].IsValid() {
			b.core.SetParameterValueByIndex(b.modelPtr, int(b.idHandles[i]), value)
		} else {
			b.core.SetParameterValue(b.modelPtr, id, value)
		}
	}
}
