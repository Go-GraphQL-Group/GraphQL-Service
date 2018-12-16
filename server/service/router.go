package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type RespData struct {
	Status bool   `json:"status"`
	Msg    string `json:"msg"`
	Token
}

// type ReqBody struct {
// 	Query         string                 `json:"query"`
// 	OperationName string                 `json:"operationName"`
// 	Variables     map[string]interface{} `json:"variables"`
// 	SW_Token      string                 `json:"sw_token"`
// }

func writeResp(status bool, msg string, token Token) []byte {
	RespData := RespData{}
	RespData.Status = status
	RespData.Msg = msg
	RespData.Token = token
	respose, err := json.Marshal(RespData)
	if err != nil {
		log.Fatalln(err)
	}
	return respose
}

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
		w.Write(writeResp(false, "Error logging in", Token{}))
		// fmt.Fprint(w, "Invalid credentials")
		return
	}

	token, err := CreateToken([]byte(SecretKey), Issuer, false)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error extracting the key")
		w.Write(writeResp(false, "Error extracting the key", Token{}))
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error marshal the token")
		w.Write(writeResp(false, "Error marshal the token", Token{}))
		log.Fatal(err)
	}
	tokens = append(tokens, token)
	w.Write(writeResp(true, "Succeed to login", token))
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	tokenStr := ""
	for k, v := range r.Header {
		if strings.ToLower(k) == TokenName {
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
	w.Write(writeResp(false, "Succeed to logout", Token{}))
}

func ApiHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{
	127.0.0.1:9090/api/login
	127.0.0.1:9090/api/query
	127.0.0.1:9090/api/logout
}
The format of query:
query{
	people(id:string){
		……
	}
	film(id:string){
		……
	}
	planet(id:string){
		……
	}
	specie(id:string){
		……
	}
	starship(id:string){
		……
	}
	vehicle(id:string){
		……
	}

	peoples(first: int, after: string){
		……
	}
	films(first: int, after: string){
		……
	}
	planets(first: int, after: string){
		……
	}
	species(first: int, after: string){
		……
	}
	starships(first: int, after: string){
		……
	}
	vehicles(first: int, after: string){
		……
	}

	peopleSearch(search: string, first: int, after: string){
		……
	}
	filmsSearch(search: string, first: int, after: string){
		……
	}
	planetsSearch(search: string, first: int, after: string){
		……
	}
	speciesSearch(search: string, first: int, after: string){
		……
	}
	starshipsSearch(search: string, first: int, after: string){
		……
	}
	vehiclesSearch(search: string, first: int, after: string){
		……
	}
}
	`))
}
