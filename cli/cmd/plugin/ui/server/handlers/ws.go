// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"syscall"

	"github.com/gorilla/websocket"

	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/log"
)

var upgrader = websocket.Upgrader{}

var (
	wsConnections []*websocket.Conn
	logs          = [][]byte{}
)

// InitWebsocketUpgrader initializes the upgrader and configures the
// CheckOrigin function using the provided host
func InitWebsocketUpgrader(hostBind string) {
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

// HandleWebsocketRequest handles the websocket request coming from clients
// upgrade normal http request to websocket request and stores the connection
func HandleWebsocketRequest(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.ForceWriteToStdErr([]byte(fmt.Sprintf("web socket upgrade error: %s", err.Error())))
		return
	}
	wsConnections = append(wsConnections, ws)

	log.ForceWriteToStdErr([]byte("web socket connection established\n"))

	sendPendingLogsOnConnection(ws)

	ws.SetCloseHandler(func(code int, text string) error {
		deleteWSConnection(ws)
		ws.Close()
		return nil
	})

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}

func sendPendingLogsOnConnection(ws *websocket.Conn) {
	log.ForceWriteToStdErr([]byte(fmt.Sprintf("sending pending %v logs to UI\n", len(logs))))
	for _, logMsg := range logs {
		err := ws.WriteMessage(1, logMsg)
		if err != nil {
			log.ForceWriteToStdErr([]byte("fail to write log message to web socket"))
			break
		}
	}

	logs = [][]byte{}
}

func deleteWSConnection(conn *websocket.Conn) {
	index := -1
	for i, ws := range wsConnections {
		if conn == ws {
			index = i
			break
		}
	}
	if index != -1 {
		wsConnections = append(wsConnections[:index], wsConnections[index+1:]...)
	}
}

// SendLog send the log message to all the connected websocket clients
func SendLog(logMsg []byte) {
	var err error

	if len(wsConnections) == 0 {
		logs = append(logs, logMsg)
		return
	}

	for _, ws := range wsConnections {
		err = ws.WriteMessage(1, logMsg)
		if err != nil {
			// when client connection is closed
			if errors.Is(err, syscall.EPIPE) || errors.Is(err, syscall.ECONNRESET) {
				deleteWSConnection(ws)
				continue
			}

			log.ForceWriteToStdErr([]byte("fail to write log message to web socket"))
			break
		}
	}
}
