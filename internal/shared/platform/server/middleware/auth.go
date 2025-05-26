package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rubenbupe/recipe-video-parser/internal/users/domain"
	"github.com/rubenbupe/recipe-video-parser/internal/users/platform/storage/sql"
)

// AuthMiddleware extrae el token Bearer del header Authorization, busca el usuario por API key y lo añade al contexto.
func AuthMiddleware(userRepo sql.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No Authorization header"})
			return
		}
		token := strings.TrimPrefix(header, "Bearer ")
		if token == "" || token == header {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Bearer token"})
			return
		}

		userApiKey, err := domain.NewUserApiKey(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key format"})
			return
		}

		user, err := userRepo.GetByApiKey(c.Request.Context(), userApiKey)
		if err != nil || user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing user for token"})
			return
		}
		// Guardar el usuario en el contexto para los handlers siguientes
		c.Set("user", user)
		c.Next()
	}
}

// GetUserFromContext obtiene el usuario autenticado del contexto Gin.
func GetUserFromContext(c *gin.Context) (*domain.User, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}
	usr, ok := user.(*domain.User)
	return usr, ok
}
