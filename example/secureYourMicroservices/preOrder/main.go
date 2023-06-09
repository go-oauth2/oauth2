package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	authServerURL  = "http://localhost:9096"
	orderServerURL = "http://localhost:8080"
	clientID       = "888888"
)

var (
	config = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: "88888888",
		Scopes:       []string{"read"},
		RedirectURL:  "http://localhost:8081",
		Endpoint: oauth2.Endpoint{
			AuthURL:  authServerURL + "/oauth/authorize",
			TokenURL: authServerURL + "/oauth/token",
		},
	}
	globalToken  *oauth2.Token // Non-concurrent security
	appUrlParams string
)

type CodeStruct struct {
	Code string `json:"code"`
}

type LineItem struct {
	ItemCode string `json:"item_code"`
	Quantity int    `json:"quantity"`
}

type Order struct {
	ID              string     `json:"id"`
	Items           []LineItem `json:"items"`
	ShippingAddress string     `json:"shipping_address"`
}

// authenticate with the authentication server
func init() {
	u := config.AuthCodeURL("xyz",
		oauth2.SetAuthURLParam("code_challenge", genCodeChallengeS256("s256example")),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"))

	// extract the url params which will be use later on
	parsedURL, err := url.Parse(u)
	if err != nil {
		log.Fatal("Error parsing URL:", err)
		return
	}

	appUrlParams = parsedURL.RawQuery

	resp, err := http.Get(u)
	if err != nil {
		// Handle error if unable to make the request
		log.Fatal("Error Authenticating with auth service: ", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Println("Authentication OK")
	} else {
		log.Fatal("Authentication failed, make sure the authentication server is running")
	}

}

func main() {

	// get the token(/apiauth is the endpoint for the backend services to get their token)
	getTheToken("apiauth")

	http.HandleFunc("/order", getOrders())

	fmt.Println("start producer-api on port 8081... !!")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func getOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// in the real world, obviously, a proxy/apiGateway would stand in front
		if carryon := allowCORS(w, r); !carryon {
			return
		}

		accessToken := r.Header.Get("Authorization")

		accessToken, _ = extractBearerToken(accessToken)

		var orderID string
		if r.Method == "POST" {
			if r.Form == nil {
				if err := r.ParseForm(); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			orderID = r.Form.Get("orderid")
		}

		// make sure the token from the caller(here the frontend) is valid
		// and make sure it have the right permission
		resp, err := http.Get(fmt.Sprintf("%s/permission?permission=read&access_token=%s", authServerURL, accessToken))
		if err != nil || resp.StatusCode == http.StatusBadRequest {
			log.Println("Http req err or Invalid oauth token")
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			} else {
				http.Error(w, "Invalid authentication", http.StatusBadRequest)
				return
			}
		}

		// Create the form data
		form := url.Values{}
		form.Set("orderid", orderID)

		// call the order service the get the order
		resp, err = http.PostForm(fmt.Sprintf("%s/order?access_token=%s", orderServerURL, globalToken.AccessToken), form)
		if err != nil || resp.StatusCode == http.StatusBadRequest {
			log.Println("Http req err or Invalid oauth token")
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			} else {
				http.Error(w, "Invalid authentication", http.StatusBadRequest)
				return
			}
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
		}

		// run the whatever logic of this service, here we change the shippingAddress
		var order Order
		err = json.Unmarshal(body, &order)
		if err != nil {
			log.Println("Error parsing response body:", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
		order.ShippingAddress = "modified by preOrder"

		// send the order
		ordeBts, err := json.Marshal(order)
		if err != nil {
			log.Println("Error marshaling resp body:", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")

		_, err = w.Write(ordeBts)
		if err != nil {
			log.Println("err sending back the order: ", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
	}
}

func allowCORS(w http.ResponseWriter, r *http.Request) bool {
	// Set the CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:9094")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Handle preflight requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return false
	}

	return true

}

func extractBearerToken(tokenWithBearer string) (string, bool) {
	bearerPrefix := "Bearer "
	if strings.HasPrefix(tokenWithBearer, bearerPrefix) {
		token := strings.TrimPrefix(tokenWithBearer, bearerPrefix)
		return token, true
	}
	return "", false
}

func getTheToken(path string) error {

	u, err := url.Parse(fmt.Sprintf("%s/%s", authServerURL, path))
	if err != nil {
		log.Println("Error parsing URL:", err)
		return err
	}

	// Add the form data to the existing query parameters
	queryParams := u.Query()

	// Add the app parmeters
	appParams, err := url.ParseQuery(appUrlParams)
	if err != nil {
		log.Println("Error parsing query string:", err)
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
		log.Println("Error request the token: ", err)
		return err
	}

	defer resp.Body.Close()

	// Handle the response as needed
	if resp.StatusCode != http.StatusOK {
		log.Println("Error StatusCode request the token: ", resp)
		return nil
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return err
	}

	// Parse the response body into CodeStruct
	var codeResponse CodeStruct
	err = json.Unmarshal(body, &codeResponse)
	if err != nil {
		log.Println("Error parsing response body:", err)
		return err
	}

	// Access the extracted code value
	token, err := config.Exchange(context.Background(), codeResponse.Code, oauth2.SetAuthURLParam("code_verifier", "s256example"))
	if err != nil {
		log.Println("Error exchange the token:", err)
		return err
	}
	globalToken = token
	return nil

}

func genCodeChallengeS256(s string) string {
	s256 := sha256.Sum256([]byte(s))
	return base64.URLEncoding.EncodeToString(s256[:])
}
