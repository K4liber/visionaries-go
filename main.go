package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"visionaries/database"
	myHandlers "visionaries/handlers"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

//IP - host IP
const IP = "46.41.149.6"

//PORT - host Port
const PORT = "80"

var db *gorm.DB

func main() {

	//Database connection
	db, err := gorm.Open("mysql", "root:Potoczek30@tcp/iidb")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.SingularTable(true)
	database.SetDB(db)

	var entry string
	var static string
	var port string
	flag.StringVar(&entry, "entry", "./build/index.html", "the entrypoint to serve.")
	flag.StringVar(&static, "static", ".", "the directory to serve static files from.")
	flag.StringVar(&port, "port", PORT, "the `port` to listen on.")
	flag.Parse()

	r := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:     []string{"*", "46.41.136.25", "visionaries.pl", IP},
		OptionsPassthrough: true,
		AllowedHeaders:     []string{"*"},
	})

	api := r.PathPrefix("/app/").Subrouter()
	api.Handle("/mems", c.Handler(myHandlers.MemsHandler))
	api.Handle("/mem/{id}", c.Handler(myHandlers.MemHandler))
	api.Handle("/profile/{nickname}", c.Handler(myHandlers.ProfileHandler))
	api.Handle("/activities/{nickname}", c.Handler(myHandlers.ActivitiesHandler))
	api.Handle("/category/{category}", c.Handler(myHandlers.CategoryHandler))
	api.Handle("/addMem", c.Handler(myHandlers.PreHandler(myHandlers.AddMemHandler)))
	api.Handle("/uploadAvatar", c.Handler(myHandlers.PreHandler(myHandlers.UploadAvatarHandler)))
	api.Handle("/addComment", c.Handler(myHandlers.PreHandler(myHandlers.AddCommentHandler)))
	api.Handle("/addMemPoint", c.Handler(myHandlers.PreHandler(myHandlers.AddMemPointHandler)))
	api.Handle("/deleteMemPoint", c.Handler(myHandlers.PreHandler(myHandlers.DeleteMemPointHandler)))
	api.Handle("/deleteMem", c.Handler(myHandlers.PreHandler(myHandlers.DeleteMemHandler)))
	api.Handle("/adminDeleteMem", c.Handler(myHandlers.PreHandler(myHandlers.AdminDeleteMemHandler)))
	api.Handle("/addCommentPoint", c.Handler(myHandlers.PreHandler(myHandlers.AddCommentPointHandler)))
	api.Handle("/deleteCommentPoint", c.Handler(myHandlers.PreHandler(myHandlers.DeleteCommentPointHandler)))
	api.Handle("/deleteComment", c.Handler(myHandlers.PreHandler(myHandlers.DeleteCommentHandler)))
	api.Handle("/adminDeleteComment", c.Handler(myHandlers.PreHandler(myHandlers.AdminDeleteCommentHandler)))

	// Serve static assets directly.
	r.PathPrefix("/resources").Handler(http.FileServer(http.Dir(static)))
	r.PathPrefix("/img").Handler(http.FileServer(http.Dir(static)))
	r.PathPrefix("/static").Handler(http.FileServer(http.Dir(static)))
	r.PathPrefix("/favicon.ico").HandlerFunc(myHandlers.IconHandler("./build/favicon.ico"))
	r.PathPrefix("/").HandlerFunc(myHandlers.IndexHandler(entry))

	srv := &http.Server{
		Handler: handlers.LoggingHandler(os.Stdout, r),
		Addr:    IP + ":" + PORT,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
