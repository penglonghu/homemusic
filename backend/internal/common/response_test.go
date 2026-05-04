package common

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSuccessResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	Success(c, "test data")
	assert.Equal(t, http.StatusOK, c.Writer.Status())

	var resp Response
	err := json.Unmarshal(c.Writer.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, CodeSuccess, resp.Code)
	assert.Equal(t, "success", resp.Msg)
	assert.Equal(t, "test data", resp.Data)
}

func TestFailResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	Fail(c, CodeParamError, "param error")
	assert.Equal(t, http.StatusOK, c.Writer.Status())

	var resp Response
	err := json.Unmarshal(c.Writer.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, CodeParamError, resp.Code)
	assert.Equal(t, "param error", resp.Msg)
}
