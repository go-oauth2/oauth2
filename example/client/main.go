package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

const (
	redirectURI = "http://localhost:9094/oauth2"
	serverURI   = "http://localhost:9096"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(serverURI + "/authorize")
		if err != nil {
			panic(err)
		}
		q := u.Query()
		q.Add("response_type", "code")
		q.Add("client_id", "222222")
		q.Add("scope", "all")
		q.Add("state", "xyz")
		q.Add("redirect_uri", url.QueryEscape(redirectURI))
		u.RawQuery = q.Encode()
		http.Redirect(w, r, u.String(), http.StatusFound)
	})

	http.HandleFunc("/oauth2", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		state := r.Form.Get("state")
		if state != "xyz" {
			http.Error(w, "State invalid", http.StatusBadRequest)
			return
		}
		code := r.Form.Get("code")
		if code == "" {
			http.Error(w, "Code not found", http.StatusBadRequest)
			return
		}
		uv := url.Values{}
		uv.Add("code", code)
		uv.Add("redirect_uri", redirectURI)
		uv.Add("grant_type", "authorization_code")
		uv.Add("client_id", "222222")
		uv.Add("client_secret", "22222222")
		resp, err := http.PostForm(serverURI+"/token", uv)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		io.Copy(w, resp.Body)
	})

	log.Println("OAuth2 client is running at 9094 port.")
	log.Fatal(http.ListenAndServe(":9094", nil))
}
