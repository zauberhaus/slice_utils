package slice_utils_test

import (
	"errors"
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zauberhaus/slice_utils"
)

var _ fmt.Stringer = MyStringer(1)

type MyStringer int

func (s MyStringer) String() string {
	return "val" + strconv.Itoa(int(s))
}

func TestHashSeq(t *testing.T) {
	data := []string{"1", "2", "3"}
	r := maps.Collect(slice_utils.HashSeq(slices.Values(data)))
	assert.NotEmpty(t, r)
	assert.Len(t, r, len(data))
	assert.IsType(t, map[uint64]string{}, r)
}

func TestAnySeq(t *testing.T) {
	data := []string{"1", "2", "3"}
	s := slice_utils.AnySeq(slices.Values(data))
	r := slices.Collect(s)
	assert.NotEmpty(t, r)
	assert.Len(t, r, len(data))
	assert.IsType(t, []any{}, r)
}

func TestFilterSeq(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	seq := slice_utils.FilterSeq(slices.Values(data), func(v int) bool {
		return v%2 == 0
	})
	got := slices.Collect(seq)
	assert.Equal(t, []int{2, 4}, got)
}

func TestRemoveSeq(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	remove := []int{2, 4}
	// Note: RemoveSeq iterates the 'remove' sequence for every element in 'data'.
	// This works for slice-backed sequences (restartable).
	seq := slice_utils.RemoveSeq(slices.Values(data), slices.Values(remove))
	got := slices.Collect(seq)
	assert.Equal(t, []int{1, 3, 5}, got)
}

func TestPatternSeq(t *testing.T) {
	data := []string{"apple", "banana", "cherry", "date"}
	re := regexp.MustCompile(`a.*e`) // matches apple, date
	seq := slice_utils.PatternSeq(slices.Values(data), re)
	got := slices.Collect(seq)
	assert.Equal(t, []string{"apple", "date"}, got)
}

func TestStringPatternSeq(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		data := []string{"apple", "banana", "cherry"}
		seq := slice_utils.StringPatternSeq(slices.Values(data), "banana")
		got := slices.Collect(seq)
		assert.Equal(t, []string{"banana"}, got)
	})

	t.Run("int", func(t *testing.T) {
		data := []int{1, 2, 3}
		seq := slice_utils.StringPatternSeq(slices.Values(data), "2")
		got := slices.Collect(seq)
		assert.Equal(t, []int{2}, got)
	})

	t.Run("Stringer", func(t *testing.T) {
		data := []MyStringer{1, 2, 3}
		seq := slice_utils.StringPatternSeq(slices.Values(data), "val2")
		got := slices.Collect(seq)
		assert.Equal(t, []MyStringer{2}, got)
	})
}

func TestDuplicateSeq(t *testing.T) {
	data := []int{1, 2, 3, 1, 4, 2, 5, 1}
	seq := slice_utils.DuplicateSeq(slices.Values(data))
	got := slices.Collect(seq)
	assert.Equal(t, []int{1, 2}, got)
}

func TestDeduplicationSeq(t *testing.T) {
	data := []int{1, 2, 2, 3, 1, 4, 1}
	seq := slice_utils.DeduplicationSeq(slices.Values(data))
	got := slices.Collect(seq)
	assert.Equal(t, []int{1, 2, 3, 4}, got)
}

func TestGroupSeq(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6}
	seq := slice_utils.GroupSeq[[]int](slices.Values(data), func(v int) string {
		if v%2 == 0 {
			return "even"
		}
		return "odd"
	})
	got := slices.Collect(seq)
	assert.Len(t, got, 2)

	// Verify contents. Order of groups is not guaranteed.
	for _, g := range got {
		if len(g) > 0 {
			if g[0]%2 == 0 {
				assert.Equal(t, []int{2, 4, 6}, g)
			} else {
				assert.Equal(t, []int{1, 3, 5}, g)
			}
		}
	}
}

func TestCountSeq(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	count := slice_utils.CountSeq(slices.Values(data))
	assert.Equal(t, 5, count)
}

func TestSumFuncSeq(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		data := []string{"1", "2", "3"}
		sum, err := slice_utils.SumFuncSeq(slices.Values(data), func(s string) (int, error) {
			return int(s[0] - '0'), nil
		})
		assert.NoError(t, err)
		assert.Equal(t, 6, sum)
	})

	t.Run("error", func(t *testing.T) {
		data := []string{"1", "a", "3"}
		_, err := slice_utils.SumFuncSeq(slices.Values(data), func(s string) (int, error) {
			if s == "a" {
				return 0, errors.New("invalid number")
			}
			return int(s[0] - '0'), nil
		})
		assert.Error(t, err)
	})
}

func TestSumSeq(t *testing.T) {
	data := []int{3, 1, 2}
	sum := slice_utils.SumSeq(slices.Values(data))
	assert.Equal(t, 6, sum)
}

func TestIsEmptySeq(t *testing.T) {
	assert.True(t, slice_utils.IsEmptySeq(slices.Values([]int{})))
	assert.False(t, slice_utils.IsEmptySeq(slices.Values([]int{1})))
}

func TestReplaceFuncSeq(t *testing.T) {
	data := []int{1, 2, 3}
	seq := slice_utils.ReplaceFuncSeq(slices.Values(data), func(v int) int {
		return v * 10
	})
	got := slices.Collect(seq)
	assert.Equal(t, []int{10, 20, 30}, got)
}

func TestReplaceSeq(t *testing.T) {
	data := []string{"a", "b", "c"}
	replacements := map[string]string{"a": "A", "c": "C"}
	seq := slice_utils.ReplaceSeq(slices.Values(data), replacements)
	got := slices.Collect(seq)
	assert.Equal(t, []string{"A", "b", "C"}, got)
}

func TestConvertSeq(t *testing.T) {
	data := []int{1, 2, 3}
	seq := slice_utils.ConvertSeq(slices.Values(data), func(v int) string {
		return string(rune('0' + v))
	})
	got := slices.Collect(seq)
	assert.Equal(t, []string{"1", "2", "3"}, got)
}
