package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/andybons/hipchat"
	"github.com/tlorens/go-ibgetkey"
)

var pom_time int //ポモドーロ合計回数を取得
var token *string
var room_name *string
var user_name *string

func main() {

	token = flag.String("t", "", "hipchat access token")
	room_name = flag.String("r", "", "send message room name")
	user_name = flag.String("u", "", "send message user name")
	flag.Parse()

	kill := make(chan bool)
	finished := make(chan bool)

	restart__key := "r"
	t := int(restart__key[0])

	finish_key := "f"
	f := int(finish_key[0])

	go pomTimerGoroutine(kill, finished)

loop:
	for {
		input := keyboard.ReadKey()
		select {
		case <-finished:
			break loop
		default:
			if input == t {
				fmt.Println("\nrestart")
				kill <- true
				go pomTimerGoroutine(kill, finished)
			}
			if input == f {
				kill <- true
				break loop
			}
		}
	}

	fmt.Println("\nfinish")
	hipchatSend(fmt.Sprintf("%dポモドーロ完了です。お疲れ様でした。", pom_time))
}

func pomTimerGoroutine(kill, finished chan bool) {
	fmt.Println("\npomodoro running")

	for i := 0; i < 4; i++ {
		hipchatSend("25分作業 (気合を入れていきましょう。)")
		//25分 1500
		fmt.Println("\n\n【work】****************************")
		for j := 0; j < 20; j++ {
			if j%10 == 0 {
				fmt.Printf("\n%dmin", j/10)
			} else {
				fmt.Print(".")
			}
			time.Sleep(1 * time.Second)
			select {
			case <-kill:
				return
			default:
			}
		}
		pom_time++
		hipchatSend("5分休憩 (歩きましょう。)")
		fmt.Println("\n\n【rest】****************************")
		//5分 300
		for j := 0; j < 20; j++ {
			if j%10 == 0 {
				fmt.Printf("\n%dmin", j/10)
			} else {
				fmt.Print(".")
			}
			time.Sleep(1 * time.Second)
			select {
			case <-kill:
				return
			default:
			}
		}
	}
	finished <- true
	return
}

func hipchatSend(msg string) {
	c := hipchat.NewClient(*token)

	req := hipchat.MessageRequest{
		RoomId:        *room_name,
		From:          "Pom",
		Message:       "@" + *user_name + " " + msg,
		Color:         hipchat.ColorGreen,
		MessageFormat: hipchat.FormatText,
		Notify:        true,
	}

	if err := c.PostMessage(req); err != nil {
		log.Printf("Expected no error, but got %q", err)
	}
}
