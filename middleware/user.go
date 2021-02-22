package middleware

import (
	"net/http"
	"webapp/appframework/app"
	"webapp/appframework/code"

	"github.com/gin-gonic/gin"
)

var auth = NewJWT()

func CheckUserToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			app.JsonResponse(c, http.StatusUnauthorized, code.ERROR_TOKEN_EMPTY, "")
			c.Abort()
			return
		}
		claims, err := auth.ParseToken(token)
		if err != nil {
			app.JsonResponse(c, http.StatusUnauthorized, code.ERROR_TOKEN_INVALID, "")
			c.Abort()
			return
		}
		if claims == nil || claims.UID == 0 {
			app.JsonResponse(c, http.StatusUnauthorized, code.ERROR_TOKEN_INVALID, "")
			c.Abort()
			return
		}
		//log.Printf("token == %v, \n uid = %v", token, claims.UID)

		c.Set("uid", claims.UID)
		c.Next()
	}
}
