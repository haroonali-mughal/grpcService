package vault

import(
	"net/http"
	httptransport "github.com/go-kit/kit/transport/http"
        //endpoints "github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
        //"github.com/go-kit/kit/endpoint"
)

func NewHTTPServer(ctx context.Context , endpoints Endpoints) http.Handler {
	m := http.NewServeMux()
	m.Handle("/hash",httptransport.NewServer(
		endpoints.HashEndpoint,
		decodeHashRequest,
		encodeResponse,
	))
	m.Handle("/validate",httptransport.NewServer(
		endpoints.ValidateEndpoint,
		decodeValidateRequest,
		encodeResponse,
	))
	return m
}
