package services

import (
	"fmt"
	"haraka-sana/config"
	"haraka-sana/orders/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

type OrderFilter struct {
	Offest   int
	PageSize int
	Page     int
	OrderBy  string
	Filter   []clause.Expression
}

func OrderFilters(c *gin.Context) OrderFilter {
	query := c.Request.URL.Query()

	m := make([]clause.Expression, 0, 4)
	m_or := make([]clause.Expression, 0, 4)

	from_date := query.Get("from_date")
	to_date := query.Get("to_date")
	delivered := query.Get("delivered")
	status := query.Get("status")
	search := query.Get("search")
	if from_date != "" {
		from_date = from_date + " 00:00:00"
		m = append(m, clause.Gte{Column: "created_at", Value: from_date})
	}

	if delivered == "true" {
		m = append(m, clause.Eq{Column: "delivered", Value: true})
	}
	if delivered == "false" {
		m = append(m, clause.Eq{Column: "delivered", Value: false})
	}

	if search != "" {
		m_or = append(m_or, clause.Like{Column: "product_name", Value: "%" + search + "%"})

		m = append(m, clause.Or(m_or...))
	}

	if status != "" {
		m = append(m, clause.Eq{Column: "status", Value: status})
	}

	if to_date != "" {
		to_date = to_date + " 23:59:59"
		m = append(m, clause.Lte{Column: "created_at", Value: to_date})
	}

	page, err := strconv.Atoi(query.Get("page"))
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(query.Get("page_size"))
	if err != nil {
		fmt.Println(err)
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	orderBy := "id DESC"
	order_by := query.Get("order_by")
	if order_by != "" {
		// with - means descending
		if strings.HasPrefix(order_by, "-") {
			orderBy = order_by[1:] + " DESC"
		} else {
			orderBy = order_by + " ASC"
		}
	}
	return OrderFilter{
		Offest:   offset,
		OrderBy:  orderBy,
		Filter:   m,
		PageSize: pageSize,
		Page:     page,
	}
}

func (filters OrderFilter) GetOrders() map[string]any {
	var orders []models.Order

	var totalCount int64

	if len(filters.Filter) == 0 {
		config.DB.
			Offset(filters.Offest).
			Limit(filters.PageSize).
			Preload("Customer").
			Preload("Seller").
			Preload("OrganizationApplication").
			Preload("Product").
			Order(filters.OrderBy).
			Find(&orders)

		config.DB.
			Model(&models.Order{}).
			Select("orders.id").
			Count(&totalCount)
	} else {
		config.DB.
			Where(clause.Where{Exprs: filters.Filter}).
			Preload("Customer").
			Preload("Seller").
			Preload("OrganizationApplication").
			Preload("Product").
			Order(filters.OrderBy).
			Find(&orders)

		config.DB.
			Where(clause.Where{Exprs: filters.Filter}).
			Model(&models.Order{}).
			Select("orders.id").
			Count(&totalCount)
	}
	totalPages := 0

	if int(totalCount)%filters.PageSize == 0 {
		totalPages = int(totalCount) / filters.PageSize
	} else {
		totalPages = (int(totalCount) / filters.PageSize) + 1
	}

	pageInfo := gin.H{
		"page":        filters.Page,
		"page_size":   filters.PageSize,
		"total_count": totalCount,
		"total_pages": totalPages,
	}
	return gin.H{
		"orders":    orders,
		"page_info": pageInfo,
	}

}
