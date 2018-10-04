package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type User struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Gamecount int    `json:"gamecount"`
	Score     int    `json:"score"`
}

func cors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
}

func leaderboards(w http.ResponseWriter, r *http.Request, users *[]User) {
	pageNumer := r.URL.Query().Get("page")
	if pageNumer == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no have correct params"))
		return
	} else {
		jsonUsers, err := json.Marshal(&users)
		if err != nil {
			log.Printf("cannot marshal:%s", err)
		}
		w.Write(jsonUsers)
	}
}

func profile(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session")
	loggedIn := (err != http.ErrNoCookie)

	if loggedIn {
		w.Write([]byte("Goooood!" + session.Value))
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	expiration := time.Now().Add(24 * time.Hour)
	cookie := http.Cookie{
		Name:    "session",
		Value:   "moleque",
		Expires: expiration,
	}
	http.SetCookie(w, &cookie)
	w.Write([]byte("Good!"))
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session")
	if err == http.ErrNoCookie {
		return
	} else {
		session.Expires = time.Now().Add(-24 * time.Hour)
		http.SetCookie(w, session)
	}
}

func main() {
	users := []User{
		{
			Username:  "Moleque",
			Email:     "moleque@mail.ru",
			Password:  "123",
			Gamecount: 5,
			Score:     72,
		},
		{
			Username:  "Molecada",
			Email:     "molecada@yandex.ru",
			Password:  "12345",
			Gamecount: 7,
			Score:     10,
		},
		{
			Username:  "Llll",
			Email:     "l@yandex.ru",
			Password:  "12345",
			Gamecount: 7,
			Score:     7,
		},
		{
			Username:  "Mlll",
			Email:     "m@yandex.ru",
			Password:  "12345",
			Gamecount: 7,
			Score:     7,
		},
		{
			Username:  "Nlll",
			Email:     "n@yandex.ru",
			Password:  "12345",
			Gamecount: 7,
			Score:     7,
		},
		{
			Username:  "Clll",
			Email:     "c@yandex.ru",
			Password:  "12345",
			Gamecount: 7,
			Score:     7,
		},
		{
			Username:  "Vlll",
			Email:     "v@yandex.ru",
			Password:  "12345",
			Gamecount: 7,
			Score:     7,
		},
	}

	http.HandleFunc("/session", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r)
		switch r.Method {
		// case http.MethodGet:
		// 	(w, r)
		case http.MethodPost:
			login(w, r)
		case http.MethodDelete:
			logout(w, r)
			//default:
		}
	})

	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r)
		switch r.Method {
		case http.MethodGet:
			page := r.URL.Query().Get("page")
			if page != "" {
				leaderboards(w, r, &users)
			}
		}
	})

	log.Printf("Try to start http server http://127.0.0.1:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("listening error:%s", err)
	}
}
