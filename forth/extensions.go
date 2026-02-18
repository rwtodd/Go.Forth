// SPDX-License-Identifier: MIT

package forth

import (
	"sort"
	"sync"
)

// extensionRegistry holds the registered extensions
var (
	extensionRegistry = make(map[string]func(*VM) error)
	registryLock      sync.RWMutex
)

// RegisterExtension registers an extension with the VM.
// This should be called in the init() function of the extension package.
func RegisterExtension(name string, init func(*VM) error) {
	registryLock.Lock()
	defer registryLock.Unlock()
	extensionRegistry[name] = init
}

// extensionList returns a list of all registered extensions
// ( -- []string )
func extensionList(vm *VM) error {
	registryLock.RLock()
	names := make([]string, 0, len(extensionRegistry))
	for name := range extensionRegistry {
		names = append(names, name)
	}
	registryLock.RUnlock()

	sort.Strings(names)
	vm.Push(names)
	return nil
}

// activateExtensions activates the specified extensions
// ( name1 ... nameN count -- )
func activateExtensions(vm *VM) error {
	countVal, err := vm.Pop()
	if err != nil {
		return err
	}
	count, ok := countVal.(int)
	if !ok || count < 0 {
		return ErrArgumentMsg("activate-extensions count must be a non-negative integer")
	}

	if len(vm.Stack) < count {
		return ErrUnderflowMsg("activate-extensions stack underflow")
	}

	names := make([]string, count)
	for i := count - 1; i >= 0; i-- {
		val, err := vm.Pop()
		if err != nil {
			return err
		}
		name, ok := val.(string)
		if !ok {
			return ErrArgumentMsg("activate-extensions expects string extension names")
		}
		names[i] = name
	}

	registryLock.RLock()
	defer registryLock.RUnlock()

	for _, name := range names {
		// Idempotency check: if already activated, skip
		if vm.ActivatedExtensions[name] {
			continue
		}

		init, ok := extensionRegistry[name]
		if !ok {
			return ErrArgumentMsg("extension not found: " + name)
		}
		if err := init(vm); err != nil {
			return err
		}
		vm.ActivatedExtensions[name] = true
	}
	return mark(vm) // if we `forget` past an activation, can't reactivate!
}

func extensionsWordsInit(vm *VM) {
	vm.Define(&NativeWord{name: "extension-list", run: extensionList, immediate: false})
	vm.Define(&NativeWord{name: "<activate-extensions>", run: activateExtensions, immediate: false})
}
