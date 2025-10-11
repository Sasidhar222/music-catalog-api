package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Sasidhar222/music-api/internal/data"
	"github.com/Sasidhar222/music-api/models"
)

// application struct to hold our dependencies
type Handlers struct {
	Artists *data.ArtistModel
}

// The handler function now uses data.Artist
// func getArtistsHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(artists)
// }

func (h *Handlers) ArtistsHandler(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case "POST":
		var newArtist models.Artist
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&newArtist)
		if err != nil {
			http.Error(w, "Bad Request:"+err.Error(), http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(newArtist.Name) == "" {
			http.Error(w, "Bad Request: name field cannot be empty", http.StatusBadRequest)
			return
		}
		_, err = h.Artists.AddArtist(&newArtist)

		if err != nil {
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newArtist)
	case "DELETE":
		split := strings.Split(r.RequestURI, "/")
		id := split[len(split)-1]

		_, err := h.Artists.DeleteArtist(id)

		if err != nil {
			if errors.Is(err, data.ErrNotFound) {
				http.Error(w, "NOT FOUND", http.StatusNotFound)
				return
			} else {
				http.Error(w, "Internal Sever Error: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusNoContent)
		return
	case "PUT":
		var artist models.Artist
		decoder := json.NewDecoder((r.Body))
		err := decoder.Decode(&artist)
		if err != nil {
			http.Error(w, "Bad Requst:"+err.Error(), http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(artist.Name) == "" {
			http.Error(w, "Bad Request: name field cannot be empty", http.StatusBadRequest)
			return
		}
		split := strings.Split(r.RequestURI, "/")
		id := split[len(split)-1]
		_, err = h.Artists.UpdateArtist(id, &artist)

		if err != nil {
			if errors.Is(err, data.ErrNotFound) {
				http.Error(w, "NOT FOUND", http.StatusNotFound)
				return
			} else {
				http.Error(w, "Internal Sever Error: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(artist)
	case "GET":
		artists, err := h.Artists.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(artists)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

}

func (h *Handlers) StatusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server is running")
}
