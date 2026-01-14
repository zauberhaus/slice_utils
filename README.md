# slice_utils

`slice_utils` is a Go library providing a collection of utility functions for working with slices and Go 1.23+ iterators (`iter.Seq`). It offers functional-style operations like filtering, mapping, grouping, and deduplication.

## Installation

```bash
go get github.com/zauberhaus/slice_utils
```

## Features

### Slice Operations

Helper functions for common slice manipulations.

*   **Filtering & Selection**: `Select`, `Delete`, `FilterStrings`, `RemoveStrings`
*   **Transformation**: `Convert`, `Change`, `ToAny`
*   **Aggregation**: `Count`, `Aggregate`, `Empty`, `Contains`
*   **Maps**: `ToMap`, `Remap`, `Group`
*   **Organization**: `SortFunc`, `Chunks`, `Pairs`
*   **Uniqueness**: `Duplicates`, `Deduplicate`

### Iterator Sequences (Go 1.23+)

Utilities for working with `iter.Seq`.

*   **Filtering**: `FilterSeq`, `RemoveSeq`, `PatternSeq`
*   **Transformation**: `ConvertSeq`, `ReplaceFuncSeq`, `ReplaceSeq`, `AnySeq`
*   **Aggregation**: `CountSeq`, `SumSeq`, `SumFuncSeq`, `IsEmptySeq`
*   **Grouping & Hashing**: `GroupSeq`, `HashSeq`
*   **Uniqueness**: `DuplicateSeq`, `DeduplicationSeq`

## Usage

### Slice Examples

```go
import (
    "fmt"
    "github.com/zauberhaus/slice_utils"
)

func main() {
    // Filtering
    nums := []int{1, 2, 3, 4, 5}
    evens := slice_utils.Select(nums, func(n int) bool {
        return n%2 == 0
    })
    fmt.Println(evens) // [2 4]

    // Conversion
    strs := slice_utils.Convert(nums, func(n int) string {
        return fmt.Sprintf("val-%d", n)
    })
    fmt.Println(strs) // [val-1 val-2 val-3 val-4 val-5]

    // Grouping
    groups := slice_utils.Group(nums, func(n int) string {
        if n%2 == 0 {
            return "even"
        }
        return "odd"
    })
    fmt.Println(groups["even"]) // [2 4]
}
```

### Sequence Examples

```go
import (
    "fmt"
    "slices"
    "github.com/zauberhaus/slice_utils"
)

func main() {
    nums := []int{1, 2, 3, 4, 5}
    
    // Create a sequence
    seq := slices.Values(nums)
    
    // Filter sequence
    filtered := slice_utils.FilterSeq(seq, func(n int) bool {
        return n > 2
    })
    
    // Collect results
    result := slices.Collect(filtered)
    fmt.Println(result) // [3 4 5]
}
```
