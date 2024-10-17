package user_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/thaironsilva/messenger/api/cognitoClient"
	"github.com/thaironsilva/messenger/api/resource/user"
)

type MockStorage struct {
	err   error
	user  user.User
	users []user.User
}

func (m *MockStorage) GetByUsername(username string) (user.User, error) {
	return m.user, m.err
}

func (m *MockStorage) GetByEmail(email string) (user.User, error) {
	return m.user, m.err
}

func (m *MockStorage) GetByString(name string) ([]user.User, error) {
	return m.users, m.err
}

func (m *MockStorage) GetAll() ([]user.User, error) {
	return m.users, m.err
}

func (m *MockStorage) Create(user user.User) error {
	return m.err
}

func (m *MockStorage) Update(user user.User) error {
	return m.err
}

func (m *MockStorage) Delete(id string) error {
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

func TestHanler_GetUsers(t *testing.T) {
	type args struct {
		cognito cognitoClient.CognitoInterface
		storage user.Storage
		r       func() *http.Request
	}

	tests := []struct {
		name           string
		args           args
		wantStatusCode int
	}{
		{
			name: "get_users_returns_200",
			args: args{
				cognito: &MockCognito{},
				storage: &MockStorage{},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, "/users/", nil)
					req.Header.Set("Authorization", "Bearer token")
					return req
				},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "get_users_returns_400_when_not_auhtorized",
			args: args{
				cognito: &MockCognito{},
				storage: &MockStorage{},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, "/users/", nil)
					return req
				},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "get_users_returns_500_when_storage_misbehaves",
			args: args{
				cognito: &MockCognito{},
				storage: &MockStorage{
					err: errors.New("something's wrong"),
				},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, "/users/", nil)
					req.Header.Set("Authorization", "Bearer token")
					return req
				},
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userHanlder := user.NewHandler(tt.args.storage, tt.args.cognito)
			handler := user.GetUsers(userHanlder)
			w := httptest.NewRecorder()
			handler(w, tt.args.r())
			result := w.Result()
			if result.StatusCode != tt.wantStatusCode {
				t.Errorf("expected '%d' but got '%d'", tt.wantStatusCode, result.StatusCode)
			}
		})
	}
}

func TestHanler_CreateUser(t *testing.T) {
	type args struct {
		cognito cognitoClient.CognitoInterface
		storage user.Storage
		r       func() *http.Request
	}

	tests := []struct {
		name           string
		args           args
		wantStatusCode int
	}{
		{
			name: "create_returns_201",
			args: args{
				cognito: &MockCognito{},
				storage: &MockStorage{},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodPost, "/users/", bytes.NewReader([]byte(`{"nickname":"john","email":"johndoe@email.com","password":"helloworld"}`)))
					return req
				},
			},
			wantStatusCode: http.StatusCreated,
		},
		{
			name: "create_returns_400_when_cognito_misbehaves",
			args: args{
				cognito: &MockCognito{err: errors.New("something's wrong")},
				storage: &MockStorage{},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodPost, "/users/", bytes.NewReader([]byte(`{"nickname":"john","email":"johndoe@email.com","password":"helloworld"}`)))
					return req
				},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "create_returns_400_when_request_body_is_invalid",
			args: args{
				cognito: &MockCognito{},
				storage: &MockStorage{},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodPost, "/users/", nil)
					return req
				},
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userHandler := user.NewHandler(tt.args.storage, tt.args.cognito)
			handler := user.CreateUser(userHandler)
			w := httptest.NewRecorder()
			handler(w, tt.args.r())
			result := w.Result()
			if result.StatusCode != tt.wantStatusCode {
				t.Errorf("expected '%d' but got '%d'", tt.wantStatusCode, result.StatusCode)
			}
		})
	}
}

func TestHanler_DeleteUser(t *testing.T) {
	type args struct {
		cognito cognitoClient.CognitoInterface
		storage user.Storage
		r       func() *http.Request
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
	}{
		{
			name: "delete_returns_200",
			args: args{
				cognito: &MockCognito{},
				storage: &MockStorage{},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodDelete, "/users/id", nil)
					req.Header.Set("Authorization", "Bearer token")
					req.SetPathValue("id", "id")
					return req
				},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "delete_returns_400_when_not_authorized",
			args: args{
				cognito: &MockCognito{},
				storage: &MockStorage{},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodDelete, "/users/id", nil)
					req.SetPathValue("id", "id")
					return req
				},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "delete_returns_500_when_storage_misbehaves",
			args: args{
				cognito: &MockCognito{},
				storage: &MockStorage{
					err: errors.New("something's wrong"),
				},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodDelete, "/users/id", nil)
					req.Header.Set("Authorization", "Bearer token")
					req.SetPathValue("id", "id")
					return req
				},
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userHanlder := user.NewHandler(tt.args.storage, tt.args.cognito)
			handler := user.DeleteUser(userHanlder)
			w := httptest.NewRecorder()
			handler(w, tt.args.r())
			result := w.Result()
			if result.StatusCode != tt.wantStatusCode {
				t.Errorf("expected '%d' but got '%d'", tt.wantStatusCode, result.StatusCode)
			}
		})
	}
}
