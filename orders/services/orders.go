package services

import (
	"fmt"
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
	status := query.Get("status")
	search := query.Get("search")
	if from_date != "" {
		from_date = from_date + " 00:00:00"
		m = append(m, clause.Gte{Column: "created_at", Value: from_date})
	}

	if search != "" {
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
