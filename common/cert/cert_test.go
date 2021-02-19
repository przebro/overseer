package cert

import (
	"os"
	"testing"
)

func TestGetClientTLSInvalid(t *testing.T) {

	if _, err := GetClientTLS("invalid path", false); err == nil {
		t.Error("unexpected result")
	}

	fpath := "cert.go"

	if _, err := GetClientTLS(fpath, false); err == nil {
		t.Error("unexpected result")
	}

}

func TestGetClientTLS(t *testing.T) {

	fpath := "../../docker/cert/root_ca.crt"
	if _, err := os.Stat(fpath); err != nil {
		t.Skip("cert file does not exists, run certgen.sh first")

	}

	if _, err := GetClientTLS(fpath, false); err != nil {
		t.Error("unexpected result:", err)
	}

}

func TestGetServerTLS(t *testing.T) {

	kpath := "../../docker/cert/overseer_srv.key"
	cpath := "../../docker/cert/overseer_srv.crt"
	if _, err := os.Stat(kpath); err != nil {
		t.Skip("key file does not exists, run certgen.sh first")

	}
	if _, err := os.Stat(cpath); err != nil {
		t.Skip("cert file does not exists, run certgen.sh first")

	}

	if _, err := GetServerTLS(cpath, kpath); err != nil {
		t.Error("unexpected result:", err)
	}
}

func TestGetServerTLSInvalid(t *testing.T) {

	if _, err := GetServerTLS("invalid path", "invlaid path"); err == nil {
		t.Error("unexpected result")
	}

	fpath := "cert.go"

	if _, err := GetServerTLS(fpath, fpath); err == nil {
		t.Error("unexpected result")
	}

}
