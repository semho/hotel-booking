package utils

import (
	"github.com/google/uuid"
)

func StringPtrToUUID(s *string) *uuid.UUID {
	if s == nil {
		return nil
	}
	id, err := uuid.Parse(*s)
	if err != nil {
		return nil
	}
	return &id
}
