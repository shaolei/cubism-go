package updater

import (
	"testing"
)

// collectUpdater records the order in which updaters are called.
type collectUpdater struct {
	order   int
	results *[]int
}

func (c *collectUpdater) OnLateUpdate(delta float64) {
	*c.results = append(*c.results, c.order)
}

func TestNewUpdateScheduler(t *testing.T) {
	s := NewUpdateScheduler()
	if s == nil {
		t.Fatal("NewUpdateScheduler returned nil")
	}
	if s.Len() != 0 {
		t.Fatalf("expected 0 updaters, got %d", s.Len())
	}
}

func TestAddUpdaterSortsByExecutionOrder(t *testing.T) {
	s := NewUpdateScheduler()
	var order []int

	// Add updaters in reverse priority order
	s.AddUpdater(UpdateOrderPose, &collectUpdater{order: UpdateOrderPose, results: &order})
	s.AddUpdater(UpdateOrderEyeBlink, &collectUpdater{order: UpdateOrderEyeBlink, results: &order})
	s.AddUpdater(UpdateOrderBreath, &collectUpdater{order: UpdateOrderBreath, results: &order})
	s.AddUpdater(UpdateOrderExpression, &collectUpdater{order: UpdateOrderExpression, results: &order})
	s.AddUpdater(UpdateOrderLook, &collectUpdater{order: UpdateOrderLook, results: &order})
	s.AddUpdater(UpdateOrderPhysics, &collectUpdater{order: UpdateOrderPhysics, results: &order})

	if s.Len() != 6 {
		t.Fatalf("expected 6 updaters, got %d", s.Len())
	}

	// Call OnLateUpdate — updaters should execute in ascending priority order
	s.OnLateUpdate(0.016)

	expected := []int{
		UpdateOrderEyeBlink,   // 200
		UpdateOrderExpression, // 300
		UpdateOrderLook,       // 400
		UpdateOrderBreath,     // 500
		UpdateOrderPhysics,    // 600
		UpdateOrderPose,       // 800
	}

	if len(order) != len(expected) {
		t.Fatalf("expected %d calls, got %d", len(expected), len(order))
	}
	for i, got := range order {
		if got != expected[i] {
			t.Errorf("call %d: expected order %d, got %d", i, expected[i], got)
		}
	}
}

func TestFuncUpdater(t *testing.T) {
	called := false
	fu := NewFuncUpdater(UpdateOrderBreath, func(delta float64) {
		called = true
		if delta != 0.033 {
			t.Errorf("expected delta 0.033, got %f", delta)
		}
	})

	if fu.GetExecutionOrder() != UpdateOrderBreath {
		t.Errorf("expected order %d, got %d", UpdateOrderBreath, fu.GetExecutionOrder())
	}

	fu.OnLateUpdate(0.033)
	if !called {
		t.Error("function was not called")
	}
}

func TestFuncUpdaterNilFunction(t *testing.T) {
	fu := NewFuncUpdater(UpdateOrderBreath, nil)
	// Should not panic
	fu.OnLateUpdate(0.016)
}

func TestRemoveUpdater(t *testing.T) {
	s := NewUpdateScheduler()
	var order []int

	u1 := &collectUpdater{order: 100, results: &order}
	u2 := &collectUpdater{order: 200, results: &order}
	u3 := &collectUpdater{order: 300, results: &order}

	s.AddUpdater(100, u1)
	s.AddUpdater(200, u2)
	s.AddUpdater(300, u3)

	if s.Len() != 3 {
		t.Fatalf("expected 3 updaters, got %d", s.Len())
	}

	// Remove u2
	s.RemoveUpdater(u2)

	if s.Len() != 2 {
		t.Fatalf("expected 2 updaters after removal, got %d", s.Len())
	}

	s.OnLateUpdate(0.016)

	expected := []int{100, 300}
	if len(order) != len(expected) {
		t.Fatalf("expected %d calls, got %d", len(expected), len(order))
	}
	for i, got := range order {
		if got != expected[i] {
			t.Errorf("call %d: expected order %d, got %d", i, expected[i], got)
		}
	}
}

