package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORSMiddleware 处理跨域资源共享(CORS)
// 允许前端应用从不同的域名或端口访问后端API
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许的来源域名，可以设置为具体的域名，如 http://localhost:8081
		// 在开发环境中，可以设置为 * 允许所有域名
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		
		// 允许的HTTP方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		
		// 允许的HTTP头
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		
		// 允许浏览器缓存预检请求结果的时间（秒）
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		
		// 允许客户端获取自定义头信息
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		
		// 允许请求携带认证信息（如cookies）
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// 处理OPTIONS请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		// 继续处理请求
		c.Next()
	}
}
