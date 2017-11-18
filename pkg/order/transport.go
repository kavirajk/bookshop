package order

import (
	"encoding/json"
	"net/http"

	"context"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/kavirajk/bookshop/transport"
	"github.com/pkg/errors"
)

var (
	ErrBadRouting = errors.New("bad routing")
)

func MakeHTTPHandler(ctx context.Context, s Service, logger log.Logger) http.Handler {
	e := MakeEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}
	placeOrderHandler := httptransport.NewServer(
		e.PlaceOrderEndpoint,
		decodePlaceOrderRequest,
		encodeResponse,
		options...,
	)
	getUserOrdersHandler := httptransport.NewServer(
		e.GetUserOrdersEndpoint,
		decodeGetUserOrdersRequest,
		encodeResponse,
		options...,
	)
	cancelOrdersHandler := httptransport.NewServer(
		e.CancelOrderEndpoint,
		decodeCancelOrderRequest,
		encodeResponse,
		options...,
	)

	r := mux.NewRouter()

	r.Handle("/orders/v1/place", placeOrderHandler).Methods("POST")
	r.Handle("/orders/v1/{user-id}", getUserOrdersHandler).Methods("GET")
	r.Handle("/orders/v1/{user-id}/cancel/{id}", cancelOrdersHandler).Methods("POST")

	return r
}
func decodePlaceOrderRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var r placeOrderRequest
	err := json.NewDecoder(req.Body).Decode(&r)
	return r, err
}

func decodeGetUserOrdersRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	vars := mux.Vars(req)
	userID, ok := vars["user-id"]
	if !ok {
		return nil, errors.Wrap(ErrBadRouting, "user-id")
	}
	return getUserOrdersRequest{
		UserID: userID,
	}, nil
}

func decodeCancelOrderRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	vars := mux.Vars(req)
	userID, ok := vars["user-id"]
	if !ok {
		return nil, errors.Wrap(ErrBadRouting, "user-id")
	}
	ID, ok := vars["id"]
	if !ok {
		return nil, errors.Wrap(ErrBadRouting, "id")
	}

	return cancelOrderRequest{
		UserID:  userID,
		OrderID: ID,
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
	case ErrOrderNotFound:
		return http.StatusNotFound
	case ErrBadRouting:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
