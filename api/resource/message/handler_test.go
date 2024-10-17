package message_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/thaironsilva/messenger/api/cognitoClient"
	"github.com/thaironsilva/messenger/api/resource/message"
	"github.com/thaironsilva/messenger/api/resource/user"
)

type MockStorage struct {
	err      error
	messages []message.Message
}

func (m *MockStorage) GetAll(sender_id string, receiver_id string) ([]message.Message, error) {
	return m.messages, m.err
}

func (m *MockStorage) Create(message message.Message) error {
	return m.err
}

type MockUserStorage struct {
	err   error
	user  user.User
	users []user.User
}

func (m *MockUserStorage) GetByUsername(username string) (user.User, error) {
	return m.user, m.err
}

func (m *MockUserStorage) GetByEmail(email string) (user.User, error) {
	return m.user, m.err
}

func (m *MockUserStorage) GetByString(name string) ([]user.User, error) {
	return m.users, m.err
}

func (m *MockUserStorage) GetAll() ([]user.User, error) {
	return m.users, m.err
}

func (m *MockUserStorage) Create(user user.User) error {
	return m.err
}

func (m *MockUserStorage) Update(user user.User) error {
	return m.err
}

func (m *MockUserStorage) Delete(id string) error {
	return m.err
}

type MockCognito struct {
	err   error
	token string
	user  cognito.GetUserOutput
}

func (m *MockCognito) SignUp(user *cognitoClient.CognitoUser) error {
	return m.err
}

func (m *MockCognito) ConfirmAccount(user *cognitoClient.UserConfirmation) error {
	return m.err
}

func (m *MockCognito) SignIn(user *cognitoClient.UserLogin) (string, error) {
	return m.token, m.err
}

func (m *MockCognito) GetUserByToken(token string) (*cognito.GetUserOutput, error) {
	return &m.user, m.err
}

func (m *MockCognito) UpdatePassword(user *cognitoClient.UserLogin) error {
	return m.err
}

func (m *MockCognito) DeleteUser(token string) error {
	return m.err
}

func TestHanler_GetMessages(t *testing.T) {
	type args struct {
		cognito     cognitoClient.CognitoInterface
		storage     message.Storage
		userStorage user.Storage
		r           func() *http.Request
	}

	tests := []struct {
		name           string
		args           args
		wantStatusCode int
	}{
		{
			name: "get_messages_returns_200",
			args: args{
				cognito:     &MockCognito{},
				storage:     &MockStorage{},
				userStorage: &MockUserStorage{},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, "/messages/username", nil)
					req.Header.Set("Authorization", "Bearer token")
					return req
				},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "get_messages_returns_400_when_not_auhtorized",
			args: args{
				cognito:     &MockCognito{},
				storage:     &MockStorage{},
				userStorage: &MockUserStorage{},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, "/messages/username", nil)
					return req
				},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "get_messages_returns_400_when_username_is_blank",
			args: args{
				cognito:     &MockCognito{},
				storage:     &MockStorage{},
				userStorage: &MockUserStorage{},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, "/messages/username", nil)
					return req
				},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "get_messages_returns_500_when_message_storage_misbehaves",
			args: args{
				cognito: &MockCognito{},
				storage: &MockStorage{
					err: errors.New("something's wrong"),
				},
				userStorage: &MockUserStorage{},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, "/messages/", nil)
					req.Header.Set("Authorization", "Bearer token")
					return req
				},
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name: "get_messages_returns_500_when_user_storage_misbehaves",
			args: args{
				cognito: &MockCognito{},
				storage: &MockStorage{},
				userStorage: &MockUserStorage{
					err: errors.New("something's wrong"),
				},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, "/messages/", nil)
					req.Header.Set("Authorization", "Bearer token")
					return req
				},
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messageHanlder := message.NewHandler(tt.args.storage, tt.args.userStorage, tt.args.cognito)
			handler := message.GetMessages(messageHanlder)
			w := httptest.NewRecorder()
			handler(w, tt.args.r())
			result := w.Result()
			if result.StatusCode != tt.wantStatusCode {
				t.Errorf("expected '%d' but got '%d'", tt.wantStatusCode, result.StatusCode)
			}
		})
	}
}
