package utilities

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestIsStringValueExistOnArray(t *testing.T) {
	// Set Gin to Test mode
	gin.SetMode(gin.TestMode)

	var testArray = []string{"value1", "target_value", "value3"}
	var lookFor = "target_value"
	assert.Equal(t, IsStringValueExistOnArray(testArray, &lookFor), true)
}
