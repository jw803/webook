package utils

import (
	"encoding/base64"
	"encoding/json"
)

type Cursor struct {
	UniqueTimeSeriesFieldValue string
}

func NewCursor(uniqueTimeSeriesFieldValue string) *Cursor {
	return &Cursor{
		UniqueTimeSeriesFieldValue: uniqueTimeSeriesFieldValue,
	}
}

func (c Cursor) Encode() string {
	data, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(data)
}
