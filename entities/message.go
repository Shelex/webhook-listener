package entities

import "encoding/json"

type Hook struct {
	ID         string `json:"id"`
	Created_at int64  `json:"created_at"`
	Channel    string `json:"channel"`
	Payload    string `json:"payload"`
	Failed     bool   `json:"failed"`
}

func (i Hook) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

func (i *Hook) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, i)
}
