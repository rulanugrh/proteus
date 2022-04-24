package controller

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	productController "github.com/ItsArul/TokoKu/controller/interfaces"
	"github.com/ItsArul/TokoKu/entity/domain"
	"github.com/ItsArul/TokoKu/entity/web"
	"github.com/ItsArul/TokoKu/services/interfaces"
	"github.com/gin-gonic/gin"
)

type productcontroler struct {
	product interfaces.ProducInterfaces
}

func StartProductController(p interfaces.ProducInterfaces) productController.ProductController {
	return &productcontroler{product: p}
}

func (p *productcontroler) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var product domain.Product
		context := context.Background()

		ctx.Bind(&product)
		productResponse, err := p.product.Create(context, product)
		if err != nil {
			response := web.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "product cannot create",
				Error:   err,
			}

			ctx.JSON(http.StatusInternalServerError, response)
		}

		response := web.SuccessResponse{
			Code:    200,
			Message: "success create  product",
			Name:    productResponse.Nama,
		}

		ctx.JSON(http.StatusOK, response)
	}
}

func (p *productcontroler) FindAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		context := context.Background()

		productResponse, err := p.product.FindAll(context)
		if err != nil {
			response := web.ErrorResponse{
				Code:    500,
				Message: "cannot findall product",
				Error:   err,
			}

			ctx.JSON(http.StatusInternalServerError, response)
		}

		response := web.SuccessFindAll{
			Code:    200,
			Message: "all product find",
			Data:    productResponse,
		}

		ctx.JSON(http.StatusOK, response)
	}
}

func (p *productcontroler) FindById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		GetId := ctx.Param("id")
		Id, _ := strconv.Atoi(GetId)

		context := context.Background()
		productresponse, err := p.product.FindById(context, uint(Id))
		if err != nil {
			response := web.ErrorResponse{
				Code:    500,
				Message: "cannot find product",
				Error:   err,
			}

			ctx.JSON(http.StatusInternalServerError, response)
		}

		response := web.SuccessResponse{
			Code:    200,
			Message: "find product",
			Name:    productresponse.Nama,
		}

		ctx.JSON(http.StatusOK, response)
	}
}

func (p *productcontroler) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var product domain.Product
		GetId := ctx.Param("id")
		Id, _ := strconv.Atoi(GetId)

		context := context.Background()

		updata, _ := ioutil.ReadAll(ctx.Request.Body)
		json.Unmarshal(updata, &product)

		productresponse, err := p.product.Update(context, uint(Id), product)
		if err != nil {
			response := web.ErrorResponse{
				Code:    500,
				Message: "cannot update product",
				Error:   err,
			}

			ctx.JSON(http.StatusInternalServerError, response)
		}

		response := web.SuccessUpdate{
			Code:    200,
			Message: "success update product",
			Data:    productresponse,
		}

		ctx.JSON(http.StatusOK, response)
	}
}

func (p *productcontroler) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		GetId := ctx.Param("id")
		Id, _ := strconv.Atoi(GetId)

		context := context.Background()
		err := p.product.Delete(context, uint(Id))
		if err != nil {
			response := web.ErrorResponse{
				Code:    500,
				Message: "cannot delete product",
				Error:   err,
			}

			ctx.JSON(http.StatusInternalServerError, response)
		}

		response := web.SuccessDelete{
			Code:    200,
			Message: "success delete product",
		}

		ctx.JSON(http.StatusOK, response)
	}
}
