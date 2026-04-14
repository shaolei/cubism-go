package motion

import "testing"

func TestGetEasingSine(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  float64
		expect float64
	}{
		{"below zero clamped", -0.5, 0.0},
		{"at zero", 0.0, 0.0},
		{"at midpoint", 0.5, 0.5},
		{"at one", 1.0, 1.0},
		{"above one clamped", 1.5, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := getEasingSine(tt.input)
			if diff := got - tt.expect; diff > 1e-10 || diff < -1e-10 {
				t.Errorf("getEasingSine(%v) = %v, want %v", tt.input, got, tt.expect)
			}
		})
	}
}

func TestLerpPoints(t *testing.T) {
	t.Parallel()

	a := Point{Time: 0, Value: 0}
	b := Point{Time: 1, Value: 10}

	t.Run("at start", func(t *testing.T) {
		got := lerpPoints(a, b, 0)
		if got.Time != 0 || got.Value != 0 {
			t.Errorf("lerpPoints(a,b,0) = %v, want {0 0}", got)
		}
	})

	t.Run("at end", func(t *testing.T) {
		got := lerpPoints(a, b, 1)
		if got.Time != 1 || got.Value != 10 {
			t.Errorf("lerpPoints(a,b,1) = %v, want {1 10}", got)
		}
	})

	t.Run("at midpoint", func(t *testing.T) {
		got := lerpPoints(a, b, 0.5)
		if got.Time != 0.5 || got.Value != 5 {
			t.Errorf("lerpPoints(a,b,0.5) = %v, want {0.5 5}", got)
		}
	})
}

func TestSegmentIntersects(t *testing.T) {
	t.Parallel()

	t.Run("linear segment in range", func(t *testing.T) {
		s := Segment{
			Type:   Linear,
			Points: []Point{{Time: 0, Value: 0}, {Time: 1, Value: 10}},
		}
		if !segmentIntersects(s, 0.5) {
			t.Error("linear segment should intersect at 0.5")
		}
	})

	t.Run("linear segment out of range", func(t *testing.T) {
		s := Segment{
			Type:   Linear,
			Points: []Point{{Time: 0.2, Value: 0}, {Time: 0.8, Value: 10}},
		}
		if segmentIntersects(s, 0.1) {
			t.Error("linear segment should not intersect at 0.1")
		}
	})

	t.Run("bezier segment in range", func(t *testing.T) {
		s := Segment{
			Type:   Bezier,
			Points: []Point{{Time: 0}, {Time: 0.33}, {Time: 0.66}, {Time: 1}},
		}
		if !segmentIntersects(s, 0.5) {
			t.Error("bezier segment should intersect at 0.5")
		}
	})

	t.Run("stepped segment in range", func(t *testing.T) {
		s := Segment{
			Type:  Stepped,
			Points: []Point{{Time: 0, Value: 5}},
			Value: 1,
		}
		if !segmentIntersects(s, 0.5) {
			t.Error("stepped segment should intersect at 0.5")
		}
	})

	t.Run("inverse stepped segment", func(t *testing.T) {
		s := Segment{
			Type:  InverseStepped,
			Points: []Point{{Time: 1, Value: 5}},
			Value: 0,
		}
		if !segmentIntersects(s, 0.5) {
			t.Error("inverse stepped segment should intersect at 0.5")
		}
	})
}

