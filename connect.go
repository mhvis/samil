package samil

import (
	"net"
	"time"
)

// Listens for incoming connections while advertising.
// Binds to given interface IP for advertising and listening.
func connect(interfaceIP net.IP) (*net.TCPConn, error) {
	done := make(chan *net.TCPConn, 1)
	fail := make(chan error, 1)
	go func() {
		conn, err := listen(interfaceIP)
		if err != nil {
			fail <- err
			return
		}
		done <- conn
	}()
	for {
		err := advertise(interfaceIP)
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
func listen(interfaceIP net.IP) (*net.TCPConn, error) {
	listener, err := net.ListenTCP("tcp4", &net.TCPAddr{
		IP:   interfaceIP,
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
func advertise(interfaceIP net.IP) error {
	socket, err := net.DialUDP("udp4", &net.UDPAddr{
		IP: interfaceIP,
	}, &net.UDPAddr{
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
