package entity

type Tier struct {
	Key       string `json:"key"`
	Name      string `json:"name"`
	Threshold int    `json:"threshold"`
}
