package model

import "testing"

func TestFirstNoFilter(t *testing.T) {
	p := FixedProvider([]uint32{1})
	r, err := First(p, Filters[uint32]())
	if err != nil {
		t.Errorf("Expected result, got err %s", err)
	}
	if r != 1 {
		t.Errorf("Expected 1, got %d", r)
	}
}

func byTwo(val uint32) (uint32, error) {
	return val * 2, nil
}

func TestMap(t *testing.T) {
	p := FixedProvider(uint32(1))
	mp := Map[uint32, uint32](byTwo)(p)

	ar, err := mp()
	if err != nil {
		t.Errorf("Expected result, got err %s", err)
	}

	er, _ := p()
	if er*2 != ar {
		t.Errorf("Expected %d, got %d", er*2, ar)
	}
}

func TestSliceMap(t *testing.T) {
	p := FixedProvider([]uint32{1, 2, 3, 4, 5})
	mp := SliceMap[uint32, uint32](p, byTwo)

	ar, err := mp()
	if err != nil {
		t.Errorf("Expected result, got err %s", err)
	}

	er, _ := p()
	for i := range er {
		if er[i]*2 != ar[i] {
			t.Errorf("Expected %d, got %d", er[i]*2, ar[i])
		}
	}
}

func TestParallelSliceMap(t *testing.T) {
	p := FixedProvider([]uint32{1, 2, 3, 4, 5})
	mp := SliceMap[uint32, uint32](p, byTwo, ParallelMap())

	ar, err := mp()
	if err != nil {
		t.Errorf("Expected result, got err %s", err)
	}

	er, _ := p()
	for i := range er {
		if er[i]*2 != ar[i] {
			t.Errorf("Expected %d, got %d", er[i]*2, ar[i])
		}
	}
}

func isTwo(val uint32) bool {
	return val == 2
}

func TestFirst(t *testing.T) {
	p := FixedProvider([]uint32{1, 2, 3, 4, 5})
	mp, err := First(p, Filters(isTwo))
	if err != nil {
		t.Errorf("Expected result, got err %s", err)
	}
	if mp != 2 {
		t.Errorf("Expected 2, got %d", mp)
	}
}

func TestThenOperator(t *testing.T) {
	p := FixedProvider(uint32(1))
	count := uint32(0)

	op1 := func(u uint32) error {
		count += u
		return nil
	}
	op2 := func(u uint32) error {
		count += u
		return nil
	}

	err := For(p, ThenOperator(op1, Operators(op2)))
	if err != nil {
		t.Errorf("Expected result, got err %s", err)
	}
	if count != 2 {
		t.Errorf("Expected 2, got %d", count)
	}
}

func TestForEachSlice(t *testing.T) {
	p := FixedProvider([]uint32{1, 2, 3, 4, 5})
	count := uint32(0)

	err := ForEachSlice(p, func(u uint32) error {
		count += u
		return nil
	})
	if err != nil {
		t.Errorf("Expected result, got err %s", err)
	}
	if count != 15 {
		t.Errorf("Expected 15, got %d", count)
	}
}

func TestForEachSliceParallel(t *testing.T) {
	p := FixedProvider([]uint32{1, 2, 3, 4, 5})
	count := uint32(0)

	err := ForEachSlice(p, func(u uint32) error {
		count += u
		return nil
	}, ParallelExecute())
	if err != nil {
		t.Errorf("Expected result, got err %s", err)
	}
	if count != 15 {
		t.Errorf("Expected 15, got %d", count)
	}
}

func TestForEachMap(t *testing.T) {
	p := FixedProvider(map[uint32][]uint32{1: {1, 2}, 2: {1, 2, 3}})
	counts := map[uint32]uint32{}

	err := ForEachMap(p, func(k uint32) Operator[[]uint32] {
		return func(vs []uint32) error {
			count := uint32(0)
			for _, v := range vs {
				count += v
			}
			counts[k] = count
			return nil
		}
	})
	if err != nil {
		t.Errorf("Expected result, got err %s", err)
	}
	if counts[1] != 3 {
		t.Errorf("Expected 3, got %d", counts[1])
	}
	if counts[2] != 6 {
		t.Errorf("Expected 6, got %d", counts[2])
	}
}

