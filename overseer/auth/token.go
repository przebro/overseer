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

var (
	methods              = map[TokenSigningMethod]jwt.SigningMethod{MethodHS256: jwt.SigningMethodHS256}
	ErrInvalidSignMethod = errors.New("invalid signing method")
	ErrInvalidIssuer     = errors.New("invalid issuer len!=0 and len <= 32")
	ErrSecretTooShort    = errors.New("secret too short")
	ErrInvalidTimeout    = errors.New("invalid timeout value")
)

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
		return nil, ErrInvalidSignMethod
	}

	if len(issuer) > 32 || len(issuer) == 0 {
		return nil, ErrInvalidIssuer
	}

	if timeout < 0 {
		return nil, ErrInvalidTimeout
	}

	if len(secret) < 16 {
		return nil, ErrSecretTooShort
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
		return tcv.secret, nil
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
