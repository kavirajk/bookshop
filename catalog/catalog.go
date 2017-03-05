package catalog

import "time"

type Book struct {
	ID              string    `json:"id"`
	ISBN            string    `json:"isbn"`
	Title           string    `json:"title"`
	Tags            []string  `json:"tag"`
	TagString       string    `json:"-"`
	Authors         []Author  `json:"authors"`
	Genres          []Genre   `json:"generes"`
	PublicationYear string    `json:"publication_year"`
	PublicationDate time.Time `json:"publication_date"`
	SampleURL       string    `json:"sample_url"`
	FullURL         string    `json:"-"`
}

type Author struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Publisher struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Genre struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
