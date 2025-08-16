package internal

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm/logger"
)

type GormZeroLogAdapter struct {
	Level int
}

func (l GormZeroLogAdapter) LogMode(logger.LogLevel) logger.Interface {
	return l
}

func (l GormZeroLogAdapter) Error(ctx context.Context, msg string, opts ...interface{}) {
	zerolog.Ctx(ctx).Error().Msg(fmt.Sprintf(msg, opts...))
}

func (l GormZeroLogAdapter) Warn(ctx context.Context, msg string, opts ...interface{}) {
	zerolog.Ctx(ctx).Warn().Msg(fmt.Sprintf(msg, opts...))
}

func (l GormZeroLogAdapter) Info(ctx context.Context, msg string, opts ...interface{}) {
	zerolog.Ctx(ctx).Info().Msg(fmt.Sprintf(msg, opts...))
}

func (l GormZeroLogAdapter) Trace(ctx context.Context, begin time.Time, f func() (string, int64), err error) {
	// zl := zerolog.Ctx(ctx)
	var event *zerolog.Event

	if err != nil {
		event = log.Debug() // zl.Debug()
	} else {
		event = log.Trace() // zl.Trace()
	}

	var dur_key string

	switch zerolog.DurationFieldUnit {
	case time.Nanosecond:
		dur_key = "elapsed_ns"
	case time.Microsecond:
		dur_key = "elapsed_us"
	case time.Millisecond:
		dur_key = "elapsed_ms"
	case time.Second:
		dur_key = "elapsed"
	case time.Minute:
		dur_key = "elapsed_min"
	case time.Hour:
		dur_key = "elapsed_hr"
	default:
		log.Error().Interface( // zl.Error().Interface(
			"zerolog.DurationFieldUnit",
			zerolog.DurationFieldUnit,
		).Msg("gormzerolog encountered a mysterious, unknown value for DurationFieldUnit")
		dur_key = "elapsed_"
	}

	event.Dur(dur_key, time.Since(begin))

	sql, rows := f()
	if sql != "" {
		event.Str("sql", sql)
	}
	if rows > -1 {
		event.Int64("rows", rows)
	}

	event.Send()
}

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
