package apis

import (
	"main/database"
	"main/middleware"
	"main/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func setupTransactionAPI(router *gin.Engine) {
	transactionAPI := router.Group("/api")
	transactionAPI.Use(middleware.JwtVerify)
	{
		transactionAPI.GET("/transaction", getTransaction)
		transactionAPI.POST("/transaction", createTransaction)
	}
}

type TransactionResult struct {
	ID            uint
	Total         float64
	Paid          float64
	Change        float64
	PaymentType   string
	PaymentDetail string
	OrderList     string
	StaffID       string
	CreatedAt     time.Time
}

func getTransaction(c *gin.Context) {
	var result []TransactionResult
	statement := `
		SELECT transactions.id, total, paid, change, payment_type, payment_detail, 
			order_list, users.username as staff, transactions.created_at
		FROM transactions
		join users on transactions.staff_id = users.id
	`
	database.GetDB().Raw(statement, nil).Scan(&result)
	c.JSON(http.StatusOK, result)
}

func createTransaction(c *gin.Context) {
	var transaction models.Transaction
	if err := c.ShouldBind(&transaction); err == nil {
		transaction.StaffID = c.GetString("jwt_staff_id")
		transaction.CreatedAt = time.Now()
		database.GetDB().Create(&transaction)
		c.JSON(http.StatusOK, map[string]interface{}{"status": "success", "data": transaction})
	} else {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"status": "Bad request"})
	}

}
