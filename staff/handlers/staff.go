package handlers

import (
	"fmt"
	"haraka-sana/config"
	"haraka-sana/staff/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func GetStaff(c *gin.Context) {

	query := c.Request.URL.Query()

	m := make([]clause.Expression, 0, 4)
	m_or := make([]clause.Expression, 0, 4)

	from_date := query.Get("from_date")
	to_date := query.Get("to_date")
	active := query.Get("active")
	search := query.Get("search")
	if from_date != "" {
		from_date = from_date + " 00:00:00"
		m = append(m, clause.Gte{Column: "created_at", Value: from_date})
	}

	if active == "true" {
		m = append(m, clause.Eq{Column: "active", Value: true})
	}
	if active == "false" {
		m = append(m, clause.Eq{Column: "active", Value: false})
	}

	if search != "" {
		m = append(m, clause.Or(m_or...))
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
	var totalCount int64

	var staff []models.Staff
	dbQuery := config.DB

	if len(m) != 0 {
		dbQuery = dbQuery.Where(clause.Where{Exprs: m})
	}

	dbQuery.
		Offset(offset).
		Limit(pageSize).
		Order(orderBy).
		Find(&staff)

	countQuery := config.DB
	if len(m) != 0 {
		countQuery = countQuery.Where(clause.Where{Exprs: m})
	}
	countQuery.
		Model(&models.Staff{}).
		Select("staffs.id").
		Count(&totalCount)
	totalPages := 0

	if int(totalCount)%pageSize == 0 {
		totalPages = int(totalCount) / pageSize
	} else {
		totalPages = (int(totalCount) / pageSize) + 1
	}

	pageInfo := gin.H{
		"page":        page,
		"page_size":   pageSize,
		"total_count": totalCount,
		"total_pages": totalPages,
	}
	c.JSON(http.StatusOK, gin.H{
		"staff":     staff,
		"page_info": pageInfo,
	})
}
