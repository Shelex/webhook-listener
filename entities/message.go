package entities

import "encoding/json"

type Hook struct {
	ID         string `json:"id"`
	Created_at int64  `json:"created_at"`
	Channel    string `json:"channel"`
	Payload    string `json:"payload"`
	Headers    string `json:"headers"`
	StatusOK   bool   `json:"statusOk"`
}

type HooksByChannel struct {
	Data  []Hook `json:"data"`
	Count int64  `json:"count"`
}

type HookExample struct {
	Hello interface{} `json:"hello" swaggertype:"string,object" example:"world"`
}

func (i Hook) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

func (i *Hook) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, i)
}
