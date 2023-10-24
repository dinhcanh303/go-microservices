package utils

import (
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

func StringToNullUUID(uuidString string) uuid.NullUUID {
	parsedUUID, err := uuid.Parse(uuidString)
	if err != nil {
		slog.Warn("Parsed UUID failed")
	}
	return uuid.NullUUID{
		UUID:  parsedUUID,
		Valid: true,
	}
}
func StringToUUID(uuidString string) uuid.UUID {
	paserdUUID, err := uuid.Parse(uuidString)
	if err != nil {
		slog.Warn("Parsed UUID failed")
	}
	return paserdUUID
}
