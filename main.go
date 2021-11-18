package main

import (
	"fmt"
	"log"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("started-service")

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(mySigningKey), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	r := mux.NewRouter()

	// 前后端分开部署时，CORS过程中， 前端先发送option request 给后端，如果后端支持浏览，再给后端发送post request
	r.Handle("/upload", jwtMiddleware.Handler(http.HandlerFunc(uploadHandler))).Methods("POST", "OPTIONS")
	r.Handle("/search", jwtMiddleware.Handler(http.HandlerFunc(searchHandler))).Methods("GET", "OPTIONS")
	r.Handle("/signup", http.HandlerFunc(signupHandler)).Methods("POST", "OPTIONS")
	r.Handle("/signin", http.HandlerFunc(signinHandler)).Methods("POST", "OPTIONS")

	log.Fatal(http.ListenAndServe(":8080", r))

}