func TestForEachMapParallel(t *testing.T) {
	p := FixedProvider(map[uint32][]uint32{1: {1, 2}, 2: {1, 2, 3}})
	counts := map[uint32]uint32{}

	err := ForEachMap(p, func(k uint32) Operator[[]uint32] {
		return func(vs []uint32) error {
			count := uint32(0)
			for _, v := range vs {
				count += v
			}
			counts[k] = count
			return nil
		}
	}, ParallelExecute())
	if err != nil {
		t.Errorf("Expected result, got err %s", err)
	}
	if counts[1] != 3 {
		t.Errorf("Expected 3, got %d", counts[1])
	}
	if counts[2] != 6 {
		t.Errorf("Expected 6, got %d", counts[2])
	}
}

func TestMerge(t *testing.T) {
	p1 := FixedProvider([]uint32{1, 2, 3, 4, 5})
	p2 := FixedProvider([]uint32{1, 2, 3, 4, 5})
	rp := MergeSliceProvider(p1, p2)

	rs, err := rp()
	if err != nil {
		t.Errorf("Expected result, got err %s", err)
	}
	if len(rs) != 10 {
		t.Errorf("Expected 10, got %d", len(rs))
	}
}

func TestApply(t *testing.T) {
	f := func(a uint32) Provider[uint32] {
		return func() (uint32, error) {
			return a + 32, nil
		}
	}

	r, err := Apply(f)(5)
	if err != nil {
		t.Errorf("Expected result, got err %s", err)
	}
	if r != 37 {
		t.Errorf("Expected 37, got %d", r)
	}
}

func TestCurry(t *testing.T) {
	f := func(a uint32, b uint32) uint32 {
		return a + b
	}
	r := Curry(f)(1)(2)
	if r != 3 {
		t.Errorf("Expected 3, got %d", r)
	}
}

func TestCompose(t *testing.T) {
	f := func(a uint64) func(b uint32) func(c uint16) Provider[uint32] {
		return func(b uint32) func(c uint16) Provider[uint32] {
			return func(c uint16) Provider[uint32] {
				return func() (uint32, error) {
					return uint32(a)*b*uint32(c) + 24, nil
				}
			}
		}
	}

	type a = uint32
	type b = func(uint16) Provider[uint32]
	type c = func(uint16) (uint32, error)

	r, err := Compose(Curry(Compose[a, b, c])(Apply[uint16, uint32]), f)(2)(3)(4)
	if err != nil {
		t.Errorf("Expected result, got err %s", err)
	}
	if r != 48 {
		t.Errorf("Expected 47, got %d", r)
	}
}

func TestParameterTransformation(t *testing.T) {
	type rt = func(uint642 uint64) Provider[uint32]

	f := func(a uint32) Provider[uint32] {
		return func() (uint32, error) {
			return a + 32, nil
		}
	}
	tf := func(uint642 uint64) uint32 {
		return uint32(uint642)
	}

	var rf rt = Compose(f, tf)
	r, err := rf(5)()
	if err != nil {
		t.Errorf("Expected result, got err %s", err)
	}
	if r != 37 {
		t.Errorf("Expected 37, got %d", r)
	}
}

func TestFilters(t *testing.T) {
	ip := FixedProvider([]uint32{1, 2, 3, 4, 5})
	f := func(i uint32) bool {
		return i > 3
	}
	op := FilteredProvider(ip, Filters(f))
	r, err := op()
	if err != nil {
		t.Errorf("Expected result, got err %s", err)
	}
	if len(r) != 2 {
		t.Errorf("Expected 2, got %d", len(r))
	}

}
