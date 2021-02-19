package cert

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"

	"google.golang.org/grpc/credentials"
)

//GetClientTLS - creates transport credentials for client
func GetClientTLS(capath string, untrusted bool) (credentials.TransportCredentials, error) {

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

//GetServerTLS - creates transport credentials for server
func GetServerTLS(certpath, keypath string) (credentials.TransportCredentials, error) {

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
