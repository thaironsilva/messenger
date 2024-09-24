package cognitoClient

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

const cognitoClientId = "qiqoemp8hjlt1ar1qf1233o7t"
const userPoolId = "us-east-2_bGTPLFgM7"

type CognitoInterface interface {
	SignUp(user *CognitoUser) error
	ConfirmAccount(user *UserConfirmation) error
	SignIn(user *UserLogin) (string, error)
	GetUserByToken(token string) (*cognito.GetUserOutput, error)
	UpdatePassword(user *UserLogin) error
}

type CognitoUser struct {
	NickName string `json:"nickname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type cognitoClient struct {
	cognitoClient *cognito.CognitoIdentityProvider
	appClientID   string
}

type UserConfirmation struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

func NewCognitoClient() CognitoInterface {
	config := &aws.Config{Region: aws.String("us-east-2")}
	sess, err := session.NewSession(config)
	if err != nil {
		panic(err)
	}
	client := cognito.New(sess)

	return &cognitoClient{
		cognitoClient: client,
		appClientID:   cognitoClientId,
	}
}

func (c *cognitoClient) SignUp(user *CognitoUser) error {
	userCognito := &cognito.SignUpInput{
		ClientId: aws.String(c.appClientID),
		Username: aws.String(user.Email),
		Password: aws.String(user.Password),
		UserAttributes: []*cognito.AttributeType{
			{
				Name:  aws.String("nickname"),
				Value: aws.String(user.NickName),
			},
			{
				Name:  aws.String("email"),
				Value: aws.String(user.Email),
			},
		},
	}
	_, err := c.cognitoClient.SignUp(userCognito)
	if err != nil {
		return err
	}
	return nil
}

func (c *cognitoClient) ConfirmAccount(user *UserConfirmation) error {
	confirmationInput := &cognito.ConfirmSignUpInput{
		Username:         aws.String(user.Email),
		ConfirmationCode: aws.String(user.Code),
		ClientId:         aws.String(c.appClientID),
	}
	_, err := c.cognitoClient.ConfirmSignUp(confirmationInput)
	if err != nil {
		return err
	}
	return nil
}

func (c *cognitoClient) SignIn(user *UserLogin) (string, error) {
	authInput := &cognito.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: aws.StringMap(map[string]string{
			"USERNAME": user.Email,
			"PASSWORD": user.Password,
		}),
		ClientId: aws.String(c.appClientID),
	}
	result, err := c.cognitoClient.InitiateAuth(authInput)
	if err != nil {
		return "", err
	}
	return *result.AuthenticationResult.AccessToken, nil
}

func (c *cognitoClient) GetUserByToken(token string) (*cognito.GetUserOutput, error) {
	input := &cognito.GetUserInput{
		AccessToken: aws.String(token),
	}
	result, err := c.cognitoClient.GetUser(input)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *cognitoClient) UpdatePassword(user *UserLogin) error {
	input := &cognito.AdminSetUserPasswordInput{
		UserPoolId: aws.String(userPoolId),
		Username:   aws.String(user.Email),
		Password:   aws.String(user.Password),
		Permanent:  aws.Bool(true),
	}
	_, err := c.cognitoClient.AdminSetUserPassword(input)
	if err != nil {
		return err
	}
	return nil
}
