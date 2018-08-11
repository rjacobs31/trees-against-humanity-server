package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	jwt.StandardClaims
	Scope     string `json:"scope"`
	populated bool
}

func (c *Claims) Populate(tokenString string) (err error) {
	token, err := jwt.ParseWithClaims(tokenString, c, nil)
	if err != nil {
		return
	}

	c, ok := token.Claims.(*Claims)
	if ok {
		c.populated = true
	} else {
		err = errors.New("could not load claims")
	}
	return
}

func (c *Claims) CheckScope(scope string) (hasScope bool) {
	if !c.populated {
		return false
	}

	for _, v := range strings.Split(c.Scope, " ") {
		if v == scope {
			return true
		}
	}

	return
}

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

func Auth0Middleware(expectedAud, expectedIss, authDomain string) (middleware *jwtmiddleware.JWTMiddleware) {
	return jwtmiddleware.New(jwtmiddleware.Options{
		SigningMethod: jwt.SigningMethodRS256,
		ValidationKeyGetter: func(token *jwt.Token) (result interface{}, err error) {
			result = token

			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(expectedAud, false)
			if !checkAud {
				err = errors.New("invalid audience")
				return
			}

			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(expectedIss, false)
			if !checkIss {
				err = errors.New(("invalid issuer"))
				return
			}

			cert, err := getPemCert(token, authDomain)
			if err != nil {
				panic(err.Error())
			}

			result, err = jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return
		},
	})
}

func getPemCert(token *jwt.Token, authDomain string) (cert string, err error) {
	authUrl, err := url.Parse(authDomain)
	if err != nil {
		return
	}

	subUrl, err := url.Parse("/.well-known/jwks.json")
	if err != nil {
		return
	}

	resp, err := http.Get(authUrl.ResolveReference(subUrl).String())
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
