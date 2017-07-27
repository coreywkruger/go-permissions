package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
)

func main() {

	config := viper.New()
	config.SetConfigFile(os.Getenv("CONFIG"))
	err := config.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	var dbconfig DbConfig
	config.UnmarshalKey("database", &dbconfig)

	db, err := InitDB(dbconfig)
	if err != nil {
		log.Fatal(err)
	}
	P := Permissionist{
		DB: db,
	}

	router := mux.NewRouter()

	router.HandleFunc("/apps", func(w http.ResponseWriter, r *http.Request) {
		id, err := P.CreateApp()
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			w.Write([]byte("Could not create app"))
			return
		}
		w.WriteHeader(200)
		w.Write([]byte(id))
	}).Methods("POST")

	router.HandleFunc("/roles/{appID}", func(w http.ResponseWriter, r *http.Request) {
		roleIDs, err := P.GetRoles(mux.Vars(r)["appID"])
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			w.Write([]byte("Could not get roles"))
			return
		}
		bytes, err := json.Marshal(roleIDs)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			w.Write([]byte("Could not get roles"))
			return
		}
		w.WriteHeader(200)
		w.Write(bytes)
	}).Methods("GET")

	router.HandleFunc("/roles/{appID}", func(w http.ResponseWriter, r *http.Request) {
		body := NewBody(w, r)
		id, err := P.CreateRole(body.GetField("role_name"), mux.Vars(r)["appID"])
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			w.Write([]byte("Could not create role"))
			return
		}
		w.WriteHeader(200)
		w.Write([]byte(id))
	}).Methods("POST")

	http.ListenAndServe(":8000", router)
}
