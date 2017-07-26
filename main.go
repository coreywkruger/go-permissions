package main

import (
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"os"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/Jeffail/gabs"
)

func main(){

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
		fmt.Println(err)
	}
	P := Permissionist{
		DB: db,
	}

	router := mux.NewRouter()
	router.HandleFunc("/roles/{appID}", func (w http.ResponseWriter, r *http.Request) {
		fmt.Println("getting roles for app")
	}).Methods("GET")

	router.HandleFunc("/roles/{appID}", func (w http.ResponseWriter, r *http.Request) {

		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			w.Write([]byte("Could not create role"))
			return
		}
		defer r.Body.Close()

		body, err := gabs.ParseJSON(bodyBytes)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			w.Write([]byte("Could not create role"))
			return
		}

		roleName, ok := body.Path("role_name").Data().(string)
		if ok != true {
			w.WriteHeader(500)
			w.Write([]byte("Could not create role"))
			return
		}

		id, err := P.CreateRole(roleName, mux.Vars(r)["appID"])
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			w.Write([]byte("Could not create role"))
			return
		}

		w.WriteHeader(200)
		w.Write([]byte(id))
	}).Methods("POST")

	router.HandleFunc("/roles/create", func (w http.ResponseWriter, r *http.Request) {
		fmt.Println("creating role")
	})

	http.ListenAndServe(":8000", router)
}
