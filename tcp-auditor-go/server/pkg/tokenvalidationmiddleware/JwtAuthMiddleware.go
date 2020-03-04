package tokenvalidationmiddleware

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"

	"bitbucket.tylertech.com/spy/scm/tcp-auditor/server/config"
	"bitbucket.tylertech.com/spy/scm/tcp-auditor/server/logging"
	"github.com/gin-gonic/gin"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

// Jwks collection of KSONWebKeys
type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

// JSONWebKeys model for the signing key
type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

// OpenIDConfig is used to parse the response from the well known endpoint, and extract the jwks_uri
type OpenIDConfig struct {
	JwksURI string `json:"jwks_uri"` // Uppercased first letter
}

// CustomClaims defines the claims to validate against during token validation
type CustomClaims struct {
	Scope string `json:"scope"`
	jwt.StandardClaims
}

// JwtAuthenticationMiddleware to add into the request pipeline
func JwtAuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		err := Validator.CheckJWT(c.Writer, c.Request)

		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}
		c.Next()

	}
}

// Validator is middleware to validate auth tokens on each request
var Validator = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		// // Verify 'aud' claim
		// aud := ""
		// checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
		// if !checkAud {
		// 	return token, errors.New("Invalid audience")
		// }
		// Verify 'iss' claim
		configuration := config.GetConfig()

		iss := configuration.OIDC.Authority
		checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, true)
		if !checkIss {
			return token, errors.New("Invalid issuer")
		}

		cert, err := getPemCert(token)
		if err != nil {
			panic(err.Error())
		}

		return cert, nil
	},
	SigningMethod: jwt.SigningMethodRS256,
})

func getPemCert(token *jwt.Token) (*rsa.PublicKey, error) {
	var openIDConfig OpenIDConfig
	nStr := ""
	eStr := ""
	var cert *rsa.PublicKey

	configuration := config.GetConfig()
	wellKnownEndpoint := configuration.OIDC.WellKnownURL
	resp, err := http.Get(wellKnownEndpoint)

	// read the payload
	body, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		return cert, err
	}

	err = json.Unmarshal(body, &openIDConfig)
	if err != nil {
		return cert, err
	}

	jwksResp, err := http.Get(openIDConfig.JwksURI)

	var jwks = Jwks{}
	err = json.NewDecoder(jwksResp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			nStr = jwks.Keys[k].N
			eStr = jwks.Keys[k].E
		}
	}

	if nStr == "" || eStr == "" {
		err := errors.New("Unable to find appropriate key")
		return cert, err
	}

	cert = genXC5(nStr, eStr)

	return cert, nil
}

// CheckScope adds to the middleware by allowing you to scope specific endpoints
var CheckScope = func(scope string, tokenString string) bool {
	token, _ := jwt.ParseWithClaims(tokenString, &CustomClaims{}, nil)

	claims, _ := token.Claims.(*CustomClaims)

	hasScope := false
	result := strings.Split(claims.Scope, " ")
	for i := range result {
		if result[i] == scope {
			hasScope = true
		}
	}

	return hasScope
}

func genXC5(nStr string, eStr string) *rsa.PublicKey {
	// decode the base64 bytes for n
	nb, err := base64.RawURLEncoding.DecodeString(nStr)
	if err != nil {
		logging.Logger.Error("Unable to decode modulus from jwks.", err)
	}

	e := 65537
	// The default exponent is usually 65537, so just compare the
	// base64 for [1,0,1] or [0,1,0,1]
	if eStr != "AQAB" && eStr != "AAEAAQ" {
		// still need to decode the big-endian int
		logging.Logger.Error("need to deocde e:", eStr)
	}

	var pubKey = &rsa.PublicKey{
		N: new(big.Int).SetBytes(nb),
		E: e,
	}

	return pubKey
}
