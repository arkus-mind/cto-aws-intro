package repository

import (
	"github.com/google/uuid"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func GenerateUUID() uuid.UUID {
	uuidNew := uuid.New()
	return uuidNew
}

func IsValidUUID(u string) bool {
	if len(u) == 0 || u == "00000000-0000-0000-0000-000000000000" {
		return false
	}
	_, err := uuid.Parse(u)
	return err == nil
}
