package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"context"
	"errors"
	"io"
	"net"
	"net/http"

	"github.com/apucontilde/illustrious-otter/transaction"
	_ "github.com/mattn/go-sqlite3"
)

const keyServerAddr = "serverAddr"

func getRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	hasFirst := r.URL.Query().Has("first")
	first := r.URL.Query().Get("first")
	hasSecond := r.URL.Query().Has("second")
	second := r.URL.Query().Get("second")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read body: %s\n", err)
	}

	fmt.Printf("%s: got / request. first(%t)=%s, second(%t)=%s, body:\n%s\n",
		ctx.Value(keyServerAddr),
		hasFirst, first,
		hasSecond, second,
		body)
	io.WriteString(w, "This is my website!\n")

}
func getHello(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Printf("%s: got /hello request\n", ctx.Value(keyServerAddr))
	method := r.Method
	fmt.Printf("method: %s\n", method)
	if method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	myName := r.PostFormValue("myName")
	if myName == "" {
		w.Header().Set("x-missing-field", "myName")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	io.WriteString(w, fmt.Sprintf("Hello, %s!\n", myName))
}

func transactionCrud(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/hello", getHello)
	mux.HandleFunc("/transaction", transactionCrud)
}

func Serve(port string, repository *transaction.SQLiteRepository) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/hello", getHello)
	mux.HandleFunc("/transaction", transactionCrud)

	// ctx, cancelCtx := context.WithCancel(context.Background())
	ctx := context.Background()
	server := &http.Server{
		Addr:    ":3333",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}
	fmt.Printf("listening on localhost%s\n", port)
	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error listening for server: %s\n", err)
	}
}

const db_file = "sqlite.db"

func main() {
	var do_migrate = true
	if _, err := os.Stat(db_file); err == nil {
		do_migrate = false
	} else if errors.Is(err, os.ErrNotExist) {
	} else {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", db_file)
	if err != nil {
		log.Fatal(err)
	}

	transactionRepository := transaction.NewSQLiteRepository(db)

	if do_migrate {
		if err := transactionRepository.Migrate(); err != nil {
			log.Fatal(err)
		}

	}

	Serve(":3333", transactionRepository)

}
