package handlers

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	services "go.simple/structure/services"
	"gorm.io/gorm"
)

const (
	deposit  = "deposit"
	withdraw = "withdraw"
)

type PaymentHandler struct {
	paymentService services.PaymentService
}

func NewPaymentHandler(db *gorm.DB) *PaymentHandler {
	return &PaymentHandler{
		paymentService: services.NewPaymentService(db),
	}
}

func (ph *PaymentHandler) CreateTransaction(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if userID <= 0 {
		c.JSON(400, "Invalid user_id")
		return
	}

	var reqBody services.CreateTransactionRequest
	err := c.BindJSON(&reqBody)
	if err != nil {
		c.JSON(400, "Parsing request error")
		return
	}
	if reqBody.TransactionType != withdraw && reqBody.TransactionType != deposit {
		c.JSON(400, "Uknown transaction type")
		return
	}
	reqBody.UserID = userID

	result, err := ph.paymentService.CreateTransaction(reqBody)
	if err != nil {
		log.Println("Create transaction error: ", err.Error())
		c.JSON(500, "Internal server error")
		return
	}

	c.JSON(200, result)
}

func (ph *PaymentHandler) GetUserTransaction(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if userID <= 0 {
		c.JSON(400, "Invalid user_id")
		return
	}

	accountID, _ := strconv.ParseInt(c.Query("account_id"), 10, 64)
	if accountID < 0 {
		c.JSON(400, "Invalid account_id")
		return
	}

	result, err := ph.paymentService.GetUserTransactions(userID, accountID)
	if err != nil {
		log.Println("Get transaction error: ", err.Error())
		c.JSON(500, "Internal server error")
		return
	}

	c.JSON(200, result)
}

func (ph *PaymentHandler) UpdateTransaction(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if userID <= 0 {
		c.JSON(400, "Invalid user_id")
		return
	}

	var reqBody services.UpdateTransactionRequest
	err := c.BindJSON(&reqBody)
	if err != nil {
		c.JSON(400, "Parsing request error")
		return
	}
	if reqBody.TransactionType != withdraw && reqBody.TransactionType != deposit {
		c.JSON(400, "Uknown transaction type")
		return
	}
	reqBody.UserID = userID

	result, err := ph.paymentService.UpdateUserTransaction(reqBody)
	if err == services.ZeroAffectedErr {
		c.JSON(400, "Zero affected row")
		return
	} else if err != nil {
		log.Println("Update transaction error: ", err.Error())
		c.JSON(500, "Internal server error")
		return
	}

	c.JSON(200, result)
}

func (ph *PaymentHandler) DeleteTransaction(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if userID <= 0 {
		c.JSON(400, "Invalid user_id")
		return
	}
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	if id <= 0 {
		c.JSON(400, "Invalid transaction id")
		return
	}

	err := ph.paymentService.DeleteUserTransaction(id, userID)
	if err == services.ZeroAffectedErr {
		c.JSON(400, "Zero affected row")
		return
	} else if err != nil {
		log.Println("Update transaction error: ", err.Error())
		c.JSON(500, "Internal server error")
		return
	}

	c.JSON(200, gin.H{"message": "success"})
}
