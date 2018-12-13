package service

import (
	"context"
	"net/http"
	"strings"
)

func TokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI[1:] != "login" {
			/*
				// token位于Authorization中，用此方法
				token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
					return []byte(SecretKey), nil
				})
			*/
			tokenStr := ""
			for k, v := range r.Header {
				if strings.ToUpper(k) == TokenName {
					tokenStr = v[0]
					break
				}
			}
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
				w.Write([]byte("Unauthorized access to this resource"))
				//fmt.Fprint(w, "Unauthorized access to this resource")
			}
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
