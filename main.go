package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"path"
	"time"

	"github.com/gorilla/websocket"
)

var port = "5000"

type SubRouter struct{}

type todoType struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Done        bool   `json:"done"`
	Order       int    `json:"order"`
}

var toDoList = map[string]todoType{}

var upgrader = websocket.Upgrader{
	CheckOrigin: CheckOriginChecker,
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

var connList []*websocket.Conn

func broadcastJson(data interface{}) {
	for conn := range connList {
		connList[conn].WriteJSON(data)
	}
}

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func CheckOriginChecker(r *http.Request) bool {
	return true
}

func socketHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("/socket")
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		fmt.Printf("WebSocket upgrade error: %s\n", err)
		return
	}
	defer conn.Close()
	connList = append(connList, conn)

	conn.WriteJSON(map[string]string{
		"type": "clear",
	})

	for item := range toDoList {
		conn.WriteJSON(map[string]todoType{item: toDoList[item]})
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("message error: %s\n", err)
			break
		}
		input := string(message)
		var inputData map[string]string
		json.Unmarshal([]byte(input), &inputData)
		fmt.Println(inputData)
		switch messageType := inputData["type"]; messageType {
		case "add":
			key := RandStringBytes(10)
			current_time := time.Now().Format("02/01/2006 3:04:05PM")
			todo := todoType{
				Name:        inputData["name"],
				Order:       1,
				Done:        false,
				Date:        current_time,
				Description: inputData["description"],
			}
			toDoList[key] = todo
			broadcastJson(map[string]todoType{key: todo})
		case "remove":
			delete(toDoList, inputData["key"])
			broadcastJson(map[string]string{
				"type": "remove",
				"key":  inputData["key"],
			})
		case "done":
			tempToDo := toDoList[inputData["key"]]
			tempToDo.Done = true
			toDoList[inputData["key"]] = tempToDo
			broadcastJson(map[string]todoType{inputData["key"]: tempToDo})
		}
	}
}

func (r *SubRouter) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var assetsPath string
	UrlPath := request.URL.Path
	switch messageType := UrlPath; messageType {
	case "/":
		assetsPath = "web/build/index.html"
	default:
		assetsPath = fmt.Sprintf("web/build%s", path.Clean(messageType))
	}
	fmt.Println(assetsPath)
	http.ServeFile(writer, request, assetsPath)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/socket", socketHandler)
	mux.Handle("/", &SubRouter{})

	fmt.Printf("Server: http://localhost:%s", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Server Shutdown")
	} else if err != nil {
		fmt.Printf("---Error---\n%s", err)
	}
}
