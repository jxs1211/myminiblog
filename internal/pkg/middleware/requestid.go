// Copyright 2022 Innkeeper Belm(孔令飞) <nosbelm@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/marmotedu/miniblog.

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/marmotedu/miniblog/internal/pkg/known"
)

// RequestID is a Gin middleware that injects the `X-Request-ID` key-value pair into the context and response of each HTTP request.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the `X-Request-ID` header exists in the request. If it does, reuse it; otherwise, create a new one.
		requestID := c.Request.Header.Get(known.XRequestIDKey)

		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Save the RequestID in the gin.Context for later use in the program.
		c.Set(known.XRequestIDKey, requestID)

		// Save the RequestID in the HTTP response header, with the key `X-Request-ID`.
		c.Writer.Header().Set(known.XRequestIDKey, requestID)
		c.Next()
	}
}
