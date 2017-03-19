package samil

import (
	"net"
	"time"
)

type Samil *net.TCPConn

// NewConnection searches for an inverter in the network and returns the
// connection if one is found.
// Calling it multiple times will connect to different inverters.
func NewConnection() (samil Samil, err error) {
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
		timer := NewTimer(5 * time.Second)
		select {
		case conn := <-done:
			return Samil(conn), nil
		case err = <-fail:
			return nil, err
		case <-timer.C:
		}
	}
}

// Listens for incoming connections, with a deadline of a minute.
func listen() (*net.TCPConn, error) {
	listener, err := net.ListenTCP("tcp4", &net.TCPAddr{net.IPv4zero, 1200})
	if err != nil {
		return nil, err
	}
	defer listener.Close()
	err = listener.SetDeadline(time.Minute)
	if err != nil {
		return nil, err
	}
	return listener.AcceptTCP()
}

func advertise() error {
	socket, err := net.DialUDP("udp4", nil, &net.UDPAddr{net.IPv4bcast, 1300})
	if err != nil {
		return err
	}
	_, err := socket.Write(advertisement)
	socket.Close()
	return err
}

/*
def _tear_down_response(data):
    """Helper function to extract header, payload and end from received response
    data."""
    response_header = data[2:5]
    # Below is actually not used
    response_payload_size = int.from_bytes(data[5:7], byteorder='big')
    response_payload = data[7:-2]
    response_end = data[-2:]
    return response_header, response_payload, response_end

*/
