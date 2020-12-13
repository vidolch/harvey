package filters

import (
	"fmt"
	"github.com/casbin/casbin"
	"github.com/dgrijalva/jwt-go"
	"github.com/revel/revel"
	"log"
	"strings"
)

func AuthFilter(c *revel.Controller, fc []revel.Filter) {
	//configURL := "https://auth.chalamov.dev/auth/realms/test"
	//ctx := context.Background()
	//provider, err := oidc.NewProvider(ctx, configURL)
	//if err != nil {
	//	panic(err)
	//}
	//
	//clientID := "test-client"
	////clientSecret := "540fb896-a614-4932-b802-327b8696e663"
	////
	////redirectURL := "http://harvey:8080/login/callback"
	////// Configure an OpenID Connect aware OAuth2 client.
	////oauth2Config := oauth2.Config{
	////	ClientID:     clientID,
	////	ClientSecret: clientSecret,
	////	RedirectURL:  redirectURL,
	////	// Discovery returns the OAuth2 endpoints.
	////	Endpoint: provider.Endpoint(),
	////	// "openid" is a required scope for OpenID Connect flows.
	////	Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	////}
	//
	//oidcConfig := &oidc.Config{
	//	ClientID: clientID,
	//}
	//
	//verifier := provider.Verifier(oidcConfig)
	//
	//res, err2 := verifier.Verify(ctx, c.Request.Header.Get("Authentication"))
	//
	//if err2 != nil {
	//	log.Fatal(err2)
	//	panic(err)
	//}
	//
	//if res.VerifyAccessToken(c.Request.Header.Get("Authentication")) == nil {
	//	c.Result = c.Forbidden("Access denied by the Auth plugin.")
	//	return
	//} else {
	//	fc[0](c, fc[1:])
	//}

	e, err := casbin.NewEnforcer("app/filters/auth_model.conf", "app/filters/auth_policy.csv")

	if err != nil {
		log.Fatal(err)
		return
	}

	if !CheckPermission(e, c.Request) {
		c.Result = c.Forbidden("Access denied by the Auth plugin.")
		return
	} else {
		fc[0](c, fc[1:])
	}
}

// GetJwtToken gets the user name from the request.
// Currently, only HTTP basic authentication is supported
func GetJwtToken(r *revel.Request) (*jwt.Token, error) {
	rawToken := r.Header.Get("Authorization")
	if rawToken == "" {
		return nil, nil
	}

	token := strings.Split(rawToken, "Bearer ")[1]

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		v := []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAm9BU67tSGahlb5xNpTuV
+3aUuI7EueJw1ubSZWp+nacO53wMeMepBkCOC9dHjbWOaghIVKWWz00tqZEgym5j
QaBLC2EHujUCjRCupBLH6s2Ejfmw70Agx816v8xoVfoqHSD3f0/hAFGslyLgdB6m
tLT9djRhqJ8u/EnFjEpKnh7a47sVJIpyxbViI2TOKwIr53TVLZzcSbiGno+VI/Ig
Z1VACfvBPOCHJYz5af2Ex/vZn/veJhLUtowXCA9UMbHFZS9pU4HOiF9vV0ZaBFhh
MmmZ7fUQGJqWpFYU9RWpmS1kk6spegmVOPx7S9RY4pDjuyjolkG7Lre9yMgdSZdy
qwIDAQAB
-----END PUBLIC KEY-----
`)
		k, _ := jwt.ParseRSAPublicKeyFromPEM(v)

		return k, nil
	})

	return parsedToken, err
}

// CheckPermission checks the user/method/path combination from the request.
// Returns true (permission granted) or false (permission forbidden)
func CheckPermission(e *casbin.Enforcer, r *revel.Request) bool {
	user, err := GetJwtToken(r)
	if err != nil {
		log.Fatal(err)
		return false
	}
	role := "anonymous"
	if user != nil {
		role = "user"
		//if user == "" {
		//	user = "anonymous"
		//}
		if claims, ok := user.Claims.(jwt.MapClaims); ok && user.Valid {
			log.Print(claims["email"])
			fmt.Println(claims["foo"], claims["nbf"])
		} else {
			log.Print(claims["email"])
			fmt.Println(err)
		}
	}

	method := r.Method
	path := r.URL.Path
	res, _ := e.Enforce(role, path, method)
	return res
}
