package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

const (
	secretKey = "This is an admin authorization"
)

type stUserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type stToken struct {
	Token string `json:"token"`
}

// LoginHandler login handler for http post
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user stUserCredentials

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		log.Println("Error in login: " + err.Error())
		return
	}

	if user.Username != "admin" || user.Password != "admin" {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Invalid credentials")
		log.Println("Error logging in")
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error extracting the key")
		return
	}

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while signing the token")
		return
	}

	response := stToken{tokenString}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// MiddleWare middleware for resource handler
func MiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		validateTokenMiddleware(w, r, next.ServeHTTP)
	})
}

func validateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

	if err == nil {
		if token.Valid {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Token is not valid")
			log.Print("Token is not valid")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		log.Print("Unauthorized access to this resource")
	}

}
