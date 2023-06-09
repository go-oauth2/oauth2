package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"order/internal/service"
	"order/internal/types"
	"strconv"
	"strings"
)

const (
	authServerURL = "http://localhost:9096"
)

func main() {

	// the storage
	orderStr := service.NewOrderSvc()

	http.HandleFunc("/orders", createOrders(&orderStr))
	http.HandleFunc("/order", getOrders(&orderStr))

	fmt.Println("start producer-api on port 8080... !!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createOrders(orderStr *service.OrderSvc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// in the real world, obviously, a proxy/apiGateway would stand in front
		if carryon := allowCORS(w, r); !carryon {
			return
		}

		accessToken := r.Header.Get("Authorization")

		// get the access token
		accessToken, _ = extractBearerToken(accessToken)

		// make sure the token from the caller(here the frontend) is valid
		// and make sure it have the right permission
		resp, err := http.Get(fmt.Sprintf("%s/permission?permission=write&access_token=%s", authServerURL, accessToken))
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

		// get the form inputs
		o := types.Order{}
		if r.Method == "POST" {

			// Parse the multipartForm
			if err := r.ParseMultipartForm(0); err != nil {
				log.Println("Failed to parse form data:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Access the form data
			form := r.MultipartForm
			values := form.Value

			iq, _ := strconv.Atoi(values["item1quantity"][0])
			li1 := types.LineItem{
				ItemCode: values["item1code"][0],
				Quantity: iq,
			}
			iq, _ = strconv.Atoi(values["item2quantity"][0])
			li2 := types.LineItem{
				ItemCode: values["item2code"][0],
				Quantity: iq,
			}

			o.Items = []types.LineItem{li1, li2}
			o.ShippingAddress = r.Form.Get("shippingaddress")
		}

		// create an uid
		newUUID := uuid.New()
		uuidString := newUUID.String()

		o.ID = uuidString

		// save the received data
		orderStr.PlaceOrder(o)

		fmt.Fprintf(w, uuidString)
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

func getOrders(orderStr *service.OrderSvc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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

		// get the access_token
		accessToken := r.URL.Query().Get("access_token")

		// make sure the token from the caller(here the preOrder service) is valid
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

		// get the data back(to make sure)
		v, _ := orderStr.GetOrder(orderID)

		ordeBts, err := json.Marshal(v)
		if err != nil {
			log.Println("err marshalling: ", err)
		}

		w.Header().Set("Content-Type", "application/json")

		_, err = w.Write(ordeBts)
		if err != nil {
			log.Println("err sending back the order: ", err)
		}
	}
}
