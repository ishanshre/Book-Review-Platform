package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/ishanshre/Book-Review-Platform/internals/config"
	"github.com/ishanshre/Book-Review-Platform/internals/driver"
	"github.com/ishanshre/Book-Review-Platform/internals/models"
	"github.com/joho/godotenv"
)

var port string
var app config.AppConfig // global config
var session *scs.SessionManager

var database string = "postgres"
var connString string

func main() {
	// checking if command has argument
	if len(os.Args) != 2 {
		log.Fatalln("usuage: command <port>")
		log.Fatalln("port must be from 1025 to 65535")
	}

	p, err := strconv.Atoi(os.Args[1]) // convert the second argument to integer
	if err != nil {
		log.Fatalln("ports must be integers")
	}
	// load the environment variable from the .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("error in loading environment files: %v\n", err)
	}
	db, err := Run()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Connected to database")

	defer db.SQL.Close()

	port = fmt.Sprintf(":%v", p)
	// create a http server with address and the handlers
	srv := http.Server{
		Addr:    port,
		Handler: router(&app),
	}
	log.Printf("Starting server at port %v", p)

	// start the server and listen to the specified port
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("error in listining to server: %s", err)
	}
}

func Run() (*driver.DB, error) {
	// store the values in the session
	gob.Register(models.User{})

	app.InProduction = false

	// Initiate a session and configure it
	session = scs.New()
	session.Lifetime = 24 * time.Hour // set time of the session
	session.Cookie.Persist = true     // true means session retains in browser even if browser is closed
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session // make session available to whole application

	// connecting to database
	log.Println("Connecting to database")
	connString = os.Getenv("postgres")
	db, err := driver.ConnectSQL(database, connString)
	if err != nil {
		return nil, fmt.Errorf("error in connecting to database: %v", err)
	}
	return db, nil
}
