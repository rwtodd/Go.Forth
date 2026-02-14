// SPDX-License-Identifier: MIT

package forth

// Dict represents a FORTH dictionary (map[string]any)
type Dict struct {
	m map[string]any
}

// emptyDict creates a new, empty dictionary
func emptyDict(vm *VM) error {
	vm.Push(Dict{m: make(map[string]any)})
	return nil
}

// dictGet looks up a value by key in the dictionary
func dictGet(vm *VM) error {
	key, err := vm.Pop()
	if err != nil {
		return err
	}
	dict, err := vm.Pop()
	if err != nil {
		return err
	}
	d, ok := dict.(Dict)
	if !ok {
		return ErrArgumentMsg("expected dictionary not found")
	}
	k, ok := key.(string)
	if !ok {
		return ErrArgumentMsg("expected string key not found")
	}
	val, found := d.m[k]
	if !found {
		return ErrKeyNotFoundMsg(k)
	}
	vm.Push(val)
	return nil
}

// dictSet sets a value in the dictionary
func dictSet(vm *VM) error {
	key, err := vm.Pop()
	if err != nil {
		return err
	}
	dict, err := vm.Pop()
	if err != nil {
		return err
	}
	val, err := vm.Pop()
	if err != nil {
		return err
	}
	d, ok := dict.(Dict)
	if !ok {
		return ErrArgumentMsg("expected dictionary not found")
	}
	k, ok := key.(string)
	if !ok {
		return ErrArgumentMsg("expected string key not found")
	}
	d.m[k] = val
	return nil
}

// dictDelete removes a key from the dictionary
func dictDelete(vm *VM) error {
	key, err := vm.Pop()
	if err != nil {
		return err
	}
	dict, err := vm.Pop()
	if err != nil {
		return err
	}
	d, ok := dict.(Dict)
	if !ok {
		return ErrArgumentMsg("expected dictionary not found")
	}
	k, ok := key.(string)
	if !ok {
		return ErrArgumentMsg("expected string key not found")
	}
	delete(d.m, k)
	return nil
}

// dictKeys returns all keys of the dictionary
func dictKeys(vm *VM) error {
	dict, err := vm.Pop()
	if err != nil {
		return err
	}
	d, ok := dict.(Dict)
	if !ok {
		return ErrArgumentMsg("expected dictionary not found")
	}
	keys := make([]string, 0, len(d.m))
	for k := range d.m {
		keys = append(keys, k)
	}
	vm.Push(keys)
	return nil
}

// dictGetOr gets a value by key, or returns default if not found
func dictGetOr(vm *VM) error {
	def, err := vm.Pop()
	if err != nil {
		return err
	}
	key, err := vm.Pop()
	if err != nil {
		return err
	}
	dict, err := vm.Pop()
	if err != nil {
		return err
	}
	d, ok := dict.(Dict)
	if !ok {
		return ErrArgumentMsg("expected dictionary not found")
	}
	k, ok := key.(string)
	if !ok {
		return ErrArgumentMsg("expected string key not found")
	}
	if val, found := d.m[k]; found {
		vm.Push(val)
	} else {
		vm.Push(def)
	}
	return nil
}

// dictGetQuery gets a value by key, returning value and -1 if found, or 0 if not found
func dictGetQuery(vm *VM) error {
	key, err := vm.Pop()
	if err != nil {
		return err
	}
	dict, err := vm.Pop()
	if err != nil {
		return err
	}
	d, ok := dict.(Dict)
	if !ok {
		return ErrArgumentMsg("expected dictionary not found")
	}
	k, ok := key.(string)
	if !ok {
		return ErrArgumentMsg("expected string key not found")
	}
	if val, found := d.m[k]; found {
		vm.Push(val)
		vm.Push(-1)
	} else {
		vm.Push(0)
	}
	return nil
}

// dictWordsInit adds dictionary-related words to the VM
func dictWordsInit(vm *VM) {
	vm.Define("empty-dict", Word{emptyDict, false})
	vm.Define("d@", Word{dictGet, false})
	vm.Define("d!", Word{dictSet, false})
	vm.Define("ddel", Word{dictDelete, false})
	vm.Define("dkeys", Word{dictKeys, false})
	vm.Define("d@|", Word{dictGetOr, false})
	vm.Define("d@?", Word{dictGetQuery, false})
}
