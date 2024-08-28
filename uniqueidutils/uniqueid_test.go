package uniqueidutils

import (
	"fmt"
	"testing"
	"time"
)

func TestGenKrandId(t *testing.T) {
	m := make(map[string]bool, 10000)
	for i := 0; i < 10000; i++ {
		sn := GenKrandId(12, RANDKINDALL)
		if _, exist := m[sn]; exist {
			fmt.Println("found repeat random...")
		}
		m[sn] = true
		fmt.Println(sn)
	}
}

func TestGenSnowflakeId(t *testing.T) {
	for i := 0; i < 100; i++ {
		time.Sleep(time.Nanosecond)
		id := GenSnowflakeId()
		fmt.Println(id.Int64(), id.Base58(), id.Base2())
	}
}
