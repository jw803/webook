package cors

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func CORSHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		//AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "HEAD", "POST", "PUT", "PATCH", "OPTIONS", "DELETE"},
		AllowHeaders: []string{"Content-Type", "Authorization", "x-shopline-session-token",
			"Dnt", "Referer", "Sec-Ch-Ua", "Sec-Ch-Ua-Mobile", "Sec-Ch-Ua-Platform", "User-Agent",
			"Accept", "Accept-Encoding", "Accept-Language", "Cache-Control", "Origin", "Pragma",
			"Refer", "Sec-Fetch-Dest", "Sec-Fetch-Mode", "Sec-Fetch-Site",
		},
		ExposeHeaders:    []string{"x-jwt-token", "x-refresh-token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.Contains(origin, "cst.shopline.io") ||
				strings.Contains(origin, "cst.shoplinetest.com") ||
				strings.Contains(origin, "shoplineapp.com") ||
				strings.Contains(origin, "localhost") {
				return true
			}
			return false
		},
		MaxAge: 12 * time.Hour,
	})
}
