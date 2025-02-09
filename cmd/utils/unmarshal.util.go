package utils

import (
	"database/sql"
	"encoding/json"
	"time"
)

type JsonNullInt64 struct {
	sql.NullInt64
}

func (v JsonNullInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	}
	return json.Marshal(nil)
}

func (v *JsonNullInt64) UnmarshalJSON(data []byte) error {
	var x *int64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Int64 = *x
	} else {
		v.Valid = false
	}
	return nil
}

// Define a custom type
type JsonNullInt16 struct {
	sql.NullInt16
}

// Implement the UnmarshalJSON method
func (v *JsonNullInt16) UnmarshalJSON(data []byte) error {
	// Attempt to unmarshal into a float64 first
	var x *float64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Int16 = int16(*x)
	} else {
		v.Valid = false
	}
	return nil
}

// Define a custom type
type JsonNullTime struct {
	sql.NullTime
}

// Implement the UnmarshalJSON method
func (v *JsonNullTime) UnmarshalJSON(data []byte) error {
	// Attempt to unmarshal into a string
	var str *string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	if str != nil {
		// Parse the time using the desired layout
		t, err := time.Parse(time.DateTime, *str)
		if err != nil {
			return err
		}
		v.Valid = true
		v.Time = t
	} else {
		v.Valid = false
	}

	return nil
}
