package internal

import (
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

func ParseUint16(number string, defaultValue uint16) uint16 {
	result, err := strconv.ParseUint(number, 0, 16)
	if err != nil {
		return defaultValue
	}
	return uint16(result)
}

func ParseBool(value string) bool {
	if value == "" {
		return false
	}
	return strings.ToUpper(value) == "TRUE" ||
		strings.ToUpper(value) == "T" ||
		value == "1" ||
		strings.ToUpper(value) == "YES" ||
		strings.ToUpper(value) == "Y"
}

func ParseInt(number string, defaultValue int) int {
	result, err := strconv.ParseInt(number, 0, 32)
	if err != nil {
		return defaultValue
	}
	return int(result)
}

func ParseUnitDuration(durationUnit string) time.Duration {
	switch strings.ToUpper(durationUnit) {
	case "H":
		return time.Hour
	case "M":
		return time.Minute
	case "S":
		return time.Second
	}
	return 0
}

var durationPattern = regexp.MustCompile(`(\d+)\s*([HhMmSs])`)

func ParseDuration(duration string, defaultValue time.Duration) time.Duration {
	var result time.Duration = 0
	for _, g := range durationPattern.FindAllStringSubmatch(duration, -1) {
		d := time.Duration(ParseInt(g[1], 1)) * ParseUnitDuration(g[2])
		result += d
	}

	return result
}

func FirstKey[K comparable, V any](dictionary map[K]V) V {
	for k := range dictionary {
		return dictionary[k]
	}
	var v V
	return v
}

func JoinUrl(base string, parts ...string) string {
	result, err := url.JoinPath(base, parts...)
	if err != nil {
		return base
	}
	return result
}

func ExistPermit(permitsList []string, permit string) bool {
	for _, x := range permitsList {
		if x == permit {
			return true
		}
	}

	return false
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func Mapping[T1 any, T2 any](slice []T1, mapperFunction func(T1) T2) []T2 {

	result := make([]T2, len(slice))

	for i, e := range slice {
		result[i] = mapperFunction(e)
	}

	return result
}

func MappingToMap[T1 any, T2 comparable, T3 any](slice []T1, mapperFunction func(T1) (T2, T3)) map[T2]T3 {

	result := make(map[T2]T3, len(slice))

	for _, e := range slice {
		key, value := mapperFunction(e)
		result[key] = value
	}

	return result
}

func MappingToMultiMap[T1 any, T2 comparable, T3 any](slice []T1, mapperFunction func(T1) (T2, T3)) map[T2][]T3 {

	result := make(map[T2][]T3, len(slice))

	for _, e := range slice {
		key, value := mapperFunction(e)
		values, ok := result[key]
		if !ok {
			values = make([]T3, 0, 4)
		}
		values = append(values, value)
		result[key] = values
	}

	return result
}

func Remove[T comparable](slice []T, element T) []T {
	pos := slices.Index(slice, element)
	if pos > -1 {
		return slices.Delete(slice, pos, pos+1)
	}
	return slice
}

func RemoveFunc[T any](slice []T, elementSelector func(T) bool) []T {
	pos := slices.IndexFunc(slice, elementSelector)
	if pos > -1 {
		return slices.Delete(slice, pos, pos+1)
	}
	return slice
}

func Any[T comparable](slice []T, predicate func(T) bool) T {
	var result T
	for _, v := range slice {
		if predicate(v) {
			return v
		}
	}
	return result
}

func Last[T any](slice []T) T {
	var result T
	if len(slice) > 0 {
		return slice[len(slice)-1]
	}
	return result
}
