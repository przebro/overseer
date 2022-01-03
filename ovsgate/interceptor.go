package ovsgate

import (
	"net/http"
	"strings"

	"github.com/przebro/overseer/common/logger"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type httpInterceptor struct {
	mux *runtime.ServeMux
	log logger.AppLogger
}

func (i *httpInterceptor) ServeHTTP(resp http.ResponseWriter, req *http.Request) {

	headers := []string{"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"}
	resp.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}
	resp.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	resp.Header().Set("Access-Control-Allow-Origin", "*")

	if req.Method == "OPTIONS" {
		return
	}

	i.log.Info(req)
	i.mux.ServeHTTP(resp, req)
}

func newHttpInterceptor(mux *runtime.ServeMux, log logger.AppLogger) http.Handler {

	return &httpInterceptor{mux: mux, log: log}
}
