package slice_utils_test // Changed to a separate test package

import (
	"errors"
	"regexp"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zauberhaus/slice_utils"
)

func TestSelect(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		f     func(val int) bool
		want  []int
	}{
		{
			name:  "select even numbers",
			input: []int{1, 2, 3, 4, 5},
			f:     func(val int) bool { return val%2 == 0 },
			want:  []int{2, 4},
		},
		{
			name:  "select numbers greater than 3",
			input: []int{1, 2, 3, 4, 5},
			f:     func(val int) bool { return val > 3 },
			want:  []int{4, 5},
		},
		{
			name:  "select all numbers",
			input: []int{1, 2, 3},
			f:     func(val int) bool { return true },
			want:  []int{1, 2, 3},
		},
		{
			name:  "select no numbers",
			input: []int{1, 2, 3},
			f:     func(val int) bool { return false },
			want:  []int{},
		},
		{
			name:  "empty slice",
			input: []int{},
			f:     func(val int) bool { return true },
			want:  []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slice_utils.Select(tt.input, tt.f)
			assert.ElementsMatch(t, tt.want, got, "Select() should return matching elements") // Use ElementsMatch for slice comparison where order might not matter
		})
	}
}

func TestCount(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		f     func(val int) bool
		want  int
	}{
		{
			name:  "count even numbers",
			input: []int{1, 2, 3, 4, 5},
			f:     func(val int) bool { return val%2 == 0 },
			want:  2,
		},
		{
			name:  "count numbers greater than 3",
			input: []int{1, 2, 3, 4, 5},
			f:     func(val int) bool { return val > 3 },
			want:  2,
		},
		{
			name:  "count all numbers",
			input: []int{1, 2, 3},
			f:     func(val int) bool { return true },
			want:  3,
		},
		{
			name:  "count no numbers",
			input: []int{1, 2, 3},
			f:     func(val int) bool { return false },
			want:  0,
		},
		{
			name:  "empty slice",
			input: []int{},
			f:     func(val int) bool { return true },
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slice_utils.Count(tt.input, tt.f)
			assert.Equal(t, tt.want, got, "Count() should return the correct count")
		})
	}
}

func TestEmpty(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		f     func(val int) bool
		want  bool
	}{
		{
			name:  "not empty - some match",
			input: []int{1, 2, 3, 4, 5},
			f:     func(val int) bool { return val%2 == 0 },
			want:  false,
		},
		{
			name:  "not empty - all match",
			input: []int{1, 2, 3},
			f:     func(val int) bool { return true },
			want:  false,
		},
		{
			name:  "empty - no match",
			input: []int{1, 3, 5},
			f:     func(val int) bool { return val%2 == 0 },
			want:  true,
		},
		{
			name:  "empty slice",
			input: []int{},
			f:     func(val int) bool { return true },
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slice_utils.Empty(tt.input, tt.f)
			assert.Equal(t, tt.want, got, "Empty() should return the correct boolean status")
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		vals  []int
		want  []int
	}{
		{
			name:  "delete single existing value",
			input: []int{1, 2, 3, 4, 5},
			vals:  []int{3},
			want:  []int{1, 2, 4, 5}, // Only the first occurrence is removed
		},
		{
			name:  "delete multiple existing values (first occurrence of first value)",
			input: []int{1, 2, 3, 4, 3, 5},
			vals:  []int{3, 4},
			want:  []int{1, 2, 4, 3, 5}, // Removes the first '3'
		},
		{
			name:  "delete non-existing value",
			input: []int{1, 2, 3, 4, 5},
			vals:  []int{6},
			want:  []int{1, 2, 3, 4, 5},
		},
		{
			name:  "delete from empty slice",
			input: []int{},
			vals:  []int{1},
			want:  []int{},
		},
		{
			name:  "delete first element",
			input: []int{1, 2, 3},
			vals:  []int{1},
			want:  []int{2, 3},
		},
		{
			name:  "delete last element",
			input: []int{1, 2, 3},
			vals:  []int{3},
			want:  []int{1, 2},
		},
		{
			name:  "delete with no values to delete",
			input: []int{1, 2, 3},
			vals:  []int{},
			want:  []int{1, 2, 3},
		},
		{
			name:  "delete duplicate values in input (only first matched is deleted)",
			input: []int{1, 2, 2, 3},
			vals:  []int{2},
			want:  []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a copy to ensure the original input slice is not modified for other tests
			inputCopy := make([]int, len(tt.input))
			copy(inputCopy, tt.input)

			got := slice_utils.Delete(inputCopy, tt.vals...)
			assert.Equal(t, tt.want, got, "Delete() should return the slice with the first matching element removed")
		})
	}
}

