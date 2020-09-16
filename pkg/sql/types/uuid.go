package types

import (
	"encoding/json"

	"github.com/gofrs/uuid"
)

type UUID struct {
	uuid.UUID
}

func (u UUID) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

func (u *UUID) UnmarshalJSON(data []byte) error {
	var id string
	if err := json.Unmarshal(data, &id); err != nil {
		return err
	}
	return u.UnmarshalText([]byte(id))
}
