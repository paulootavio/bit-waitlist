package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

type Message struct {
	Text string `json:"text"`
}

type Waitlist struct {
	Email string `json:"email"`
}

func main() {

	http.HandleFunc("/waitlist", helloHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

	defer db.Close()

}

func helloHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var waitlist Waitlist
	body, errres := ioutil.ReadAll(r.Body)

	if errres != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &waitlist); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	//Verificar se o parâmetro foi fornecido
	if waitlist.Email == "" {
		http.Error(w, "Parâmetro 'nome' não encontrado", http.StatusBadRequest)
		return
	}

	err2 := godotenv.Load(".env")
	if err2 != nil {
		log.Fatal("Error loading .env file")
	}

	port := 5432
	user := os.Getenv("username")
	password := os.Getenv("password")
	dbname := os.Getenv("dbname")
	hostnamedb := os.Getenv("hostnamedb")

	//postgres
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hostnamedb, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	_, errDb := db.Exec("INSERT INTO waitlist (email) VALUES ($1)", waitlist.Email)

	if errDb != nil {
		http.Error(w, errDb.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data saved successfully"))

	message := Message{Text: "Hello, world!"}
	json.NewEncoder(w).Encode(message)
}
