package mongox

import (
	"encoding/base64"
	"encoding/json"
	"errors"
)

type Cursor struct {
	UniqueTimeSeriesFieldValue string
}

func NewCursor(uniqueTimeSeriesFieldValue string) *Cursor {
	return &Cursor{
		UniqueTimeSeriesFieldValue: uniqueTimeSeriesFieldValue,
	}
}

func (c Cursor) EncodeCursor() (string, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data), err
}

type CursorString string

func (c CursorString) Decode() (Cursor, error) {
	var cursor Cursor
	if c == "" {
		return cursor, errors.New("invalid cursor string")
	}

	data, err := base64.RawURLEncoding.DecodeString(string(c))
	if err != nil {
		return cursor, err
	}

	err = json.Unmarshal(data, &cursor)
	if err != nil {
		return cursor, err
	}
	return cursor, err
}