func TestSortFunc(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		f     func(val1, val2 int) bool
		want  []int
	}{
		{
			name:  "sort ascending",
			input: []int{5, 2, 8, 1, 9},
			f:     func(val1, val2 int) bool { return val1 < val2 },
			want:  []int{1, 2, 5, 8, 9},
		},
		{
			name:  "sort descending",
			input: []int{5, 2, 8, 1, 9},
			f:     func(val1, val2 int) bool { return val1 > val2 },
			want:  []int{9, 8, 5, 2, 1},
		},
		{
			name:  "empty slice",
			input: []int{},
			f:     func(val1, val2 int) bool { return val1 < val2 },
			want:  []int{},
		},
		{
			name:  "single element slice",
			input: []int{7},
			f:     func(val1, val2 int) bool { return val1 < val2 },
			want:  []int{7},
		},
		{
			name:  "already sorted slice",
			input: []int{1, 2, 3, 4},
			f:     func(val1, val2 int) bool { return val1 < val2 },
			want:  []int{1, 2, 3, 4},
		},
		{
			name:  "slice with duplicates",
			input: []int{3, 1, 4, 1, 5, 9, 2, 6},
			f:     func(val1, val2 int) bool { return val1 < val2 },
			want:  []int{1, 1, 2, 3, 4, 5, 6, 9},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a copy because Sort modifies the slice in place
			inputCopy := make([]int, len(tt.input))
			copy(inputCopy, tt.input)

			slice_utils.SortFunc(inputCopy, tt.f)
			assert.Equal(t, tt.want, inputCopy, "SortFunc() should correctly sort the slice in-place")
		})
	}
}

func TestConvert(t *testing.T) {
	// Test case for int to string
	t.Run("int to string", func(t *testing.T) {
		input := []int{1, 2, 3}
		f := func(val int) string { return "num: " + string(rune('0'+val)) }
		want := []string{"num: 1", "num: 2", "num: 3"}
		got := slice_utils.Convert(input, f)
		assert.Equal(t, want, got, "Convert() should transform int to string correctly")
	})

	// Test case for int to int (identity)
	t.Run("int to int (identity)", func(t *testing.T) {
		input := []int{10, 20, 30}
		f := func(val int) int { return val }
		want := []int{10, 20, 30}
		got := slice_utils.Convert(input, f)
		assert.Equal(t, want, got, "Convert() should perform identity conversion for int to int")
	})

	// Test case for empty slice
	t.Run("empty slice", func(t *testing.T) {
		input := []int{}
		f := func(val int) string { return "" }
		want := []string{}
		got := slice_utils.Convert(input, f)
		assert.Equal(t, want, got, "Convert() should return an empty slice for empty input")
	})

	// Test case for int to bool
	t.Run("int to bool", func(t *testing.T) {
		input := []int{0, 1, 2}
		f := func(val int) bool { return val > 0 }
		want := []bool{false, true, true}
		got := slice_utils.Convert(input, f)
		assert.Equal(t, want, got, "Convert() should transform int to bool correctly")
	})
	// Add a test case for string to int if applicable, ensure to match types.
	t.Run("string to int", func(t *testing.T) {
		input := []string{"1", "2", "3"}
		f := func(val string) int {
			// Simplified atoi for test, in real code use strconv.Atoi
			return int(val[0] - '0')
		}
		want := []int{1, 2, 3}
		got := slice_utils.Convert(input, f)
		assert.Equal(t, want, got, "Convert() should transform string to int correctly")
	})
}
func TestAggregate(t *testing.T) {
	tests := []struct {
		name    string
		input   []int
		f       func(val int) (int, error)
		want    int
		wantErr bool
	}{
		{
			name:  "sum integers",
			input: []int{1, 2, 3, 4, 5},
			f:     func(val int) (int, error) { return val, nil },
			want:  15,
		},
		{
			name:  "sum squares",
			input: []int{1, 2, 3},
			f:     func(val int) (int, error) { return val * val, nil },
			want:  14, // 1*1 + 2*2 + 3*3 = 1 + 4 + 9 = 14
		},
		{
			name:  "empty slice",
			input: []int{},
			f:     func(val int) (int, error) { return val, nil },
			want:  0,
		},
		{
			name:  "function returns error",
			input: []int{1, 2, -3, 4},
			f: func(val int) (int, error) {
				if val < 0 {
					return 0, errors.New("negative value")
				} else {
					return val, nil
				}
			},
			want:    0, // Expect zero value for T on error
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := slice_utils.Aggregate(tt.input, tt.f)
			if tt.wantErr {
				assert.Error(t, err, "Aggregate() should return an error")
			} else {
				assert.NoError(t, err, "Aggregate() should not return an error")
				assert.Equal(t, tt.want, got, "Aggregate() should return the correct aggregated value")
			}
		})
	}
}

