package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
)

func main() {
	r := gin.Default()
	r.Static("/statics", "./statics")
	r.LoadHTMLGlob("templates/*")
	r.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", gin.H{
			"code": 0,
		})
	})
	r.GET("/login/:id", func(c *gin.Context) {
		// 获取 URL 中的参数 :id
		id := c.Param("id")

		// 这里你可以根据用户ID去数据库中获取用户信息或者执行其他逻辑

		// 渲染 mainpage.html 模板
		c.HTML(http.StatusOK, "mainpage.html", gin.H{
			"code": 2,
			"user": id,
		})
	})
	r.DELETE("/user/:id", func(c *gin.Context) {
		// 获取 URL 参数中的用户 ID
		userID := c.Param("id")

		// 连接数据库
		db := connect()

		// 在数据库中删除用户
		if err := db.Where("username = ?", userID).Delete(&UserInfo{}).Error; err != nil {
			// 如果删除出错，返回错误信息
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 删除成功，返回成功信息
		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	})

	r.POST("/login", func(c *gin.Context) {
		db := connect()
		username := c.PostForm("username")
		password := c.PostForm("password")
		var uu = UserInfo{}
		result := db.Find(&uu, "username=? AND password=?", username, password)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.HTML(200, "login.html", gin.H{
				"code": 1,
			})
		} else {
			c.Redirect(http.StatusFound, "/login/"+username)
		}
	})
	r.GET("/index", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})
	r.GET("/register", func(c *gin.Context) {
		c.HTML(200, "register.html", gin.H{
			"code": 0,
		})
	})
	r.POST("/register", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		var uu = UserInfo{}
		db := connect()
		result := db.Find(&uu, "username=?", username)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			userinfo := UserInfo{username, password}
			db.Create(&userinfo)
			c.HTML(200, "register.html", gin.H{
				"code": 2,
			})
		} else {
			c.HTML(200, "register.html", gin.H{
				"code": 1,
			})
		}
		defer db.Close()
	})
	r.POST("/changeName/:id", func(c *gin.Context) {
		db := connect()
		// 从URL参数中获取用户ID
		username := c.Param("id")

		// 从请求体中解析新的名字
		var newName struct {
			NewName string `json:"newName"`
		}
		if err := c.ShouldBindJSON(&newName); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 查询数据库，更新用户信息
		var existingUser UserInfo
		result := db.Where("username = ?", newName.NewName).First(&existingUser)
		if result.RowsAffected != 0 {
			// 如果存在相同的记录，说明新名字已经被占用
			c.JSON(http.StatusBadRequest, gin.H{"error": "New name is already taken"})
			return
		}

		// 更新数据库，将旧名字修改为新名字
		var user UserInfo
		result = db.Where("username = ?", username).First(&user)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find user"})
			return
		}
		user.Username = newName.NewName
		db.Model(&user).Where("username = ?", username).Update("username", newName.NewName)
		c.JSON(http.StatusOK, gin.H{"message": "Name updated successfully"})
	})
	r.Run("127.0.0.1:8080")
}
