package utils

import (
	"assignment/internal/constants"
	employeeerror "assignment/internal/errors"

	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitLogClient() {
	Logger, _ = zap.NewDevelopment()
}

func RespondWithError(c *gin.Context, statusCode int, message string) {

	c.AbortWithStatusJSON(statusCode, employeeerror.EmployeeError{
		Trace:   c.Request.Header.Get(constants.TransactionID),
		Code:    statusCode,
		Message: message,
	})
}
