package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

type response struct {
	Message string `json:"message"`
}

type jwks struct {
	Keys []jsonWebKeys `json:"keys"`
}

type jsonWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

func Auth0Middleware(aud, iss string) (middleware *jwtmiddleware.JWTMiddleware) {
	return jwtmiddleware.New(jwtmiddleware.Options{
		SigningMethod: jwt.SigningMethodRS256,
		ValidationKeyGetter: func(token *jwt.Token) (result interface{}, err error) {
			result = token
			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
			if !checkAud {
				err = errors.New("invalid audience")
				return
			}

			checkIss := token.Claims.(jwt.MapClaims).VerifyAudience(iss, false)
			if !checkIss {
				err = errors.New(("invalid issuer"))
				return
			}

			cert, err := getPemCert(token)
			if err != nil {
				panic(err.Error())
			}

			result, err = jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return
		},
	})
}

func getPemCert(token *jwt.Token) (cert string, err error) {
	resp, err := http.Get("https://YOUR_AUTH0_DOMAIN/.well-known/jwks.json")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	jwks := jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)
	if err != nil {
		return
	}

	for _, v := range jwks.Keys {
		if token.Header["kid"] == v.Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + v.X5c[0] + "\n-----END CERTIFICATE-----"
			return
		}
	}

	err = errors.New("unable to find appropriate key")
	return
}
