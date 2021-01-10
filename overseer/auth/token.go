package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenSigningMethod string

const (
	MethodHS256 TokenSigningMethod = "hs256"
)

var methods = map[TokenSigningMethod]jwt.SigningMethod{MethodHS256: jwt.SigningMethodES256}

type TokenCreatorVerifier struct {
	method  jwt.SigningMethod
	timeout int
	issuer  string
	secret  []byte
}

func NewTokenCreatorVerifier(method TokenSigningMethod, issuer string, timeout int, secret []byte) (*TokenCreatorVerifier, error) {

	var exists bool
	var m jwt.SigningMethod

	if m, exists = methods[method]; !exists {
		return nil, errors.New("invalid signing method")
	}

	if len(issuer) > 32 {
		return nil, errors.New("issuer too long")
	}

	if len(secret) < 16 {
		return nil, errors.New("secret too short")
	}
	if timeout < 0 {
		return nil, errors.New("invalid timeout value")
	}

	return &TokenCreatorVerifier{method: m, timeout: timeout, secret: secret, issuer: issuer}, nil
}

func (tcv *TokenCreatorVerifier) Create(username string, userdata map[string]interface{}) (string, error) {

	now := time.Now()

	if username == "" {
		return "", errors.New("empty username")
	}

	claims := jwt.MapClaims{
		"iat": now.Unix(),
		"nbf": now.Unix(),
		"sub": username,
	}
	if tcv.timeout > 0 {

		exp := now.Add(time.Duration(tcv.timeout) * time.Second)
		claims["exp"] = exp.Unix()
	}

	if tcv.issuer != "" {
		claims["iss"] = tcv.issuer
	}

	for k, v := range userdata {

		claims[k] = v
	}

	tk := jwt.NewWithClaims(tcv.method, claims)
	return tk.SignedString(tcv.secret)

}
func (tcv *TokenCreatorVerifier) Verify(token string) (string, error) {

	var tk *jwt.Token = nil
	var err error
	var ok bool
	var claims jwt.MapClaims

	method := func(t *jwt.Token) (interface{}, error) {
		if t.Method != tcv.method {
			return nil, errors.New("")
		}
		return tcv.method, nil
	}

	if tk, err = jwt.Parse(token, method); err != nil {
		return "", err
	}

	if !tk.Valid {
		return "", errors.New("invalid token")
	}
	if claims, ok = tk.Claims.(jwt.MapClaims); !ok {
		return "", err
	}

	sub := claims["sub"].(string)

	return sub, nil
}
