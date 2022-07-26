package cert

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/przebro/overseer/common/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var caPool = x509.NewCertPool()

//getClientTLS - creates transport credentials for client
func getClientTLS(untrusted bool) (credentials.TransportCredentials, error) {

	config := &tls.Config{
		RootCAs:            caPool,
		InsecureSkipVerify: untrusted,
	}

	return credentials.NewTLS(config), nil
}

//getClientWithOwnCertificateTLS - creates transport credentials for client
func getClientWithOwnCertificateTLS(certkey, certpath string, untrusted bool) (credentials.TransportCredentials, error) {

	var clientCerts []tls.Certificate = nil

	cert, err := tls.LoadX509KeyPair(certpath, certkey)
	if err != nil {
		return nil, err
	}

	clientCerts = []tls.Certificate{}
	clientCerts = append(clientCerts, cert)

	config := &tls.Config{
		RootCAs:            caPool,
		InsecureSkipVerify: untrusted,
		Certificates:       clientCerts,
	}

	return credentials.NewTLS(config), nil
}

//getServerTLS - creates transport credentials for server
func getServerTLS(certpath, keypath string) (credentials.TransportCredentials, error) {

	cert, err := tls.LoadX509KeyPair(certpath, keypath)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}

//getServerWithClientCArootTLS - creates transport credentials for server
func getServerWithClientCArootTLS(certpath, keypath string, policy types.CertPolicy) (credentials.TransportCredentials, error) {

	cert, err := tls.LoadX509KeyPair(certpath, keypath)
	if err != nil {
		return nil, err
	}

	var authType tls.ClientAuthType
	switch policy {
	case types.CertPolicyNone:
		{
			authType = tls.NoClientCert
		}
	case types.CertPolicyRequired:
		{
			authType = tls.RequireAnyClientCert
		}
	case types.CertPolicyVerify:
		{
			authType = tls.RequireAndVerifyClientCert
		}
	}

	var cca *x509.CertPool

	if authType != tls.NoClientCert {
		cca = caPool
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   authType,
		ClientCAs:    cca,
	}

	return credentials.NewTLS(config), nil
}

//BuildClientCredentials - factory method, creates grpc credentials for client
func BuildClientCredentials(certificate, certificateKey string, clientPolicy types.CertPolicy, level types.ConnectionSecurityLevel) (grpc.DialOption, error) {

	var creds credentials.TransportCredentials
	var err error

	if !filepath.IsAbs(certificateKey) {
		certificateKey, _ = filepath.Abs(certificateKey)
	}

	if !filepath.IsAbs(certificate) {
		certificate, _ = filepath.Abs(certificate)
	}

	trusted := clientPolicy != types.CertPolicyVerify

	if level == types.ConnectionSecurityLevelServeOnly {
		creds, err = getClientTLS(trusted)
	} else {
		creds, err = getClientWithOwnCertificateTLS(certificateKey, certificate, trusted)
	}

	if err != nil {
		return nil, err
	}
	return grpc.WithTransportCredentials(creds), nil
}

//BuildServerCredentials - factory method, creates grpc credentials for server
func BuildServerCredentials(certpath, keypath string, clientPolicy types.CertPolicy, level types.ConnectionSecurityLevel) (grpc.ServerOption, error) {

	var creds credentials.TransportCredentials
	var err error

	if !filepath.IsAbs(certpath) {
		certpath, _ = filepath.Abs(certpath)
	}

	if !filepath.IsAbs(keypath) {
		keypath, _ = filepath.Abs(keypath)
	}

	if level == types.ConnectionSecurityLevelServeOnly {
		creds, err = getServerTLS(certpath, keypath)
	} else {
		creds, err = getServerWithClientCArootTLS(certpath, keypath, clientPolicy)
	}

	if err != nil {
		return nil, err
	}
	return grpc.Creds(creds), nil
}

func RegisterCA(path string) error {

	data, err := ioutil.ReadFile(path)

	if err != nil {
		return err
	}
	ok := caPool.AppendCertsFromPEM(data)
	if !ok {
		return errors.New("unable to read root certificate")
	}

	return nil
}
