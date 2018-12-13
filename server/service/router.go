package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	fmt.Println(r.Form.Get("username"))
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Error in request")
		return
	}

	if strings.ToLower(r.Form.Get("username")) != "admin" || r.Form.Get("password") != "password" {
		w.WriteHeader(http.StatusForbidden)
		fmt.Println("Error logging in")
		fmt.Fprint(w, "Invalid credentials")
		return
	}

	token, err := CreateToken([]byte(SecretKey), Issuer, false)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error extracting the key")
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	tokenBytes, err := json.Marshal(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error marshal the token")
		log.Fatal(err)
	}
	tokens = append(tokens, token)
	w.Write(tokenBytes)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	tokenStr := ""
	for k, v := range r.Header {
		if strings.ToUpper(k) == TokenName {
			tokenStr = v[0]
			break
		}
	}
	for i, token := range tokens {
		if token.SW_TOKEN == tokenStr {
			tokens = append(tokens[:i], tokens[i+1:]...)
			break
		}
	}
	w.Write([]byte("logout"))
}
