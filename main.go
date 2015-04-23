package main

import (
  "code.google.com/p/go.net/websocket"
  "html/template"
  "log"
  "net/http"
  "os"
)

var (
  pwd, _        = os.Getwd()
  RootTemp      = template.Must(template.ParseFiles(pwd + "/main.html"))
  JSON          = websocket.JSON
  Message       = websocket.Message
  ActiveClients = make(map[ClientConn]int)
)

// Initialize handlers
func init() {
  http.HandleFunc("/", RootHandler)
  // Get static CSS and JavaScript files
  http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
  http.Handle("/socket", websocket.Handler(SocketServer))
}

// Client connection consists of the WS and client IP
type ClientConn struct {
  websocket *websocket.Conn
  clientIP  string
}

// WS server to handle chat between clients
func SocketServer(ws *websocket.Conn) {
  var err error
  var clientMessage string

  // Executes after surrounding func returns
  defer func() {
    if err = ws.Close(); err != nil {
      log.Println("WebSocket could not be closed.", err.Error())
    }
  }()

  client := ws.Request().RemoteAddr
  log.Println("Client connected: ", client)
  socketCli := ClientConn{ws, client}
  ActiveClients[socketCli] = 0
  log.Println("Number of clients currently connected: ", len(ActiveClients))

  // Keeps WebSocket open, rather than closing after one Send/Receieve
  for {
    if err = Message.Receive(ws, &clientMessage); err != nil {
      // Connection closes if there is an error reading
      log.Println("Websocket disconnected", err.Error())
      // Remove active clients
      delete(ActiveClients, socketCli)
      log.Println("Number of clients currently connected: ", len(ActiveClients))
      return
    }

    clientMessage = socketCli.clientIP + " Said: " + clientMessage
    for cs, _ := range ActiveClients {
      if err = Message.Send(cs.websocket, clientMessage); err != nil {
        log.Println("Could not send message to ", cs.clientIP, err.Error())
      }
    }
  }
}

// Renders the template for the root page
func RootHandler(w http.ResponseWriter, req *http.Request) {
  err := RootTemp.Execute(w, ":"+os.Getenv("PORT"))
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

// The main fuction to setup routes and start the web server
func main() {
  // Start an HTTP server with given address and DefaultServeMux handler
  // DefaultServeMux is like a HTTP request router that is instantiated by default when the HTTP package is used
  err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
  if err != nil {
    // panic creates run-time error and stops the program
    panic("ListenAndServe: " + err.Error())
  }
}