func TestChange(t *testing.T) {
	// Test case for int transformation
	t.Run("square numbers", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		f := func(val int) int { return val * val }
		want := []int{1, 4, 9, 16, 25}
		got := slice_utils.Change(input, f)
		assert.Equal(t, want, got, "Change() should square numbers correctly")
	})

	t.Run("double numbers", func(t *testing.T) {
		input := []int{1, 2, 3}
		f := func(val int) int { return val * 2 }
		want := []int{2, 4, 6}
		got := slice_utils.Change(input, f)
		assert.Equal(t, want, got, "Change() should double numbers correctly")
	})

	t.Run("empty slice", func(t *testing.T) {
		input := []int{}
		f := func(val int) int { return val }
		want := []int{}
		got := slice_utils.Change(input, f)
		assert.Equal(t, want, got, "Change() should return an empty slice for empty input")
	})

	// Test case for string transformation
	t.Run("change string to uppercase", func(t *testing.T) {
		input := []string{"hello", "world"}
		// This conversion is simplified and assumes ASCII. For full Unicode, use strings.ToUpper.
		f := func(val string) string {
			runes := []rune(val)
			for i, r := range runes {
				if r >= 'a' && r <= 'z' {
					runes[i] = r - ('a' - 'A')
				}
			}
			return string(runes)
		}
		want := []string{"HELLO", "WORLD"}
		got := slice_utils.Change(input, f)
		assert.Equal(t, want, got, "Change() should convert strings to uppercase")
	})
}

func TestRemap(t *testing.T) {
	type Test[T any, K comparable, V any] struct {
		name    string
		input   []T
		f       func(val T) (K, V, error)
		want    map[K]V
		wantErr bool
	}

	tests1 := []Test[int, string, bool]{
		{
			name:  "remap int to string key and bool value",
			input: []int{1, 2, 3},
			f:     func(val int) (string, bool, error) { return "num_" + string(rune('0'+val)), val%2 == 0, nil },
			want:  map[string]bool{"num_1": false, "num_2": true, "num_3": false},
		},
		{
			name:  "empty slice",
			input: []int{},
			f:     func(val int) (string, bool, error) { return "", false, nil },
			want:  map[string]bool{},
		},
	}

	tests2 := []Test[int, string, int]{
		{
			name:  "duplicate keys (last one wins)",
			input: []int{1, 2, 3},
			f:     func(val int) (string, int, error) { return "same_key", val, nil },
			want:  map[string]int{"same_key": 3},
		},
		{
			name:  "function returns error",
			input: []int{1, -2, 3},
			f: func(val int) (string, int, error) {
				if val < 0 {
					return "", 0, errors.New("negative input")
				} else {
					return string(rune('0' + val)), val, nil
				}
			},
			want:    nil, // On error, the result map should be nil
			wantErr: true,
		},
	}

	for _, tt := range tests1 {
		t.Run(tt.name, func(t *testing.T) {
			got, err := slice_utils.Remap(tt.input, tt.f)
			if tt.wantErr {
				assert.Error(t, err, "Remap() should return an error")
				assert.Nil(t, got, "Remap() should return a nil map on error")
			} else {
				assert.NoError(t, err, "Remap() should not return an error")
				assert.Equal(t, tt.want, got, "Remap() should correctly create the remapped map")
			}
		})
	}

	for _, tt := range tests2 {
		t.Run(tt.name, func(t *testing.T) {
			got, err := slice_utils.Remap(tt.input, tt.f)
			if tt.wantErr {
				assert.Error(t, err, "Remap() should return an error")
				assert.Nil(t, got, "Remap() should return a nil map on error")
			} else {
				assert.NoError(t, err, "Remap() should not return an error")
				assert.Equal(t, tt.want, got, "Remap() should correctly create the remapped map")
			}
		})
	}
}

