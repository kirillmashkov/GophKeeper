package client

import (
	"encoding/json"
	"errors"
)

type ISecret interface {
	String() string
	Kind() int
}

type secret struct {
	Kind int `json:"type"`
	Data json.RawMessage `json:"data"`
}

func Encode(s ISecret) ([]byte, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return json.Marshal(secret{
		Kind: s.Kind(),
		Data: data,
	})
}

func Decode(data []byte) (ISecret, error) {
	var s secret
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}

	if s.Kind == 1 {
		var bin Bin
		if err := json.Unmarshal(s.Data, &bin); err != nil {
			return nil, err
		}
		return bin, nil
	} else if s.Kind == 2 {
		var cred Credential
		if err := json.Unmarshal(s.Data, &cred); err != nil {
			return nil, err
		}
		return cred, nil

	} else if s.Kind == 3 {
		var card Card
		if err := json.Unmarshal(s.Data, &card); err != nil {
			return nil, err
		}
		return card, nil
	} else {
		return nil, errors.New("unknown secret")
	}
}