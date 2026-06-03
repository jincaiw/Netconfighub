package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	ErrCodeSuccess       = 0
	ErrCodeBadRequest    = 40001
	ErrCodeDuplicateName = 40002
	ErrCodeInvalidIP     = 40003
	ErrCodeUnauthorized  = 40101
	ErrCodeDatabase      = 50001
	ErrCodeSSHConnect    = 50002
	ErrCodeTelnetConnect = 50003
	ErrCodeCommandFailed = 50004
	ErrCodeGitCommit     = 50005
)

type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PageData struct {
	Items    interface{} `json:"items"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, ApiResponse{
		Code:    ErrCodeSuccess,
		Message: "success",
		Data:    data,
	})
}

func Error(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, ApiResponse{
		Code:    code,
		Message: message,
	})
}

func PageResult(c *gin.Context, items interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, ApiResponse{
		Code:    ErrCodeSuccess,
		Message: "success",
		Data: PageData{
			Items:    items,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

func BadRequest(c *gin.Context, msg string) {
	Error(c, http.StatusBadRequest, ErrCodeBadRequest, msg)
}

func DuplicateName(c *gin.Context, msg string) {
	Error(c, http.StatusBadRequest, ErrCodeDuplicateName, msg)
}

func InvalidIP(c *gin.Context, msg string) {
	Error(c, http.StatusBadRequest, ErrCodeInvalidIP, msg)
}

func Unauthorized(c *gin.Context, msg string) {
	Error(c, http.StatusUnauthorized, ErrCodeUnauthorized, msg)
}

func DatabaseError(c *gin.Context, msg string) {
	Error(c, http.StatusInternalServerError, ErrCodeDatabase, msg)
}
