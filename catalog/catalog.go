package catalog

import (
	"strings"
	"time"
)

type Book struct {
	ID              string     `json:"id"`
	ISBN            string     `json:"isbn"`
	Title           string     `json:"title"`
	TagString       string     `json:"-"`
	Authors         []Author   `json:"authors" gorm:"many_to_many"`
	Genres          []Genre    `json:"generes" gorm:"many_to_many"`
	Publisher       *Publisher `json:"publisher,omitempty"`
	PublisherID     string     `json:"-"`
	PublicationYear string     `json:"publication_year"`
	PublicationDate time.Time  `json:"publication_date"`
	SampleURL       string     `json:"sample_url"`
	FullURL         string     `json:"-"`
	Price           float64    `json:"price"`
}

func (b *Book) Tags() []string {
	tags := strings.Split(b.TagString, ",")
	for i := range tags {
		tags[i] = strings.TrimSpace(tags[i])
	}
	return tags
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
