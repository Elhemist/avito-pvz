package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	roleCtx             = "role"
)

func (h *Handler) JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader(authorizationHeader)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
			c.Abort()
			return
		}

		if len(tokenString) < (len("Bearer ")) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}
		tokenString = tokenString[len("Bearer "):]

		userId, role, err := h.services.Authorization.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set(userCtx, userId)
		c.Set(roleCtx, role)

		c.Next()
	}
}

// func getUserId(c *gin.Context) (uuid.UUID, error) {
// 	id, ok := c.Get(userCtx)
// 	if !ok {
// 		return uuid.Nil, fmt.Errorf("user id not found")
// 	}

// 	idUuid, ok := id.(uuid.UUID)
// 	if !ok {
// 		return uuid.Nil, fmt.Errorf("user id is of invalid type")
// 	}

// 	return idUuid, nil
// }
