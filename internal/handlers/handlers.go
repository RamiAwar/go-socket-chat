package handlers

import (
	"log"
	"net/http"
	"sort"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
)

var wsChan = make(chan WsPayload)

var clients = make(map[WebSocketConnection]string)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(),
)

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Home(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "home.jet", nil)
	if err != nil {
		log.Println(err)
	}
}

type WebSocketConnection struct {
	*websocket.Conn
}

// WsJsonResponse defines the response sent back from websokcet
type WsJsonResponse struct {
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	User           string   `json:"user"`
	MessageType    string   `json:"message_type"`
	ConnectedUsers []string `json:"connected_users"`
}

type WsPayload struct {
	Action   string              `json:"action"`
	Username string              `json:"username"`
	Message  string              `json:"message"`
	Conn     WebSocketConnection `json:"-"`
}

// Connect upgrades connection to websocket
func Connect(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client connected")

	var response WsJsonResponse
	response.Message = `<em><small>Connected to server</small></em>`

	// Add to list of clients
	conn := WebSocketConnection{Conn: ws}
	clients[conn] = ""

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err)
	}

	go ListenForWs(&conn)
}

func ListenForWs(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("error: ", r)
		}
	}()

	var payload WsPayload
	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
		} else {
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}

func ListenToWsChannel() {
	var res WsJsonResponse

	for {
		e := <-wsChan

		switch e.Action {

		case "username":
			// Send back list of all users
			clients[e.Conn] = e.Username
			users := getUserList()

			res.Action = "list_users"
			res.ConnectedUsers = users
			broadcast(res)

		case "left":
			delete(clients, e.Conn)
			users := getUserList()

			res.Action = "list_users"
			res.ConnectedUsers = users
			broadcast(res)

		case "get_users":
			res.Action = "list_users"
			res.ConnectedUsers = getUserList()
			send(res, &e.Conn)

		case "broadcast":
			res.Message = e.Message
			res.Action = "broadcast"
			res.User = e.Username
			broadcast(res)
		}
	}
}

func getUserList() []string {
	userList := make([]string, 0)
	for _, x := range clients {
		if len(x) > 0 { // Only append users with username
			userList = append(userList, x)
		}
	}

	sort.Strings(userList)
	return userList
}

func send(response WsJsonResponse, conn *WebSocketConnection) {
	err := conn.WriteJSON(response)
	if err != nil {
		log.Println("Websocket err", err)
		_ = conn.Close()
		delete(clients, *conn)

		response.ConnectedUsers = getUserList()
		broadcast(response)
	}
}

func broadcast(response WsJsonResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println("Websocket err", err)
			_ = client.Close()
			delete(clients, client)
		}
	}
}

func renderPage(w http.ResponseWriter, template string, variables jet.VarMap) error {
	view, err := views.GetTemplate(template)
	if err != nil {
		log.Println(err)
		return err
	}

	err = view.Execute(w, variables, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
