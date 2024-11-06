package routes

import "github.com/gin-gonic/gin"

func UserGroup(rg *gin.RouterGroup) {
	rg.GET("/")
}

func AdminGroup(rg *gin.RouterGroup) {

}
