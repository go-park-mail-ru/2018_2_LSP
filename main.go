package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/satori/go.uuid"
)

type User struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Gamecount int    `json:"gamecount"`
	Score     int    `json:"score"`
}

var session = make(map[string]string)

var users = map[string]User{
	"moleque@mail.ru": User{
		Username:  "Moleque",
		Email:     "moleque@mail.ru",
		Password:  "123",
		Gamecount: 5,
		Score:     72,
	},
	"molecada@yandex.ru": User{
		Username:  "Molecada",
		Email:     "molecada@yandex.ru",
		Password:  "12345",
		Gamecount: 7,
		Score:     10,
	},
	"l@yandex.ru": User{
		Username:  "Llll",
		Email:     "l@yandex.ru",
		Password:  "12345",
		Gamecount: 7,
		Score:     7,
	},
	"m@yandex.ru": User{
		Username:  "Mlll",
		Email:     "m@yandex.ru",
		Password:  "12345",
		Gamecount: 7,
		Score:     7,
	},
	"n@yandex.ru": User{
		Username:  "Nlll",
		Email:     "n@yandex.ru",
		Password:  "12345",
		Gamecount: 7,
		Score:     7,
	},
	"c@yandex.ru": User{
		Username:  "Clll",
		Email:     "c@yandex.ru",
		Password:  "12345",
		Gamecount: 7,
		Score:     7,
	},
	"v@yandex.ru": User{
		Username:  "Vlll",
		Email:     "v@yandex.ru",
		Password:  "12345",
		Gamecount: 7,
		Score:     7,
	},
}

//========
func cors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
}

func leaderboards(w http.ResponseWriter, r *http.Request) {
	pageNumer := r.URL.Query().Get("page")
	if pageNumer == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no have correct params"))
		return
	} else {
		sliceUser := make([]User, 0, 10)
		for _, user := range users {
			sliceUser = append(sliceUser, user)
		}

		jsonUsers, err := json.Marshal(&sliceUser)
		if err != nil {
			log.Printf("cannot marshal:%s", err)
		}
		w.Write(jsonUsers)
	}
}

func profile(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err == http.ErrNoCookie {
		log.Println("coook")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	email, ok := session[cookie.Value]
	if !ok {
		log.Println("ses")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user := users[email]
	response, err := json.Marshal(user)
	if err != nil {
		log.Println("cannot marshal user:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func login(w http.ResponseWriter, r *http.Request) {

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

func signup(w http.ResponseWriter, r *http.Request) {
	request := &User{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(request)
	if err != nil {
		log.Printf("err %s", request)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, exists := users[request.Email]; request.Email == "" || exists == true {
		log.Printf("user %s is exist", request.Email)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users[request.Email] = User{
		Username:  request.Username,
		Email:     request.Email,
		Password:  request.Password,
		Gamecount: 0,
		Score:     0,
	}
	sessionId := uuid.NewV4().String()
	session[sessionId] = request.Email

	expiration := time.Now().Add(24 * time.Hour)
	cookie := http.Cookie{
		Name:    "session",
		Value:   sessionId,
		Expires: expiration,
	}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}

func main() {
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
				leaderboards(w, r)
			} else {
				profile(w, r)
			}
		case http.MethodPost:
			signup(w, r)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := fmt.Sprintf("cannot get %s", r.URL)
		log.Printf(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err))
	})

	log.Printf("Try to start http server http://127.0.0.1:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("listening error:%s", err)
	}
}
