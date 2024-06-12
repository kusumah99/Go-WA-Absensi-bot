package msghandler

import (
	"context"
	"fmt"
	"strings"
	"time"

	// _ "ram/msgbard"
	// _ "../ai"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

var eventsMessage *events.Message
var client *whatsmeow.Client

const menu string = "Perintah yang dapat digunakan adalah: \n" +
	"/pesan => pesan yang akan di kirim\n" +
	"/kirim => kirimkan ke semua anggota\n" +
	"/menu => menampilkan pesan ini\n" +
	"/list => menampilkan daftar anggota\n" +
	"/tambah => tambah anggota\n" +
	"/edit => edit anggota\n" +
	"/hapus => menghapus anggota"

func Initialize(waClient *whatsmeow.Client) {
	client = waClient
}

func OnMessage(v *events.Message) {
	eventsMessage = v

	pesan := strings.TrimSpace(v.Message.GetConversation())

	lPesan := strings.ToLower(pesan)

	fmt.Println("PESAN MASUK: ", pesan)
	fmt.Println("Panjang PESAN: ", len(pesan))
	if lPesan == "/menu" {
		time.Sleep(time.Second * 2)

		fmt.Println("MENU")
		sender := v.Info.Sender.User
		fmt.Println("Sendernya diatas : " + sender)
		// targetJID := v.Message. .Info.DeviceSentMeta.DestinationJID
		// targetJID, _ := types.ParseJID(v.Message. )
		// client.SendMessage(context.Background(), v.Info.Sender, &waProto.Message{
		// 	Conversation: proto.String("Menu Anda hari ini adalah: Nasgor"),
		// })
		// sendResp, err := kirimPesan(sender, "Menu Anda hari ini adalah: Nasgor")
		sendResp, err := kirimPesan(sender, menu)

		fmt.Println(sendResp)
		// time.Sleep(time.Second * 5)
		if err != nil {
			fmt.Println("Ada error: " + err.Error())
			// client.SendMessage(context.Background(), v.Info.Sender, &waProto.Message{
			// 	Conversation: proto.String("Ada Error: " + err.Error()),
			// })
		} else {
			fmt.Println("Kirim pesan berhasil ")
			// client.SendMessage(context.Background(), v.Info.Sender, &waProto.Message{
			// 	Conversation: proto.String("Pesan terkirim"),
			// })
		}
		// } else if strings.EqualFold(lPesan, "/msg") {
	} else if strings.HasPrefix(lPesan, "/msg ") {
		// sintak: msg noHp Pesan
		// a := regexp.MustCompile(` `)
		time.Sleep(time.Millisecond * 2330)
		a := strings.SplitN(lPesan, " ", 3)
		if len(a) == 3 {
			fmt.Println("Proses KIRIM Pesan ke " + a[1] + " isinya: \"" + a[2] + "\"")
			sendResp, err := kirimPesan(a[1], a[2])
			fmt.Println(sendResp)
			time.Sleep(time.Millisecond * 2132)
			if err != nil {
				fmt.Println("Ada error: " + err.Error())
				kirimPesanJid2(v.Info.Sender, "Ada Error: "+err.Error())
			} else {
				fmt.Println("Kirim pesan ke: " + a[1])
				kirimPesanJid2(v.Info.Sender, "Pesan terkirim")
			}
		} else {
			fmt.Println("GAGAL FORMAT KIRIM Pesan")
			time.Sleep(time.Millisecond * 2132)
			kirimPesanJid2(v.Info.Sender, "Format yang digunakan harus sepertin ini: /msg noHP Pesan")
		}
		fmt.Println(len(a))

		// } else if strings.HasPrefix(lPesan, "/ai ") {
		// 	// a := regexp.MustCompile(` `)
		// 	a := strings.SplitN(lPesan, " ", 2)
		// 	if len(a) == 2 {
		// 		fmt.Println("Permintaan ke AI Bard: " + a[1])
		// 		jawab, _, err := msgbard.TanyaBard(a[1])
		// 		time.Sleep(time.Millisecond * 2530)
		// 		if err != nil {
		// 			fmt.Println("Ada error: " + err.Error())
		// 			kirimPesanJid2(v.Info.Sender, "Ada Error: "+err.Error())
		// 		} else {
		// 			fmt.Println("Dapat jawaban dari AI")
		// 			fmt.Println("Jawabannya: " + jawab)
		// 			kirimPesanJid2(v.Info.Sender, jawab)
		// 			// kirimPesanJid2("6281210007733", jawab)
		// 		}
		// 	} else {
		// 		fmt.Println("GAGAL FORMAT KIRIM Pesan")
		// 		time.Sleep(time.Millisecond * 1723)
		// 		kirimPesanJid2(v.Info.Sender, "Format yang digunakan harus sepertin ini: /ai pertanyaan")
		// 	}
		// 	fmt.Println(len(a))

		// } else {
		// fmt.Println("NOT MENU")
	} else if strings.HasPrefix(lPesan, "/tambah ") {
		// sintak: member nohp nama
		// nama gk boleh ada spasi
		time.Sleep(time.Millisecond * 2213)
		a := strings.SplitN(lPesan, " ", 3)
		if len(a) == 3 {
			fmt.Println("Proses nambah anggota " + a[1] + " nama: \"" + a[2] + "\"")
			memberResp, err := addMember(a[1], a[2])
			fmt.Println(memberResp)
			time.Sleep(time.Millisecond * 2232)
			if err != nil {
				fmt.Println("Ada error: " + err.Error())
				kirimPesanJid2(v.Info.Sender, "Ada Error: "+err.Error())
			} else {
				fmt.Println("Berhasil nambah anggota")
				kirimPesanJid2(v.Info.Sender, "Berhasil nambah anggota")
			}
		} else {
			fmt.Println("GAGAL FORMAT tambah anggota")
			time.Sleep(time.Millisecond * 2132)
			kirimPesanJid2(v.Info.Sender, "Format yang digunakan harus sepertin ini: /tambah noHp Nama")
		}
		fmt.Println(len(a))

	}
}

func kirimPesan(nomorTujuan string, pesan string) (resp whatsmeow.SendResponse, err error) {
	return client.SendMessage(context.Background(), types.JID{
		User:   nomorTujuan,
		Server: types.DefaultUserServer,
	}, &waE2E.Message{
		Conversation: proto.String(pesan),
	})
}

func kirimPesanJid2(Jid types.JID, pesan string) (resp whatsmeow.SendResponse, err error) {
	nomorTujuan := Jid.User
	return client.SendMessage(context.Background(), types.JID{
		User:   nomorTujuan,
		Server: types.DefaultUserServer,
	}, &waE2E.Message{
		Conversation: proto.String(pesan),
	})
}

func kirimPesanJid(Jid types.JID, pesan string) (resp whatsmeow.SendResponse, err error) {
	return client.SendMessage(context.Background(), Jid, &waE2E.Message{
		Conversation: proto.String(pesan),
	})
}
