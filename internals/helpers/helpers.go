package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/ishanshre/Book-Review-Platform/internals/config"
)

var app *config.AppConfig

// NewHelpers sets up access to gloabal app config
func NewHelpers(a *config.AppConfig) {
	app = a
}

// ClientError handles the client errors
func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status of", status)
	http.Error(w, http.StatusText(status), status)
}

// ServerError handles the server error
func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
