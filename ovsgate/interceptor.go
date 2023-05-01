package ovsgate

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/rs/zerolog"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type httpInterceptor struct {
	mux *runtime.ServeMux
	log *zerolog.Logger
}

func (i *httpInterceptor) ServeHTTP(resp http.ResponseWriter, req *http.Request) {

	payload, _ := ioutil.ReadAll(req.Body)
	req.Body = ioutil.NopCloser(bytes.NewBuffer(payload))

	i.log.Info().Str("path", req.URL.Path).Str("body", string(payload)).Msg("request received")

	headers := []string{"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"}
	resp.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}
	resp.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	resp.Header().Set("Access-Control-Allow-Origin", "*")

	if req.Method == "OPTIONS" {
		return
	}

	log := zerolog.Ctx(req.Context()).With().Str("service", "exec").Logger()

	nctx := log.WithContext(req.Context())
	req = req.WithContext(nctx)

	i.mux.ServeHTTP(resp, req)
}

func newHttpInterceptor(mux *runtime.ServeMux, log *zerolog.Logger) http.Handler {

	return &httpInterceptor{mux: mux, log: log}
}
