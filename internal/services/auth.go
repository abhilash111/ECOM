package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/abhilash111/ecom/config"
	"github.com/abhilash111/ecom/internal/models"
	"github.com/abhilash111/ecom/internal/repository"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	SignUp(email, phoneNumber, password string, roles []models.Role) (*models.User, error)
	LoginWithPassword(username, password string) (*AuthResponse, error)
	InitiatePhoneLogin(phoneNumber string) (string, error)
	VerifyPhoneLogin(phoneNumber, otp string) (*AuthResponse, error)
	RefreshToken(refreshToken string) (*AuthResponse, error)
}

type authService struct {
	userService UserService
	otpService  OTPService
}

func NewAuthService(userService UserService, otpService OTPService) AuthService {
	return &authService{
		userService: userService,
		otpService:  otpService,
	}
}

func calculateSecretHash(clientID, clientSecret, username string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(username + clientID))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func (s *authService) SignUp(email, phoneNumber, password string, roles []models.Role) (*models.User, error) {
	// Create user in Cognito
	sess := session.Must(session.NewSession())
	cognitoSvc := cognitoidentityprovider.New(sess, aws.NewConfig().WithRegion(config.Envs.CognitoRegion))

	// Convert roles to comma-separated string
	var rolesStr string
	for i, role := range roles {
		if i > 0 {
			rolesStr += ","
		}
		rolesStr += string(role)
	}
	fmt.Println("rolesStr", rolesStr)
	fmt.Println("email", email)
	fmt.Println("cognitoAppID", config.Envs.CognitoAppID)
	fmt.Println("CognitoAppSecret", config.Envs.CognitoAppSecret)
	fmt.Println("JWTSECRET", config.Envs.JWTSecret)

	fmt.Println("cognitoAppSecret1", config.Envs.CognitoAppID, config.Envs.CognitoAppSecret, email, aws.String(calculateSecretHash(
		config.Envs.CognitoAppID,
		config.Envs.CognitoAppSecret,
		email,
	)))
	// Create user in Cognito
	signUpInput := &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(config.Envs.CognitoAppID),
		SecretHash: aws.String(calculateSecretHash(
			config.Envs.CognitoAppID,
			config.Envs.CognitoAppSecret,
			email,
		)),

		Username: aws.String(email),
		Password: aws.String(password),
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(email),
			},
			{
				Name:  aws.String("phone_number"),
				Value: aws.String(phoneNumber),
			},
			{
				Name:  aws.String("custom:roles"),
				Value: aws.String(rolesStr),
			},
		},
	}

	_, err := cognitoSvc.SignUp(signUpInput)
	if err != nil {
		fmt.Println("Error signing up user:", err)
		return nil, err
	}

	// Confirm user (in production, you'd want email/phone verification)
	_, err = cognitoSvc.AdminConfirmSignUp(&cognitoidentityprovider.AdminConfirmSignUpInput{
		UserPoolId: aws.String(config.Envs.CognitoPoolID),
		Username:   aws.String(email),
	})
	if err != nil {
		return nil, err
	}

	// Get user details to store in our DB
	user, err := cognitoSvc.AdminGetUser(&cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(config.Envs.CognitoPoolID),
		Username:   aws.String(email),
	})
	if err != nil {
		return nil, err
	}

	// Create user in our database
	dbUser, err := s.userService.CreateUser(email, phoneNumber, *user.Username, roles)
	if err != nil {
		return nil, err
	}

	return dbUser, nil
}

func (s *authService) LoginWithPassword(username, password string) (*AuthResponse, error) {
	sess := session.Must(session.NewSession())
	cognitoSvc := cognitoidentityprovider.New(sess, aws.NewConfig().WithRegion(config.Envs.CognitoRegion))
	fmt.Println("cognitoAppID", config.Envs.CognitoAppID)
	authParams := map[string]*string{
		"USERNAME": aws.String(username),
		"PASSWORD": aws.String(password),
	}

	// Add SECRET_HASH if client secret is configured
	if config.Envs.CognitoAppSecret != "" {
		authParams["SECRET_HASH"] = aws.String(calculateSecretHash(
			config.Envs.CognitoAppID,
			config.Envs.CognitoAppSecret,
			username,
		))
	}

	authInput := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow:       aws.String("USER_PASSWORD_AUTH"),
		ClientId:       aws.String(config.Envs.CognitoAppID),
		AuthParameters: authParams,
	}

	result, err := cognitoSvc.InitiateAuth(authInput)
	if err != nil {
		fmt.Println("Error during authentication:", err)
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  *result.AuthenticationResult.AccessToken,
		RefreshToken: *result.AuthenticationResult.RefreshToken,
		ExpiresIn:    *result.AuthenticationResult.ExpiresIn,
		TokenType:    *result.AuthenticationResult.TokenType,
	}, nil
}

