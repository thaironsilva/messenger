package connectionManager_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/thaironsilva/messenger/api/cognitoClient"
	"github.com/thaironsilva/messenger/api/connectionManager"
	"github.com/thaironsilva/messenger/api/resource/message"
	"github.com/thaironsilva/messenger/api/resource/user"

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type MockMessageStorage struct {
	err      error
	messages []message.Message
}

func (m *MockMessageStorage) GetAll(sender_id string, receiver_id string) ([]message.Message, error) {
	return m.messages, m.err
}

func (m *MockMessageStorage) Create(message message.Message) error {
	return m.err
}

type MockUserStorage struct {
	err   error
	user  user.User
	users []user.User
}

func (m *MockUserStorage) GetByUsername(username string) (user.User, error) {
	if username == "user1" {
		return user.User{Username: "user1"}, nil
	}
	if username == "user2" {
		return user.User{Username: "user2"}, nil
	}
	return m.user, m.err
}

func (m *MockUserStorage) GetByEmail(email string) (user.User, error) {
	if email == "email1" {
		return user.User{Username: "user1"}, nil
	}
	if email == "email2" {
		return user.User{Username: "user2"}, nil
	}
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
	if token == "token1" {
		name := "email"
		value := "email1"
		cogUser := &cognito.GetUserOutput{}
		cogUser.SetUserAttributes([]*cognito.AttributeType{{Name: &name, Value: &value}})
		return cogUser, nil
	}
	if token == "token2" {
		name := "email"
		value := "email2"
		cogUser := &cognito.GetUserOutput{}
		cogUser.SetUserAttributes([]*cognito.AttributeType{{Name: &name, Value: &value}})
		return cogUser, nil
	}
	return &m.user, m.err
}

func (m *MockCognito) UpdatePassword(user *cognitoClient.UserLogin) error {
	return m.err
}

func (m *MockCognito) DeleteUser(token string) error {
	return m.err
}

func TestConnectionManager_testHandleConnections(t *testing.T) {
	t.Run("stabishes_double_sided_connection_and_exchange_messages", func(t *testing.T) {
		wantCount := 100
		connHandler := connectionManager.NewConnectionHandler(&MockMessageStorage{}, &MockUserStorage{}, &MockCognito{})
		s := httptest.NewServer(http.HandlerFunc(connHandler.HandleConnections))
		defer s.Close()

		u := "ws" + strings.TrimPrefix(s.URL, "http") + "/api/v0/chat/user2"

		header := http.Header{}
		header.Set("Authorization", "Bearer token1")
		ws1, _, err := websocket.DefaultDialer.DialContext(context.TODO(), u, header)
		if err != nil {
			t.Fatalf("%v", err)
		}

		defer ws1.Close()

		u = "ws" + strings.TrimPrefix(s.URL, "http") + "/api/v0/chat/user1"

		header = http.Header{}
		header.Set("Authorization", "Bearer token2")
		ws2, _, err := websocket.DefaultDialer.DialContext(context.TODO(), u, header)
		if err != nil {
			t.Fatalf("%v", err)
		}

		defer ws2.Close()

		go func() {
			for i := 0; i < 100; i++ {
				ws1.WriteJSON("test message 1")
			}
		}()

		go func() {
			for i := 0; i < 100; i++ {
				ws2.WriteJSON("test message 2")
			}
		}()

		var wg1 sync.WaitGroup
		wg1.Add(wantCount)

		counter1 := 0
		go func() {
			for {
				var receive string
				ws1.ReadJSON(&receive)
				switch string(receive) {
				case "test message 2":
					counter1++
					wg1.Done()
				default:
					t.Errorf("Received unexpected message: %s", string(receive))
				}
				if counter1 == 100 {
					return
				}
			}
		}()
		wg1.Wait()

		var wg2 sync.WaitGroup
		wg2.Add(wantCount)

		counter2 := 0
		go func() {
			for {
				var receive string
				ws2.ReadJSON(&receive)
				switch string(receive) {
				case "test message 1":
					counter2++
					wg2.Done()
				default:
					t.Errorf("Received unexpected message: %s", string(receive))
				}
				if counter2 == 100 {
					return
				}
			}
		}()
		wg2.Wait()

		if counter1 != wantCount {
			t.Errorf("expected '%d' but got '%d'", wantCount, counter1)
		}
		if counter2 != wantCount {
			t.Errorf("expected '%d' but got '%d'", wantCount, counter2)
		}
	})

	t.Run("establishes_one_sided_connection_and_dont_fail", func(t *testing.T) {
		wantCount := 100
		connHandler := connectionManager.NewConnectionHandler(&MockMessageStorage{}, &MockUserStorage{}, &MockCognito{})
		s := httptest.NewServer(http.HandlerFunc(connHandler.HandleConnections))
		defer s.Close()

		u := "ws" + strings.TrimPrefix(s.URL, "http") + "/messages/user2"

		header := http.Header{}
		header.Set("Authorization", "Bearer token1")
		ws, _, err := websocket.DefaultDialer.DialContext(context.TODO(), u, header)
		if err != nil {
			t.Fatalf("%v", err)
		}

		defer ws.Close()

		go func() {
			for i := 0; i < 10; i++ {
				ws.WriteJSON("test message 1")
			}
		}()

		var wg sync.WaitGroup
		wg.Add(wantCount)
		go func() {
			for i := 0; i < wantCount; i++ {
				var receive string
				ws.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
				ws.ReadJSON(&receive)
				if string(receive) != "" {
					t.Errorf("Expected no message but got '%s'", receive)
				}
				wg.Done()
			}
		}()
		wg.Wait()
	})
}
