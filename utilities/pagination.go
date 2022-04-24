package utilities

import (
	"strconv"

	"github.com/ItsArul/TokoKu/entity/domain"
	"github.com/gin-gonic/gin"
)

func PaginationRequest(ctx *gin.Context) domain.PaginationProduct {
	limit := 5
	page := 1
	sort := "nama asc"

	query := ctx.Request.URL.Query()
	for key, value := range query {
		queryValue := value[len(value)-1]
		switch key {
		case "limit":
			limit, _ = strconv.Atoi(queryValue)
		case "page":
			page, _ = strconv.Atoi(queryValue)

		case "sort":
			sort = queryValue

		}
	}

	return domain.PaginationProduct{
		Limit: limit,
		Page:  page,
		Sort:  sort,
	}
}
