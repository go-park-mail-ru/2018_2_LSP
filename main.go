package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type User struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Gamecount int    `json:"gamecount"`
	Score     int    `json:"score"`
}

func main() {
	users := map[string]User{
		"moleque@mail.ru": User{
			Username:  "Moleque",
			Email:     "moleque@mail.ru",
			Password:  "123",
			Gamecount: 5,
			Score:     72,
		},
		"molecada@yandex.ru": {
			Username:  "Molecada",
			Email:     "molecada@yandex.ru",
			Password:  "12345",
			Gamecount: 7,
			Score:     10,
		},
		"l@yandex.ru": {
			Username:  "Llll",
			Email:     "l@yandex.ru",
			Password:  "12345",
			Gamecount: 7,
			Score:     7,
		},
	}

	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Println("bad request: wrong method")
		}

		log.Println(r.FormValue("email"), r.FormValue("password"))
	})

	// http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
	// 	log.Println("request", r.URL)
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.Write(response)
	// })

	// http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
	// 	log.Println("request", r.URL)
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.Write(response)
	// })

	http.HandleFunc("/leaderboard", func(w http.ResponseWriter, r *http.Request) {
		pageNumer := r.URL.Query().Get("page")
		if pageNumer == "" {
			return
		} else {
			jsonUsers, err := json.Marshal(&users)
			if err != nil {
				log.Printf("cannot marshal:%s", err)
			}
			w.Header().Set("Content-Type", "application/json")
			log.Println(r.Header.Get("Origin"))
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			if r.Method == http.MethodOptions {
				return
			}
			w.Write(jsonUsers)
		}
	})

	log.Printf("Try to start http server http://127.0.0.1:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("listening error:%s", err)
	}
}
