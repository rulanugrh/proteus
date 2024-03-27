package middleware

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rulanugrh/order/internal/config"
	"github.com/rulanugrh/order/internal/util/constant"
	"google.golang.org/grpc/metadata"
)

func GetID() (*uint, error) {
	meta, ok := metadata.FromIncomingContext(context.Background())

	if !ok {
		return nil, constant.Unauthorized("you are not logged")
	}

	token := meta["Authorization"]
	if len(token) == 0 {
		return nil, constant.NotFound("you token not found")
	}

	claim, err := getToken(token[0])
	if err != nil {
		return nil, constant.Unauthorized(err.Error())
	}

	return &claim.ID, nil
}

func getToken(token string) (*jwtclaim, error) {
	conf := config.GetConfig()
	tkn, err := jwt.ParseWithClaims(token, &jwtclaim{}, func(t *jwt.Token) (data interface{}, err error) {
		data = []byte(conf.Server.Secret)
		return data, constant.BadRequest("Secret token not valid", err)
	})

	if err != nil {
		return nil, constant.BadRequest("Invalid parse something error", err)
	}

	claim, valid := tkn.Claims.(*jwtclaim)
	if !valid {
		return nil, constant.Forbidden("sorry your token invalid")
	}

	return claim, nil
}