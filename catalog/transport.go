package book

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"context"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

var (
	ErrEmptyQuery = errors.New("empty query")
	ErrNoNextPage = errors.New("no next page")
	ErrNoPrevPage = errors.New("no prev page")
)

const (
	defaultPageLimit = 20
)

func MakeHTTPHandler(ctx context.Context, s Service, logger log.Logger) http.Handler {
	e := MakeEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}
	searchHandler := httptransport.NewServer(
		ctx,
		e.SearchEndpoint,
		decodeSearchRequest,
		encodeResponse,
		options...,
	)
	getHandler := httptransport.NewServer(
		ctx,
		e.GetEndpoint,
		decodeGetRequest,
		encodeResponse,
		options...,
	)
	r := mux.NewRouter()

	r.Handle("/books/v1/search", searchHandler).Methods("GET")
	r.Handle("/books/v1/{id}", getHandler).Methods("GET")

	return r
}
func decodeSearchRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	q := req.FormValue("q")
	if strings.TrimSpace(q) == "" {
		return nil, ErrEmptyQuery
	}
	return searchRequest{
		Q: q,
	}, nil
}

func decodeGetRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getRequest{
		ID: id,
	}
}

type errorer interface {
	error() error
}

type statuser interface {
	status() int
}

type pager interface {
	page() (total int, previous, next string)
}

// formatResponse is the uniform response format used throughout the books service,
// for every endpoint response.
type formatResponse struct {
	Data interface{}  `json:"data,omitempty"`
	Meta metaResponse `json:"meta"`
}

// metaResponse is part of response json that tells about basic meta information.
type metaResponse struct {
	Status   int    `json:"status"`
	Error    string `json:"error,omitempty"`
	Previous string `json:"previous,omitempty"`
	Next     string `json:"next,omitempty"`
	Total    int    `json:"total,omitempty"`
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, d interface{}) error {
	if e, ok := d.(errorer); ok && e.error() != nil {
		// Now its a business logic error.
		// Extract base domain error.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	status := http.StatusOK
	if s, ok := d.(statbook); ok && s.status() != 0 {
		status = s.status()
	}

	f := formatResponse{
		Data: d,
		Meta: metaResponse{Status: status},
	}

	if page, ok := d.(pager); ok {
		t, p, n := page.page()
		f.Meta.Total = t
		f.Meta.Previous = p
		f.Meta.Next = n
	}

	return json.NewEncoder(w).Encode(f)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// Its important to pass errors.Cause() as we decide status code based on
	// root error which is domain specific
	code := codeFrom(errors.Cause(err))
	w.WriteHeader(code)
	f := formatResponse{Meta: metaResponse{Status: code, Error: err.Error()}}
	json.NewEncoder(w).Encode(f)
}

func codeFrom(err error) int {
	switch err {
	case ErrBookNotFound:
		return http.StatusNotFound
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrInvalidPassword, ErrInvalidResetKey, ErrMissingField, ErrPasswordMismatch:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func nextLimitOffset(total, currentLimit, currentOffset int) (limit, offset int, err error) {
	if currentLimit+currentOffset <= total {
		// there exists next page
		return currentLimit, currentOffset + currentLimit, nil
	}
	return 0, 0, ErrNoNextPage
}

func prevLimitOffset(total, currentLimit, currentOffset int) (limit, offset int, err error) {
	if total > 0 && currentOffset > 0 {
		limit = currentLimit

		// there exists prev page
		if currentOffset-currentLimit <= 0 {
			offset = 0
		} else {
			offset = currentOffset - currentLimit
		}

		return
	}
	return 0, 0, ErrNoNextPage
}

func appendLimitOffset(values url.Values, limit, offset int) url.Values {
	values.Set("limit", strconv.Itoa(limit))
	values.Set("offset", strconv.Itoa(offset))
	return values
}
