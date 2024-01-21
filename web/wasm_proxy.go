package web

import (
	"fmt"
	"mtui/types"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var MinetestMagic = []byte{0x4f, 0x45, 0x74, 0x03}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (api *Api) HandleProxy(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	api.app.CreateUILogEntry(&types.Log{
		Username: c.Username,
		Event:    "ws-proxy",
		Message:  fmt.Sprintf("User '%s' connects via websocket-proxy", c.Username),
	}, r)

	go func() {
		err = api.handleProxyConnection(conn)
		if err != nil {
			logrus.WithError(err).Error("ws handler error")
		}
	}()
}

func (api *Api) handleProxyConnection(conn *websocket.Conn) error {
	_, data, err := conn.ReadMessage()
	if err != nil {
		return err
	}
	defer conn.Close()

	parts := strings.Split(string(data), " ")
	if len(parts) != 5 {
		return fmt.Errorf("invalid command: '%s'", data)
	}
	if parts[0] != "PROXY" {
		return fmt.Errorf("command not implemented: '%s'", parts[0])
	}
	if parts[1] != "IPV4" {
		return fmt.Errorf("ip version not implemented: '%s'", parts[1])
	}
	protocol := parts[2]
	host := parts[3]
	port, _ := strconv.ParseInt(parts[4], 10, 32)

	if port < 1 || port >= 65536 {
		return fmt.Errorf("invalid port: %d", port)
	}

	logrus.WithFields(logrus.Fields{
		"host": host,
		"port": port,
	}).Info("WASM WS Proxy connecting")

	// only allow dns requests and minetest-protocol forwarding
	if host == "10.0.0.1" && port == 53 && protocol == "TCP" {
		err = resolveDNS(conn)
	} else if protocol == "UDP" {
		// override port/host for local minetest connection
		host = api.app.Config.WASMMinetestHost
		if host == "" {
			// fallback
			host = "engine"
		}
		port = int64(api.app.Config.DockerMinetestPort)
		if port == 0 {
			// fallback
			port = 30000
		}

		err = forwardData(conn, host, port)
	} else {
		return fmt.Errorf("unsupported command: '%s'", data)
	}
	return err
}

func resolveDNS(conn *websocket.Conn) error {
	conn.WriteMessage(websocket.TextMessage, []byte("PROXY OK"))

	_, data, err := conn.ReadMessage()
	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"host": string(data),
	}).Debug("WASM WS: Resolving host")

	ips, err := net.LookupIP(string(data))
	if err != nil {
		return err
	}

	if len(ips) == 0 {
		return fmt.Errorf("host not found")
	}

	err = conn.WriteMessage(websocket.BinaryMessage, []byte(ips[0]))
	if err != nil {
		return err
	}

	return conn.Close()
}

func forwardData(conn *websocket.Conn, host string, port int64) error {
	uaddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}

	udpconn, err := net.DialUDP("udp", nil, uaddr)
	if err != nil {
		return err
	}
	errchan := make(chan error, 1)
	run := atomic.Bool{}
	run.Store(true)

	conn.WriteMessage(websocket.TextMessage, []byte("PROXY OK"))

	go func() {
		buf := make([]byte, 3000)
		for run.Load() {
			len, err := udpconn.Read(buf)
			if err != nil {
				errchan <- err
				return
			}
			err = conn.WriteMessage(websocket.BinaryMessage, buf[:len])
			if err != nil {
				errchan <- err
				return
			}
		}
	}()

	go func() {
		for run.Load() {
			_, data, err := conn.ReadMessage()
			if err != nil {
				errchan <- err
				return
			}
			if len(data) < 9 {
				errchan <- fmt.Errorf("invalid packet size: %d", len(data))
				return
			}

			// ensure that we are using the minetest protocol
			for i, b := range MinetestMagic {
				if data[i] != b {
					errchan <- fmt.Errorf("invalid magic at offset %d: %d", i, data[i])
					return
				}
			}

			_, err = udpconn.Write(data)
			if err != nil {
				errchan <- err
				return
			}
		}
	}()

	err = <-errchan
	run.Store(false)
	conn.Close()
	return err
}
