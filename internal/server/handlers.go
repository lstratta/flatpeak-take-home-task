package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (s *serveMux) testHandler() http.Handler {
return http.HandlerFunc( func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("test handler reached")
		d := struct{
			Name string `json:"name"`
			Date time.Time `json:"date"`
		}{
			Name: "Luke",
			Date: time.Now(),
		}

		b, err := json.Marshal(d)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
    
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	})
}
