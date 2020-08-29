package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/tebro/albion-mapper-backend/albion"

	"github.com/tebro/albion-mapper-backend/db"

	"github.com/gorilla/mux"
)

type apiPortal struct {
	Source  string `json:"source"`
	Target  string `json:"target"`
	Size    int    `json:"size"`
	Hours   int    `json:"hours"`
	Minutes int    `json:"minutes"`
}

var password = os.Getenv("AUTH_PASSWORD")

func isAuth(r *http.Request) bool {
	header := r.Header["X-Tebro-Auth"]
	for _, s := range header {
		if s == password {
			return true
		}
	}
	return false
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	_, err := db.Hello()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error")
		log.Printf("Health error: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func send401(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	fmt.Fprint(w, "Authenticate")
}

func send400AndLog(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, "Bad request")
	log.Printf("400 error: %v\n", err)
}

func send500AndLog(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, "Error")
	log.Printf("500 error: %v\n", err)
}

func getZonesHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuth(r) {
		send401(w)
		return
	}
	zones, err := albion.GetZones()
	if err != nil {
		send500AndLog(w, err)
		return
	}
	json, err := json.Marshal(zones)
	if err != nil {
		send500AndLog(w, err)
		return
	}
	fmt.Fprint(w, string(json))
}

func setZoneHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuth(r) {
		send401(w)
		return
	}
	var zone albion.Zone
	err := json.NewDecoder(r.Body).Decode(&zone)
	if err != nil {
		send400AndLog(w, err)
		return
	}
	if !albion.IsValidZone(zone) {
		send400AndLog(w, fmt.Errorf("Not valid zone: %v", zone))
		return
	}
	err = albion.SetZone(zone)
	if err != nil {
		send500AndLog(w, err)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprint(w, "ACCEPTED")
}

func getPortalsHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuth(r) {
		send401(w)
		return
	}
	portals, err := albion.GetPortals()
	if err != nil {
		send500AndLog(w, err)
		return
	}
	json, err := json.Marshal(portals)
	if err != nil {
		send500AndLog(w, err)
		return
	}
	fmt.Fprint(w, string(json))
}

func addPortalHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuth(r) {
		send401(w)
		return
	}

	var portal apiPortal
	err := json.NewDecoder(r.Body).Decode(&portal)
	if err != nil {
		send400AndLog(w, err)
		return
	}

	expires := time.Now()
	expires = expires.Add(time.Hour * time.Duration(portal.Hours))
	expires = expires.Add(time.Minute * time.Duration(portal.Minutes))

	dbPortal := albion.Portal{
		Source:  portal.Source,
		Target:  portal.Target,
		Size:    portal.Size,
		Expires: expires,
	}

	isValid, err := albion.IsValidPortal(dbPortal)
	if err != nil {
		send500AndLog(w, err)
		return
	}

	if !isValid {
		send400AndLog(w, fmt.Errorf("Invalid portal: %v", dbPortal))
		return
	}

	err = albion.AddPortal(dbPortal)
	if err != nil {
		send500AndLog(w, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprint(w, "ACCEPTED")
}

func setupRoutes(r *mux.Router) {
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/health", healthHandler)
	r.HandleFunc("/api/zone", setZoneHandler).Methods("POST")
	r.HandleFunc("/api/zone", getZonesHandler)
	r.HandleFunc("/api/portal", addPortalHandler).Methods("POST")
	r.HandleFunc("/api/portal", getPortalsHandler)
}

// StartServer starts the HTTP server
func StartServer() error {
	router := mux.NewRouter()
	setupRoutes(router)

	log.Println("Server starting on port 8080")
	err := http.ListenAndServe(":8080", router)
	return err
}
