package helpers

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
)

func ShuffleInt64s(slice []int64) []int64 {

	for i := range slice {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}

	return slice
}

func StringToSlice(s string, glue string) (ret []string) {

	if len(s) > 0 {
		ret = strings.Split(s, glue)
	}
	return ret
}

func StringSliceToIntSlice(in []string) (ret []int) {

	for _, v := range in {
		v = strings.TrimSpace(v)
		i, err := strconv.Atoi(v)
		if err == nil {
			ret = append(ret, i)
		}
	}
	return ret
}

func StringSliceToInt64Slice(in []string) (ret []int64) {

	for _, v := range in {
		v = strings.TrimSpace(v)
		i, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			ret = append(ret, i)
		}
	}
	return ret
}

func IntSliceToStringSlice(in []int) (ret []string) {

	for _, v := range in {
		ret = append(ret, strconv.Itoa(v))
	}
	return ret
}

func SliceHasInt(slice []int, i int) bool {
	for _, v := range slice {
		if v == i {
			return true
		}
	}
	return false
}

func SliceHasInt64(slice []int64, i int64) bool {
	for _, v := range slice {
		if v == i {
			return true
		}
	}
	return false
}

func SliceHasString(i string, slice []string) bool {
	for _, v := range slice {
		if v == i {
			return true
		}
	}
	return false
}

// todo, keep order by doing https://codereview.stackexchange.com/questions/191238/return-unique-items-in-a-go-slice
func UniqueInt(arg []int) []int {

	tempMap := make(map[int]bool)

	for idx := range arg {
		tempMap[arg[idx]] = true
	}

	tempSlice := make([]int, 0)
	for key := range tempMap {
		tempSlice = append(tempSlice, key)
	}
	return tempSlice
}

func UniqueInt64(arg []int64) []int64 {

	tempMap := make(map[int64]bool)

	for idx := range arg {
		tempMap[arg[idx]] = true
	}

	tempSlice := make([]int64, 0)
	for key := range tempMap {
		tempSlice = append(tempSlice, key)
	}
	return tempSlice
}

func UniqueString(strings []string) []string {

	var keys = make(map[string]bool, len(strings))
	var list []string

	for _, v := range strings {
		if _, value := keys[v]; !value {
			keys[v] = true
			list = append(list, v)
		}
	}

	return list
}

func FirstInts(slice []int, x int) []int {

	x = int(math.Min(float64(x), float64(len(slice))))
	return slice[0:x]
}

func Filter(ss []string, filter func(string) bool) (ret []string) {
	for _, s := range ss {
		if filter(s) {
			ret = append(ret, s)
		}
	}
	return
}
