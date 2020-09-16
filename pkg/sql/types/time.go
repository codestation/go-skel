package types

import (
	"database/sql"
	"encoding/json"
)

type NullTime struct {
	sql.NullTime
}

func (s NullTime) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(s.Time)
}

func (s *NullTime) UnmarshalJSON(data []byte) error {
	s.Valid = string(data) != "null"
	err := json.Unmarshal(data, &s.Time)
	return err
}
