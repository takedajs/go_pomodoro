package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/andybons/hipchat"
	"github.com/tlorens/go-ibgetkey"
)

var pom_time int
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
	hipchatSend(fmt.Sprintf("%dポモドーロ完了です。お疲れ様でした。", pom_time), "red")
}

func pomTimerGoroutine(kill, finished chan bool) {
	fmt.Print("\npomodoro running")

	for i := 0; i < 4; i++ {
		hipchatSend("25分作業 (気合を入れていきましょう。)", "gray")
		fmt.Print("\n【25min work】****************************")
		//25m = 1500s
		for j := 0; j < 1500; j++ {
			if j%60 == 0 {
				fmt.Printf("\n%dmin", j/60)
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
		hipchatSend("5分休憩 (歩きましょう。)", "green")
		fmt.Print("\n【5min rest】****************************")
		//5m = 300s
		for j := 0; j < 300; j++ {
			if j%60 == 0 {
				fmt.Printf("\n%dmin", j/60)
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

func hipchatSend(msg string, bg_color string) {
	c := hipchat.NewClient(*token)

	req := hipchat.MessageRequest{
		RoomId:        *room_name,
		From:          "Pom",
		Message:       "@" + *user_name + " " + msg,
		Color:         bg_color,
		MessageFormat: hipchat.FormatText,
		Notify:        true,
	}

	if err := c.PostMessage(req); err != nil {
		log.Printf("Expected no error, but got %q", err)
	}
}
