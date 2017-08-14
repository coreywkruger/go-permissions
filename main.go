package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

func handleCreateApp(P *Permissionist) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := NewBody(w, r)
		app, err := P.CreateApp(body.GetField("name"))
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			w.Write([]byte("Could not create app"))
			return
		}
		bytes, err := json.Marshal(&app)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			w.Write([]byte("Could not parse json"))
			return
		}
		w.WriteHeader(200)
		w.Write(bytes)
	})
}

func handleGetApp(P *Permissionist) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		appID, err := P.GetApp(mux.Vars(r)["appID"])
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			w.Write([]byte("Could not get roles"))
			return
		}
		w.WriteHeader(200)
		w.Write([]byte(appID))
	})
}

func handleGetRoles(P *Permissionist) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		permissionNames, err := P.GetRoles(mux.Vars(r)["appID"])
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			w.Write([]byte("Could not get roles"))
			return
		}
		bytes, err := json.Marshal(permissionNames)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			w.Write([]byte("Could not get roles"))
			return
		}
		w.WriteHeader(200)
		w.Write(bytes)
	})
}

func handleCreateRole(P *Permissionist) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := NewBody(w, r)
		role, err := P.CreateRole(body.GetField("role_name"), mux.Vars(r)["appID"])
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			w.Write([]byte("Could not create role"))
			return
		}
		bytes, err := json.Marshal(&role)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			w.Write([]byte("Could not parse json"))
			return
		}
		w.WriteHeader(200)
		w.Write(bytes)
	})
}

func handleGetRole(P *Permissionist) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, err := P.GetRoleByID(mux.Vars(r)["roleID"])
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			w.Write([]byte("Could not get role"))
			return
		}
		bytes, err := json.Marshal(&role)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			w.Write([]byte("Could not parse json"))
			return
		}
		w.WriteHeader(200)
		w.Write(bytes)
	})
}

func handleGetPermissionsByRoleID(P *Permissionist) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roles, err := P.GetPermissionsByRoleID(mux.Vars(r)["roleID"])
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			w.Write([]byte("Could not get permissions"))
			return
		}
		bytes, err := json.Marshal(&roles)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			w.Write([]byte("Could not get permissions"))
			return
		}
		w.WriteHeader(200)
		w.Write(bytes)
	})
}

func handleCreatePermission(P *Permissionist) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := NewBody(w, r)
		permission, err := P.CreatePermission(body.GetField("name"), mux.Vars(r)["appID"])
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			w.Write([]byte("Could not create permission"))
			return
		}
		bytes, err := json.Marshal(&permission)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			w.Write([]byte("Could not create permission"))
			return
		}
		w.WriteHeader(200)
		w.Write(bytes)
	})
}

func handleGrantPermissionToRole(P *Permissionist) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := P.GrantPermissionToRole(mux.Vars(r)["roleID"], mux.Vars(r)["permissionID"])
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			w.Write([]byte("Could not grant permission"))
			return
		}
		w.WriteHeader(200)
	})
}

func main() {

	config := InitConfig()
	db := InitDb(config.GetString("database"))

	schemas, err := ioutil.ReadFile("schemas.sql")
	if err != nil {
		log.Fatal(err)
	}

	// Create schema
	_, err = db.Exec(string(schemas))
	if err != nil {
		log.Fatal(err)
	}

	P := Permissionist{
		DB: db,
	}

	router := mux.NewRouter()

	router.HandleFunc("/apps", handleCreateApp(&P)).Methods("POST")
	router.HandleFunc("/apps/{appID}", handleGetApp(&P)).Methods("GET")
	router.HandleFunc("/apps/{appID}/roles", handleGetRoles(&P)).Methods("GET")
	router.HandleFunc("/apps/{appID}/roles", handleCreateRole(&P)).Methods("POST")
	router.HandleFunc("/roles/{roleID}/permissions", handleGetPermissionsByRoleID(&P)).Methods("GET")
	router.HandleFunc("/roles/{roleID}/permissions/{permissionID}", handleGrantPermissionToRole(&P)).Methods("POST")
	router.HandleFunc("/permissions", handleCreatePermission(&P)).Methods("POST")
	http.ListenAndServe(":8000", router)
}
