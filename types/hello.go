package types

type Message struct {
	Message string `json:"message"`
	Id      int    `json:"id"`
}

type Me struct {
	Name  string `json:"name"`
	IsGay bool   `json:"is_gay"`
}
