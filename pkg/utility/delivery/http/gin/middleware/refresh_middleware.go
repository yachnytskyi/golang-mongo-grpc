package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/internal/common/interfaces"
	utility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/utility"
	http "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

// RefreshTokenAuthenticationMiddleware is a Gin middleware for handling user authentication using refresh tokens.
func RefreshTokenAuthenticationMiddleware(config interfaces.Config, logger interfaces.Logger) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
		defer cancel()

		// Extract the refresh token from the request headers or cookies.
		refreshToken := extractRefreshToken(ginContext, location+"RefreshTokenAuthenticationMiddleware")
		if validator.IsError(refreshToken.Error) {
			abortWithStatusJSON(ginContext, logger, refreshToken.Error, constants.StatusUnauthorized)
			return
		}

		// Extract the refresh token from the request headers or cookies.
		refreshTokenConfig := config.GetConfig()
		userTokenPayload := utility.ValidateJWTToken(
			logger,
			location+"RefreshTokenAuthenticationMiddleware",
			refreshToken.Data,
			refreshTokenConfig.RefreshToken.PublicKey,
		)
		if validator.IsError(userTokenPayload.Error) {
			httpAuthorizationError := http.NewHTTPAuthorizationError(location+"RefreshTokenAuthenticationMiddleware.ValidateJWTToken", constants.LoggingErrorNotification)
			abortWithStatusJSON(ginContext, logger, httpAuthorizationError, constants.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, constants.ID, userTokenPayload.Data.UserID)
		ctx = context.WithValue(ctx, constants.UserRole, userTokenPayload.Data.Role)
		ginContext.Request = ginContext.Request.WithContext(ctx)
		ginContext.Next()
	}
}
