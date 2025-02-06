package handlers

import (
	"haraka-sana/config"
	"haraka-sana/orders/models"
	"haraka-sana/orders/objects"
	staffModel "haraka-sana/staff/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

// Generate a unique tracking code based on time
func CreateTrackingCode() string {
	// Get the current timestamp in milliseconds
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	// Generate a random 4-byte value
	bytes := make([]byte, 4)
	_, err := rand.Read(bytes)
	if err != nil {
		return "" // Handle error properly in production
	}

	// Convert to hexadecimal string
	randomPart := hex.EncodeToString(bytes)

	// Combine timestamp and random part
	trackingCode := fmt.Sprintf("%x-%s", timestamp, randomPart)

	return trackingCode
}

func CreateShippingBatch(c *gin.Context) {
	contextStaff, _ := c.Get("staff")
	fmt.Print(contextStaff)
	staff := contextStaff.(staffModel.Staff)
	var batchInfo objects.CreateBatch

	err := c.ShouldBindJSON(&batchInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error persing information, check fiedls",
		})
		return
	}

	batch := models.Batch{
		OpenTracking:  batchInfo.OpenTracking,
		StartLocation: batchInfo.StartLocation,
		StopCountry:   batchInfo.StartCountry,
		VehicleNumber: batchInfo.VehicleNumber,
		VehicleType:   batchInfo.VehicleType,
		Status:        batchInfo.Status,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		StaffId:       staff.Id,
	}
	config.DB.Create(&batch)
	c.JSON(http.StatusOK, gin.H{
		"success": "batch created successfully, list orders under it",
	})
}

func GetBatches(c *gin.Context) {
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

	var batches []models.Batch
	dbQuery := config.DB

	if len(m) != 0 {
		dbQuery = dbQuery.Where(clause.Where{Exprs: m})
	}

	dbQuery.
		Offset(offset).
		Order(orderBy).
		Find(&batches)

	countQuery := config.DB
	if len(m) != 0 {
		countQuery = countQuery.Where(clause.Where{Exprs: m})
	}
	countQuery.
		Model(&models.Batch{}).
		Select("batches.id").
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
		"batches":   batches,
		"page_info": pageInfo,
	})
}

func AddProductToBatch(c *gin.Context) {
	contextStaff, _ := c.Get("staff")
	fmt.Print(contextStaff)
	staff := contextStaff.(staffModel.Staff)

	var orderBatch objects.AddOrderBatch

	err := c.ShouldBindJSON(&orderBatch)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error persing information, check fiedls",
		})
		return
	}
	var batchInfo models.Batch
	config.DB.
		Where(&models.Batch{
			Id: orderBatch.BatchId,
		}).First(&batchInfo)

	if batchInfo.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "batch not found",
		})
		return
	}

	if batchInfo.StaffId != staff.Id {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "only staff who created batch can add items",
		})
		return
	}

	batchProd := models.BatchOrder{
		BatchId:   orderBatch.BatchId,
		OrderId:   orderBatch.OrderId,
		CreatedAt: time.Now(),
	}
	config.DB.Create(&batchProd)
	c.JSON(http.StatusOK, gin.H{
		"success": "Order added to batch",
	})
}
