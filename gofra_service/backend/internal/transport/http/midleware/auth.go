package midleware

import (
	"context"
	"net/http"
	"time"

	"Gofra_Market/internal/domain"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ctxKey string 

type SessionStore interface {
	BySID(ctx context.Context, sid string) (*domain.Session, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}

func Auth(sess SessionStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("sid")
		if err != nil {
			c.Next()
			return
		}

		s, err := sess.BySID(c.Request.Context(), cookie)
		if err != nil {
			c.Next()
			return
		}

		if s.ExpiredAt.Before(time.Now()) {
			_ = sess.Delete(c.Request.Context(), s.ID)
			http.SetCookie(c.Writer, &http.Cookie{Name: "sid", Value: "", MaxAge: -1, Path: "/"})
			c.Next()
			return
		}

		c.Set(string(ctxKey("userID")), s.UserID)
		c.Next()
	}
}
