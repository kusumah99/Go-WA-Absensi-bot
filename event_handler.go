package main

import (
	"context"
	"fmt"
	"ksa/dbservice"
	"ksa/msghandler"
	"os"
	"time"

	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"rsc.io/qr"
)

var WaClient *whatsmeow.Client
var dbContainer *sqlstore.Container
var deviceStore *store.Device
var dbLog waLog.Logger

func init() {
	dbLog = waLog.Stdout("Database", "DEBUG", true)
}

func EventsHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.LoggedOut:
		// delay 10 detik, request qrcode ulang pake thread berbeda, lempar ka browser
		fmt.Println("Logged out, attempting to relogin in 10 seconds....")
		time.Sleep(10 * time.Second)
		go func() {
			WaLogin()
		}()
		break
	case *events.Disconnected:
		// delay 10 detik, request qrcode ulang pake thread berbeda, lempar ka browser
		fmt.Println("Disconnected, attempting to reconnect in 10 seconds...")
		time.Sleep(10 * time.Second)
		go func() {
			WaLogin()
		}()
		break
	case *events.Message:

		msghandler.OnMessage(v)

		// if !v.Info.IsFromMe {
		// 	fmt.Println("Pesan: " + v.Message.String())
		// 	// targetJID, err := types.ParseJID(v.Info.ID)
		// 	targetJID, _ := types.ParseJID(v.Info.ID)
		// 	client.SendMessage(context.Background(), targetJID, &waProto.Message{
		// 		Conversation: proto.String("Ini pesan balasan otomatis dari pesan Anda: " + v.Message.GetConversation()),
		// 	})
		// }
	}
}

func WaDisconnect() {
	WaClient.Disconnect()
}

func WaLogin() {
	var err error
	if WaClient != nil && WaClient.IsConnected() {
		WaClient.Disconnect()
	}
	WaClient = nil // garbage
	if deviceStore != nil {
		// deviceStore.DeleteAllSessions()
		deviceStore = nil // garbage
	}
	if dbContainer != nil {
		dbContainer.Close()
		dbContainer = nil // garbage
	}
	// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
	dbContainer, err = sqlstore.New("sqlite3", "file:absensi.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err = dbContainer.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	WaClient = whatsmeow.NewClient(deviceStore, clientLog)
	WaClient.AddEventHandler(EventsHandler)

	msghandler.Initialize(WaClient)

	dbservice.Init("sqlite3", "file:kontak.db?_foreign_keys=on")
	dbservice.PrepareTable()

	if WaClient.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := WaClient.GetQRChannel(context.Background())
		err := WaClient.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qr.L, os.Stdout)
				fmt.Println("QR code:", evt.Code)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err := WaClient.Connect()
		if err != nil {
			panic(err)
		}
	}
}
