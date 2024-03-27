package middleware

import (
	"context"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rulanugrh/order/internal/config"
	"github.com/rulanugrh/order/internal/util/constant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type jwtclaim struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	RoleID   uint   `json:"role_id"`
	Avatar   string `json:"avatar"`
	jwt.RegisteredClaims
}

type GRPCInterceptors struct {
	conf *config.App
}

func NewGRPCInterceptors(conf *config.App) *GRPCInterceptors {
	return &GRPCInterceptors{conf: conf}
}

func (g *GRPCInterceptors) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		log.Println("[*] Unary Interceptors: ", info.FullMethod)
		err = g.verify(ctx)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func (g *GRPCInterceptors) Stream() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		log.Println("[*] Stream Interceptors", info.FullMethod)
		err := g.verify(ss.Context())
		if err != nil {
			return err
		}

		return handler(srv, ss)
	}
}

func (g *GRPCInterceptors) verify(ctx context.Context) error {
	meta, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return constant.Unauthorized("you are not logged")
	}

	token := meta["Authorization"]
	if len(token) == 0 {
		return constant.NotFound("you token not found")
	}

	tkn, err := jwt.ParseWithClaims(token[0], &jwtclaim{}, func(t *jwt.Token) (data interface{}, err error) {
		data = []byte(g.conf.Server.Secret)
		return data, constant.BadRequest("Secret token not valid", err)
	})

	if err != nil {
		return constant.BadRequest("Invalid parse something error", err)
	}

	claim, valid := tkn.Claims.(*jwtclaim)
	if !valid {
		return constant.Forbidden("sorry your token invalid")
	}

	if claim.ExpiresAt.Unix() < time.Now().Unix() {
		return constant.Unauthorized("sorry you token expired")
	}

	return nil
}
