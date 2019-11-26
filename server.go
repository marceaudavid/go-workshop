package main

import (
		"fmt"
		"math/rand"
		"net/http"
		"encoding/json"
		"github.com/gorilla/sessions"
		"github.com/gorilla/securecookie"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var key = securecookie.GenerateRandomKey(32)
var store = sessions.NewCookieStore([]byte(key))


func main() {
    http.HandleFunc("/random", func (w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "%f", rand.Float64())
		})
		
		http.HandleFunc("/login", func (w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {
				var user User
				dec := json.NewDecoder(r.Body)
				dec.DisallowUnknownFields()
				err := dec.Decode(&user)
				if err != nil {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					fmt.Fprintf(w, "%s", err.Error())
				} else {
					session, _ := store.Get(r, "auth")
					session.Values["authenticated"] = true
					session.Save(r, w)
					fmt.Fprintf(w, "username: %s \npassword: %s", user.Username, user.Password)
				}
			} else {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
		})

		http.HandleFunc("/status", func (w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, "auth")
			auth, ok := session.Values["authenticated"].(bool)
			if !auth || !ok {
				fmt.Fprintf(w, "Visiteur")
    	} else {
				fmt.Fprintf(w, "Authentifié")
			}
		})

		http.HandleFunc("/logout", func (w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, "auth")
			session.Values["authenticated"] = false
			session.Save(r, w)
			fmt.Fprintf(w, "Vous avez été deconnecté")
		})

		http.ListenAndServe(":3000", nil)
}