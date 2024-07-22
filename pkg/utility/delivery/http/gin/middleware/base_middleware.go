package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/common"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	correlationIDHeader = "X-Correlation-ID"
)

// CorrelationIDMiddleware adds a correlation ID to requests and responses.
func CorrelationIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationID := c.GetHeader(correlationIDHeader)
		if len(correlationID) == 0 {
			correlationID = uuid.New().String()
		}

		c.Set(correlationIDHeader, correlationID)
		c.Writer.Header().Set(correlationIDHeader, correlationID)
		c.Next()
	}
}

// SecureHeadersMiddleware adds secure headers to HTTP responses.
func SecureHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		securityConfig := config.GetSecurityConfig()

		c.Header(securityConfig.ContentSecurityPolicyHeader.Key, securityConfig.ContentSecurityPolicyHeader.Value)
		c.Header(securityConfig.StrictTransportSecurityHeader.Key, securityConfig.StrictTransportSecurityHeader.Value)
		c.Header(securityConfig.XContentTypeOptionsHeader.Key, securityConfig.XContentTypeOptionsHeader.Value)

		c.Next()
	}
}

// CSPMiddleware adds Content Security Policy (CSP) headers to HTTP responses.
func CSPMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		securityConfig := config.GetSecurityConfig()

		c.Writer.Header().Set(
			securityConfig.ContentSecurityPolicyHeaderFull.Key,
			securityConfig.ContentSecurityPolicyHeaderFull.Value,
		)

		c.Next()
	}
}

// RateLimitMiddleware implements rate limiting to control the number of requests from clients.
func RateLimitMiddleware() gin.HandlerFunc {
	securityConfig := config.GetSecurityConfig()

	limiterOptions := &limiter.ExpirableOptions{
		DefaultExpirationTTL: time.Second,
	}

	limiter := tollbooth.NewLimiter(securityConfig.RateLimit, limiterOptions)
	return func(c *gin.Context) {
		tollbooth_gin.LimitHandler(limiter)(c)
		c.Next()
	}
}

// ValidateInputMiddleware allows specific HTTP methods and checks for the content type.
func ValidateInputMiddleware(logger applicationModel.Logger) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		contentType := ginContext.GetHeader(constants.ContentType)
		securityConfig := config.GetSecurityConfig()

		// Check if the request method is in the list of allowed methods.
		if validator.IsSliceNotContains(securityConfig.AllowedHTTPMethods, ginContext.Request.Method) {
			allowedMethods := strings.Join(config.GetSecurityConfig().AllowedHTTPMethods, ", ")
			notification := constants.InvalidHTTPMethodNotification + allowedMethods
			httpRequestError := httpError.NewHTTPRequestError(location+"ValidateInputMiddleware.AllowedHTTPMethods", ginContext.Request.Method, notification)
			abortWithStatusJSON(ginContext, logger, httpRequestError, constants.StatusBadRequest)
			return
		}

		// Check if the content type is in the list of allowed content types.
		if len(contentType) > 0 && validator.IsSliceNotContains(securityConfig.AllowedContentTypes, contentType) {
			allowedContentTypes := strings.Join(config.GetSecurityConfig().AllowedContentTypes, ", ")
			notification := constants.InvalidHTTPMethodNotification + allowedContentTypes
			httpRequestError := httpError.NewHTTPRequestError(location+"ValidateInputMiddleware.AllowedContentTypes", contentType, notification)
			abortWithStatusJSON(ginContext, logger, httpRequestError, constants.StatusBadRequest)
			return
		}

		ginContext.Next()
	}
}

// TimeoutMiddleware sets a timeout for each request.
func TimeoutMiddleware(logger applicationModel.Logger) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
		defer cancel()

		ginContext.Request = ginContext.Request.WithContext(ctx)
		ch := make(chan struct{})
		go func() {
			defer close(ch)
			ginContext.Next()
		}()

		select {
		case <-ch:
		case <-ctx.Done():
			httpInternalError := httpError.NewHTTPInternalError(location+"TimeOutMiddleware", ctx.Err().Error())
			abortWithStatusJSON(ginContext, logger, httpInternalError, constants.StatusBadGateway)
		}
	}
}

// LoggerMiddleware logs incoming requests and outgoing responses with additional context.
func LoggerMiddleware(logger applicationModel.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		correlationID := c.GetString("X-Correlation-ID")

		httpIncomingLog := commonModel.NewHTTPIncomingLog(
			location+"LoggerMiddleware",
			correlationID,
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			c.Request.UserAgent(),
		)

		c.Next()

		httpOutgoingLog := commonModel.NewHTTPOutgoingLog(
			location+"LoggerMiddleware",
			correlationID,
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			c.Request.UserAgent(),
			c.Writer.Status(),
			time.Since(start),
		)

		logger.Info(httpIncomingLog)
		logger.Info(httpOutgoingLog)
	}
}
