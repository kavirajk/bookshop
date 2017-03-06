// tranport contains common tranport utils for all the services.
package transport

// formatResponse is the uniform response format used throughout the books service,
// for every endpoint response.
type FormatResponse struct {
	Data interface{}  `json:"data,omitempty"`
	Meta MetaResponse `json:"meta"`
}

// metaResponse is part of response json that tells about basic meta information.
type MetaResponse struct {
	Status   int    `json:"status"`
	Error    string `json:"error,omitempty"`
	Previous string `json:"previous,omitempty"`
	Next     string `json:"next,omitempty"`
	Total    int    `json:"total,omitempty"`
}
