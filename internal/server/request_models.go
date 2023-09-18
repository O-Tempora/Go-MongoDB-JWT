package server

type TokenPair struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}
