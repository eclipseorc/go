package login

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_util/api/v1/models"
	"net/http"
)

func UserLogin(c *gin.Context) {
	fmt.Println("i'm login")
	var user models.User

	user.Name = c.PostForm("username")

	err := user.Add(user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": 1, "msg": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"err": 0, "msg": "登录成功"})
}
