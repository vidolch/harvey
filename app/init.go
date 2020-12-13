package app

import (
	"github.com/revel/revel"
	"harvey/app/filters"
	"net/http"
)

var (
	// AppVersion revel app version (ldflags)
	AppVersion string

	// BuildTime revel app build-time (ldflags)
	BuildTime string
)

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		HeaderFilter,                  // Add some security based headers
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.BeforeAfterFilter,       // Call the before and after filter functions
		filters.AuthFilter,	   		   // Authentication filter
		revel.ActionInvoker,           // Invoke the action.
	}

	// Register startup functions with OnAppStart
	// revel.DevMode and revel.RunMode only work inside of OnAppStart. See Example Startup Script
	// ( order dependent )
	// revel.OnAppStart(ExampleStartupScript)
	// revel.OnAppStart(InitDB)
	// revel.OnAppStart(FillCache)
}

// HeaderFilter adds common security headers
// There is a full implementation of a CSRF filter in
// https://github.com/revel/modules/tree/master/csrf
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	if c.Request.Method == "OPTIONS" {
		c.Response.Out.Header().Add("Access-Control-Allow-Origin", "*")
		c.Response.Out.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization") //自定义 Header
		c.Response.Out.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Response.Out.Header().Add("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Response.Out.Header().Add("Access-Control-Allow-Credentials", "true")
		c.Response.SetStatus(http.StatusNoContent)
		// 截取复杂请求下post变成options请求后台处理方法(针对跨域请求检测)
	} else {
		c.Response.Out.Header().Add("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		c.Response.Out.Header().Add("Access-Control-Allow-Origin", "*")
		c.Response.Out.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Response.Out.Header().Add("Content-Type", "application/json; charset=UTF-8")
		c.Response.Out.Header().Add("X-Frame-Options", "SAMORIGIN")
		c.Response.Out.Header().Add("Vary", "Origin, Access-Control-Request-Method, Access-Control-Request-Headers")

		fc[0](c, fc[1:]) // Execute the next filter stage.
	}
}

//func ExampleStartupScript() {
//	// revel.DevMod and revel.RunMode work here
//	// Use this script to check for dev mode and set dev/prod startup scripts here!
//	if revel.DevMode == true {
//		// Dev mode
//	}
//}
