// Copyright 2026 Zauberhaus
// Licensed to Zauberhaus under one or more agreements.
// Zauberhaus licenses this file to you under the Apache 2.0 License.
// See the LICENSE file in the project root for more information.

package slice_utils

import (
	"cmp"
	"maps"
	"reflect"
	"slices"

	"regexp"
	"sort"
)

func Select[Slice ~[]V, V any](slice Slice, f func(val V) bool) Slice {
	return slices.Collect(FilterSeq(slices.Values(slice), f))
}

func Count[Slice ~[]V, V any](slice Slice, f func(val V) bool) int {
	return CountSeq(FilterSeq(slices.Values(slice), f))
}

func Empty[Slice ~[]V, V any](slice Slice, f func(val V) bool) bool {
	return IsEmptySeq(FilterSeq(slices.Values(slice), f))
}

func Delete[Slice ~[]V, V comparable](slice Slice, vals ...V) Slice {

	for i, v := range slice {
		for _, val := range vals {
			if v == val {
				return append(slice[:i], slice[i+1:]...)
			}
		}
	}

	return slice
}

func SortFunc[Slice ~[]V, V any](slice Slice, f func(val1 V, val2 V) bool) {
	sort.Slice(slice, func(i, j int) bool {
		v1 := slice[i]
		v2 := slice[j]

		return f(v1, v2)
	})
}

func To[T any, V any, Slice ~[]V](slice Slice) []T {
	return slices.Collect(ConvertSeq(slices.Values(slice), func(val V) T {
		t := reflect.TypeFor[T]()
		v := reflect.ValueOf(val)

		if v.Type() == t {
			return any(val).(T)
		}

		isPtr := t.Kind() == reflect.Pointer
		if isPtr {
			t = t.Elem()
		}

		if v.Kind() == reflect.Pointer {
			v = v.Elem()
		}

		if v.CanConvert(t) {
			v2 := v.Convert(t)
			if isPtr {
				p := reflect.New(t)
				p.Elem().Set(v2)

				if p.CanInterface() {
					return p.Interface().(T)
				}
			} else {
				if v2.CanInterface() {
					return v2.Interface().(T)
				}
			}
		}

		return *new(T)
	}))
}

func ToAny[Slice ~[]V, V any](slice Slice) []any {
	r := slices.Collect(AnySeq(slices.Values(slice)))
	if r == nil {
		return []any{}
	}

	return r
}

func Convert[Slice ~[]V, V any, T any](slice Slice, f func(val1 V) T) []T {
	r := slices.Collect(ConvertSeq(slices.Values(slice), f))
	if r == nil {
		return []T{}
	}

	return r
}

func Aggregate[Slice ~[]V, V any, T cmp.Ordered](slice Slice, f func(val1 V) (T, error)) (T, error) {
	return SumFuncSeq(slices.Values(slice), f)
}

func Change[Slice ~[]V, V any](slice Slice, f func(val1 V) V) Slice {
	r := slices.Collect(ReplaceFuncSeq(slices.Values(slice), f))
	if r == nil {
		return Slice{}
	}

	return r
}

func Remap[Slice ~[]V, V any, K comparable, T any](slice Slice, f func(val V) (K, T, error)) (map[K]T, error) {
	result := map[K]T{}

	for _, v := range slice {
		k, v, err := f(v)
		if err != nil {
			return nil, err
		}

		result[k] = v
	}

	return result, nil
}

func ToMap[Slice ~[]V, K comparable, V any](slice Slice, f func(val V) K) map[K]V {
	result := map[K]V{}

	for _, v := range slice {
		k := f(v)
		result[k] = v
	}

	return result
}

func Duplicates[Slice ~[]V, V comparable](slice Slice) Slice {
	tmp := map[V]int{}
	for _, v := range slice {
		if cnt, ok := tmp[v]; ok {
			tmp[v] = cnt + 1
		} else {
			tmp[v] = 1
		}
	}

	result := Slice{}

	for k, v := range tmp {
		if v > 1 {
			result = append(result, k)
		}
	}

	return result
}

func Deduplicate[Slice ~[]V, V comparable](s Slice) Slice {
	r := slices.Collect(DeduplicationSeq(slices.Values(s)))
	if r == nil {
		return Slice{}
	}
	return r
}

func Groups[Slice ~[]V, V any, K cmp.Ordered](s Slice, f func(v V) K) []Slice {
	return slices.Collect(GroupSeq[Slice](slices.Values(s), f))
}

func FilterStrings(s []string, p *regexp.Regexp) []string {
	values := slices.Values(s)

	r := slices.Collect(PatternSeq(values, p))
	if r == nil {
		return []string{}
	}

	return r
}

func RemoveStrings(s []string, p *regexp.Regexp) []string {
	values := slices.Values(s)

	r := slices.Collect(RemoveSeq(values, PatternSeq(values, p)))
	if r == nil {
		return []string{}
	}

	return r
}

func Chunks[Slice ~[]V, V any](slice Slice, size int) []Slice {
	if size < 1 {
		if len(slice) == 0 {
			return []Slice{}
		}

		return []Slice{slice}
	}

	r := slices.Collect(slices.Chunk(slice, size))
	if r == nil {
		return []Slice{}
	}
	return r
}

func Group[S ~[]E, E any, H cmp.Ordered](s S, f func(v E) H) map[H]S {
	groups := map[H]S{}

	for _, v := range s {
		key := f(v)

		if group, ok := groups[key]; ok {
			groups[key] = append(group, v)
		} else {
			groups[key] = []E{v}
		}
	}

	keys := slices.Collect(maps.Keys(groups))
	sort.Slice(keys, func(i, j int) bool {
		k1 := keys[i]
		k2 := keys[j]

		return k1 > k2
	})

	return groups
}

func Contains[V any](slice []V, f func(val V) bool) bool {
	return slices.ContainsFunc(slice, f)
}

func Pairs[T any](values ...T) [][2]T {
	result := [][2]T{}

	for i := 0; i < len(values); i += 2 {
		key := values[i]
		value := *new(T)

		if i+1 < len(values) {
			value = values[i+1]
		}

		result = append(result, [2]T{key, value})
	}

	return result
}
