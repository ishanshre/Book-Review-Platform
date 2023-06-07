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
	"github.com/ishanshre/Book-Review-Platform/internals/handler"
	"github.com/ishanshre/Book-Review-Platform/internals/helpers"
	"github.com/ishanshre/Book-Review-Platform/internals/models"
	"github.com/ishanshre/Book-Review-Platform/internals/render"
	"github.com/joho/godotenv"
)

var port string
var app config.AppConfig // global config
var session *scs.SessionManager

var database string = "postgres"
var connString string

var infoLog *log.Logger
var errorLog *log.Logger

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

	// close the channel at last
	defer close(app.MailChan)

	log.Println("Starting the mail listener")
	// starting the mail listener
	listenForMail()

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

	// create a mail channel and assign it to app.MailChan
	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// change to true in production
	app.InProduction = false

	// create a logger with log.New()
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// assign logs to global configs
	app.InfoLog = infoLog
	app.ErrorLog = errorLog

	// Initiate a session and configure it
	session = scs.New()
	session.Lifetime = 24 * time.Hour // set time of the session
	session.Cookie.Persist = true     // true means session retains in browser even if browser is closed
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session // make session available to whole application

	// initiate the template cache
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// store the templates into global app config
	app.TemplateCache = tc
	app.UseCache = false

	// pass the global app config reference to render app
	render.NewRenderer(&app)

	// pass the global config to handler

	log.Println("Connecting to database")
	connString = os.Getenv("postgres")
	db, err := driver.ConnectSQL(database, connString)
	if err != nil {
		return nil, fmt.Errorf("error in connecting to database: %v", err)
	}

	// handlers connecting to database
	repo := handler.NewRepo(&app, db)
	handler.NewHandler(repo)

	helpers.NewHelpers(&app)

	return db, nil
}
