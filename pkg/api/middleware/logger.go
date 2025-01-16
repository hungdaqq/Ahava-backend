package middleware

import (
	"ahava/pkg/domain"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func DefaultStructuredLogger() gin.HandlerFunc {
	return StructuredLogger(&log.Logger)
}

func StructuredLogger(logger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		param := gin.LogFormatterParams{}

		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)
		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}

		param.ClientIP = c.Request.Header.Get("X-Forwarded-For")
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		param.BodySize = c.Writer.Size()
		if raw != "" {
			path = path + "?" + raw
		}
		param.Path = path

		// Log using the params
		var logEvent *zerolog.Event
		if c.Writer.Status() >= 500 {
			logEvent = logger.Error()
		} else {
			logEvent = logger.Info()
		}

		logEvent.
			Str("client_ip", param.ClientIP).
			Str("method", param.Method).
			Int("status_code", param.StatusCode).
			Int("body_size", param.BodySize).
			Str("path", param.Path).
			Str("latency", param.Latency.String()).
			Msg(param.ErrorMessage)

		go func() {
			tranCh <- domain.RequestTransaction{
				Method:       param.Method,
				Path:         param.Path,
				StatusCode:   param.StatusCode,
				ClientIP:     param.ClientIP,
				Latency:      param.Latency.String(),
				BodySize:     param.BodySize,
				ErrorMessage: param.ErrorMessage,
			}
		}()
	}
}

func SaveRequestTransaction(db *gorm.DB) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	trans := make([]domain.RequestTransaction, 0)
	for {
		select {
		case <-ticker.C:
			if len(trans) > 0 {
				fmt.Printf("Saving request transactions, count: %d", len(trans))
				db.Create(&trans)
				trans = make([]domain.RequestTransaction, 0)
			}
		case tran := <-tranCh:
			trans = append(trans, tran)
		}
	}
}

var (
	tranCh = make(chan domain.RequestTransaction, 100)
)
