package filters

import (
	"github.com/casbin/casbin"
	"github.com/revel/revel"
	"log"
	"net/http"
)

func AuthFilter(c *revel.Controller, fc []revel.Filter) {
	e, err := casbin.NewEnforcer("app/filters/auth_model.conf", "app/filters/auth_policy.csv")

	log.Print("comming in")

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

// GetUserName gets the user name from the request.
// Currently, only HTTP basic authentication is supported
func GetUserName(r *revel.Request) string {
	req := r.In.GetRaw().(*http.Request)
	username, _, _ := req.BasicAuth()
	return username
}

// CheckPermission checks the user/method/path combination from the request.
// Returns true (permission granted) or false (permission forbidden)
func CheckPermission(e *casbin.Enforcer, r *revel.Request) bool {
	user := GetUserName(r)
	method := r.Method
	path := r.URL.Path
	res, _ := e.Enforce(user, path, method)
	return res
}
