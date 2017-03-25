package samil

import (
	"net"
	"time"
)

// Listens for incoming connections while advertising.
func connect() (*net.TCPConn, error) {
	done := make(chan *net.TCPConn, 1)
	fail := make(chan error, 1)
	go func() {
		conn, err := listen()
		if err != nil {
			fail <- err
			return
		}
		done <- conn
	}()
	for {
		err := advertise()
		if err != nil {
			return nil, err
		}
		timer := time.NewTimer(5 * time.Second)
		select {
		case conn := <-done:
			return conn, nil
		case err = <-fail:
			return nil, err
		case <-timer.C:
		}
	}
}

// Listens for incoming connections with a deadline of a minute.
func listen() (*net.TCPConn, error) {
	listener, err := net.ListenTCP("tcp4", &net.TCPAddr{
		IP:   net.IPv4zero,
		Port: 1200,
	})
	if err != nil {
		return nil, err
	}
	defer listener.Close()
	err = listener.SetDeadline(time.Now().Add(time.Minute))
	if err != nil {
		return nil, err
	}
	return listener.AcceptTCP()
}

// Advertises the existence of this client.
func advertise() error {
	socket, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: 1300,
	})
	if err != nil {
		return err
	}
	_, err = socket.Write(advertisement)
	socket.Close()
	return err
}
