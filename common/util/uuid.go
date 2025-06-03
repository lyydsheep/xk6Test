package util

import (
	"github.com/google/uuid"
	"strings"
)

func NewUUIDWithoutDash() string {
	Uuid, _ := uuid.NewUUID()
	uuids := strings.Split(Uuid.String(), "-")
	return strings.Join(uuids, "")
}
