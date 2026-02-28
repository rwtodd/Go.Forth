package random

import (
	"math/rand/v2"
	"time"

	"github.com/rwtodd/Go.Forth/forth"
)

var rng = rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano())))

func init() {
	forth.RegisterExtension("random", func(vm *forth.VM) error {
		vm.Define(forth.NewNativeWord("randint", randint))
		vm.Define(forth.NewNativeWord("randfloat", randfloat))
		vm.Define(forth.NewNativeWord("@shuffle", arrayShuffle))
		vm.Define(forth.NewNativeWord("@select", arraySelect))
		vm.Define(forth.NewNativeWord("randseed!", seedWord))
		return nil
	})
}

// randint ( max min -- n ) gives a random integer in the half-closed range [min,max).
func randint(vm *forth.VM) error {
	minVal, err := vm.Pop()
	if err != nil {
		return err
	}
	maxVal, err := vm.Pop()
	if err != nil {
		return err
	}

	min, ok1 := minVal.(int)
	max, ok2 := maxVal.(int)
	if !ok1 || !ok2 {
		return forth.ErrArgumentMsg("randint requires integer arguments")
	}
	if max <= min {
		return forth.ErrArgumentMsg("randint requires max > min")
	}

	n := rng.IntN(max-min) + min
	vm.Push(n)
	return nil
}

// randfloat ( -- n ) a random number in the range [0,1).
func randfloat(vm *forth.VM) error {
	vm.Push(rng.Float64())
	return nil
}

// seedWord ( n -- ) sets the random seed
func seedWord(vm *forth.VM) error {
	seedVal, err := vm.Pop()
	if err != nil {
		return err
	}
	seed, ok := seedVal.(int)
	if !ok {
		return forth.ErrArgumentMsg("@seed requires an integer argument")
	}

	rng = rand.New(rand.NewPCG(uint64(seed), uint64(seed)))
	return nil
}

func getSlice(arrVal any) (any, error) {
	if v, ok := arrVal.(*forth.Variable); ok {
		return v.Value(), nil
	}
	return arrVal, nil
}

// arrayShuffle ( arr -- ) expects an array (or a variable that points to an array) and shuffles it.
func arrayShuffle(vm *forth.VM) error {
	arrVal, err := vm.Pop()
	if err != nil {
		return err
	}

	slice, err := getSlice(arrVal)
	if err != nil {
		return err
	}

	switch a := slice.(type) {
	case []byte:
		rng.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	case []int:
		rng.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	case []float64:
		rng.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	case []string:
		rng.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	case []any:
		rng.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	default:
		return forth.ErrArgumentMsg("@shuffle expects an array or variable pointing to an array")
	}
	return nil
}

// arraySelect ( arr -- value ) selects an element from the array arr at random leaving it on the stack.
func arraySelect(vm *forth.VM) error {
	arrVal, err := vm.Pop()
	if err != nil {
		return err
	}

	slice, err := getSlice(arrVal)
	if err != nil {
		return err
	}

	switch a := slice.(type) {
	case []byte:
		if len(a) == 0 {
			return forth.ErrArgumentMsg("@select from empty array")
		}
		vm.Push(int(a[rng.IntN(len(a))]))
	case []int:
		if len(a) == 0 {
			return forth.ErrArgumentMsg("@select from empty array")
		}
		vm.Push(a[rng.IntN(len(a))])
	case []float64:
		if len(a) == 0 {
			return forth.ErrArgumentMsg("@select from empty array")
		}
		vm.Push(a[rng.IntN(len(a))])
	case []string:
		if len(a) == 0 {
			return forth.ErrArgumentMsg("@select from empty array")
		}
		vm.Push(a[rng.IntN(len(a))])
	case []any:
		if len(a) == 0 {
			return forth.ErrArgumentMsg("@select from empty array")
		}
		vm.Push(a[rng.IntN(len(a))])
	default:
		return forth.ErrArgumentMsg("@select expects an array or variable pointing to an array")
	}
	return nil
}
