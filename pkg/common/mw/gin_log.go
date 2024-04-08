// Copyright © 2023 OpenIM open source community. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mw

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/openimsdk/tools/log"
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
		req, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			c.Abort()
			return
		}
		start := time.Now()
		log.ZDebug(c, "gin request", "method", c.Request.Method, "uri", c.Request.RequestURI, "req", string(req))
		c.Request.Body = io.NopCloser(bytes.NewReader(req))
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			buf:            bytes.NewBuffer(nil),
		}
		c.Writer = writer
		c.Next()
		resp := writer.buf.Bytes()
		log.ZDebug(c, "gin response", "time", time.Since(start), "status", c.Writer.Status(), "resp", string(resp))
	}
}
