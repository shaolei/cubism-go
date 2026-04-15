package motion

// CubismMotionQueueEntry represents a single motion in the playback queue
// with full lifecycle and fade state tracking, matching the official SDK's
// CubismMotionQueueEntry design.
type CubismMotionQueueEntry struct {
	motion Motion
	id     int
	loop   bool // Override loop behavior

	// Lifecycle state
	available bool // Whether this entry is valid
	finished  bool // Whether playback has completed
	started   bool // Whether playback has begun

	// Timing (in global/accumulated seconds)
	startTimeSeconds       float64 // Global time when entry started playing
	fadeInStartTimeSeconds float64 // Global time when fade-in started
	endTimeSeconds         float64 // Global time when entry should end

	// Fade-out state
	isTriggeredFadeOut bool    // Whether fade-out has been explicitly triggered
	fadeOutSeconds     float64 // Fade-out duration (negative = use motion default)

	// Computed fade weight
	fadeWeight float64
}

func newCubismMotionQueueEntry(mtn Motion, id int, loop bool) *CubismMotionQueueEntry {
	return &CubismMotionQueueEntry{
		motion:          mtn,
		id:              id,
		loop:            loop,
		available:       true,
		fadeOutSeconds:  -1.0,
	}
}

// IsAvailable returns whether this entry is still valid
func (e *CubismMotionQueueEntry) IsAvailable() bool { return e.available }

// IsFinished returns whether playback has completed
func (e *CubismMotionQueueEntry) IsFinished() bool { return e.finished }

// IsStarted returns whether playback has begun
func (e *CubismMotionQueueEntry) IsStarted() bool { return e.started }

// GetId returns the motion queue entry ID
func (e *CubismMotionQueueEntry) GetId() int { return e.id }

// GetMotion returns a pointer to the motion data
func (e *CubismMotionQueueEntry) GetMotion() *Motion { return &e.motion }

// GetFadeWeight returns the current computed fade weight
func (e *CubismMotionQueueEntry) GetFadeWeight() float64 { return e.fadeWeight }

// IsLoop returns whether this entry should loop
func (e *CubismMotionQueueEntry) IsLoop() bool {
	if e.loop {
		return true
	}
	return e.motion.Meta.Loop
}

// GetLocalTime returns the time elapsed since this entry started playing
func (e *CubismMotionQueueEntry) GetLocalTime(userTimeSeconds float64) float64 {
	return userTimeSeconds - e.startTimeSeconds
}

// StartFadeout triggers an explicit fade-out on this entry starting from the given global time
func (e *CubismMotionQueueEntry) StartFadeout(fadeOutSeconds float64, userTimeSeconds float64) {
	e.isTriggeredFadeOut = true
	e.fadeOutSeconds = fadeOutSeconds
	e.endTimeSeconds = userTimeSeconds + fadeOutSeconds
}

// SetFadeOut sets the fade-out duration without triggering it immediately
func (e *CubismMotionQueueEntry) SetFadeOut(fadeOutSeconds float64) {
	e.fadeOutSeconds = fadeOutSeconds
}

// Start initializes the entry timing on first update
func (e *CubismMotionQueueEntry) Start(userTimeSeconds float64) {
	e.started = true
	e.startTimeSeconds = userTimeSeconds
	e.fadeInStartTimeSeconds = userTimeSeconds

	duration := e.motion.Meta.Duration
	if duration > 0.0 {
		e.endTimeSeconds = userTimeSeconds + duration
	} else {
		e.endTimeSeconds = 1e30 // Effectively infinite
	}
}

// Restart resets timing for looping motions
func (e *CubismMotionQueueEntry) Restart(userTimeSeconds float64) {
	e.startTimeSeconds = userTimeSeconds
	e.fadeInStartTimeSeconds = userTimeSeconds

	duration := e.motion.Meta.Duration
	if duration > 0.0 {
		e.endTimeSeconds = userTimeSeconds + duration
	}
}

// UpdateFadeWeight calculates the current fade weight based on global time.
// Uses (userTimeSeconds - fadeInStartTimeSeconds) / fadeInSeconds for fade-in,
// matching the official SDK's calculation rather than the previous currentTime/fadeInTime approach.
func (e *CubismMotionQueueEntry) UpdateFadeWeight(userTimeSeconds float64) {
	var fadeInWeight float64 = 1.0
	var fadeOutWeight float64 = 1.0

	// Calculate fade-in weight
	if e.motion.FadeInTime > 0.0 {
		fadeInWeight = getEasingSine((userTimeSeconds - e.fadeInStartTimeSeconds) / e.motion.FadeInTime)
	}

	// Calculate fade-out weight
	if e.isTriggeredFadeOut {
		if e.fadeOutSeconds > 0.0 {
			fadeOutWeight = getEasingSine((e.endTimeSeconds - userTimeSeconds) / e.fadeOutSeconds)
		}
	} else if e.motion.FadeOutTime > 0.0 && e.motion.Meta.Duration > 0.0 {
		// Natural fade-out at end of motion
		fadeOutWeight = getEasingSine((e.endTimeSeconds - userTimeSeconds) / e.motion.FadeOutTime)
	}

	e.fadeWeight = fadeInWeight * fadeOutWeight
}
