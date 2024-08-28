package uniqueidutils

import (
	"math/rand"
)

const (
	RandKindNum   = 0 // pure number
	RANDKINDLOWER = 1 // lowercase letter
	RANDKINDUPPER = 2 // uppercase letter
	RANDKINDALL   = 3 // mix number letter
)

// GenKrandId Generate random string
func GenKrandId(size int, kind int) string {
	// If SEED is a not deterministic value (example: SEED is derived from the current time),
	//then there's no need to seed. The default random source is seeded automatically in Go 1.20 and later.
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	isAll := kind > 2 || kind < 0
	for i := 0; i < size; i++ {
		if isAll { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return string(result)
}
