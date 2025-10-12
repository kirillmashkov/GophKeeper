package model

type CreateSecretRequest struct {
	Name     string `json:"name"`
	Data     string `json:"data"`
	Metadata string `json:"metadata"`
	Kind     int    `json:"type"`
}

type UpdateSecretRequest struct {
	Name     string `json:"name"`
	Data     string `json:"data"`
	Metadata string `json:"metadata"`
	Kind     int    `json:"type"`
}

type SecretDB struct {
	ID       string
	Name     string
	Data     []byte
	Metadata []byte
	Kind     int
}

type GetSecretResponse struct {
	Name     string `json:"name"`
	Data     string `json:"data"`
	Metadata string `json:"metadata"`
	Kind     int    `json:"type"`
}
