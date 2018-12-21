package service

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func jsonDecode(r io.Reader, val interface{}) error {
	dec := json.NewDecoder(r)
	dec.UseNumber()
	return dec.Decode(val)
}

func TokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, sw_token,sign")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		if r.RequestURI[1:] != "api/login" {
			/*
				// token位于Authorization中，用此方法
				token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
					return []byte(SecretKey), nil
				})
			*/
			tokenStr := ""
			for k, v := range r.Header {
				if strings.ToLower(k) == TokenName {
					tokenStr = v[0]
					break
				}
			}
			// var reqbody ReqBody
			// if tokenStr == "" {
			// 	if err := jsonDecode(r.Body, &reqbody); err == nil {
			// 		tokenStr = reqbody.SW_TOKEN
			// 	}
			// }
			validToken := false
			for _, token := range tokens {
				if token.SW_TOKEN == tokenStr {
					validToken = true
				}
			}

			if validToken {
				ctx := context.WithValue(r.Context(), TokenName, tokenStr)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write(writeResp(false, "Unauthorized access to this resource", Token{}))
				//fmt.Fprint(w, "Unauthorized access to this resource")
			}
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
