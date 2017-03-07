package catalog

import (
	"encoding/json"
	"net/http"
	"strings"

	"context"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/kavirajk/bookshop/transport"
	"github.com/pkg/errors"
)

var (
	ErrEmptyQuery = errors.New("empty query")
	ErrBadRouting = errors.New("bad routing")
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

	r.Handle("/catalog/v1/search", searchHandler).Methods("GET")
	r.Handle("/catalog/v1/{id}", getHandler).Methods("GET")

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
	vars := mux.Vars(req)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getRequest{
		ID: id,
	}, nil
}

// errorer interface should be implemented by all the doman specific errors.
// easy to set different status code in case of different errors.
type errorer interface {
	error() error
}

// statuser allows any response to get customer status code
// e.g: 201 for successfull resource creation.
type statuser interface {
	status() int
}

// pager used to paginate any transport response.
type pager interface {
	page() (total int, previous, next string)
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
	if s, ok := d.(statuser); ok && s.status() != 0 {
		status = s.status()
	}

	f := transport.FormatResponse{
		Data: d,
		Meta: transport.MetaResponse{Status: status},
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
	f := transport.FormatResponse{Meta: transport.MetaResponse{Status: code, Error: err.Error()}}
	json.NewEncoder(w).Encode(f)
}

func codeFrom(err error) int {
	switch err {
	case ErrBookNotFound:
		return http.StatusNotFound
	case ErrEmptyQuery, ErrBadRouting:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
