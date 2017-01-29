package models

import "encoding/json"

type CheckRequestV1 struct {
	Metadata CheckMetadaV1    `json:"metadata"`
	Param    *json.RawMessage `json:"param"`
}

type CheckMetadaV1 struct {
	Type string `json:"type"`
	Freq int    `json:"freq"`
	Timestamp int64 `json:"timestamp"`
	Id   int    `json:"id"`
	ServiceId int `json:"service_id"`
}
