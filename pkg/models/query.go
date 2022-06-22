package models

import (
	"encoding/json"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

type QueryPayload struct {
	Action  string `json:"action"`
	GroupId string `json:"groupId"`
	Topic   string `json:"topic"`
	From    int    `json:"from"`
	To      int    `json:"to"`
	Hide    bool   `json:"hide,omitempty"`
}

func ParsePayload(query backend.DataQuery) (*QueryPayload, error) {
	var payload QueryPayload

	// Unmarshal the JSON into QueryPayload.
	err := json.Unmarshal(query.JSON, &payload)
	if err != nil {
		return nil, err
	}

	log.DefaultLogger.Info("ParsePayload", "GroupId", payload.GroupId, "Action", payload.Action)

	return &payload, nil
}
