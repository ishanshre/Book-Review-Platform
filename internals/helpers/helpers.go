package helpers

import (
	"encoding/json"
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

type Message struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func WriteJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func ApiStatusOk(w http.ResponseWriter, message string) {
	WriteJson(w, http.StatusOK, Message{
		Status:  "success",
		Message: message,
	})
}

func ApiStatusOkData(w http.ResponseWriter, data any) {
	WriteJson(w, http.StatusOK, Message{
		Status: "success",
		Data:   data,
	})
}

func StatusInternalServerError(w http.ResponseWriter, message string) {
	WriteJson(w, http.StatusInternalServerError, Message{
		Status:  "error",
		Message: message,
	})
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

func PageNotFound(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

// IsAuthenticated return true if authenticated else false
func IsAuthenticated(r *http.Request) bool {
	exists := app.Session.Exists(r.Context(), "user_id")
	return exists
}

// IsAdmin returns true if authenticated user is admin else return false
func IsAdmin(r *http.Request) bool {
	access_level := app.Session.GetInt(r.Context(), "access_level")
	return access_level == 1
}
