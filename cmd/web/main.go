package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/ishanshre/Book-Review-Platform/internals/config"
	"github.com/ishanshre/Book-Review-Platform/internals/driver"
	"github.com/ishanshre/Book-Review-Platform/internals/handler"
	"github.com/ishanshre/Book-Review-Platform/internals/helpers"
	"github.com/ishanshre/Book-Review-Platform/internals/middleware"
	"github.com/ishanshre/Book-Review-Platform/internals/models"
	"github.com/ishanshre/Book-Review-Platform/internals/render"
	"github.com/ishanshre/Book-Review-Platform/internals/router"
	"github.com/joho/godotenv"
)

var app config.AppConfig // global config
var session *scs.SessionManager
var host string = "127.0.0.1:6379"
var database string = "postgres"
var connString string

var infoLog *log.Logger
var errorLog *log.Logger

type FileLogger struct {
	file *os.File
}

func (f *FileLogger) Write(p []byte) (n int, err error) {
	return f.file.Write(p)
}

func main() {
	// using flag for command line arguments
	port := flag.Int("port", 8000, "The port to run the web application")

	flag.Parse()

	db, err := Run()
	if err != nil {
		app.ErrorLog.Println(err)
	}
	app.InfoLog.Println("Connected to database")

	defer db.SQL.Close()

	// close the channel at last
	defer close(app.MailChan)

	app.InfoLog.Println("Starting the mail listener")
	// starting the mail listener
	listenForMail()

	// pass app config to middleware
	middleware.NewMiddlewareApp(&app)

	addr := fmt.Sprintf(":%d", *port)
	// create a http server with address and the handlers
	srv := http.Server{
		Addr:    addr,
		Handler: router.Router(&app),
	}
	app.InfoLog.Printf("Starting server at port %d", *port)

	// start the server and listen to the specified port
	if err := srv.ListenAndServe(); err != nil {
		app.ErrorLog.Fatalf("error in listining to server: %s", err)
	}
}

func Run() (*driver.DB, error) {

	// create a logger with log.New()
	infoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	// assign logs to global configs
	app.InfoLog = infoLog
	app.ErrorLog = errorLog
	// load the environment variable from the .env file
	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("error in loading environment files")
	}

	// store the values in the session
	gob.Register(models.User{})

	// create a mail channel and assign it to app.MailChan
	mailChan := make(chan models.MailData, 10)
	app.MailChan = mailChan

	// change to true in production
	app.InProduction = false
	app.UseRedis = true

	// Initiate a session and configure it
	session = scs.New()

	// Establish a pool to Redis if UseRedis config is true
	if app.UseRedis {
		pool := &redis.Pool{
			MaxIdle: 10,
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", host)
			},
		}
		session.Store = redisstore.New(pool)
	}
	session.Lifetime = 24 * time.Hour // set time of the session
	session.Cookie.Persist = true     // true means session retains in browser even if browser is closed
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session // make session available to whole application

	// initiate the template cache
	tc, err := render.CreateTemplateCache()
	if err != nil {
		infoLog.Println(err)
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