func TestRemoveNonExistentUpdater(t *testing.T) {
	s := NewUpdateScheduler()
	var order []int

	u1 := &collectUpdater{order: 100, results: &order}
	u2 := &collectUpdater{order: 200, results: &order}

	s.AddUpdater(100, u1)

	// Removing u2 (not added) should be a no-op
	s.RemoveUpdater(u2)

	if s.Len() != 1 {
		t.Fatalf("expected 1 updater, got %d", s.Len())
	}
}

func TestAddNilUpdater(t *testing.T) {
	s := NewUpdateScheduler()
	s.AddUpdater(100, nil)

	if s.Len() != 0 {
		t.Fatalf("expected 0 updaters when adding nil, got %d", s.Len())
	}
}

func TestClear(t *testing.T) {
	s := NewUpdateScheduler()
	var order []int

	s.AddUpdater(100, &collectUpdater{order: 100, results: &order})
	s.AddUpdater(200, &collectUpdater{order: 200, results: &order})

	if s.Len() != 2 {
		t.Fatalf("expected 2 updaters, got %d", s.Len())
	}

	s.Clear()

	if s.Len() != 0 {
		t.Fatalf("expected 0 updaters after clear, got %d", s.Len())
	}

	// OnLateUpdate after clear should do nothing
	s.OnLateUpdate(0.016)
	if len(order) != 0 {
		t.Fatalf("expected 0 calls after clear, got %d", len(order))
	}
}

func TestSortStability(t *testing.T) {
	// Multiple updaters with the same priority should maintain insertion order
	s := NewUpdateScheduler()
	var order []string

	s.AddUpdater(100, NewFuncUpdater(100, func(delta float64) { order = append(order, "first") }))
	s.AddUpdater(100, NewFuncUpdater(100, func(delta float64) { order = append(order, "second") }))
	s.AddUpdater(100, NewFuncUpdater(100, func(delta float64) { order = append(order, "third") }))

	s.OnLateUpdate(0.016)

	expected := []string{"first", "second", "third"}
	if len(order) != len(expected) {
		t.Fatalf("expected %d calls, got %d", len(expected), len(order))
	}
	for i, got := range order {
		if got != expected[i] {
			t.Errorf("call %d: expected %q, got %q", i, expected[i], got)
		}
	}
}

func TestLazySortOnlyWhenNeeded(t *testing.T) {
	s := NewUpdateScheduler()
	var order []int

	u1 := &collectUpdater{order: 100, results: &order}
	u2 := &collectUpdater{order: 200, results: &order}

	// First update triggers sort
	s.AddUpdater(200, u2)
	s.AddUpdater(100, u1)
	s.OnLateUpdate(0.016)

	order = order[:0]

	// Second update should NOT re-sort (already sorted)
	s.OnLateUpdate(0.016)

	// Verify order is still correct
	expected := []int{100, 200}
	for i, got := range order {
		if got != expected[i] {
			t.Errorf("call %d: expected order %d, got %d", i, expected[i], got)
		}
	}

	// Add a new updater — should mark for re-sort
	s.AddUpdater(150, &collectUpdater{order: 150, results: &order})
	order = order[:0]
	s.OnLateUpdate(0.016)

	expected = []int{100, 150, 200}
	for i, got := range order {
		if got != expected[i] {
			t.Errorf("call %d: expected order %d, got %d", i, expected[i], got)
		}
	}
}

func TestUpdateOrderConstants(t *testing.T) {
	// Verify constants match official SDK values and are in ascending order
	orders := []int{
		UpdateOrderEyeBlink,
		UpdateOrderExpression,
		UpdateOrderLook,
		UpdateOrderBreath,
		UpdateOrderPhysics,
		UpdateOrderLipSync,
		UpdateOrderPose,
	}

	expected := []int{200, 300, 400, 500, 600, 700, 800}
	for i, got := range orders {
		if got != expected[i] {
			t.Errorf("order constant %d: expected %d, got %d", i, expected[i], got)
		}
	}

	// Verify ascending order
	for i := 1; i < len(orders); i++ {
		if orders[i] <= orders[i-1] {
			t.Errorf("order constants not ascending: %d <= %d", orders[i], orders[i-1])
		}
	}
}