func TestSegmentInterpolate(t *testing.T) {
	t.Parallel()

	t.Run("linear interpolation", func(t *testing.T) {
		s := Segment{
			Type:   Linear,
			Points: []Point{{Time: 0, Value: 0}, {Time: 1, Value: 10}},
		}
		got := segmentInterpolate(s, 0.5)
		if got != 5 {
			t.Errorf("linear interpolate(0.5) = %v, want 5", got)
		}
	})

	t.Run("linear clamps negative k", func(t *testing.T) {
		s := Segment{
			Type:   Linear,
			Points: []Point{{Time: 0.5, Value: 0}, {Time: 1, Value: 10}},
		}
		got := segmentInterpolate(s, 0.3)
		if got != 0 {
			t.Errorf("linear interpolate with k<0 should clamp to base value, got %v", got)
		}
	})

	t.Run("stepped returns base value", func(t *testing.T) {
		s := Segment{
			Type:   Stepped,
			Points: []Point{{Time: 0, Value: 7}},
		}
		got := segmentInterpolate(s, 0.5)
		if got != 7 {
			t.Errorf("stepped interpolate = %v, want 7", got)
		}
	})

	t.Run("inverse stepped returns base value", func(t *testing.T) {
		s := Segment{
			Type:   InverseStepped,
			Points: []Point{{Time: 1, Value: 3}},
		}
		got := segmentInterpolate(s, 0.5)
		if got != 3 {
			t.Errorf("inverse stepped interpolate = %v, want 3", got)
		}
	})

	t.Run("bezier at endpoints", func(t *testing.T) {
		s := Segment{
			Type:   Bezier,
			Points: []Point{{Time: 0, Value: 0}, {Time: 0.33, Value: 3}, {Time: 0.66, Value: 7}, {Time: 1, Value: 10}},
		}
		start := segmentInterpolate(s, 0)
		end := segmentInterpolate(s, 1)
		if start != 0 {
			t.Errorf("bezier start = %v, want 0", start)
		}
		if end != 10 {
			t.Errorf("bezier end = %v, want 10", end)
		}
	})

	t.Run("unknown type returns 0", func(t *testing.T) {
		s := Segment{Type: 99}
		got := segmentInterpolate(s, 0.5)
		if got != 0 {
			t.Errorf("unknown type should return 0, got %v", got)
		}
	})
}

func TestGetFade(t *testing.T) {
	t.Parallel()

	t.Run("no fade times returns full weight", func(t *testing.T) {
		m := Motion{FadeInTime: 0, FadeOutTime: 0, Meta: Meta{Duration: 5}}
		_, _, fadeWeight := getFade(m, 1.0, 2.5)
		if fadeWeight != 1.0 {
			t.Errorf("no fade: weight = %v, want 1.0", fadeWeight)
		}
	})

	t.Run("with fade in", func(t *testing.T) {
		m := Motion{FadeInTime: 1.0, FadeOutTime: 0, Meta: Meta{Duration: 5}}
		fadeIn, _, _ := getFade(m, 1.0, 0.5)
		if fadeIn >= 1.0 {
			t.Errorf("fade in at t=0.5 with fadeInTime=1.0 should be < 1.0, got %v", fadeIn)
		}
	})

	t.Run("with fade out", func(t *testing.T) {
		m := Motion{FadeInTime: 0, FadeOutTime: 1.0, Meta: Meta{Duration: 5}}
		_, fadeOut, _ := getFade(m, 1.0, 4.5)
		if fadeOut >= 1.0 {
			t.Errorf("fade out near end should be < 1.0, got %v", fadeOut)
		}
	})

	t.Run("negative duration skips fade out", func(t *testing.T) {
		m := Motion{FadeInTime: 0, FadeOutTime: 1.0, Meta: Meta{Duration: -1}}
		_, fadeOut, _ := getFade(m, 1.0, 5.0)
		if fadeOut != 1.0 {
			t.Errorf("negative duration should skip fade out, got %v", fadeOut)
		}
	})
}

func TestEntryUpdate(t *testing.T) {
	t.Parallel()

	t.Run("not finished before duration", func(t *testing.T) {
		e := &Entry{
			motion:      Motion{Meta: Meta{Duration: 5.0}},
			currentTime: 0,
		}
		finished := e.Update(1.0)
		if finished {
			t.Error("should not be finished before duration")
		}
		if e.currentTime != 1.0 {
			t.Errorf("currentTime = %v, want 1.0", e.currentTime)
		}
	})

	t.Run("finished at duration", func(t *testing.T) {
		e := &Entry{
			motion:      Motion{Meta: Meta{Duration: 2.0}},
			currentTime: 1.5,
		}
		finished := e.Update(1.0)
		if !finished {
			t.Error("should be finished at duration")
		}
	})
}
