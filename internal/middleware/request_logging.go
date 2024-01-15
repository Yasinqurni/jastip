package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"jastip-app/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestLog struct {
	Method string `json:"method"`
	URL    string `json:"url"`
	Body   string `json:"body"`
}

type ResponseLog struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Body   string `json:"body"`
}

type GinLog struct {
	Request  *RequestLog  `json:"request"`
	Response *ResponseLog `json:"response"`
}

func GinLogger(ginLogger *logger.Logger, errorLogger *logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := ctx.Request
		body, err := io.ReadAll(request.Body)
		if err != nil {
			errorLogger.Info(err.Error())
		}
		request.Body = io.NopCloser(bytes.NewBuffer(body))

		writer := &responseWriter{ResponseWriter: ctx.Writer}
		ctx.Writer = writer

		ctx.Next()

		req := RequestLog{
			Method: request.Method,
			URL:    ctx.Request.RequestURI,
			Body:   string(body),
		}

		res := ResponseLog{
			Code:   writer.status,
			Status: http.StatusText(writer.status),
			Body:   writer.responseBody.String(),
		}

		ginlog := GinLog{Request: &req, Response: &res}
		ginlogJSON, err := json.Marshal(&ginlog)
		if err != nil {
			errorLogger.Info(err.Error())
			return
		}

		ginLogger.Info(string(ginlogJSON))
	}
}

type responseWriter struct {
	gin.ResponseWriter
	status       int
	responseBody *bytes.Buffer
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriter) Write(data []byte) (int, error) {
	if w.responseBody == nil {
		w.responseBody = &bytes.Buffer{}
	}
	w.responseBody.Write(data)

	return w.ResponseWriter.Write(data)
}
