package model

// struct for *.exp3.json files
type ExpJson struct {
	Name        string `json:"-"`
	Type        string `json:"Type"`
	FadeInTime  float64 `json:"FadeInTime"`
	FadeOutTime float64 `json:"FadeOutTime"`
	Parameters []struct {
		Id    string  `json:"Id"`
		Value float64 `json:"Value"`
		Blend string  `json:"Blend"`
	} `json:"Parameters"`
}
