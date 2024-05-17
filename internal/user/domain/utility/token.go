package utility

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	signingMethod       = "RS256"
	userIDClaim         = "user_id"
	userRoleClaim       = "user_role"
	expirationClaim     = "exp"
	issuedAtClaim       = "iat"
	notBeforeClaim      = "nbf"
	location            = "internal.user.domain.utility."
	unexpectedMethod    = "unexpected method: %s"
	invalidTokenMessage = "validate: invalid token"
)

// GenerateJWTToken generates a JWT token with the provided UserTokenPayload, using the given private key,
// and sets the token's expiration based on the specified token lifetime.
func GenerateJWTToken(ctx context.Context, tokenLifeTime time.Duration, userTokenPayload commonModel.UserTokenPayload, privateKey string) (string, error) {
	// Decode the private key from base64-encoded string.
	decodedPrivateKey, decodeStringError := decodeBase64String(privateKey)
	if validator.IsError(decodeStringError) {
		return constants.EmptyString, decodeStringError
	}

	// Parse the private key for signing.
	key, parsePrivateKeyError := parsePrivateKey(decodedPrivateKey)
	if validator.IsError(parsePrivateKeyError) {
		return constants.EmptyString, parsePrivateKeyError
	}

	// Generate claims for the JWT token.
	// Create the signed token using the private key and claims.
	now := time.Now().UTC()
	claims := generateClaims(tokenLifeTime, now, userTokenPayload)
	token, newWithClaimsError := createSignedToken(key, claims)
	if validator.IsError(newWithClaimsError) {
		return constants.EmptyString, newWithClaimsError
	}

	return token, nil
}

// ValidateJWTToken validates a JWT token using the provided public key and returns the claims
// extracted from the token if it's valid.
func ValidateJWTToken(token string, publicKey string) (commonModel.UserTokenPayload, error) {
	// Decode the public key from a base64-encoded string.
	decodedPublicKey, decodeStringError := decodeBase64String(publicKey)
	if validator.IsError(decodeStringError) {
		return commonModel.UserTokenPayload{}, decodeStringError
	}

	// Parse the public key for verification.
	key, parsePublicKeyError := parsePublicKey(decodedPublicKey)
	if validator.IsError(parsePublicKeyError) {
		return commonModel.UserTokenPayload{}, parsePublicKeyError
	}

	// Parse and verify the token using the public key.
	parsedToken, parseTokenError := parseToken(token, key)
	if validator.IsError(parseTokenError) {
		return commonModel.UserTokenPayload{}, parseTokenError
	}

	// Extract and validate the claims from the parsed token.
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if ok && parsedToken.Valid {
		return commonModel.UserTokenPayload{
			UserID: fmt.Sprint(claims[userIDClaim]),
			Role:   fmt.Sprint(claims[userRoleClaim]),
		}, nil
	}

	internalError := domainError.NewInternalError(location+"parsedToken.Claims.NotOk", invalidTokenMessage)
	logging.Logger(internalError)
	return commonModel.UserTokenPayload{}, domainError.HandleError(internalError)
}

// decodeBase64String decodes a base64-encoded string into a byte slice.
func decodeBase64String(base64String string) ([]byte, error) {
	decodedPrivateKey, decodeStringError := base64.StdEncoding.DecodeString(base64String)
	if validator.IsError(decodeStringError) {
		internalError := domainError.NewInternalError(location+"decodeString.StdEncoding.DecodeString", decodeStringError.Error())
		logging.Logger(internalError)
		return []byte{}, domainError.HandleError(internalError)
	}

	return decodedPrivateKey, nil
}

// parsePrivateKey parses the RSA private key from the provided byte slice.
func parsePrivateKey(decodedPrivateKey []byte) (*rsa.PrivateKey, error) {
	key, parsePrivateKeyError := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if validator.IsError(parsePrivateKeyError) {
		internalError := domainError.NewInternalError(location+"parsePrivateKey.ParseRSAPrivateKeyFromPEM", parsePrivateKeyError.Error())
		logging.Logger(internalError)
		return nil, domainError.HandleError(internalError)
	}

	return key, nil
}

// generateClaims generates JWT claims with the specified token lifetime and UserTokenPayload.
func generateClaims(tokenLifeTime time.Duration, now time.Time, userTokenPayload commonModel.UserTokenPayload) jwt.MapClaims {
	return jwt.MapClaims{
		userIDClaim:     userTokenPayload.UserID,
		userRoleClaim:   userTokenPayload.Role,
		expirationClaim: now.Add(tokenLifeTime).Unix(),
		issuedAtClaim:   now.Unix(),
		notBeforeClaim:  now.Unix(),
	}
}

// createSignedToken creates a signed JWT token using the provided private key and claims.
func createSignedToken(key *rsa.PrivateKey, claims jwt.MapClaims) (string, error) {
	token, newWithClaimsError := jwt.NewWithClaims(jwt.GetSigningMethod(signingMethod), claims).SignedString(key)
	if validator.IsError(newWithClaimsError) {
		internalError := domainError.NewInternalError(location+"createSignedToken.NewWithClaims", newWithClaimsError.Error())
		logging.Logger(internalError)
		return constants.EmptyString, domainError.HandleError(internalError)
	}

	return token, nil
}

// parsePublicKey parses the RSA Public key from the provided byte slice.
func parsePublicKey(decodedPublicKey []byte) (*rsa.PublicKey, error) {
	key, parsePublicKeyError := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if validator.IsError(parsePublicKeyError) {
		internalError := domainError.NewInternalError(location+"parsePublicKey.ParseRSAPublicKeyFromPEM", parsePublicKeyError.Error())
		logging.Logger(internalError)
		return nil, domainError.HandleError(internalError)
	}

	return key, nil
}

// parseToken parses and verifies the JWT token using the provided public key.
func parseToken(token string, key *rsa.PublicKey) (*jwt.Token, error) {
	parsedToken, parseError := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		_, ok := t.Method.(*jwt.SigningMethodRSA)
		if ok {
			return key, nil
		}
		internalError := domainError.NewInternalError(location+"parseToken.jwt.Parse.NotOk", unexpectedMethod+" t.Header[alg]")
		logging.Logger(internalError)
		return nil, domainError.HandleError(internalError)
	})
	if validator.IsError(parseError) {
		internalError := domainError.NewInternalError(location+"parseToken.jwt.Parse", parseError.Error())
		logging.Logger(internalError)
		return nil, domainError.HandleError(internalError)
	}

	return parsedToken, nil
}
