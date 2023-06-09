package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-oauth2/oauth2/v4/server"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"
)

var dumpvar bool

const (
	authServerURL string = "http://localhost:9096"
)

type Handlers struct {
	Permissions
	Authentication
}

func NewHandlers(dv bool, srv *server.Server) *Handlers {
	dumpvar = dv

	perm := NewPermissions(srv)
	auth := NewAuthentication(srv)

	return &Handlers{
		Permissions:    perm,
		Authentication: auth,
	}

}

type Permissions struct {
	srv *server.Server
}

func NewPermissions(srv *server.Server) Permissions {
	return Permissions{
		srv: srv,
	}
}

// Endpoint to validate token and permission
func (p Permissions) ValidPermission(w http.ResponseWriter, r *http.Request) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "validPermission", r) // Ignore the error
	}

	// validate the token
	token, err := p.srv.ValidationBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	permission := r.URL.Query().Get("permission")

	// validate the permission
	switch permission {
	case "read":
		log.Println("In read permission")
		if !strings.Contains(token.GetScope(), "read") && !strings.Contains(token.GetScope(), "all") {
			http.Error(w, "Unauthorized", http.StatusBadRequest)
			return
		}

	case "write":
		log.Println("In write permission")
		if !strings.Contains(token.GetScope(), "write") && !strings.Contains(token.GetScope(), "all") {
			fmt.Println("do not have Write permission.")
			http.Error(w, "Unauthorized", http.StatusBadRequest)
			return
		}

	case "all":
		log.Println("In all permission")
		if !strings.Contains(token.GetScope(), "all") {
			fmt.Println("do not have All permission.")
			http.Error(w, "Unauthorized", http.StatusBadRequest)
			return
		}
	default:
		log.Println("In default permission")
		http.Error(w, "Unauthorized", http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
		"client_id":  token.GetClientID(),
		"user_id":    token.GetUserID(),
		"permission": token.GetScope(),
	}
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(data)
}

func dumpRequest(writer io.Writer, header string, r *http.Request) error {
	data, err := httputil.DumpRequest(r, true)
	if err != nil {
		return err
	}
	writer.Write([]byte("\n" + header + ": \n"))
	writer.Write(data)
	return nil
}
