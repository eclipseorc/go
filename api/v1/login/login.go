package login

import (
	"github.com/gin-gonic/gin"
	"go_util/modes"
	"net/http"
)

func UserLogin(c *gin.Context) {
	var user modes.AdminUser
	user.Phone = c.PostForm("username")

	if err := c.PostForm("username"); err == "" {
		c.JSON(http.StatusOK, gin.H{"err": 2, "msg": "用户名不能为空"})
		return
	}
	if fage, err := user.Get(); nil == err {
		if !fage {
			c.JSON(http.StatusOK, gin.H{"err": 2, "msg": "此账号不属于教学系统"})
			return
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"err": 1, "msg": "账号不存在"})
		return
	}
	if user.Pass != c.PostForm("password") {
		c.JSON(http.StatusOK, gin.H{"err": 1, "msg": "密码不正确"})
		return
	}
}
