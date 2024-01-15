package utils

import (
	"log/slog"
	"reflect"

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
