package customer

import (
	"store-api/internal/dto"
	"store-api/internal/service"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ProductController struct {
	Logger         *logrus.Logger
	ProductService service.ProductService
}

func NewProductController(logger *logrus.Logger, services *service.Services) *ProductController {
	return &ProductController{
		Logger:         logger,
		ProductService: services.ProductService,
	}
}

func (ct *ProductController) GetAll(c *fiber.Ctx) error {
	request := new(dto.FindAndSearchProductRequest)
	request.Page, _ = strconv.Atoi(c.Query("page"))
	if request.Page < 1 {
		request.Page = 1
	}
	request.PageSize, _ = strconv.Atoi(c.Query("page_size"))
	if request.PageSize < 1 {
		request.PageSize = 10
	}
	request.QueryName = c.Query("name")
	request.QueryPriceMin, _ = strconv.ParseFloat(strings.TrimSpace(c.Query("min_price")), 64)
	request.QueryPriceMax, _ = strconv.ParseFloat(strings.TrimSpace(c.Query("max_price")), 64)
	request.QueryCategoryID = c.Query("category_id")
	request.QueryStoreID = c.Query("store_id")
	if strings.ToLower(c.Query("sort_max_price")) == "true" {
		request.SortByPriceHighest = true
	} else if strings.ToLower(c.Query("sort_min_price")) == "true" {
		request.SortByPriceLowest = true
	}

	res, pagination, err := ct.ProductService.FindAll(c.UserContext(), request)
	if err != nil {
		return err
	}

	response := dto.Response[[]dto.ProductItemDTO]{
		Status: "SUCCESS",
		Messages: map[string]string{
			"_success": "success getting categories",
		},
		Data:       res,
		Pagination: pagination,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
