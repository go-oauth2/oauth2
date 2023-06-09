package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/go-oauth2/oauth2/v4/generates"
	"server/internal/handlers"

	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/jackc/pgx/v4"
	pg "github.com/vgarvardt/go-oauth2-pg/v4"
	"github.com/vgarvardt/go-pg-adapter/pgx4adapter"
)

var (
	dumpvar   bool
	idvar     string
	secretvar string
	domainvar string
	portvar   int
)

func init() {
	// credential for the client
	flag.BoolVar(&dumpvar, "d", true, "Dump requests and responses")
	flag.StringVar(&idvar, "i", "222222", "The client id being passed in")
	flag.StringVar(&secretvar, "s", "22222222", "The client secret being passed in")
	flag.StringVar(&domainvar, "r", "http://localhost:9094", "The domain of the redirect url")
	flag.IntVar(&portvar, "p", 9096, "the base port for the server")
}

const (
	// credential for the preOrder service
	idPreorder     string = "888888"
	secretPreorder string = "88888888"
	domainPreorder string = "http://localhost:8081"

	dbUser     = "postgres"
	dbHost     = "localhost"
	dbPassword = "password"
	dbDatabase = "users"
	dbSSL      = "disable"
	dbPort     = "5432"
)

func main() {
	flag.Parse()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbDatabase,
		dbSSL,
	)

	pgxConn, _ := pgx.Connect(context.TODO(), dsn)

	manager := manage.NewDefaultManager()

	// use PostgreSQL token store with pgx.Connection adapter
	adapter := pgx4adapter.NewConn(pgxConn)
	tokenStore, _ := pg.NewTokenStore(adapter, pg.WithTokenStoreGCInterval(time.Minute))
	defer tokenStore.Close()

	clientStore, _ := pg.NewClientStore(adapter)

	manager.MapTokenStorage(tokenStore)
	manager.MapClientStorage(clientStore)

	// generate jwt access token
	// manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte("00000000"), jwt.SigningMethodHS512))
	manager.MapAccessGenerate(generates.NewAccessGenerate())

	manager.MapClientStorage(clientStore)

	// register the front-end
	clientStore.Create(&models.Client{
		ID:     idvar,
		Secret: secretvar,
		Domain: domainvar,
	})

	// register prePost service
	clientStore.Create(&models.Client{
		ID:     idPreorder,
		Secret: secretPreorder,
		Domain: domainPreorder,
	})

	srv := server.NewServer(server.NewConfig(), manager)

	// set the oauth package to work without browser
	// the token will be return as a json payload
	srv.SetModeAPI()

	// handlers will handle all handlers
	handler := handlers.NewHandlers(dumpvar, srv)

	srv.SetUserAuthorizationHandler(handler.UserAuthorizeHandler)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	// Endpoints for the front-end
	// (use this service for the example but a specific users' service may be better)
	http.HandleFunc("/signup", handler.SignupHandler)
	http.HandleFunc("/login", handler.LoginHandler)

	// Endpoints for the backend services to authenticate and get their token
	http.HandleFunc("/apiauth", handler.ApiAuthHandler)

	// Endpoints specific to validate the authorization
	http.HandleFunc("/oauth/authorize", handler.Authorize)
	http.HandleFunc("/oauth/token", handler.Token)

	// Endpoint which validate a client's token and the given permission
	http.HandleFunc("/permission", handler.ValidPermission)

	log.Printf("Server is running at %d port.\n", portvar)
	log.Printf("Point your OAuth client Auth endpoint to %s:%d%s", "http://localhost", portvar, "/oauth/authorize")
	log.Printf("Point your OAuth client Token endpoint to %s:%d%s", "http://localhost", portvar, "/oauth/token")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", portvar), nil))
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

func createDsn() string {
	// dbHost := os.Getenv("DATABASE_HOST")
	// dbPort := os.Getenv("DATABASE_PORT")
	// dbUser := os.Getenv("DATABASE_USER")
	// dbPass := os.Getenv("DATABASE_PASS")
	// databaseName := os.Getenv("DATABASE_NAME")

	dbHost := "localhost"
	dbPort := "5432"
	dbUser := "postgres"
	dbPass := "password"
	databaseName := "users"

	dsnString := ""
	if dbPass == "" {
		dsnString = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s timezone=UTC connect_timeout=5",
			dbHost,
			dbPort,
			dbUser,
			databaseName,
			"disable")
	} else {
		dsnString = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s timezone=UTC connect_timeout=5",
			dbHost,
			dbPort,
			dbUser,
			dbPass,
			databaseName,
			"disable")
	}
	return dsnString
}
