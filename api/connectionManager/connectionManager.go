package connectionManager

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/thaironsilva/messenger/api/cognitoClient"
	"github.com/thaironsilva/messenger/api/resource/user"

	"github.com/gorilla/websocket"
)

type ConnectionHandler struct {
	userRepository user.Storage
	cognito        cognitoClient.CognitoInterface
	Clients        map[string]*websocket.Conn
	Channels       map[string]chan string
	mu             sync.Mutex
}

func NewConnectionHandler(userRepository user.Storage, cognito cognitoClient.CognitoInterface) *ConnectionHandler {
	return &ConnectionHandler{
		userRepository: userRepository,
		cognito:        cognito,
		Clients:        make(map[string]*websocket.Conn),
		Channels:       make(map[string]chan string),
	}
}

func (h *ConnectionHandler) HandleConnections(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if token == "" {
		fmt.Println("message: unauthorized")
		return
	}

	cognitoUser, err := h.cognito.GetUserByToken(token)
	if err != nil {
		fmt.Println(err)
		return
	}

	var email string

	for _, attribute := range cognitoUser.UserAttributes {
		if *attribute.Name == "email" {
			email = *attribute.Value
		}
	}

	sender, err := h.userRepository.GetByEmail(email)
	if err != nil {
		fmt.Println(err)
		return
	}

	username := strings.TrimPrefix(r.URL.Path, "/messages/")
	if username == "" {
		fmt.Println("message: not found")
		return
	}

	receiver, err := h.userRepository.GetByUsername(username)
	if err != nil {
		fmt.Println(err)
		return
	}

	h.mu.Lock()
	h.Clients[sender.Username+"-"+receiver.Username] = conn
	h.Channels[receiver.Username+"-"+sender.Username] = make(chan string)
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.Clients, sender.Username+"-"+receiver.Username)
		delete(h.Channels, receiver.Username+"-"+sender.Username)
		h.mu.Unlock()
	}()

	// receive messages
	go func() {
		for {
			receiveChannel, open := h.Channels[receiver.Username+"-"+sender.Username]
			if !open {
				return
			}
			msg := <-receiveChannel
			myConn, connected := h.Clients[sender.Username+"-"+receiver.Username]
			if !connected {
				return
			}

			err := myConn.WriteJSON(msg)
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	// send messages
	for {
		var msg string
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println(err)
			return
		}

		sendChannel, open := h.Channels[sender.Username+"-"+receiver.Username]
		if open {
			sendChannel <- msg
		}
	}
}
