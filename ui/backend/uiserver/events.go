package uiserver

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func (srv *uiServer) getEvents(w http.ResponseWriter, r *http.Request) {
	log.Printf("configured config path %s", srv.cfgPath)

	hookerUrl := os.Getenv("HOOKER_UI_UPDATE_URL")
	if len(hookerUrl) <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("No Hooker URL configured, set HOOKER_UI_UPDATE_URL to the Hooker URL")
		return
	}

	resp, err := http.Get(hookerUrl + "/events")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Unable to reach Hooker at URL: " + hookerUrl + "/events" + " err: " + err.Error())
		return
	}

	currentEvents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Failed to read events: " + err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(currentEvents)
}
