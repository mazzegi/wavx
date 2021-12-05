package wavx

import "math"

func RoundInt(f float64) int {
	return int(math.Round(f))
}

func AbsInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func MaxAbsInts(ns []int) int {
	if len(ns) == 0 {
		return 0
	}
	m := AbsInt(ns[0])
	for i, n := range ns {
		if i == 0 {
			continue
		}
		an := AbsInt(n)
		if an > m {
			m = an
		}
	}
	return m
}
