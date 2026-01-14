// Copyright 2026 Zauberhaus
// Licensed to Zauberhaus under one or more agreements.
// Zauberhaus licenses this file to you under the Apache 2.0 License.
// See the LICENSE file in the project root for more information.

package slice_utils

import (
	"cmp"
	"fmt"
	"iter"
	"regexp"
	"slices"

	"hash/maphash"
)

func FilterSeq[S any](s iter.Seq[S], fn func(S) bool) iter.Seq[S] {
	return func(yield func(s S) bool) {
		for v := range s {
			if fn(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func RemoveSeq[S comparable](s iter.Seq[S], g iter.Seq[S]) iter.Seq[S] {
	return func(yield func(s S) bool) {
		for v1 := range s {
			found := false

			for v2 := range g {
				if v1 == v2 {
					found = true
					break
				}
			}

			if !found {
				if !yield(v1) {
					return
				}
			}
		}

	}
}

func PatternSeq[S any](s iter.Seq[S], pattern *regexp.Regexp) iter.Seq[S] {
	return func(yield func(s S) bool) {
		for v := range s {
			var txt string
			switch o := any(v).(type) {
			case string:
				txt = o
			case fmt.Stringer:
				txt = o.String()
			default:
				txt = fmt.Sprintf("%v", o)
			}

			if pattern.MatchString(txt) {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func StringPatternSeq[S any](s iter.Seq[S], pattern string) iter.Seq[S] {
	return func(yield func(s S) bool) {
		for v := range s {
			var txt string
			switch o := any(v).(type) {
			case string:
				txt = o
			case fmt.Stringer:
				txt = o.String()
			default:
				txt = fmt.Sprintf("%v", o)
			}

			if txt == pattern {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func DuplicateSeq[V comparable](s iter.Seq[V]) iter.Seq[V] {
	m := map[V]int{}

	return func(yield func(s V) bool) {
		for v := range s {
			if cnt, ok := m[v]; ok {
				m[v] = cnt + 1
				if cnt == 1 {
					if !yield(v) {
						return
					}
				}
			} else {
				m[v] = 1
			}
		}
	}
}

func DeduplicationSeq[V comparable](s iter.Seq[V]) iter.Seq[V] {
	m := map[V]bool{}

	return func(yield func(s V) bool) {
		for v := range s {
			if _, ok := m[v]; ok {
				continue
			} else {
				m[v] = true
				if !yield(v) {
					return
				}
			}
		}
	}
}

func HashSeq[E comparable](s iter.Seq[E]) iter.Seq2[uint64, E] {
	var h maphash.Hash

	return func(yield func(uint64, E) bool) {
		for v := range s {
			h.Reset()
			maphash.WriteComparable(&h, v)
			if !yield(h.Sum64(), v) {
				return
			}
		}
	}
}

func GroupSeq[S ~[]E, E any, H comparable](s iter.Seq[E], fn func(v E) H) iter.Seq[S] {
	groups := map[H]S{}

	for v := range s {
		h := fn(v)
		g := groups[h]
		g = append(g, v)
		groups[h] = g
	}

	return func(yield func(S) bool) {
		for _, v := range groups {
			if !yield(v) {
				return
			}
		}
	}
}

func CountSeq[S any](s iter.Seq[S]) int {
	var result int

	for range s {
		result++
	}

	return result
}

func SumFuncSeq[S any, T cmp.Ordered](s iter.Seq[S], fn func(S) (T, error)) (T, error) {
	var result T

	for v := range s {
		val, err := fn(v)
		if err != nil {
			return *new(T), err
		}

		result += val
	}

	return result, nil
}

func SumSeq[S cmp.Ordered](s iter.Seq[S]) S {
	var result S

	items := slices.Collect(s)
	slices.Sort(items)

	for _, v := range items {
		result += v
	}

	return result
}

func IsEmptySeq[S any](s iter.Seq[S]) bool {
	for range s {
		return false
	}

	return true
}

func ReplaceFuncSeq[S any](s iter.Seq[S], fn func(val S) S) iter.Seq[S] {
	return func(yield func(s S) bool) {
		for v := range s {
			if !yield(fn(v)) {
				return
			}
		}
	}
}

func ReplaceSeq[S comparable](s iter.Seq[S], g map[S]S) iter.Seq[S] {
	return func(yield func(s S) bool) {
		for v := range s {
			if r, ok := g[v]; ok {
				if !yield(r) {
					return
				}
			} else {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func ConvertSeq[S any, T any](s iter.Seq[S], fn func(val S) T) iter.Seq[T] {
	return func(yield func(s T) bool) {
		for v := range s {
			if !yield(fn(v)) {
				return
			}
		}
	}
}

func AnySeq[S any](s iter.Seq[S]) iter.Seq[any] {
	return func(yield func(s any) bool) {
		for v := range s {
			if !yield(any(v)) {
				return
			}
		}
	}
}
