package handler

// Message representa a resposta simples do endpoint hello.
type Message struct {
	Message string `json:"message"`
	Id      int    `json:"id"`
}

// Me representa a resposta com o nome extraído da rota.
type Me struct {
	Name string `json:"name"`
}
