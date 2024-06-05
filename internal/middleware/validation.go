package middleware

import (
	"assignment/internal/constants"
	"assignment/internal/models"
	"assignment/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

// This function gets the unique transactionID
func GetTransactionID(c *gin.Context) string {

	transactionID := c.GetHeader(constants.TransactionID)
	_, err := uuid.Parse(transactionID)
	if err != nil {
		transactionID = uuid.New().String()
		c.Request.Header.Set(constants.TransactionID, transactionID)
	}
	return transactionID
}

func ValidateCreateEmployeeRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// validate the body params
		var employee models.Employee
		err := ctx.ShouldBindBodyWith(&employee, binding.JSON)
		if err != nil {
			utils.RespondWithError(ctx, http.StatusBadRequest, constants.InvalidBody)
			return
		}

		if employee.Name == "" {
			utils.RespondWithError(ctx, http.StatusBadRequest, "employee name is missing")
		}

		if employee.Position == "" {
			utils.RespondWithError(ctx, http.StatusBadRequest, "employee position is missing")
		}

		if employee.Salary == nil {
			utils.RespondWithError(ctx, http.StatusBadRequest, "employee salary is missing")
		}

		ctx.Next()
	}
}

func ValidateUpdateEmployeeRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// validate the body params
		var employee models.Employee
		err := ctx.ShouldBindBodyWith(&employee, binding.JSON)
		if err != nil {
			utils.RespondWithError(ctx, http.StatusBadRequest, constants.InvalidBody)
			return
		}

		if employee.ID == "" {
			utils.RespondWithError(ctx, http.StatusBadRequest, "employee Id is missing")
		}

		ctx.Next()
	}
}

func ValidateEmployeeID() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		employeeID := ctx.Param("id")
		// Validate request body
		if employeeID == "" || employeeID == ":" {
			utils.RespondWithError(ctx, http.StatusBadRequest, "employee Id is missing the request")
			return
		}

		ctx.Next()
	}
}
