package cubismmath

import "math"

// TargetPoint handles the direction the model is facing, providing smooth
// interpolation with acceleration and deceleration.
// Matches the official SDK's CubismTargetPoint.
type TargetPoint struct {
	faceTargetX      float32
	faceTargetY      float32
	faceX            float32
	faceY            float32
	faceVX           float32
	faceVY           float32
	lastTimeSeconds  float32
	userTimeSeconds  float32
}

const (
	frameRate      int     = 30
	tpEpsilon              = float32(0.01)
	faceParamMaxV           = float32(40.0 / 10.0)                     // Average speed for head turning
	maxVPerFrame            = faceParamMaxV * float32(1.0) / float32(frameRate) // Max velocity per frame
	timeToMaxSpeed          = float32(0.15)                              // Time to reach max speed (seconds)
)

// NewTargetPoint creates a new TargetPoint with all values initialized to zero.
func NewTargetPoint() *TargetPoint {
	return &TargetPoint{}
}

// Update advances the face direction interpolation by deltaTimeSeconds.
// Uses acceleration-based smoothing: the face accelerates toward the target
// direction and decelerates smoothly as it approaches.
func (t *TargetPoint) Update(deltaTimeSeconds float32) {
	t.userTimeSeconds += deltaTimeSeconds

	if t.lastTimeSeconds == 0.0 {
		t.lastTimeSeconds = t.userTimeSeconds
		return
	}

	deltaTimeWeight := (t.userTimeSeconds - t.lastTimeSeconds) * float32(frameRate)
	t.lastTimeSeconds = t.userTimeSeconds

	frameToMaxSpeed := timeToMaxSpeed * float32(frameRate)
	maxA := deltaTimeWeight * maxVPerFrame / frameToMaxSpeed

	// Direction to target
	dx := t.faceTargetX - t.faceX
	dy := t.faceTargetY - t.faceY

	if absF(dx) <= tpEpsilon && absF(dy) <= tpEpsilon {
		return // No significant change needed
	}

	// Distance to target
	d := sqrtF(dx*dx + dy*dy)

	// Max velocity in the direction of the target
	vx := maxVPerFrame * dx / d
	vy := maxVPerFrame * dy / d

	// Acceleration needed to change from current velocity to target velocity
	ax := vx - t.faceVX
	ay := vy - t.faceVY

	a := sqrtF(ax*ax + ay*ay)

	// Clamp acceleration
	if a < -maxA || a > maxA {
		ax *= maxA / a
		ay *= maxA / a
	}

	// Apply acceleration
	t.faceVX += ax
	t.faceVY += ay

	// Smooth deceleration as we approach the target
	// Using the relationship between acceleration, velocity, and distance:
	// maxV = 0.5 * (sqrt(maxA^2 + 16*maxA*d - 8*maxA*d) - maxA)
	maxV := 0.5 * (sqrtF(maxA*maxA+16.0*maxA*d-8.0*maxA*d) - maxA)
	curV := sqrtF(t.faceVX*t.faceVX + t.faceVY*t.faceVY)

	if curV > maxV {
		t.faceVX *= maxV / curV
		t.faceVY *= maxV / curV
	}

	t.faceX += t.faceVX
	t.faceY += t.faceVY
}

// Set sets the target direction. X and Y should be in the range [-1.0, 1.0].
func (t *TargetPoint) Set(x, y float32) {
	t.faceTargetX = x
	t.faceTargetY = y
}

// GetX returns the current X direction value in the range [-1.0, 1.0].
func (t *TargetPoint) GetX() float32 { return t.faceX }

// GetY returns the current Y direction value in the range [-1.0, 1.0].
func (t *TargetPoint) GetY() float32 { return t.faceY }

// sqrtF returns the square root of a float32.
func sqrtF(x float32) float32 {
	return float32(math.Sqrt(float64(x)))
}
