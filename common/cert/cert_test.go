package cert

import (
	"os"
	"testing"

	"github.com/przebro/overseer/common/types"
)

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
	RegisterCA("../../docker/cert/root_ca.crt")

	if _, err := os.Stat(kpath); err != nil {
		t.Skip("key file does not exists, run certgen.sh first")

	}
	if _, err := os.Stat(cpath); err != nil {
		t.Skip("cert file does not exists, run certgen.sh first")

	}

	if _, err := getServerWithClientCArootTLS(cpath, kpath, types.CertPolicyNone); err != nil {
		t.Error("unexpected result:", err)
	}

	if _, err := getServerWithClientCArootTLS(cpath, kpath, types.CertPolicyRequired); err != nil {
		t.Error("unexpected result:", err)
	}

	if _, err := getServerWithClientCArootTLS(cpath, kpath, types.CertPolicyVerify); err != nil {
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

func TestBuildServerCredentials_Errors(t *testing.T) {

	kpath := "cert.go"
	cpath := "../../docker/cert/overseer_srv.crt"

	RegisterCA("cert.go")

	if _, err := os.Stat(kpath); err != nil {
		t.Skip("key file does not exists, run certgen.sh first")

	}
	if _, err := os.Stat(cpath); err != nil {
		t.Skip("cert file does not exists, run certgen.sh first")

	}

	_, err := BuildServerCredentials(cpath, kpath, types.CertPolicyRequired, types.ConnectionSecurityLevelClientAndServer)
	if err == nil {
		t.Error("unexpected result,err is nil")
	}

	_, err = BuildServerCredentials(cpath, kpath, types.CertPolicyRequired, types.ConnectionSecurityLevelServeOnly)
	if err == nil {
		t.Error("unexpected result,err is nil")
	}
}

func TestBuildClientCredentials_Errors(t *testing.T) {

	kpath := "cert.go"
	cpath := "../../docker/cert/overseer_srv.crt"

	if _, err := os.Stat(kpath); err != nil {
		t.Skip("key file does not exists, run certgen.sh first")

	}
	if _, err := os.Stat(cpath); err != nil {
		t.Skip("cert file does not exists, run certgen.sh first")

	}

	_, err := BuildClientCredentials(cpath, kpath, types.CertPolicyRequired, types.ConnectionSecurityLevelClientAndServer)
	if err == nil {
		t.Error("unexpected result,err is nil")
	}
}
