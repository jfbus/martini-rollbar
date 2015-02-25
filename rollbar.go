// Package martini-rollbar is a middleware for Martini that reports panics to rollbar.com.
//
//  package main
//
//  import (
//    "github.com/go-martini/martini"
//    "github.com/jfbus/martini-rollbar"
//  )
//
//
//  func main() {
//    m := martini.Classic()
//    m.Use(rollbar.Report(rollbar.Config{Token: ROLLBAR_TOKEN}))
//
//    m.Get("/panic", func() {
//      panic("an error occured")
//    })
//
//    m.Run()
//  }
package rollbar

import (
	"fmt"
	"log"
	"net/http"

	"github.com/codegangsta/inject"
	"github.com/go-martini/martini"
	rb "github.com/stvp/rollbar"
)

type Config struct {
	Token string
}

// Recovery returns a middleware that recovers from any panics, sends the error to rollbar and writes a HTTP 500 response.
// It only triggers in production environments.
func Recovery(cfg Config) martini.Handler {
	if martini.Env != martini.Prod {
		return func() {}
	}
	rb.Token = cfg.Token
	rb.Environment = "production"

	return func(c martini.Context, r *http.Request, log *log.Logger) {
		defer func() {
			if err := recover(); err != nil {

				stack := rb.BuildStack(3)
				if e, ok := err.(error); ok {
					rb.RequestErrorWithStack(rb.CRIT, r, e, stack)
				} else {
					rb.RequestErrorWithStack(rb.CRIT, r, fmt.Errorf("%s", err), stack)
				}
				var str string
				for _, f := range stack {
					str += fmt.Sprintf("File \"%s\" line %d in %s\n", f.Filename, f.Line, f.Method)
				}

				log.Printf("PANIC: %s\n%s", err, str)

				val := c.Get(inject.InterfaceOf((*http.ResponseWriter)(nil)))
				res := val.Interface().(http.ResponseWriter)

				body := []byte("500 Internal Server Error")

				res.WriteHeader(http.StatusInternalServerError)
				res.Write(body)
			}
		}()

		c.Next()
	}
}
