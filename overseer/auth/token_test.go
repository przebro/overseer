package auth

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

func TestNewCreator(t *testing.T) {
	_, err := NewTokenCreatorVerifier(TokenSigningMethod("abc"), "testissuer", 1000, []byte("a_secret_key"))
	if err != ErrInvalidSignMethod {
		t.Error("unexpected result:", err, "expected:", ErrInvalidSignMethod)
	}

	_, err = NewTokenCreatorVerifier(MethodHS256, "", 1000, []byte("a_secret_key"))
	if err != ErrInvalidIssuer {
		t.Error("unexpected result:", err, "expected:", ErrInvalidIssuer)
	}

	_, err = NewTokenCreatorVerifier(MethodHS256, "123456789123456789123456789ABCDEFGHIJ", 1000, []byte("a_secret_key"))
	if err != ErrInvalidIssuer {
		t.Error("unexpected result:", err, "expected:", ErrInvalidIssuer)
	}

	_, err = NewTokenCreatorVerifier(MethodHS256, "testissuer", -1, []byte("a_secret_key"))
	if err != ErrInvalidTimeout {
		t.Error("unexpected result:", err, "expected:", ErrInvalidTimeout)
	}

	_, err = NewTokenCreatorVerifier(MethodHS256, "testissuer", 0, []byte("a_secret_key"))
	if err != ErrSecretTooShort {
		t.Error("unexpected result:", err, "expected:", ErrSecretTooShort)
	}

	_, err = NewTokenCreatorVerifier(MethodHS256, "testissuer", 0, []byte("a_secret_key_at_least_16"))
	if err != nil {
		t.Error("unexpected result:", err, "expected:", nil)
	}

}

func TestCreatorCreate(t *testing.T) {

	exp := 123
	issuer := "testissuer"
	sub := "testuser"
	tcv, err := NewTokenCreatorVerifier(MethodHS256, issuer, exp, []byte("a_secret_key_at_least_16"))
	if err != nil {
		t.Error("unexpected result:", err, "expected:", nil)
	}

	_, err = tcv.Create("", map[string]interface{}{"testkey": "testvalue"})
	if err == nil {
		t.Error("unexpected result, invalid username")
	}

	token, err := tcv.Create(sub, map[string]interface{}{"testkey": "testvalue"})
	if err != nil {
		t.Error("unexpected result:", err)
	}
	tb := strings.Split(token, ".")

	br := bytes.NewReader([]byte(tb[1]))

	rdr := base64.NewDecoder(base64.RawStdEncoding, br)
	data, _ := ioutil.ReadAll(rdr)
	result := struct {
		Iat     int
		Iss     string
		Nbf     int
		Sub     string
		testkey string
		Exp     int
	}{}

	json.Unmarshal(data, &result)

	if result.Exp-result.Iat != 123 {
		t.Error("unexpected result, invalid iat")
	}
	if sub != result.Sub {
		t.Error("unexpected result, invalid sub")
	}

	if issuer != result.Iss {
		t.Error("unexpected result, invalid iss")
	}
}

func TestVerifyToken(t *testing.T) {

	exp := 1
	issuer := "testissuer"
	sub := "testuser"
	tcv, err := NewTokenCreatorVerifier(MethodHS256, issuer, exp, []byte("a_secret_key_at_least_16"))
	if err != nil {
		t.Error("unexpected result:", err, "expected:", nil)
	}

	token, err := tcv.Create(sub, map[string]interface{}{"testkey": "testvalue"})

	if err != nil {
		t.Error(err)
	}

	result, err := tcv.Verify(token)
	if err != nil {
		t.Error("unexpected result:", err)
	}
	if result != sub {
		t.Error("unexpected result, sub value:", result, "expected:", sub)
	}

	time.Sleep(2 * time.Second)

	_, err = tcv.Verify(token)

	if err == nil {
		t.Error("unexpected result, token should be invalid after ", exp, "seconds")
	}

}
