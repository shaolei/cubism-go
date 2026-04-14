package moc

type Moc struct {
	MocPtr      uintptr
	MocBuffer   []byte
	ModelPtr    uintptr
	ModelBuffer []byte
	closed      bool
}

// Close releases the resources held by Moc.
// After calling Close, the Moc must not be used anymore.
func (m *Moc) Close() {
	if m.closed {
		return
	}
	m.MocBuffer = nil
	m.ModelBuffer = nil
	m.MocPtr = 0
	m.ModelPtr = 0
	m.closed = true
}
