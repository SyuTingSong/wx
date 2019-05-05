package login

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func L2(c *gin.Context) {
	l2, ok := c.GetQuery("l2")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "l2 is required",
		})
		return
	}
}

func UserInfo(c *gin.Context) {
}
