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

func (m *MockStorage) GetById(id string) (user.User, error) {
	return m.user, m.err
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

func TestHanler_GetUser(t *testing.T) {
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
			name: "get_by_id_returns_200",
			args: args{
				storage: &MockStorage{},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, "/users/id", nil)
					req.SetPathValue("id", "id")
					return req
				},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "get_by_id_returns_500_when_storage_misbehaves",
			args: args{
				storage: &MockStorage{
					err: errors.New("something's wrong"),
				},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, "/users/id", nil)
					req.SetPathValue("id", "id")
					return req
				},
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name: "get_by_id_returns_404_when_user_doesnt_exist",
			args: args{
				storage: &MockStorage{
					err: errors.New("sql: no rows in result set"),
				},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, "/users/id", nil)
					req.SetPathValue("id", "id")
					return req
				},
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "get_by_id_returns_400_when_id_is_empty",
			args: args{
				storage: &MockStorage{},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, "/users/", nil)
					return req
				},
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := user.GetUser(tt.args.storage)
			w := httptest.NewRecorder()
			handler(w, tt.args.r())
			result := w.Result()
			if result.StatusCode != tt.wantStatusCode {
				t.Errorf("expected '%d' but got '%d'", tt.wantStatusCode, result.StatusCode)
			}
		})
	}
}

func TestHanler_GetUsers(t *testing.T) {
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
					req, _ := http.NewRequest(http.MethodGet, "/users/", nil)
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
					req, _ := http.NewRequest(http.MethodGet, "/users/", nil)
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

func TestHanler_CreateUser(t *testing.T) {
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
			name: "create_returns_201",
			args: args{
				storage: &MockStorage{},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodPost, "/users/", bytes.NewReader([]byte(`{"username":"john","email":"johndoe@email.com","password":"helloworld"}`)))
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
					req, _ := http.NewRequest(http.MethodPost, "/users/", bytes.NewReader([]byte(`{"username":"john","email":"johndoe@email.com","password":"helloworld"}`)))
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
					req, _ := http.NewRequest(http.MethodPost, "/users/", nil)
					return req
				},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "create_returns_400_when_request_params_are_invalid",
			args: args{
				storage: &MockStorage{},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodPost, "/users/", bytes.NewReader([]byte(`{"password":"helloworld"}`)))
					return req
				},
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := user.CreateUser(tt.args.storage)
			w := httptest.NewRecorder()
			handler(w, tt.args.r())
			result := w.Result()
			if result.StatusCode != tt.wantStatusCode {
				t.Errorf("expected '%d' but got '%d'", tt.wantStatusCode, result.StatusCode)
			}
		})
	}
}

func TestHanler_UpdateUser(t *testing.T) {
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
			name: "update_returns_200",
			args: args{
				storage: &MockStorage{},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodPut, "/users/id", bytes.NewReader([]byte(`{"username":"john","email":"johndoe@email.com","password":"helloworld"}`)))
					req.SetPathValue("id", "id")
					return req
				},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "update_returns_500_when_storage_misbehaves",
			args: args{
				storage: &MockStorage{
					err: errors.New("something's wrong"),
				},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodPut, "/users/id", bytes.NewReader([]byte(`{"username":"john","email":"johndoe@email.com","password":"helloworld"}`)))
					req.SetPathValue("id", "id")
					return req
				},
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name: "update_returns_400_when_request_body_is_invalid",
			args: args{
				storage: &MockStorage{},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodPut, "/users/id", nil)
					req.SetPathValue("id", "id")
					return req
				},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "update_returns_400_when_id_is_empty",
			args: args{
				storage: &MockStorage{},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodPut, "/users/", bytes.NewReader([]byte(`{"username":"john","email":"johndoe@email.com","password":"helloworld"}`)))
					return req
				},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "update_returns_400_when_request_params_are_invalid",
			args: args{
				storage: &MockStorage{},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodPut, "/users/id", bytes.NewReader([]byte(`{"password":"helloworld"}`)))
					req.SetPathValue("id", "id")
					return req
				},
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := user.UpdateUser(tt.args.storage)
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
				storage: &MockStorage{},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodDelete, "/users/id", nil)
					req.SetPathValue("id", "id")
					return req
				},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "delete_returns_500_when_storage_misbehaves",
			args: args{
				storage: &MockStorage{
					err: errors.New("something's wrong"),
				},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodDelete, "/users/id", nil)
					req.SetPathValue("id", "id")
					return req
				},
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name: "delete_returns_400_when_id_is_empty",
			args: args{
				storage: &MockStorage{},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodDelete, "/users/", nil)
					return req
				},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "get_by_id_returns_404_when_user_doesnt_exist",
			args: args{
				storage: &MockStorage{
					err: errors.New("sql: no rows in result set"),
				},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodDelete, "/users/id", nil)
					req.SetPathValue("id", "id")
					return req
				},
			},
			wantStatusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := user.DeleteUser(tt.args.storage)
			w := httptest.NewRecorder()
			handler(w, tt.args.r())
			result := w.Result()
			if result.StatusCode != tt.wantStatusCode {
				t.Errorf("expected '%d' but got '%d'", tt.wantStatusCode, result.StatusCode)
			}
		})
	}
}
