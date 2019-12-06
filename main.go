package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ciiiii/iconServiceBackend/cos"
	"github.com/ciiiii/iconServiceBackend/config"
	"github.com/golang/groupcache/lru"
	"net/http"
	"strings"
)

func main() {
	cosService := cos.Init()
	cache := lru.New(100)
	r := gin.New()
	r.Use(cors.New(cors.Config{
		AllowMethods: []string{"GET", "OPTIONS"},
		AllowOrigins: []string{"http://localhost:8000"},
		AllowOriginFunc: func(origin string) bool {
			return strings.Contains(origin, "chrome-extension://")
		},
	}))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	icons := r.Group("icons")
	{
		icons.GET("/", func(c *gin.Context) {
			query := c.Query("query")
			tag := c.Query("tag")
			marker := c.Query("marker")
			var prefix string
			search := false
			prefix = strings.Join([]string{"svgs", tag, query}, "/")
			if query != "" {
				search = true
			}
			key := prefix + marker
			value, ok := cache.Get(key)
			if ok {
				c.JSON(200, gin.H{
					"success": true,
					"data":    value,
				})
				return
			}
			iconList, err := cosService.List(prefix, marker, search)
			if err != nil {
				c.JSON(400, gin.H{
					"success": false,
					"message": "cos service error",
				})
				return
			}
			cache.Add(key, iconList)
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
