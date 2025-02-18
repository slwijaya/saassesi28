package models

type Rating struct {
	Rate  float64 `json:"rate"`
	Count int     `json:"count"`
}

type FakeStoreProduct struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	ImageURL    string  `json:"image"`
	Description string  `json:"description"`
	Rating      Rating  `json:"rating"`
}
