package test

import (
	"encoding/json"
	"fmt"
	"io"
	"store-api/internal/dto"
	"store-api/internal/entity"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestCustomer_Products(t *testing.T) {
	var product entity.Product
	TestCfg.DB.Preload("Category").Preload("Store").First(&product)

	testItems := map[string]TestItem{
		"get_all_success_query": {
			"query_name":            product.Name,
			"query_store_id":        product.StoreID,
			"query_category_id":     product.CategoryID,
			"response_code":         fiber.StatusOK,
			"response_status":       "SUCCESS",
			"response_current_page": 1,
			"response_total_page":   1,
			"response_page_size":    10,
			"count_item":            1,
		},
		"get_all_success_paginate": {
			"query_page":            "2",
			"query_page_size":       "5",
			"response_code":         fiber.StatusOK,
			"response_status":       "SUCCESS",
			"response_current_page": 2,
			"response_page_size":    5,
			"count_item":            5,
		},
		"get_all_success_record_empty": {
			"query_name":            "thisIsSoRandomName",
			"response_code":         fiber.StatusOK,
			"response_status":       "SUCCESS",
			"response_current_page": 1,
			"response_total_page":   0,
			"response_page_size":    10,
			"count_item":            0,
		},
		"get_all_success_paginate_empty": {
			"query_page":            "1000",
			"query_page_size":       "1000",
			"response_code":         fiber.StatusOK,
			"response_status":       "SUCCESS",
			"response_current_page": 1000,
			"response_page_size":    1000,
			"count_item":            0,
		},
	}

	for testName, testItem := range testItems {
		t.Run(testName, func(t *testing.T) {
			queryUrl := "?"
			for k, v := range testItem {
				if strings.Contains(k, "query_") {
					queryUrl += fmt.Sprintf("%s=%s&",
						strings.Replace(k, "query_", "", 1), v)
				}
			}
			queryUrl = queryUrl[:len(queryUrl)-1] // remove trailing &
			request := NewRequestWithToken(fiber.MethodGet, CustomerProductURL+queryUrl, "", ExistingCustomer["token"])

			response, err := TestCfg.App.Test(request)
			require.Nil(t, err)
			require.Equal(t, testItem["response_code"], response.StatusCode)

			responseBodyByte, err := io.ReadAll(response.Body)
			require.Nil(t, err)

			responseBody := new(dto.Response[[]dto.ProductDTO])
			err = json.Unmarshal(responseBodyByte, responseBody)
			require.Nil(t, err)

			require.Equal(t, testItem["response_status"], responseBody.Status)
			require.Equal(t, testItem["response_current_page"], responseBody.Pagination.CurrentPage)
			if _, ok := testItem["response_total_page"]; ok {
				require.Equal(t, testItem["response_total_page"], responseBody.Pagination.TotalPage)
			}
			require.Equal(t, testItem["response_page_size"], responseBody.Pagination.PageSize)
			require.Equal(t, testItem["count_item"], len(responseBody.Data))

		})
	}
}
