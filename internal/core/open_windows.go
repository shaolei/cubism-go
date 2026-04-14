//go:build windows

package core

import (
	"fmt"
	"sync"

	"golang.org/x/sys/windows"
)

var (
	loadedDLLs   = make(map[string]*windows.DLL)
	loadedDLLsMu sync.Mutex
)

func openLibrary(name string) (uintptr, error) {
	loadedDLLsMu.Lock()
	defer loadedDLLsMu.Unlock()

	if dll, ok := loadedDLLs[name]; ok {
		return uintptr(dll.Handle), nil
	}

	handle, err := windows.LoadDLL(name)
	if err != nil {
		return 0, fmt.Errorf("failed to load DLL %s: %w", name, err)
	}
	loadedDLLs[name] = handle
	return uintptr(handle.Handle), nil
}

// CloseLibrary releases the loaded DLL. Call this when the application exits.
func CloseLibrary(name string) error {
	loadedDLLsMu.Lock()
	defer loadedDLLsMu.Unlock()

	dll, ok := loadedDLLs[name]
	if !ok {
		return nil
	}
	delete(loadedDLLs, name)
	return dll.Release()
}
