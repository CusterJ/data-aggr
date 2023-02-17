package domain

type Dataset struct {
	Dataset []FooData `json:"dataset"`
}

type FooData struct {
	ID     string `json:"id"`
	Time   int    `json:"time"`
	Signal string `json:"signal,omitempty"`
	Data   string `json:"data,omitempty"`
}