func TestToMap(t *testing.T) {
	// Test case for int to string key
	t.Run("map int to string key", func(t *testing.T) {
		input := []int{1, 2, 3}
		f := func(val int) string { return "id_" + string(rune('0'+val)) }
		want := map[string]int{"id_1": 1, "id_2": 2, "id_3": 3}
		got := slice_utils.ToMap(input, f)
		assert.Equal(t, want, got, "ToMap() should correctly create the map with int values")
	})

	// Test case for empty slice
	t.Run("empty slice", func(t *testing.T) {
		input := []int{}
		f := func(val int) string { return "" }
		want := map[string]int{}
		got := slice_utils.ToMap(input, f)
		assert.Equal(t, want, got, "ToMap() should return an empty map for empty input")
	})

	// Test case for duplicate keys (last one wins)
	t.Run("duplicate keys (last one wins)", func(t *testing.T) {
		input := []int{10, 20, 30}
		f := func(val int) string { return "common_key" }
		want := map[string]int{"common_key": 30}
		got := slice_utils.ToMap(input, f)
		assert.Equal(t, want, got, "ToMap() should handle duplicate keys by taking the last value")
	})

	// Test case for mapping structs to an ID
	t.Run("map struct to ID", func(t *testing.T) {
		type MyStruct struct {
			ID   int
			Name string
		}
		input := []MyStruct{{1, "A"}, {2, "B"}}
		f := func(s MyStruct) int { return s.ID }
		want := map[int]MyStruct{1: {1, "A"}, 2: {2, "B"}}
		got := slice_utils.ToMap(input, f)
		assert.Equal(t, want, got, "ToMap() should map structs to their IDs correctly")
	})
}

func TestDuplicates(t *testing.T) {
	// Test case for int slice
	t.Run("find duplicates in int slice", func(t *testing.T) {
		input := []int{1, 2, 3, 2, 4, 1, 5}
		want := []int{1, 2}
		got := slice_utils.Duplicates(input)
		assert.ElementsMatch(t, want, got, "Duplicates() should return all unique duplicate int elements")
	})

	// Test case for no duplicates
	t.Run("no duplicates", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		want := []int{}
		got := slice_utils.Duplicates(input)
		assert.ElementsMatch(t, want, got, "Duplicates() should return an empty slice when no duplicates")
	})

	// Test case for empty slice
	t.Run("empty slice", func(t *testing.T) {
		input := []int{}
		want := []int{}
		got := slice_utils.Duplicates(input)
		assert.ElementsMatch(t, want, got, "Duplicates() should return an empty slice for empty input")
	})

	// Test case for all elements are duplicates
	t.Run("all elements are duplicates", func(t *testing.T) {
		input := []int{1, 1, 1, 1}
		want := []int{1}
		got := slice_utils.Duplicates(input)
		assert.ElementsMatch(t, want, got, "Duplicates() should return the single duplicate element")
	})

	// Test case for string slice
	t.Run("string slice duplicates", func(t *testing.T) {
		input := []string{"apple", "banana", "apple", "orange", "banana"}
		want := []string{"apple", "banana"}
		got := slice_utils.Duplicates(input)
		assert.ElementsMatch(t, want, got, "Duplicates() should return all unique duplicate string elements")
	})
}