func (s *authService) InitiatePhoneLogin(phoneNumber string) (string, error) {
	// Check if user exists with this phone number
	user, err := s.userService.GetUserByPhone(phoneNumber)
	if err != nil {
		return "", fmt.Errorf("user not found")
	}

	// Generate and send OTP
	otp, err := s.otpService.SendOTP(phoneNumber)
	if err != nil {
		return "", err
	}

	fmt.Println("user", user)

	return otp, nil
}

func (s *authService) VerifyPhoneLogin(phoneNumber, otp string) (*AuthResponse, error) {
	// Verify OTP
	valid, err := s.otpService.VerifyOTP(phoneNumber, otp)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, fmt.Errorf("invalid OTP")
	}

	// Clean up OTP
	_ = repository.DeleteOTP(phoneNumber)

	// Get user by phone number
	user, err := s.userService.GetUserByPhone(phoneNumber)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Initiate custom auth flow with Cognito
	sess := session.Must(session.NewSession())
	cognitoSvc := cognitoidentityprovider.New(sess, aws.NewConfig().WithRegion(config.Envs.CognitoRegion))

	// Initiate custom auth flow
	authParams := map[string]*string{
		"USERNAME": aws.String(user.Email),
	}

	fmt.Println("user.Email", user.Email)
	// Add SECRET_HASH if client secret exists
	if config.Envs.CognitoAppSecret != "" {
		authParams["SECRET_HASH"] = aws.String(calculateSecretHash(
			config.Envs.CognitoAppID,
			config.Envs.CognitoAppSecret,
			user.Email,
		))
	}

	initAuthResp, err := cognitoSvc.InitiateAuth(&cognitoidentityprovider.InitiateAuthInput{
		AuthFlow:       aws.String("CUSTOM_AUTH"),
		ClientId:       aws.String(config.Envs.CognitoAppID),
		AuthParameters: authParams,
	})
	if err != nil {
		fmt.Println("Error initiating auth:", err)
		return nil, fmt.Errorf("initiate auth failed: %w", err)
	}

	// Prepare challenge response
	challengeResponses := map[string]*string{
		"USERNAME": aws.String(user.Email),
		"ANSWER":   aws.String("1234"),
	}

	// Add SECRET_HASH to challenge if needed
	if config.Envs.CognitoAppSecret != "" {
		challengeResponses["SECRET_HASH"] = authParams["SECRET_HASH"]
	}

	respondToAuthResp, err := cognitoSvc.RespondToAuthChallenge(&cognitoidentityprovider.RespondToAuthChallengeInput{
		ClientId:           aws.String(config.Envs.CognitoAppID),
		ChallengeName:      aws.String("CUSTOM_CHALLENGE"),
		Session:            initAuthResp.Session,
		ChallengeResponses: challengeResponses,
	})

	if err != nil {
		fmt.Println("Error responding to auth challenge:", err)
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  *respondToAuthResp.AuthenticationResult.AccessToken,
		RefreshToken: *respondToAuthResp.AuthenticationResult.RefreshToken,
		ExpiresIn:    *respondToAuthResp.AuthenticationResult.ExpiresIn,
		TokenType:    *respondToAuthResp.AuthenticationResult.TokenType,
	}, nil
}

func (s *authService) RefreshToken(refreshToken string) (*AuthResponse, error) {
	sess := session.Must(session.NewSession())
	cognitoSvc := cognitoidentityprovider.New(sess, aws.NewConfig().WithRegion(config.Envs.CognitoRegion))

	authInput := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("REFRESH_TOKEN_AUTH"),
		ClientId: aws.String(config.Envs.CognitoAppID),
		AuthParameters: map[string]*string{
			"REFRESH_TOKEN": aws.String(refreshToken),
		},
	}

	result, err := cognitoSvc.InitiateAuth(authInput)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken: *result.AuthenticationResult.AccessToken,
		ExpiresIn:   *result.AuthenticationResult.ExpiresIn,
		TokenType:   *result.AuthenticationResult.TokenType,
	}, nil
}

// CustomClaims extends the standard JWT claims with our custom fields
type CustomClaims struct {
	Roles []string `json:"roles"`
	jwt.RegisteredClaims
}

func parseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the token with Cognito's public keys
		// In production, you should fetch and cache Cognito's public keys
		return []byte(""), nil // This is simplified - actual implementation needs proper key verification
	})

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
