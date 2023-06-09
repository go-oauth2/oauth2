package handlers

import (
	"fmt"
	"github.com/go-oauth2/oauth2/v4/server"
	"net/http"
	"os"
	"sync"
)

// List of the allowed services
var apiWhiteList = map[string]string{
	"888888": "88888888",
}

// DBRepo is the db repo
type Authentication struct {
	srv           *server.Server
	extStore      map[string]interface{} // like a redis store
	databaseUsers map[string]interface{} // use a db
	// one mutex for 2 stores, ok for the demo
	sync.RWMutex
}

// NewPostgresqlHandlers creates db repo for postgres
func NewAuthentication(srv *server.Server) Authentication {

	return Authentication{
		srv:           srv,
		extStore:      make(map[string]interface{}),
		databaseUsers: make(map[string]interface{}),
	}
}

func (a Authentication) Authorize(w http.ResponseWriter, r *http.Request) {
	if dumpvar {
		dumpRequest(os.Stdout, "authorize", r)
	}

	err := a.srv.HandleAuthorizeRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

}

func (a Authentication) UserAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "userAuthorizeHandler", r) // Ignore the error
	}

	clientID := r.Form.Get("client_id")

	switch clientID {
	case "222222":

		a.RLock()
		uid, ok := a.extStore[fmt.Sprintf("LoggedInUserID-%v", r.Form.Get("email"))]
		a.RUnlock()
		if !ok {
			if r.Form == nil {
				r.ParseForm()
			}

			w.WriteHeader(http.StatusOK)
			return
		}

		userID = uid.(string)

		a.Lock()
		delete(a.extStore, fmt.Sprintf("LoggedInUserID-%v", r.Form.Get("email")))
		a.Unlock()
		return
	case "888888":

		a.RLock()
		uid, ok := a.extStore[fmt.Sprintf("LoggedInUserID-%v", r.Form.Get("client_id"))]
		a.RUnlock()
		if !ok {
			if r.Form == nil {
				r.ParseForm()
			}

			w.WriteHeader(http.StatusOK)
			return
		}

		userID = uid.(string)

		a.Lock()
		delete(a.extStore, fmt.Sprintf("LoggedInUserID-%v", r.Form.Get("client_id")))
		a.Unlock()
		return
	default:
		userID = ""
		return
	}
}

func (a Authentication) Token(w http.ResponseWriter, r *http.Request) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "token", r) // Ignore the error
	}

	err := a.srv.HandleTokenRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Endpoint specific for the APIs
func (a Authentication) ApiAuthHandler(w http.ResponseWriter, r *http.Request) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "apiAuthHandler", r) // Ignore the error
	}

	if r.Method == "POST" {

		if r.Form == nil {
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// make sure the client's api is allow
		_, ok := apiWhiteList[r.Form.Get("client_id")]
		if ok {
			// save user in a temporary store for the user to be reconized later on
			a.Lock()
			a.extStore[fmt.Sprintf("LoggedInUserID-%v", r.Form.Get("client_id"))] = r.Form.Get("client_id")
			a.Unlock()

			a.Authorize(w, r)
			return
		} else {

			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

	}

	http.Error(w, "Bad Request", http.StatusBadRequest)
}

func (a Authentication) SignupHandler(w http.ResponseWriter, r *http.Request) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "signup", r) // Ignore the error
	}

	if r.Method == "POST" {

		if r.Form == nil {
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// some logic
		if len(r.Form.Get("email")) < 1 && len(r.Form.Get("password")) < 1 {

			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		} else {
			// Save the user in db(simplistic, for example)
			a.Lock()
			a.databaseUsers[fmt.Sprintf(r.Form.Get("email"))] = r.Form.Get("password")
			a.Unlock()
		}

		// save user in a temporary store for the user to be reconized later on
		a.Lock()
		a.extStore[fmt.Sprintf("LoggedInUserID-%v", r.Form.Get("email"))] = r.Form.Get("email")
		a.Unlock()

		a.Authorize(w, r)
		return
	}

	http.Error(w, "Bad Request", http.StatusBadRequest)
}

func (a Authentication) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "login", r) // Ignore the error
	}

	if r.Method == "POST" {

		if r.Form == nil {
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if len(r.Form.Get("email")) < 1 && len(r.Form.Get("password")) < 1 {

			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		} else {
			a.RLock()
			password, ok := a.databaseUsers[fmt.Sprintf(r.Form.Get("email"))]
			a.RUnlock()
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if password != r.Form.Get("password") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// save user in a temporary store for the user to be reconized later on
			a.Lock()
			a.extStore[fmt.Sprintf("LoggedInUserID-%v", r.Form.Get("email"))] = r.Form.Get("email")
			a.Unlock()

		}

		a.Authorize(w, r)
		return
	}
	http.Error(w, "Bad Request", http.StatusBadRequest)
}
