package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var tpl = template.Must(template.ParseFiles("index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := u.Query()
	searchQuery := params.Get("q")
	page := params.Get("page")
	if page == "" {
		page = "1"
	}

	entries, err := GetEntireWardrobe(w, r)
	if err != nil {
		fmt.Printf("Unable to fetch entries from wardrobe table")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	search := &Search{
		Query:   searchQuery,
		Entries: entries,
	}

	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf.WriteTo(w)

}

func main() {

	/*
		cred := new(credentials)

		cred.auth.username = os.Getenv("AUTH_USERNAME")
		cred.auth.password = os.Getenv("AUTH_PASSWORD")

		if cred.auth.username == "" {
			log.Fatal("basic auth username must be provided")
		}

		if cred.auth.password == "" {
			log.Fatal("basic auth password must be provided")
		}
	*/

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	cassInit()
	defer closeCass()

	fs := http.FileServer(http.Dir("frontend"))

	m := mux.NewRouter()
	m.Handle("/frontend/", http.StripPrefix("/frontend/", fs))
	m.HandleFunc("/", indexHandler)

	// Login API
	//m.HandleFunc("/login", cred.basicAuth(cred.loginHandler))
	m.HandleFunc("/login", userLogin)
	m.HandleFunc("/search", searchHandler) // returns html page

	// DB APIs

	m.HandleFunc("/item", func(w http.ResponseWriter, r *http.Request) { AddItem(w, r) }).Methods("POST")
	m.HandleFunc("/item", func(w http.ResponseWriter, r *http.Request) { UpdateItem(w, r) }).Methods("PUT")
	//m.HandleFunc("/item", func(w http.ResponseWriter, r *http.Request) { GetItem(w,r) }).Methods("GET")
	//m.HandleFunc("/item", func(w http.ResponseWriter, r *http.Request) { DeleteItem(w,r) }).Methods("DELETE")

	m.HandleFunc("/allitems", func(w http.ResponseWriter, r *http.Request) { AddAllItems(w, r) }).Methods("POST")
	m.HandleFunc("/allitems", func(w http.ResponseWriter, r *http.Request) { GetAllItems(w, r) }).Methods("GET")
	m.HandleFunc("/allitems", func(w http.ResponseWriter, r *http.Request) { DeleteAllItems(w, r) }).Methods("DELETE")

	m.HandleFunc("/count", func(w http.ResponseWriter, r *http.Request) { CountAllItems(w, r) }).Methods("GET")
	m.HandleFunc("/getfor", func(w http.ResponseWriter, r *http.Request) { GetForPerson(w, r) }).Methods("GET")
	m.HandleFunc("/pickootd", func(w http.ResponseWriter, r *http.Request) { PickOOTD(w, r) }).Methods("GET")

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	log.Fatal(http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), handlers.CORS(headers, methods, origins)(m)))

	//closeCass()

}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

/*
func (cred *credentials) basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pwd, ok := r.BasicAuth()
		if ok {
			uHash := sha256.Sum256([]byte(user))
			pHash := sha256.Sum256([]byte(pwd))
			expUHash := sha256.Sum256([]byte(cred.auth.username))
			expPHash := sha256.Sum256([]byte(cred.auth.password))

			uMatch := (subtle.ConstantTimeCompare(uHash[:], expUHash[:]) == 1)
			pMatch := (subtle.ConstantTimeCompare(pHash[:], expPHash[:]) == 1)

			if uMatch && pMatch {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

func (cred *credentials) loginHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintln(w, "This is the protected handler")
}
*/

func userLogin(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var user Credentials
	dbUser := &Credentials{
		Username: "User1",
		Password: "passwd",
	}

	json.NewDecoder(r.Body).Decode(&user)
	if user.Username != dbUser.Username {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"response":"Failure"}`))
		return
	}

	if user.Password != dbUser.Password {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"response":"Failure"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"response":"Success"}`))
}
