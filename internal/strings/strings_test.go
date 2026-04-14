package strings_test

import (
	"testing"
	"unsafe"

	"github.com/shaolei/cubism-go/internal/strings"
	"github.com/stretchr/testify/assert"
)

func TestGoString(t *testing.T) {
	t.Parallel()

	t.Run("normal string", func(t *testing.T) {
		t.Parallel()
		cs := "hello\x00"
		ptr := uintptr(unsafe.Pointer(&[]byte(cs)[0]))
		got := strings.GoString(ptr)
		assert.Equal(t, "hello", got)
	})

	t.Run("empty string", func(t *testing.T) {
		t.Parallel()
		cs := "\x00"
		ptr := uintptr(unsafe.Pointer(&[]byte(cs)[0]))
		got := strings.GoString(ptr)
		assert.Equal(t, "", got)
	})

	t.Run("single character", func(t *testing.T) {
		t.Parallel()
		cs := "A\x00"
		ptr := uintptr(unsafe.Pointer(&[]byte(cs)[0]))
		got := strings.GoString(ptr)
		assert.Equal(t, "A", got)
	})

	t.Run("unicode string", func(t *testing.T) {
		t.Parallel()
		cs := "こんにちは\x00"
		ptr := uintptr(unsafe.Pointer(&[]byte(cs)[0]))
		got := strings.GoString(ptr)
		assert.Equal(t, "こんにちは", got)
	})

	t.Run("nil pointer returns empty", func(t *testing.T) {
		t.Parallel()
		var nilPtr uintptr
		got := strings.GoString(nilPtr)
		assert.Equal(t, "", got)
	})

	t.Run("long string", func(t *testing.T) {
		t.Parallel()
		cs := "abcdefghijklmnopqrstuvwxyz0123456789\x00"
		ptr := uintptr(unsafe.Pointer(&[]byte(cs)[0]))
		got := strings.GoString(ptr)
		assert.Equal(t, "abcdefghijklmnopqrstuvwxyz0123456789", got)
	})
}
