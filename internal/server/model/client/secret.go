package client

import "encoding/json"

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