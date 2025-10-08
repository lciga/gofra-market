package midleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Достаёт cookie cfg.CookieName ("sid"), sess.BySID, проверяет expires_at,
// кладёт userID в контекст: c.Set("userID", primitive.ObjectID)
type ctxKey string // userID

// SessionStore is the minimal interface used by the middleware. Using an interface
// allows tests to inject a fake implementation.
type SessionStore interface {
	BySID(ctx context.Context, sid string) (*struct {
		ID        primitive.ObjectID
		UserID    primitive.ObjectID
		ExpiredAt time.Time
	}, error)
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
			// invalid session — ignore and continue as anonymous
			c.Next()
			return
		}

		if s.ExpiredAt.Before(time.Now()) {
			// session expired
			_ = sess.Delete(c.Request.Context(), s.ID)
			http.SetCookie(c.Writer, &http.Cookie{Name: "sid", Value: "", MaxAge: -1, Path: "/"})
			c.Next()
			return
		}

		c.Set(string(ctxKey("userID")), s.UserID)
		c.Next()
	}
}
