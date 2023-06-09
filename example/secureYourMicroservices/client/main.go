package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"sync"

	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
	// "golang.org/x/oauth2/clientcredentials"
)

const (
	authServerURL = "http://localhost:9096"
	clientID      = "222222"
)

var (
	config = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: "22222222",
		Scopes:       []string{"write, read"},
		RedirectURL:  "http://localhost:9094",
		Endpoint: oauth2.Endpoint{
			AuthURL:  authServerURL + "/oauth/authorize",
			TokenURL: authServerURL + "/oauth/token",
		},
	}
	globalTokens = make(map[string]*oauth2.Token) // should be persisted outside, keep it stateless
	// globalToken  *oauth2.Token // for mobile app that is ok as each app run its own instance
	appUrlParams string
	mu           = sync.RWMutex{}
)

type CodeStruct struct {
	Code string `json:"code"`
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		u := config.AuthCodeURL("xyz",
			oauth2.SetAuthURLParam("code_challenge", genCodeChallengeS256("s256example")),
			oauth2.SetAuthURLParam("code_challenge_method", "S256"))

		// extract the url params which will be use later on
		parsedURL, err := url.Parse(u)
		if err != nil {
			log.Println("Error parsing URL:", err)
			return
		}

		appUrlParams = parsedURL.RawQuery

		resp, err := http.Get(u)
		if err != nil {
			// Handle error if unable to make the request
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			http.Redirect(w, r, fmt.Sprintf("/portail?id=%v", clientID), http.StatusFound)
		} else {

			http.Error(w, "An err occur", http.StatusInternalServerError)
		}
	})

	// Portail shows the signup or login buttons
	http.HandleFunc("/portail", func(w http.ResponseWriter, r *http.Request) {

		outputHTML(w, r, "static/portail.html")
	})

	// Render the form for signup
	http.HandleFunc("/signupfront", func(w http.ResponseWriter, r *http.Request) {

		outputHTML(w, r, "static/signup.html")
	})

	// Handle the data from the form signup
	http.HandleFunc("/signupdata", func(w http.ResponseWriter, r *http.Request) {

		err := helperHandleAuth(w, r, "signup")
		if err != nil {
			log.Println("Error in signupdata: ", err)
		}
	})

	// Render the form for login
	http.HandleFunc("/loginfront", func(w http.ResponseWriter, r *http.Request) {

		outputHTML(w, r, "static/login.html")
	})

	// Handle the data for login
	http.HandleFunc("/logindata", func(w http.ResponseWriter, r *http.Request) {
		err := helperHandleAuth(w, r, "login")
		if err != nil {
			log.Println("Error in logindata: ", err)
		}
	})

	// After login/signup show the menu
	http.HandleFunc("/welcome", func(w http.ResponseWriter, r *http.Request) {

		outputHTML(w, r, "static/welcome.html")

	})

	// Render the page to create order
	http.HandleFunc("/createorder", func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")

		mu.RLock()
		_, ok := globalTokens[email]
		mu.RUnlock()
		if ok {

			outputHTML(w, r, "static/createOrders.html")
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	})

	// Render the page to get an order
	http.HandleFunc("/getorder", func(w http.ResponseWriter, r *http.Request) {

		email := r.URL.Query().Get("email")

		mu.RLock()
		_, ok := globalTokens[email]
		mu.RUnlock()
		if ok {

			outputHTML(w, r, "static/getOrders.html")
		} else {

			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	})

	http.HandleFunc("/signout", func(w http.ResponseWriter, r *http.Request) {

		email := r.URL.Query().Get("email")

		mu.RLock()
		_, ok := globalTokens[email]
		mu.RUnlock()
		if ok {
			// delete the token
			mu.Lock()
			delete(globalTokens, email)
			mu.Unlock()

			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {

			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	})

	// Other functionalities
	// NOTE for the reader to finish the implementation if needed
	// http.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
	// 	if globalToken == nil {
	// 		http.Redirect(w, r, "/", http.StatusFound)
	// 		return
	// 	}

	// 	globalToken.Expiry = time.Now()
	// 	token, err := config.TokenSource(context.Background(), globalToken).Token()
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}

	// 	globalToken = token
	// 	e := json.NewEncoder(w)
	// 	e.SetIndent("", "  ")
	// 	e.Encode(token)
	// })

	// http.HandleFunc("/pwd", func(w http.ResponseWriter, r *http.Request) {
	// 	token, err := config.PasswordCredentialsToken(context.Background(), "test", "test")
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}

	// 	globalToken = token
	// 	e := json.NewEncoder(w)
	// 	e.SetIndent("", "  ")
	// 	e.Encode(token)
	// })

	// http.HandleFunc("/client", func(w http.ResponseWriter, r *http.Request) {
	// 	cfg := clientcredentials.Config{
	// 		ClientID:     config.ClientID,
	// 		ClientSecret: config.ClientSecret,
	// 		TokenURL:     config.Endpoint.TokenURL,
	// 	}

	// 	token, err := cfg.Token(context.Background())
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}

	// 	e := json.NewEncoder(w)
	// 	e.SetIndent("", "  ")
	// 	e.Encode(token)
	// })

	log.Println("Client is running at 9094 port.Please open http://localhost:9094")
	log.Fatal(http.ListenAndServe(":9094", nil))
}

// Target the  auth server to get a token, and extract it
func helperHandleAuth(w http.ResponseWriter, r *http.Request, path string) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	u, err := url.Parse(fmt.Sprintf("%s/%s", authServerURL, path))
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return err
	}

	// Add the form data to the existing query parameters
	queryParams := u.Query()
	for key, values := range r.Form {
		for _, value := range values {
			queryParams.Add(key, value)
		}
	}

	// Add the app parmeters
	appParams, err := url.ParseQuery(appUrlParams)
	if err != nil {
		fmt.Println("Error parsing query string:", err)
		return err
	}

	for key, values := range appParams {
		for _, value := range values {
			queryParams.Add(key, value)
		}
	}

	u.RawQuery = queryParams.Encode()

	resp, err := http.Post(u.String(), "application/x-www-form-urlencoded", strings.NewReader(""))
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Redirect(w, r, fmt.Sprintf("/portail?id=%v", clientID), http.StatusFound)
		return nil
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return err
	}

	// Parse the response body into CodeStruct
	var codeResponse CodeStruct
	err = json.Unmarshal(body, &codeResponse)
	if err != nil {
		fmt.Println("Error parsing response body:", err)
		return err
	}

	// extract the token
	token, err := config.Exchange(context.Background(), codeResponse.Code, oauth2.SetAuthURLParam("code_verifier", "s256example"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	// save the token for a specific user
	email := r.Form.Get("email")
	mu.Lock()
	globalTokens[email] = token
	mu.Unlock()

	outputHTML(w, r, "static/welcome.html")
	return nil

}

func outputHTML(w http.ResponseWriter, req *http.Request, filename string) {
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	email := req.URL.Query().Get("email")

	tk := ""
	t, ok := globalTokens[email]
	if ok {
		tk = t.AccessToken
	}

	data := struct {
		Token string
	}{
		Token: tk,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func genCodeChallengeS256(s string) string {
	s256 := sha256.Sum256([]byte(s))
	return base64.URLEncoding.EncodeToString(s256[:])
}