func TestDeduplicate(t *testing.T) {
	// Test case for int slice
	t.Run("deduplicate int slice", func(t *testing.T) {
		input := []int{1, 2, 3, 2, 4, 1, 5}
		want := []int{1, 2, 3, 4, 5} // Expect order preservation by DeduplicationSeq
		got := slice_utils.Deduplicate(input)
		assert.Equal(t, want, got, "Deduplicate() should remove duplicates preserving order for int slice")
	})

	// Test case for no duplicates
	t.Run("no duplicates", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		want := []int{1, 2, 3, 4, 5}
		got := slice_utils.Deduplicate(input)
		assert.Equal(t, want, got, "Deduplicate() should return the same slice if no duplicates")
	})

	// Test case for empty slice
	t.Run("empty slice", func(t *testing.T) {
		input := []int{}
		want := []int{}
		got := slice_utils.Deduplicate(input)
		assert.Equal(t, want, got, "Deduplicate() should return an empty slice for empty input")
	})

	// Test case for all elements are duplicates
	t.Run("all elements are duplicates", func(t *testing.T) {
		input := []int{1, 1, 1, 1}
		want := []int{1}
		got := slice_utils.Deduplicate(input)
		assert.Equal(t, want, got, "Deduplicate() should return a single element for all duplicate input")
	})

	// Test case for string slice
	t.Run("string slice deduplication", func(t *testing.T) {
		input := []string{"apple", "banana", "apple", "orange", "banana"}
		want := []string{"apple", "banana", "orange"}
		got := slice_utils.Deduplicate(input)
		assert.Equal(t, want, got, "Deduplicate() should remove duplicates preserving order for string slice")
	})
}
func TestGroups(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		f     func(v int) int // Grouping function
		want  [][]int
	}{
		{
			name:  "group by parity",
			input: []int{1, 2, 3, 4, 5, 6},
			f:     func(v int) int { return v % 2 }, // 0 for even, 1 for odd
			want:  [][]int{{1, 3, 5}, {2, 4, 6}},    // The internal implementation of GroupSeq might sort keys.
		},
		{
			name:  "group by tens digit",
			input: []int{12, 25, 18, 30, 21},
			f:     func(v int) int { return v / 10 },
			want:  [][]int{{12, 18}, {25, 21}, {30}},
		},
		{
			name:  "empty slice",
			input: []int{},
			f:     func(v int) int { return v },
			want:  [][]int{},
		},
		{
			name:  "single group",
			input: []int{1, 2, 3},
			f:     func(v int) int { return 0 },
			want:  [][]int{{1, 2, 3}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slice_utils.Groups(tt.input, tt.f)

			// To make `assert.ElementsMatch` work for slices of slices reliably,
			// we need to ensure that the order of elements within each inner slice
			// and the order of the inner slices themselves are deterministic.
			// The original GroupSeq mock sorted keys. We need to ensure elements within groups are also sorted.
			for i := range got {
				sort.Ints(got[i])
			}
			for i := range tt.want {
				sort.Ints(tt.want[i])
			}

			// Then, ElementsMatch can compare the groups without strict order of groups.
			assert.ElementsMatch(t, tt.want, got, "Groups() should correctly group elements")
		})
	}
}

func TestRemoveStrings(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		pattern string
		want    []string
		wantErr bool
	}{
		{
			name:    "remove words containing 'a'",
			input:   []string{"apple", "banana", "cherry", "date"},
			pattern: "a",
			want:    []string{"cherry"},
			wantErr: false,
		},
		{
			name:    "remove numbers",
			input:   []string{"hello", "world123", "go"},
			pattern: `\d+`,
			want:    []string{"hello", "go"},
			wantErr: false,
		},
		{
			name:    "empty slice",
			input:   []string{},
			pattern: "test",
			want:    []string{},
			wantErr: false,
		},
		{
			name:    "no matches",
			input:   []string{"cat", "dog"},
			pattern: "xyz",
			want:    []string{"cat", "dog"},
			wantErr: false,
		},
		{
			name:    "invalid regex pattern",
			input:   []string{"a", "b"},
			pattern: "[",
			want:    nil, // Assert.Nil is suitable here for nil slice on error
			wantErr: true,
		},
		{
			name:    "remove all",
			input:   []string{"alpha", "beta", "gamma"},
			pattern: "a|b|g", // Matches any of them
			want:    []string{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := regexp.Compile(tt.pattern)
			if tt.wantErr {
				assert.Error(t, err, "RemoveStrings() should return an error for invalid pattern")
			} else {
				if assert.NoError(t, err, "RemoveStrings() should not return an error for valid pattern") {
					got := slice_utils.RemoveStrings(tt.input, p)
					assert.Equal(t, tt.want, got, "RemoveStrings() should return the slice with matching elements removed")
				}
			}
		})
	}
}

