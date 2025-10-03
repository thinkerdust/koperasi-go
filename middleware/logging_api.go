package middleware

import (
	"bytes"
	"io"
	"strings"
	"koperasi-go/db"
	"koperasi-go/model"

	"github.com/gin-gonic/gin"
)

// custom response writer
type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LogRouteAPI() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Hanya log route yang prefix /api/
		if !strings.HasPrefix(c.Request.RequestURI, "/api/") {
			c.Next()
			return
		}

		// Copy request body
		reqBody, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody)) // reset body biar bisa dibaca ulang handler

		// Wrap response writer
		bw := &bodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bw

		// Process request
		c.Next()

		// Ambil user_id kalau ada
		var userID uint = 0
		if val, exists := c.Get("user_id"); exists {
			if id, ok := val.(uint); ok {
				userID = id
			}
		}

		// Insert log ke DB
		log := model.LoggingAPI{
			URI:      c.Request.RequestURI,
			Method:   c.Request.Method,
			Request:  string(reqBody),
			Response: bw.body.String(),
			IP:       c.ClientIP(),
			UserID:   userID,
		}

		db.DB.Create(&log)
	}
}
