package main

import (
	"context"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/NainVictorin1/homework2/Internal/data"
	_ "github.com/lib/pq"
)

type application struct {
	addr     *string
	feedback *data.FeedbackModel

	logger        *slog.Logger
	templateCache map[string]*template.Template
	todos         *data.TodoModel
	journals      *data.JournalModel
}

func main() {
	initDatabase() // Ensure this function is properly implemented

	addr := flag.String("addr", ":8080", "HTTP network address") // Default port if empty
	dsn := flag.String("dsn", "", "PostgreSQL DSN")
	flag.Parse()

	// Set up logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Open database connection
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("database connection pool established")

	// Create template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Initialize the application struct
	app := &application{
		addr:          addr,
		feedback:      &data.FeedbackModel{DB: db},
		logger:        logger,
		templateCache: templateCache,
		journals:      &data.JournalModel{DB: db},
		todos:         &data.TodoModel{DB: db},
	}

	// Register routes
	registerRoutes(app)

	// Start the server
	err = app.serve()
	if err != nil {
		// Log the error if server fails to start
		logger.Error("Server failed", "error", err)
		os.Exit(1)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func registerRoutes(app *application) {
	// Register handlers with logging middleware
	http.HandleFunc("/", loggingMiddleware(http.HandlerFunc(app.homeHandler)).ServeHTTP)

	http.HandleFunc("/journal", loggingMiddleware(http.HandlerFunc(app.submitJournalHandler)).ServeHTTP)
	http.HandleFunc("/journals", loggingMiddleware(http.HandlerFunc(app.viewJournalsHandler)).ServeHTTP)

	http.HandleFunc("/todos", loggingMiddleware(http.HandlerFunc(app.viewTodosHandler)).ServeHTTP)
	http.HandleFunc("/todo", loggingMiddleware(http.HandlerFunc(app.addTodoHandler)).ServeHTTP)

	http.HandleFunc("/feedbacks", loggingMiddleware(http.HandlerFunc(app.viewFeedbacksHandler)).ServeHTTP)
	http.HandleFunc("/feedback", loggingMiddleware(http.HandlerFunc(app.submitFeedbackHandler)).ServeHTTP)

}

// loggingMiddleware is an example middleware that logs HTTP requests.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log request details
		log.Printf("Received request: %s %s", r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

// Example handler for home page
func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, "home.tmpl", nil)
}

func (app *application) addJournalHandler(w http.ResponseWriter, r *http.Request) {
	// Your handler logic for submitting a journal
	app.renderTemplate(w, "submit_journal.tmpl", nil)
}

// Example handler for submitting feedback
func (app *application) addFeedbackHandler(w http.ResponseWriter, r *http.Request) {
	// Your handler logic for submitting feedback
	app.renderTemplate(w, "feedback.tmpl", nil)
}

// Example handler for adding todo
func (app *application) addTodoHandler(w http.ResponseWriter, r *http.Request) {
	// Your handler logic for adding todo
	app.renderTemplate(w, "add_todo.tmpl", nil)
}

// Example handler for viewing journals
func (app *application) listJournalsHandler(w http.ResponseWriter, r *http.Request) {
	// Your handler logic for viewing journals
	app.renderTemplate(w, "view_journal.tmpl", nil)
}

// Example handler for viewing feedbacks
func (app *application) listFeedbacksHandler(w http.ResponseWriter, r *http.Request) {
	// Your handler logic for viewing feedback
	app.renderTemplate(w, "view_feedback.tmpl", nil)
}

// Example handler for viewing todos
func (app *application) listTodosHandler(w http.ResponseWriter, r *http.Request) {
	// Your handler logic for viewing todos
	app.renderTemplate(w, "view_todo.tmpl", nil)
}

// renderTemplate is a utility function to render templates.
func (app *application) renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	tmplObj, ok := app.templateCache[tmpl]
	if !ok {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	err := tmplObj.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// initDatabase is a placeholder for your database initialization function (if needed).
func initDatabase() {
	// Initialize the database connection or setup (if required).
}
