package api

import (
	"encoding/json"
	"errors"
	"github.com/bewiwi/mta/models"
	"github.com/bewiwi/mta/check"
)

type CheckApi struct {
	Type  string                 `json:"type"`
	Freq  int                    `json:"freq"`
	Id    int                    `json:"id"`
	Param map[string]interface{} `json:"param"`
}

func (c *CheckApi) GetCheck() (models.CheckV1, error) {
	var checkModel models.CheckV1
	var err error
	checkModel.Metadata.Id = c.Id
	checkModel.Metadata.Freq = c.Freq
	checkModel.Metadata.Type = c.Type

	param, err := json.Marshal(c.Param)
	if err != nil {
		return checkModel, err
	}

	var current_check check.CheckRunInterface
	if checkModel.Metadata.Type == "ping" {
		var config check.Ping
		err = json.Unmarshal(param, &config)
		if err != nil {
			return checkModel, err
		}
		param, _ = json.Marshal(config)
		current_check = &config

	} else if checkModel.Metadata.Type == "http" {
		var config check.HttpCheck
		err = json.Unmarshal(param, &config)
		if err != nil {
			return checkModel, err
		}
		param, _ = json.Marshal(config)
		current_check = &config
	} else {
		return checkModel, errors.New("Type not found")
	}

	err = current_check.ValidParam()
	if err != nil {
		return checkModel, err
	}

	finalParam := json.RawMessage(param)
	checkModel.Param = &finalParam
	return checkModel, nil
}
