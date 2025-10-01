package collections

import (
	"github.com/google/uuid"
)

func ToUUIDs(ids []string) ([]uuid.UUID, error) {
	uuids := make([]uuid.UUID, len(ids))
	for i, id := range ids {
		uuid, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}
		uuids[i] = uuid
	}
	return uuids, nil
}

func ToStrings(uuids []uuid.UUID) []string {
	strings := make([]string, len(uuids))
	for i, uuid := range uuids {
		strings[i] = uuid.String()
	}
	return strings
}
