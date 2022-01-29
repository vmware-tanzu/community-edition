// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	"github.com/spf13/cobra"

	cliv1alpha1 "github.com/vmware-tanzu/tanzu-framework/apis/cli/v1alpha1"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/buildinfo"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/command/plugin"
)

// Variables injected by ldflags by the bundler.
var (
	AppName            string
	BuiltAt            string
	VersionAstilectron string
	VersionElectron    string
)

var descriptor = cliv1alpha1.PluginDescriptor{
	Name:        "installer",
	Description: "Launch the Tanzu installer.",
	Group:       cliv1alpha1.RunCmdGroup,
	Version:     buildinfo.Version,
}

var (
	// logLevel is the verbosity to print
	logLevel int32

	// debug controls whether to log debug level messages
	debug bool

	// w is a pointer to the UI window
	w *astilectron.Window

	// default build version
	defaultVersion = "v0.0.1-unversioned"
)

func main() {
	if descriptor.Version == "" {
		descriptor.Version = defaultVersion
	}

	p, err := plugin.NewPlugin(&descriptor)
	if err != nil {
		log.Fatal(err, "unable to initialize new plugin")
	}

	p.Cmd.PersistentFlags().Int32VarP(&logLevel, "verbose", "v", 0, "Number for the log level verbosity (0-9)")
	p.Cmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug level logs")
	p.Cmd.Run = launch

	if err := p.Execute(); err != nil {
		os.Exit(1)
	}
}

func launch(cmd *cobra.Command, args []string) {
	l := log.New(log.Writer(), log.Prefix(), log.Flags())

	// Run bootstrap
	l.Printf("Running app built at %s\n", BuiltAt)
	if err := bootstrap.Run(bootstrap.Options{
		Asset:    Asset,
		AssetDir: AssetDir,
		AstilectronOptions: astilectron.Options{
			AppName:            AppName,
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
			SingleInstance:     true,
			VersionAstilectron: VersionAstilectron,
			VersionElectron:    VersionElectron,
		},
		Debug:  debug,
		Logger: l,
		MenuOptions: []*astilectron.MenuItemOptions{{
			Label: astikit.StrPtr("File"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Label: astikit.StrPtr("About"),
					OnClick: func(e astilectron.Event) (deleteListener bool) {
						if err := bootstrap.SendMessage(w, "about", "TCE UI - under development", func(m *bootstrap.MessageIn) {
							// Unmarshal payload
							var s string
							if err := json.Unmarshal(m.Payload, &s); err != nil {
								l.Println(fmt.Errorf("unmarshaling payload failed: %w", err))
								return
							}
							l.Printf("About modal has been displayed and payload is %s!\n", s)
						}); err != nil {
							l.Println(fmt.Errorf("sending about event failed: %w", err))
						}
						return
					},
				},
				{Role: astilectron.MenuItemRoleClose},
			},
		}},
		OnWait: func(_ *astilectron.Astilectron, ws []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
			w = ws[0]
			go func() {
				time.Sleep(5 * time.Second)
				if err := bootstrap.SendMessage(w, "check.out.menu", "Don't forget to check out the menu!"); err != nil {
					l.Println(fmt.Errorf("sending check.out.menu event failed: %w", err))
				}
			}()
			return nil
		},
		RestoreAssets: RestoreAssets,
		Windows: []*bootstrap.Window{{
			Homepage:       "index.html",
			MessageHandler: handleMessages,
			Options: &astilectron.WindowOptions{
				BackgroundColor: astikit.StrPtr("#333"),
				Center:          astikit.BoolPtr(true),
				Height:          astikit.IntPtr(700),
				Width:           astikit.IntPtr(700),
			},
		}},
	}); err != nil {
		l.Println(fmt.Errorf("running bootstrap failed: %w", err))
	}
}

// handleMessages handles messages
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "explore":
		// Unmarshal payload
		var path string
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &path); err != nil {
				payload = err.Error()
				return
			}
		}

		// Explore
		// if payload, err = explore(path); err != nil {
		// 	payload = err.Error()
		// 	return
		// }

	case "event.name":
		// Unmarshal payload
		var s string
		if err = json.Unmarshal(m.Payload, &s); err != nil {
			payload = err.Error()
			return
		}
		payload = s + " world"
	}
	return
}

// func Asset(name string) ([]byte, error) {
// 	return []byte{}, nil
// }

// func AssetDir(name string) ([]string, error) {
// 	return []string{}, nil
// }

// func RestoreAssets(dir, name string) error {
// 	return nil
// }
