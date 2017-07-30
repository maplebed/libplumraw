package libplumraw

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
)

type LightpadAnnouncement struct {
	ID   string
	IP   net.IP
	Port int
}

type DefaultLightpadHeartbeat struct{}

func (d *DefaultLightpadHeartbeat) Listen(ctx context.Context) chan LightpadAnnouncement {
	logrus.Debug("about to listen for broadcast heartbeats")
	/* Lets prepare a address at any address at port */
	ServerAddr, err := net.ResolveUDPAddr("udp", ":43770")
	if err != nil {
		panic(err)
	}

	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 1024)

	responses := make(chan LightpadAnnouncement, 0)

	go func() {
		defer ServerConn.Close()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("got cancel; bailing.")
				return
			default:
			}
			ServerConn.SetReadDeadline(time.Now().Add(1 * time.Second))
			n, addr, err := ServerConn.ReadFromUDP(buf)
			if err != nil {
				if err.(*net.OpError).Timeout() {
					continue
				}
				panic(err)
			}
			msg := string(buf[0:n])
			// 2017-07-29 15:34:41.655324026 -0700 PDT Received  PLUM 8888 8429176c-bf88-4aee-be07-b6a9064cf1ab 8443  from  192.168.1.91:54209
			fmt.Println(time.Now(), "Received ", msg, " from ", addr)

			if strings.HasPrefix(msg, "PLUM 8888") {
				bits := strings.Split(msg, " ")
				if len(bits) == 4 {
					port, err := strconv.Atoi(bits[3])
					if err != nil {
						fmt.Printf("couldn't parse port from lightpad announcement: %s", bits[3])
						continue
					}
					la := LightpadAnnouncement{
						ID:   bits[2],
						IP:   addr.IP,
						Port: port,
					}
					responses <- la
					// make a new buf to erase the contents of the old one from memory
					buf = make([]byte, 1024)
				}
			}
			if err != nil {
				fmt.Println("Error: ", err)
			}
		}
	}()
	return responses
}
