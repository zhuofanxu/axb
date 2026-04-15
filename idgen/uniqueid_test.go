package idgen

import (
	"sync"
	"testing"
	"unicode"
)

// ---- Snowflake ----

func TestGenSnowflakeId_Unique(t *testing.T) {
	const count = 10000
	seen := make(map[int64]struct{}, count)
	for i := 0; i < count; i++ {
		id := GenSnowflakeId().Int64()
		if _, exists := seen[id]; exists {
			t.Fatalf("duplicate snowflake ID: %d", id)
		}
		seen[id] = struct{}{}
	}
}

func TestGenSnowflakeId_ConcurrentUnique(t *testing.T) {
	const count = 10000
	ids := make([]int64, count)
	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(idx int) {
			defer wg.Done()
			id := GenSnowflakeId().Int64()
			mu.Lock()
			ids[idx] = id
			mu.Unlock()
		}(i)
	}
	wg.Wait()

	seen := make(map[int64]struct{}, count)
	for _, id := range ids {
		if _, exists := seen[id]; exists {
			t.Fatalf("duplicate snowflake ID in concurrent generation: %d", id)
		}
		seen[id] = struct{}{}
	}
}

func TestGenSnowflakeId_Positive(t *testing.T) {
	for i := 0; i < 100; i++ {
		id := GenSnowflakeId().Int64()
		if id <= 0 {
			t.Errorf("snowflake ID should be positive, got %d", id)
		}
	}
}

// ---- KrandId ----

func TestGenKrandId_Length(t *testing.T) {
	tests := []struct {
		size int
		kind int
	}{
		{8, RandKindNum},
		{12, RandKindLower},
		{16, RandKindUpper},
		{10, RandKindNumLower},
		{10, RandKindNumUpper},
		{10, RandKindLowerUpper},
		{20, RandKindAll},
	}
	for _, tt := range tests {
		got := GenKrandId(tt.size, tt.kind)
		if len(got) != tt.size {
			t.Errorf("kind=%d: expected length %d, got %d (value=%q)", tt.kind, tt.size, len(got), got)
		}
	}
}

func TestGenKrandId_CharacterSet(t *testing.T) {
	tests := []struct {
		name    string
		kind    int
		allowed func(r rune) bool
	}{
		{"num only", RandKindNum, func(r rune) bool { return unicode.IsDigit(r) }},
		{"lower only", RandKindLower, func(r rune) bool { return unicode.IsLower(r) }},
		{"upper only", RandKindUpper, func(r rune) bool { return unicode.IsUpper(r) }},
		{"num+lower", RandKindNumLower, func(r rune) bool { return unicode.IsDigit(r) || unicode.IsLower(r) }},
		{"num+upper", RandKindNumUpper, func(r rune) bool { return unicode.IsDigit(r) || unicode.IsUpper(r) }},
		{"lower+upper", RandKindLowerUpper, func(r rune) bool { return unicode.IsLetter(r) }},
		{"all", RandKindAll, func(r rune) bool { return unicode.IsLetter(r) || unicode.IsDigit(r) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := GenKrandId(200, tt.kind)
			for _, r := range s {
				if !tt.allowed(r) {
					t.Errorf("unexpected character %q in result %q", r, s)
				}
			}
		})
	}
}

func TestGenKrandId_NoAmbiguousChars(t *testing.T) {
	// 文档注释：排除易混淆字符 0,1,O,o,l
	ambiguous := []rune{'0', '1', 'O', 'o', 'l'}
	for i := 0; i < 50; i++ {
		s := GenKrandId(100, RandKindAll)
		for _, r := range s {
			for _, bad := range ambiguous {
				if r == bad {
					t.Errorf("ambiguous character %q should not appear in output", r)
				}
			}
		}
	}
}

func TestGenKrandId_Unique(t *testing.T) {
	const count = 10000
	seen := make(map[string]struct{}, count)
	for i := 0; i < count; i++ {
		s := GenKrandId(10, RandKindAll)
		if _, exists := seen[s]; exists {
			t.Fatalf("duplicate krand ID after %d iterations: %q", i, s)
		}
		seen[s] = struct{}{}
	}
}
