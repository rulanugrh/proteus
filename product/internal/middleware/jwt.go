package middleware

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rulanugrh/tokoku/product/internal/config"
	"github.com/rulanugrh/tokoku/product/internal/entity/web"
)

type jwtclaim struct {
	ID       uint   `json:"id"`
	RoleID   uint   `json:"role_id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	jwt.RegisteredClaims
}

func CheckToken(tkn string) (*jwtclaim, error) {
	conf := config.GetConfig()
	token, _ := jwt.ParseWithClaims(tkn, &jwtclaim{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(conf.Server.Secret), web.Forbidden("this page strict")
	})

	claim, err := token.Claims.(*jwtclaim)
	if !err {
		return nil, web.Unauthorized("token invalid or missing")
	}

	return claim, nil
}

func ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token == "" {
			response, err := json.Marshal(web.Unauthorized("sorry token missing"))
			if err != nil {
				w.WriteHeader(500)
				return
			}

			w.WriteHeader(401)
			w.Write(response)
			return
		}

		claim, err := CheckToken(token)
		if err != nil {
			response, err := json.Marshal(web.Unauthorized(err.Error()))
			if err != nil {
				w.WriteHeader(500)
				return
			}

			w.WriteHeader(401)
			w.Write(response)
			return
		}

		if claim.ExpiresAt.Unix() < time.Now().Unix() {
			response, err := json.Marshal(web.Unauthorized("token expire"))
			if err != nil {
				w.WriteHeader(500)
				return
			}

			w.WriteHeader(401)
			w.Write(response)
			return
		}

		next.ServeHTTP(w, r)
	})
}
