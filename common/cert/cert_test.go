package cert

import (
	"os"
	"testing"

	"github.com/przebro/overseer/common/types"
)

func TestGetClientTLSInvalid(t *testing.T) {

	if _, err := getClientTLS("invalid path", false); err == nil {
		t.Error("unexpected result")
	}

	fpath := "cert.go"

	if _, err := getClientTLS(fpath, false); err == nil {
		t.Error("unexpected result")
	}

}

func TestGetClientTLS(t *testing.T) {

	fpath := "../../docker/cert/root_ca.crt"
	if _, err := os.Stat(fpath); err != nil {
		t.Skip("cert file does not exists, run certgen.sh first")

	}

	if _, err := getClientTLS(fpath, false); err != nil {
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

	if _, err := getServerTLS(cpath, kpath); err != nil {
		t.Error("unexpected result:", err)
	}
}

func TestGetServerWithClientCArootTLS(t *testing.T) {

	kpath := "../../docker/cert/overseer_srv.key"
	cpath := "../../docker/cert/overseer_srv.crt"
	capath := "../../docker/cert/root_ca.crt"
	if _, err := os.Stat(kpath); err != nil {
		t.Skip("key file does not exists, run certgen.sh first")

	}
	if _, err := os.Stat(cpath); err != nil {
		t.Skip("cert file does not exists, run certgen.sh first")

	}

	if _, err := getServerWithClientCArootTLS(cpath, kpath, capath, types.CertPolicyNone); err != nil {
		t.Error("unexpected result:", err)
	}

	if _, err := getServerWithClientCArootTLS(cpath, kpath, capath, types.CertPolicyRequired); err != nil {
		t.Error("unexpected result:", err)
	}

	if _, err := getServerWithClientCArootTLS(cpath, kpath, capath, types.CertPolicyVerify); err != nil {
		t.Error("unexpected result:", err)
	}
}

func TestGetServerTLSInvalid(t *testing.T) {

	if _, err := getServerTLS("invalid path", "invlaid path"); err == nil {
		t.Error("unexpected result")
	}

	fpath := "cert.go"

	if _, err := getServerTLS(fpath, fpath); err == nil {
		t.Error("unexpected result")
	}
}

func TestGetServerTLSInvalidCAPath(t *testing.T) {

	kpath := "../../docker/cert/overseer_srv.key"
	cpath := "../../docker/cert/overseer_srv.crt"
	fpath := "cert.go"
	if _, err := os.Stat(kpath); err != nil {
		t.Skip("key file does not exists, run certgen.sh first")

	}
	if _, err := os.Stat(cpath); err != nil {
		t.Skip("cert file does not exists, run certgen.sh first")

	}

	if _, err := getServerWithClientCArootTLS(cpath, kpath, fpath, types.CertPolicyRequired); err == nil {
		t.Error("unexpected result")
	}
}
