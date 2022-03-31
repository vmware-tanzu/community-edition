// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package server provides backend api for UI
package server

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-openapi/loads"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations"
)

// Content is set from main since we can't embed a path starting with ../
// We may want to move the web folder under server to avoid this.
var Content embed.FS

// Serve provides the backend REST API for the UI.
func Serve(bind, browser string) error {
	swaggerSpec, err := loads.Analyzed(restapi.FlatSwaggerJSON, "2.0")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewTanzuUIAPI(swaggerSpec)
	server := restapi.NewServer(api)

	server.EnabledListeners = []string{"http"}
	host, port, err := net.SplitHostPort(bind)
	if err != nil {
		return errors.New("invalid binding address provided, please use address in the the form '127.0.0.1:8080'")
	}
	server.Port, err = strconv.Atoi(port)
	if err != nil {
		return errors.New("invalid binding port provided, please provide a valid number (e.g. 8080)")
	}
	server.Host = host
	server.Browser = browser

	// TODO: Need to determine what we actually need here.
	// ws.InitWebsocketUpgrader(server.Host)

	// TODO: Define and wire up handlers for each API call processing.
	// app := &handlers.App{InitOptions: initOptions, AppConfig: appConfig, TKGTimeout: tkgTimeOut, TKGConfigReaderWriter: tkgConfigReaderWriter}
	// app.ConfigureHandlers(api)
	server.SetAPI(api)
	server.SetHandler(api.Serve(FileServerMiddleware))

	// Configure out static page handling for the UI
	server.SetHandler(http.HandlerFunc(staticHandler))

	// check if the port is already in use, if so exit gracefully
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Host, server.Port))
	if err != nil {
		server.Logf("Failed to start the kickstart UI Server[Host:%s, Port:%d], error: %s\n", server.Host, server.Port, err)
		os.Exit(1)
	}
	l.Close()

	defer func() {
		err := server.Shutdown()
		if err != nil {
			fmt.Printf("Error shutting down server: %s\n", err.Error())
		}
	}()

	mux := http.NewServeMux()
	mux.Handle("/api", http.HandlerFunc(staticHandler))

	if err := server.Serve(); err != nil {
		return err
	}
	return nil
}

// FileServerMiddleware serves ui resource
func FileServerMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/ws") {
			// TODO: Check what we need here.
			// ws.HandleWebsocketRequest(w, r)
		} else if strings.HasPrefix(r.URL.Path, "/api") {
			handler.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, "/ui", http.StatusMovedPermanently)
		}
	})
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	if strings.HasSuffix(r.URL.Path, ".css") {
		w.Header().Add("Content-Type", "text/css")
	}

	// get static content from go embed
	fsys := fs.FS(Content)
	staticContent, _ := fs.Sub(fsys, "web/tanzu-ui/build")

	http.StripPrefix("/ui", Logger(http.FileServer(http.FS(staticContent)), "UI")).ServeHTTP(w, r)
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s %s %s %s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
