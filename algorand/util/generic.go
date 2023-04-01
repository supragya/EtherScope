package util

import (
	"reflect"

	"golang.org/x/exp/constraints"
)

func Map[T any, U any](vals []T, fn func(T) U) []U {
	var us []U
	for _, v := range vals {
		us = append(us, fn(v))
	}
	return us
}

func Filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func Cast[T any, S any](arr []T) []S {
	result := make([]S, len(arr))
	targetType := reflect.TypeOf(arr[0])
	for i, v := range arr {
		result[i] = reflect.ValueOf(v).Convert(targetType).Interface().(S)
	}
	return result
}

func Contains[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func MapMax[T comparable, U constraints.Ordered](m map[T]U) (T, U) {
	var maxValue U
	var maxKey T
	for k, v := range m {
		if v > maxValue {
			maxValue = v
			maxKey = k
		}
	}

	return maxKey, maxValue
}
