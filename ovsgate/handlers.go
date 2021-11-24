package ovsgate

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func grpcMetadataHandler(c context.Context, r *http.Request) metadata.MD {

	meta := make(map[string]string)

	val := r.Header.Get("Authorization")

	meta["authorization"] = val
	return metadata.New(meta)
}

func grpcErrorHandler(c context.Context, sm *runtime.ServeMux, m runtime.Marshaler, rw http.ResponseWriter, r *http.Request, e error) {

	fallback := `{"error": "failed to marshal error message"}`

	rw.Header().Set("Content-type", "application/json")

	rw.WriteHeader(runtime.HTTPStatusFromCode(status.Code(e)))

	jErr := json.NewEncoder(rw).Encode(
		struct {
			Err string `json:"error,omitempty"`
		}{
			Err: status.Convert(e).Message(),
		})

	if jErr != nil {
		rw.Write([]byte(fallback))
	}

}
