package utilities

import (
	"database/sql"
	"encoding/json"
	"time"
)

type NullTime struct {
	sql.NullTime
}

func (nt *NullTime) MarshalJSON() ([]byte, error) {
	if nt != nil {
		if nt.Valid {
			return json.Marshal(nt.Time)
		}
		return json.Marshal(nil)
	} else {
		return json.Marshal(nil)
	}
}

func (nt *NullTime) UnmarshalJSON(data []byte) error {
	var t *time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	if t != nil {
		nt.Valid = true
		nt.Time = *t
	} else {
		nt.Valid = false
	}
	return nil
}

func ValueOfNullNullTime(value *string, layout string, timeZone string) *NullTime {
	// Parse the datetime string into a time.Time object
	// Layout string that matches the format of datetimeStr
	if value != nil {
		if timeZone != "None" {
			valueToUTC, err := DatetimeUTC(*value, layout, timeZone)
			if err != nil {
				return &NullTime{NullTime: sql.NullTime{Valid: false}}
			}

			return &NullTime{NullTime: sql.NullTime{Time: valueToUTC, Valid: true}}
		} else {
			timeLayout, err := time.Parse(layout, *value)
			if err != nil {
				return &NullTime{NullTime: sql.NullTime{Valid: false}}
			}

			return &NullTime{NullTime: sql.NullTime{Time: timeLayout, Valid: true}}
		}
	}

	return &NullTime{NullTime: sql.NullTime{Valid: false}}
}

type NullBool struct {
	sql.NullBool
}

func (nb *NullBool) MarshalJSON() ([]byte, error) {
	if nb != nil {
		if nb.Valid {
			return json.Marshal(nb.Bool)
		}
		return json.Marshal(nil)
	} else {
		return json.Marshal(nil)
	}
}

func (nb *NullBool) UnmarshalJSON(data []byte) error {
	var b *bool
	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}
	if b != nil {
		nb.Valid = true
		nb.Bool = *b
	} else {
		nb.Valid = false
	}
	return nil
}

func ValueOfNullNullBool(value *bool) *NullBool {
	if value != nil {
		return &NullBool{NullBool: sql.NullBool{Bool: *value, Valid: true}}
	}
	return &NullBool{NullBool: sql.NullBool{Valid: false}}
}

type NullFloat64 struct {
	sql.NullFloat64
}

func (nf *NullFloat64) MarshalJSON() ([]byte, error) {
	if nf != nil {
		if nf.Valid {
			return json.Marshal(nf.Float64)
		}
		return json.Marshal(nil)
	} else {
		return json.Marshal(nil)
	}
}

func (nf *NullFloat64) UnmarshalJSON(data []byte) error {
	var f *float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	if f != nil {
		nf.Valid = true
		nf.Float64 = *f
	} else {
		nf.Valid = false
	}
	return nil
}

func ValueOfNullFloat64(value *float64) *NullFloat64 {
	if value != nil {
		return &NullFloat64{NullFloat64: sql.NullFloat64{Float64: *value, Valid: true}}
	}
	return &NullFloat64{NullFloat64: sql.NullFloat64{Valid: false}}
}

type NullInt64 struct {
	sql.NullInt64
}

func (ni *NullInt64) MarshalJSON() ([]byte, error) {
	if ni != nil {
		if ni.Valid {
			return json.Marshal(ni.Int64)
		}
		return json.Marshal(nil)
	} else {
		return json.Marshal(nil)
	}
}

func (ni *NullInt64) UnmarshalJSON(data []byte) error {
	var i *int64
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	if i != nil {
		ni.Valid = true
		ni.Int64 = *i
	} else {
		ni.Valid = false
	}
	return nil
}

func ValueOfNullInt64(value *int64) *NullInt64 {
	if value != nil {
		return &NullInt64{NullInt64: sql.NullInt64{Int64: *value, Valid: true}}
	}
	return &NullInt64{NullInt64: sql.NullInt64{Valid: false}}
}

type NullString struct {
	sql.NullString
}

func (ns *NullString) MarshalJSON() ([]byte, error) {
	if ns != nil {
		if ns.Valid {
			return json.Marshal(ns.String)
		}
		return json.Marshal(nil)
	} else {
		return json.Marshal(nil)
	}
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		ns.Valid = true
		ns.String = *s
	} else {
		ns.Valid = false
	}
	return nil
}

func ValueOfNullString(value *string) *NullString {
	if value != nil {
		return &NullString{NullString: sql.NullString{String: *value, Valid: true}}
	}
	return &NullString{NullString: sql.NullString{Valid: false}}
}

func ParseValueOfNullNullTime(value *string, layout string) *NullTime {
	// Parse the datetime string into a time.Time object
	// Layout string that matches the format of datetimeStr
	if value != nil {
		timeLayout, err := time.Parse(layout, *value)
		if err != nil {
			return &NullTime{NullTime: sql.NullTime{Valid: false}}
		}

		return &NullTime{NullTime: sql.NullTime{Time: timeLayout, Valid: true}}
	}

	return &NullTime{NullTime: sql.NullTime{Valid: false}}
}
