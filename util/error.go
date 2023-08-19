package util

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func ErrorResponse(err error) map[string]interface{} {

	return gin.H{"error": err.Error()}

}

func ErrorResponseString(err string) map[string]interface{} {

	return gin.H{"error": fmt.Errorf("%s", err)}

}
