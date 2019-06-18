package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"

	"github.com/lib/pq"

	"github.com/jinzhu/gorm"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type User struct {
	Username string         `gorm:"unique,not null" json:"username,omitempty"`
	Password string         `gorm:"not null" json:"password,omitempty"`
	Contacts pq.StringArray `gorm:"type:varchar(255)[]" json:"contacts,omitempty"`
}

type Message struct {
	ID        uint      `gorm:"primary_key,auto_increment" json:"id,omitempty"`
	Sender    string    `gorm:"not null" json:"sender,omitempty"`
	Recipient string    `gorm:"not null" json:"recipient,omitempty"`
	Timestamp time.Time `gorm:"not null" json:"timestamp,omitempty"`
	Data      string    `json:"data,omitempty"`
}

type Config struct {
	Hostname    string `json:"hostname"`
	Port        string `json:"port"`
	DbDialect   string `json:"db_dialect"`
	DbArgs      string `json:"db_args"`
	ActivateLog bool   `json:"activate_log"`
}

var db *gorm.DB

func main() {
	defer shutdown()
	start(readConfig("config.json"))
}

func readConfig(filename string) Config {
	var data []byte
	var err error
	if data, err = ioutil.ReadFile(filename); err != nil {
		panic(err)
	}

	var conf Config
	if err = json.Unmarshal(data, &conf); err != nil {
		panic(err)
	}
	return conf
}

func start(conf Config) {
	var err error
	db, err = gorm.Open(conf.DbDialect, conf.DbArgs)
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{}, &Message{})

	r := mux.NewRouter()
	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users", postUser).Methods("POST")
	r.HandleFunc("/users/{username}", getUser).Methods("GET", "HEAD")
	r.HandleFunc("/users/{username}/contacts", postContact).Methods("POST")
	r.HandleFunc("/messages", postMessage).Methods("POST")
	r.HandleFunc("/messages", getMessages).Methods("GET").
		Queries("sender", "", "recipient", "")

	if conf.ActivateLog {
		loggedRouter := handlers.LoggingHandler(os.Stdout, r)
		http.ListenAndServe(":"+conf.Port, loggedRouter)
	} else {
		log.Fatal(http.ListenAndServe(":"+conf.Port, r))
	}
}

func shutdown() {
	if err := db.Close(); err != nil {
		log.Fatal(err.Error())
	}
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i := range users {
		users[i].Password = ""
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func postUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !db.Where("username = ?", user.Username).First(&user).RecordNotFound() {
		http.Error(w, "user "+user.Username+" already exists", http.StatusConflict)
		return
	}

	if err = db.Create(&user).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	var user User
	if db.Where("username = ?", username).First(&user).RecordNotFound() {
		http.Error(w, "user "+username+" not found", http.StatusNotFound)
		return
	}

	authUsername, authPassword, authOK := r.BasicAuth()

	if !authOK {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if authUsername != user.Username ||
		authPassword != user.Password {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if r.Method != "HEAD" {
		if err := json.NewEncoder(w).Encode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func postMessage(w http.ResponseWriter, r *http.Request) {
	type msgWithTimeAsString struct {
		ID        uint   `gorm:"primary_key,auto_increment" json:"id,omitempty"`
		Sender    string `gorm:"not null" json:"sender,omitempty"`
		Recipient string `gorm:"not null" json:"recipient,omitempty"`
		Timestamp string `gorm:"not null" json:"timestamp,omitempty"`
		Data      string `json:"data,omitempty"`
	}

	var msgWithStringTime msgWithTimeAsString

	if err := json.NewDecoder(r.Body).Decode(&msgWithStringTime); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	timestamp, err := time.Parse("2006-01-02T15:04:05-07:00", msgWithStringTime.Timestamp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	msg := Message{
		msgWithStringTime.ID,
		msgWithStringTime.Sender,
		msgWithStringTime.Recipient,
		timestamp,
		msgWithStringTime.Data,
	}

	var sender User
	if db.Where("username = ?", msg.Sender).First(&sender).RecordNotFound() {
		http.Error(w, "user "+msg.Sender+" doesn't exist", http.StatusNotFound)
		return
	}

	authUsername, authPassword, authOK := r.BasicAuth()

	if !authOK {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if authUsername != sender.Username ||
		authPassword != sender.Password {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var recipient User
	if db.Where("username = ?", msg.Recipient).First(&recipient).RecordNotFound() {
		http.Error(w, "user "+msg.Sender+" doesn't exist", http.StatusNotFound)
		return
	}

	if err := db.Create(&msg).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	senderUsername := r.FormValue("sender")
	recipientUsername := r.FormValue("recipient")

	var sender, recipient User
	if db.First(&sender, "username = ?", senderUsername).RecordNotFound() {
		http.Error(w, "user "+senderUsername+" not found", http.StatusNotFound)
		return
	}
	if db.First(&recipient, "username = ?", recipientUsername).RecordNotFound() {
		http.Error(w, "user "+recipientUsername+" not found", http.StatusNotFound)
		return
	}

	var msg []Message
	if err := db.Find(&msg, "sender = ? AND recipient = ?", senderUsername, recipientUsername).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authUsername, authPassword, authOK := r.BasicAuth()

	if !authOK {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if authUsername == senderUsername && authPassword == sender.Password ||
		authUsername == recipientUsername && authPassword == recipient.Password {

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(&msg); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	} else {
		w.WriteHeader(http.StatusForbidden)
		return
	}
}

func postContact(w http.ResponseWriter, r *http.Request) {

	var newContact User
	if err := json.NewDecoder(r.Body).Decode(&newContact); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	username := mux.Vars(r)["username"]

	var user User
	if db.First(&user, "username = ?", username).RecordNotFound() {
		http.Error(w, "user "+username+" not found", http.StatusNotFound)
		return
	}

	authUsername, authPassword, authOK := r.BasicAuth()

	if !authOK {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if authUsername != user.Username ||
		authPassword != user.Password {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var tempContact User
	if db.First(&tempContact, "username = ?", newContact.Username).RecordNotFound() {
		http.Error(w, "user "+newContact.Username+" not found", http.StatusNotFound)
		return
	}

	for _, contactFromList := range user.Contacts {
		if contactFromList == newContact.Username {
			http.Error(w, "user "+newContact.Username+" already contact of "+username, http.StatusConflict)
			return
		}
	}

	db.Model(&user).Where("username = ?", username).Update("contacts", append(user.Contacts, newContact.Username))
}
