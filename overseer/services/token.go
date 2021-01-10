package services

import (
	"encoding/base64"
	"goscheduler/overseer/auth"
	"goscheduler/overseer/config"
)

//NewTokenCreatorVerifier - Creates a new token creator verifier
func NewTokenCreatorVerifier(sec config.SecurityConfiguration) (*auth.TokenCreatorVerifier, error) {

	b, err := base64.StdEncoding.DecodeString(sec.Secret)
	if err != nil {
		return nil, err
	}

	return auth.NewTokenCreatorVerifier(auth.MethodHS256, sec.Issuer, sec.Timeout, b)

}
