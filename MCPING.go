package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	raknet "github.com/beito123/go-raknet"
	"github.com/beito123/go-raknet/protocol"
)

type ServerInfo struct {
	Name   string
	Online int
	Max    int
	Addr   string
}

func main() {
	var (
		addr = flag.String("s", "play.lbsg.net:19132", "Ip Addr")
	)
	flag.Parse()
	s, err := Ping(*addr)
	if err != nil {
		fmt.Println("Offline")
		os.Exit(1)
	}
	print(s)
}

func Ping(addr string) (*ServerInfo, error) {
	conn, cer := net.Dial("udp", addr)
	if cer != nil {
		return nil, cer
	}
	defer conn.Close()
	pl := new(protocol.Protocol)
	pl.RegisterPackets()
	pac := &protocol.UnconnectedPing{
		Timestamp:  time.Now().Unix(),
		Magic:      true,
		PingID:     0,
		Connection: raknet.ConnectionVanilla,
	}
	pac.Encode()
	_, err := conn.Write(pac.Bytes())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	buffer := make([]byte, 2048)

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, err = conn.Read(buffer)
	if err != nil {
		//IO TimeOut
		return nil, err
	}
	pk, ok := pl.Packet(buffer[0])
	if !ok {
		log.Println(ok)
	}
	pk.SetBytes(buffer)
	upong := pk.(*protocol.UnconnectedPong)
	upong.Decode()
	spl := strings.Split(upong.Identifier.Build(), ";")
	online, err := strconv.Atoi(spl[4])
	if err != nil {
		return nil, err
	}
	max, err := strconv.Atoi(spl[5])
	if err != nil {
		return nil, err
	}
	spl2 := strings.Split(spl[1], "§")
	name := ""
	for i := 0; i < len(spl2); i++ {
		if len(spl2[i]) >= 2 {
			name += string([]rune(spl2[i])[1:])
		}
	}
	sinfo := &ServerInfo{
		Name:   name,
		Online: online,
		Max:    max,
		Addr:   addr,
	}
	return sinfo, nil
}

func print(si *ServerInfo) {
	fmt.Printf("[%s] %s (%d/%d)\n", si.Addr, si.Name, si.Online, si.Max)
}
