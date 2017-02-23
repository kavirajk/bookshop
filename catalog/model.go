package catalog

type Album struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Tags      []string `json:"tag"`
	Tracks    []Track  `json:"-"`
	TagString string   `json:"-"`
	SampleURL string   `json:"sample_url"`
	FullURL   string   `json:"-"`
}

type Track struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	SampleURL string `json:"sample_url"`
	FullURL   string `json:"-"`
}
