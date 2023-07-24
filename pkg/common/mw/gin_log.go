package mw

import (
	"bytes"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

type responseWriter struct {
	gin.ResponseWriter
	buf *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.buf.Write(b)
	return w.ResponseWriter.Write(b)
}

func GinLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		req, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			c.Abort()
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewReader(req))
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			buf:            bytes.NewBuffer(nil),
		}
		c.Writer = writer
		c.Next()
		resp := writer.buf.Bytes()
		log.ZDebug(c, "gin request", "method", c.Request.Method, "time", time.Since(start).String(), "uri", c.Request.RequestURI, "req header", c.Request.Header, "req body", string(req), "code", c.Writer.Status(), "resp header", c.Writer.Header(), "resp", string(resp))
	}
}
