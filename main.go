package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ciiiii/iconServiceBackend/cos"
	"github.com/ciiiii/iconServiceBackend/config"
	"net/http"
	"strings"
)

func main() {
	cosService := cos.Init()
	r := gin.New()
	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "OPTIONS"},
		AllowOrigins:     []string{"http://localhost:8000"},
		AllowOriginFunc: func(origin string) bool {
			return strings.Contains(origin, "chrome-extension://")
		},
	}))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	icons := r.Group("icons")
	{
		icons.GET("/", func(c *gin.Context) {
			iconList, err := cosService.List("")
			if err != nil {
				c.JSON(400, gin.H{
					"success": false,
					"message": "cos service errror",
				})
			}
			c.JSON(200, gin.H{
				"success": true,
				"data":    iconList,
			})
		})
	}

	gin.SetMode(config.Parser().Config.Mode)
	gin.DisableConsoleColor()
	server := &http.Server{
		Addr:    ":" + config.Parser().Config.Port,
		Handler: r,
	}
	panic(server.ListenAndServe())
}
