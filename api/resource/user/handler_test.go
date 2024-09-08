package user_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thaironsilva/messenger/api/resource/user"
)

type MockStorage struct {
	err   error
	user  user.User
	users []user.User
}

func (m *MockStorage) GetAll() ([]user.User, error) {
	return m.users, m.err
}

func (m *MockStorage) Create(user user.User) error {
	return m.err
}

func TestHanler_GetAll(t *testing.T) {
	type args struct {
		storage user.Storage
		r       func() *http.Request
	}

	tests := []struct {
		name           string
		args           args
		wantStatusCode int
	}{
		{
			name: "get_all_returns_200",
			args: args{
				storage: &MockStorage{},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, "/", bytes.NewReader([]byte{}))
					return req
				},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "get_all_returns_500_when_storage_misbehaves",
			args: args{
				storage: &MockStorage{
					err: errors.New("something's wrong"),
				},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, "/", bytes.NewReader([]byte{}))
					return req
				},
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := user.GetUsers(tt.args.storage)
			w := httptest.NewRecorder()
			handler(w, tt.args.r())
			result := w.Result()
			if result.StatusCode != tt.wantStatusCode {
				t.Errorf("expected '%d' but got '%d'", tt.wantStatusCode, result.StatusCode)
			}
		})
	}
}

func TestHanler_Create(t *testing.T) {
	type args struct {
		storage user.Storage
		r       func() *http.Request
	}

	tests := []struct {
		name           string
		args           args
		wantStatusCode int
	}{
		{
			name: "create_returns_200",
			args: args{
				storage: &MockStorage{},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"username":"john","email":"johndoe@email.com","password":"helloworld"}`)))
					return req
				},
			},
			wantStatusCode: http.StatusCreated,
		},
		{
			name: "create_returns_500_when_storage_misbehaves",
			args: args{
				storage: &MockStorage{
					err: errors.New("something's wrong"),
				},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"username":"john","email":"johndoe@email.com","password":"helloworld"}`)))
					return req
				},
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name: "create_returns_400_when_request_body_is_invalid",
			args: args{
				storage: &MockStorage{},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodPost, "/", nil)
					return req
				},
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := user.CreateUsers(tt.args.storage)
			w := httptest.NewRecorder()
			handler(w, tt.args.r())
			result := w.Result()
			if result.StatusCode != tt.wantStatusCode {
				t.Errorf("expected '%d' but got '%d'", tt.wantStatusCode, result.StatusCode)
			}
		})
	}
}
