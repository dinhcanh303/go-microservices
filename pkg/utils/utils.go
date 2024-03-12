package utils

import (
	"fmt"
	"log/slog"
	"reflect"
	"strconv"

	st "github.com/golang/protobuf/ptypes/struct"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

func StringToNullUUID(uuidString string) (uuid.NullUUID, error) {
	parsedUUID, err := uuid.Parse(uuidString)
	if err != nil {
		return uuid.NullUUID{
			UUID:  parsedUUID,
			Valid: false,
		}, errors.Wrap(err, "Parsed UUID failed")
	}
	return uuid.NullUUID{
		UUID:  parsedUUID,
		Valid: true,
	}, nil
}
func StringToUUID(uuidString string) (uuid.UUID, error) {
	parsedUUID, err := uuid.Parse(uuidString)
	if err != nil {
		return uuid.UUID{}, errors.Wrap(err, "Parsed UUID failed")
	}
	return parsedUUID, nil
}

func StringToNullUUIDNormal(uuidString string) uuid.NullUUID {
	parsedUUID, err := uuid.Parse(uuidString)
	if err != nil {
		slog.Warn("Parsed UUID failed")
	}
	return uuid.NullUUID{
		UUID:  parsedUUID,
		Valid: true,
	}
}
func StringToUUIDNormal(uuidString string) uuid.UUID {
	parsedUUID, err := uuid.Parse(uuidString)
	if err != nil {
		slog.Warn("Parsed UUID failed")
	}
	return parsedUUID
}
func ConvertArUUIDToArString(uuids []uuid.UUID) []string {
	results := make([]string, 0)
	for _, uuid := range uuids {
		results = append(results, uuid.String())
	}
	return results
}

func ConvertArStringToArUUID(strings []string) ([]uuid.UUID, error) {
	uuids := make([]uuid.UUID, 0)
	for _, str := range strings {
		uuid, err := StringToUUID(str)
		if err != nil {
			return nil, errors.Wrap(err, "Converted UUID failed")
		}
		uuids = append(uuids, uuid)
	}
	return uuids, nil
}
func ConvertArStringToArNullUUID(strings []string) ([]uuid.NullUUID, error) {
	uuids := make([]uuid.NullUUID, 0)
	for _, str := range strings {
		uuid, err := StringToNullUUID(str)
		if err != nil {
			return nil, errors.Wrap(err, "Converted UUID failed")
		}
		uuids = append(uuids, uuid)
	}
	return uuids, nil
}
func LoadFileEnvOnLocal() error {

	err := godotenv.Load("../../.test.env")
	if err != nil {
		return err
	}
	return nil
}
func Contains[T any](arr []T, x T) bool {
	for _, v := range arr {
		if reflect.ValueOf(v) == reflect.ValueOf(x) {
			return true
		}
	}
	return false
}
func ContainsFunc[T any](arr []T, predicate func(T) bool) bool {
	for _, v := range arr {
		if predicate(v) {
			return true
		}
	}
	return false
}
func String(n int32) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(n)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}
func FormatInt(n int32) string {
	return strconv.FormatInt(int64(n), 10)
}

func UniqueSlice[T comparable](inputSlice []T) []T {
	uniqueSlice := make([]T, 0, len(inputSlice))
	seen := make(map[T]bool, len(inputSlice))
	for _, ele := range inputSlice {
		if !seen[ele] {
			uniqueSlice = append(uniqueSlice, ele)
			seen[ele] = true
		}
	}
	return uniqueSlice
}
func RemoveSlice[T comparable](inputSlice *[]T, removeItems ...T) []T {
	results := make([]T, 0, len(*inputSlice))
	// var end bool
	return results
}
func ToStruct(v map[string]interface{}) *st.Struct {
	size := len(v)
	if size == 0 {
		return nil
	}
	fields := make(map[string]*st.Value, size)
	for k, v := range v {
		fields[k] = ToValue(v)
	}
	return &st.Struct{
		Fields: fields,
	}
}

// ToValue converts an interface{} to a ptypes.Value
func ToValue(v interface{}) *st.Value {
	switch v := v.(type) {
	case bool:
		return &st.Value{
			Kind: &st.Value_BoolValue{
				BoolValue: v,
			},
		}
	case int:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case int8:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case int32:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case int64:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint8:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint32:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint64:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case float32:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case float64:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: v,
			},
		}
	case string:
		return &st.Value{
			Kind: &st.Value_StringValue{
				StringValue: v,
			},
		}
	case error:
		return &st.Value{
			Kind: &st.Value_StringValue{
				StringValue: v.Error(),
			},
		}
	case nil:
		return &st.Value{
			Kind: &st.Value_NullValue{
				NullValue: st.NullValue_NULL_VALUE,
			},
		}
	default:
		// Fallback to reflection for other types
		return toValue(reflect.ValueOf(v))
	}
}
func toValue(v reflect.Value) *st.Value {
	switch v.Kind() {
	case reflect.Bool:
		return &st.Value{
			Kind: &st.Value_BoolValue{
				BoolValue: v.Bool(),
			},
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: float64(v.Int()),
			},
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: float64(v.Uint()),
			},
		}
	case reflect.Float32, reflect.Float64:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: v.Float(),
			},
		}
	case reflect.Ptr:
		if v.IsNil() {
			return nil
		}
		return toValue(reflect.Indirect(v))
	case reflect.Array, reflect.Slice:
		size := v.Len()
		if size == 0 {
			return nil
		}
		values := make([]*st.Value, size)
		for i := 0; i < size; i++ {
			values[i] = toValue(v.Index(i))
		}
		return &st.Value{
			Kind: &st.Value_ListValue{
				ListValue: &st.ListValue{
					Values: values,
				},
			},
		}
	case reflect.Struct:
		t := v.Type()
		size := v.NumField()
		if size == 0 {
			return nil
		}
		fields := make(map[string]*st.Value, size)
		for i := 0; i < size; i++ {
			name := t.Field(i).Name
			// Better way?
			if len(name) > 0 && 'A' <= name[0] && name[0] <= 'Z' {
				fields[name] = toValue(v.Field(i))
			}
		}
		if len(fields) == 0 {
			return nil
		}
		return &st.Value{
			Kind: &st.Value_StructValue{
				StructValue: &st.Struct{
					Fields: fields,
				},
			},
		}
	case reflect.Map:
		keys := v.MapKeys()
		if len(keys) == 0 {
			return nil
		}
		fields := make(map[string]*st.Value, len(keys))
		for _, k := range keys {
			if k.Kind() == reflect.String {
				fields[k.String()] = toValue(v.MapIndex(k))
			}
		}
		if len(fields) == 0 {
			return nil
		}
		return &st.Value{
			Kind: &st.Value_StructValue{
				StructValue: &st.Struct{
					Fields: fields,
				},
			},
		}
	default:
		// Last resort
		return &st.Value{
			Kind: &st.Value_StringValue{
				StringValue: fmt.Sprint(v),
			},
		}
	}
}
