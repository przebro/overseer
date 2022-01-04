package cert

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/przebro/overseer/common/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

//getClientTLS - creates transport credentials for client
func getClientTLS(capath string, untrusted bool) (credentials.TransportCredentials, error) {

	pool := x509.NewCertPool()
	data, err := ioutil.ReadFile(capath)

	if err != nil {
		return nil, err
	}
	ok := pool.AppendCertsFromPEM(data)
	if !ok {
		return nil, errors.New("unable to read certificate")
	}
	config := &tls.Config{
		RootCAs:            pool,
		InsecureSkipVerify: untrusted,
	}

	return credentials.NewTLS(config), nil
}

//getClientWithOwnCertificateTLS - creates transport credentials for client
func getClientWithOwnCertificateTLS(capath, certkey, certpath string, untrusted bool) (credentials.TransportCredentials, error) {

	pool := x509.NewCertPool()
	data, err := ioutil.ReadFile(capath)

	if err != nil {
		return nil, err
	}
	ok := pool.AppendCertsFromPEM(data)
	if !ok {
		return nil, errors.New("unable to read certificate")
	}

	var clientCerts []tls.Certificate = nil

	cert, err := tls.LoadX509KeyPair(certpath, certkey)
	if err != nil {
		return nil, err
	}

	clientCerts = []tls.Certificate{}
	clientCerts = append(clientCerts, cert)

	config := &tls.Config{
		RootCAs:            pool,
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
func getServerWithClientCArootTLS(certpath, keypath, clientCA string, policy types.CertPolicy) (credentials.TransportCredentials, error) {

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

		data, err := os.ReadFile(clientCA)
		if err != nil {
			return nil, err
		}

		cca = x509.NewCertPool()
		if ok := cca.AppendCertsFromPEM(data); !ok {
			return nil, errors.New("failed to add certificate")
		}
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   authType,
		ClientCAs:    cca,
	}

	return credentials.NewTLS(config), nil
}

//BuildClientCredentials - factory method, creates grpc credentials for client
func BuildClientCredentials(remoteCA, certificate, certificateKey string, clientPolicy types.CertPolicy, level types.ConnectionSecurityLevel) (grpc.DialOption, error) {

	var creds credentials.TransportCredentials
	var err error

	if !filepath.IsAbs(remoteCA) {
		remoteCA, _ = filepath.Abs(remoteCA)
	}

	if !filepath.IsAbs(certificateKey) {
		certificateKey, _ = filepath.Abs(certificateKey)
	}

	if !filepath.IsAbs(certificate) {
		certificate, _ = filepath.Abs(certificate)
	}

	trusted := clientPolicy != types.CertPolicyVerify

	if level == types.ConnectionSecurityLevelServeOnly {
		creds, err = getClientTLS(remoteCA, trusted)
	} else {
		creds, err = getClientWithOwnCertificateTLS(remoteCA, certificateKey, certificate, trusted)
	}

	if err != nil {
		return nil, err
	}
	return grpc.WithTransportCredentials(creds), nil
}

//BuildServerCredentials - factory method, creates grpc credentials for server
func BuildServerCredentials(clientCA, certpath, keypath string, clientPolicy types.CertPolicy, level types.ConnectionSecurityLevel) (grpc.ServerOption, error) {

	var creds credentials.TransportCredentials
	var err error

	if !filepath.IsAbs(certpath) {
		certpath, _ = filepath.Abs(certpath)
	}

	if !filepath.IsAbs(keypath) {
		keypath, _ = filepath.Abs(keypath)
	}

	if !filepath.IsAbs(clientCA) {
		clientCA, _ = filepath.Abs(clientCA)
	}

	if level == types.ConnectionSecurityLevelServeOnly {
		creds, err = getServerTLS(certpath, keypath)
	} else {
		creds, err = getServerWithClientCArootTLS(certpath, keypath, clientCA, clientPolicy)
	}

	if err != nil {
		return nil, err
	}
	return grpc.Creds(creds), nil
}