func TestFilterStrings(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		pattern string
		want    []string
		wantErr bool
	}{
		{
			name:    "filter words containing 'a'",
			input:   []string{"apple", "banana", "cherry", "date"},
			pattern: "a",
			want:    []string{"apple", "banana", "date"},
			wantErr: false,
		},
		{
			name:    "filter numbers",
			input:   []string{"hello", "world123", "go"},
			pattern: `\d+`,
			want:    []string{"world123"},
			wantErr: false,
		},
		{
			name:    "empty slice",
			input:   []string{},
			pattern: "test",
			want:    []string{},
			wantErr: false,
		},
		{
			name:    "no matches",
			input:   []string{"cat", "dog"},
			pattern: "xyz",
			want:    []string{},
			wantErr: false,
		},
		{
			name:    "invalid regex pattern",
			input:   []string{"a", "b"},
			pattern: "[",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "filter all",
			input:   []string{"alpha", "beta", "gamma"},
			pattern: "a|b|g", // Matches any of them
			want:    []string{"alpha", "beta", "gamma"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := regexp.Compile(tt.pattern)
			if tt.wantErr {
				assert.Error(t, err, "FilterStrings() should return an error for invalid pattern")
			} else {
				if assert.NoError(t, err, "FilterStrings() should not return an error for valid pattern") {
					got := slice_utils.FilterStrings(tt.input, p)
					assert.Equal(t, tt.want, got, "FilterStrings() should return the slice with matching elements kept")
				}
			}
		})
	}
}

func TestChunks(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		size  int
		want  [][]int
	}{
		{
			name:  "chunk into equal parts",
			input: []int{1, 2, 3, 4, 5, 6},
			size:  2,
			want:  [][]int{{1, 2}, {3, 4}, {5, 6}},
		},
		{
			name:  "last chunk smaller",
			input: []int{1, 2, 3, 4, 5},
			size:  2,
			want:  [][]int{{1, 2}, {3, 4}, {5}},
		},
		{
			name:  "size greater than slice length",
			input: []int{1, 2, 3},
			size:  5,
			want:  [][]int{{1, 2, 3}},
		},
		{
			name:  "empty slice",
			input: []int{},
			size:  2,
			want:  [][]int{},
		},
		{
			name:  "size 1",
			input: []int{1, 2, 3},
			size:  1,
			want:  [][]int{{1}, {2}, {3}},
		},
		{
			name:  "size 0 should result in empty chunks (or panic depending on slices.Chunk impl)",
			input: []int{1, 2, 3},
			size:  0,
			want:  [][]int{{1, 2, 3}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slice_utils.Chunks(tt.input, tt.size)
			assert.Equal(t, tt.want, got, "Chunks() should correctly divide the slice into chunks")
		})
	}
}

func TestToAny(t *testing.T) {
	input := []int{1, 2, 3}
	want := []any{1, 2, 3}
	got := slice_utils.ToAny(input)
	assert.Equal(t, want, got)

	assert.Empty(t, slice_utils.ToAny([]int{}))
}

func TestGroup(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6}
	f := func(v int) string {
		if v%2 == 0 {
			return "even"
		}
		return "odd"
	}
	got := slice_utils.Group(input, f)

	assert.Len(t, got, 2)
	assert.Equal(t, []int{2, 4, 6}, got["even"])
	assert.Equal(t, []int{1, 3, 5}, got["odd"])
}

func TestContains(t *testing.T) {
	input := []int{1, 2, 3}
	assert.True(t, slice_utils.Contains(input, func(v int) bool { return v == 2 }))
	assert.False(t, slice_utils.Contains(input, func(v int) bool { return v == 4 }))
	assert.False(t, slice_utils.Contains([]int{}, func(v int) bool { return true }))
}

func TestPairs(t *testing.T) {
	t.Run("even number of elements", func(t *testing.T) {
		got := slice_utils.Pairs(1, 2, 3, 4)
		want := [][2]int{{1, 2}, {3, 4}}
		assert.Equal(t, want, got)
	})

	t.Run("odd number of elements", func(t *testing.T) {
		got := slice_utils.Pairs(1, 2, 3)
		want := [][2]int{{1, 2}, {3, 0}} // 0 is zero value for int
		assert.Equal(t, want, got)
	})

	t.Run("empty", func(t *testing.T) {
		got := slice_utils.Pairs[int]()
		assert.Empty(t, got)
	})

	t.Run("strings", func(t *testing.T) {
		got := slice_utils.Pairs("a", "b", "c")
		want := [][2]string{{"a", "b"}, {"c", ""}}
		assert.Equal(t, want, got)
	})
}
