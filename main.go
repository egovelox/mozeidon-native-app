package main

// TODO: add in readme
// To learn about native-messaging protocol (common to browsers like Chrome or Firefox)
// see https://developer.chrome.com/docs/extensions/develop/concepts/native-messaging#native-messaging-host-protocol

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/egovelox/mozeidon-native-app/models"
	"github.com/james-barrow/golang-ipc"
	"github.com/rickypc/native-messaging-host"
)

func main() {
	/*
		Listen to, and handle incoming ipc message :
		- forward each incoming message to the web-browser
		- send back each browser response as outgoing ipc message

		It acts like a web-browser proxy
	*/
	if err := webBrowserProxy(); err != nil {
		// write to stderr
		log.Printf("Error in mozeidon_native_app: %v", err)
	}
}

type IpcIncomingMessage struct {
	Command string `json:"command"        binding:"required"`
	Args    string `json:"args,omitempty"`
}

func webBrowserProxy() error {
	browserMessagingClient := (&host.Host{}).Init()

	// Step 1. Register this running native-app profile into the ProfileDirectory
	var nativeAppProfile *models.NativeAppProfile
	firstMessage := &host.H{}
	if err := browserMessagingClient.OnMessage(os.Stdin, firstMessage); err != nil {
		return fmt.Errorf("Error receiving message from browser: %w", err)
	}
	response, err := json.Marshal(firstMessage)
	if err != nil {
		return fmt.Errorf("error parsing registration response message: %w", err)
	}
	var registrationData models.RegistrationInfoResponse
	if err := json.Unmarshal(response, &registrationData); err != nil {
		return fmt.Errorf("error parsing registration message: %w", err)
	}
	nativeAppProfile = models.GetNativeAppProfile(&registrationData)

	profileDataDir, err := models.GetProfileDirectory()
	if err != nil {
		return fmt.Errorf("Error getting the profile directory: %w", err)
	}

	jsonProfile, err := json.MarshalIndent(nativeAppProfile, "", "  ")
	jsonProfilePath := filepath.Join(profileDataDir, nativeAppProfile.FileName)

	if err := os.WriteFile(jsonProfilePath, jsonProfile, 0644); err != nil {
		return fmt.Errorf("error writing profile file: %w", err)
	}

	/*
		On exits, unregister this running native-app profile from the ProfileDirectory
		for exits triggered by the browser-extension, handle SIGTERM sent from browser-extension.
		will not work for windows.
		@see https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/Native_messaging#closing_the_native_app
	*/
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sigChan
		os.Remove(jsonProfilePath)
		os.Exit(0)
	}()

	/*
		Also unregister the profile when the process ends
		with an error
	*/
	defer os.Remove(jsonProfilePath)

	// Step 2. Start IPC server
	ipcConfig := &ipc.ServerConfig{

		Encryption:        true, // allows encryption to be switched off (bool - default is true)
		UnmaskPermissions: true, // make the socket writeable for other users (default is false)
	}

	ipcServer, err := ipc.StartServer(nativeAppProfile.IpcName, ipcConfig)
	if err != nil {
		return fmt.Errorf("Error starting %s ipc-server: %w", nativeAppProfile.IpcName, err)
	}

	// Listen to client, and handle incoming message
	for {
		message, _ := ipcServer.Read()
		if message.MsgType > 0 {

			// Parse incoming message
			incomingMessage := IpcIncomingMessage{}
			json.Unmarshal(message.Data, &incomingMessage)

			// Send incoming message to browser
			request := &host.H{"payload": incomingMessage}
			if err := browserMessagingClient.PostMessage(os.Stdout, request); err != nil {
				return fmt.Errorf("Error posting message to browser: %w", err)
			}

			for {
				// Wait for browser messages
				// browser may send many messages before the data:end message
				response := &host.H{}
				if err := browserMessagingClient.OnMessage(os.Stdin, response); err != nil {
					return fmt.Errorf("Error receiving message from browser: %w", err)
				}

				// send back browser message to client
				responseMessage, _ := json.Marshal(response)
				err = ipcServer.Write(1, responseMessage)
				if err != nil {
					return fmt.Errorf("Error writing to ipc server: %w", err)
				}
				// end of browser response for the incoming message
				if string(responseMessage) == `{"data":"end"}` {
					break
				}
			}
		}
	}
}
