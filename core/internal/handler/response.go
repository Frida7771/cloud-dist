package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func respondError(c *gin.Context, err error) {
	if err == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unknown error"})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

func respondUnauthorized(c *gin.Context, err error) {
	if err == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
}
