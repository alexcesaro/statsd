package statsd

import (
	"errors"
	"io"
	"net"
	"testing"
	"time"
)

type mockNetConn struct {
	read             func(p []byte) (n int, err error)
	write            func(p []byte) (n int, err error)
	close            func() error
	localAddr        func() net.Addr
	remoteAddr       func() net.Addr
	setDeadline      func(t time.Time) error
	setReadDeadline  func(t time.Time) error
	setWriteDeadline func(t time.Time) error
}

func (m *mockNetConn) Read(p []byte) (n int, err error) {
	if m.read != nil {
		return m.read(p)
	}
	panic("implement me")
}

func (m *mockNetConn) Write(p []byte) (n int, err error) {
	if m.write != nil {
		return m.write(p)
	}
	panic("implement me")
}

func (m *mockNetConn) Close() error {
	if m.close != nil {
		return m.close()
	}
	panic("implement me")
}

func (m *mockNetConn) LocalAddr() net.Addr {
	if m.localAddr != nil {
		return m.localAddr()
	}
	panic("implement me")
}

func (m *mockNetConn) RemoteAddr() net.Addr {
	if m.remoteAddr != nil {
		return m.remoteAddr()
	}
	panic("implement me")
}

func (m *mockNetConn) SetDeadline(t time.Time) error {
	if m.setDeadline != nil {
		return m.setDeadline(t)
	}
	panic("implement me")
}

func (m *mockNetConn) SetReadDeadline(t time.Time) error {
	if m.setReadDeadline != nil {
		return m.setReadDeadline(t)
	}
	panic("implement me")
}

func (m *mockNetConn) SetWriteDeadline(t time.Time) error {
	if m.setWriteDeadline != nil {
		return m.setWriteDeadline(t)
	}
	panic("implement me")
}

func TestSafeConn_FailsToWriteIfCannotRead(t *testing.T) {
	c := &mockNetConn{
		setReadDeadline: func(t time.Time) error {
			return nil
		},
		read: func(b []byte) (int, error) {
			return 0, io.EOF
		},
	}

	s := SafeConn{
		netConn: c,
	}

	p := []byte("test_key:1|c\n")
	n, err := s.Write(p)
	if n != 0 {
		t.Error("Write() did not return 0 bytes when it failed")
	}
	if err == nil {
		t.Error("Error should have been connection is closed")
	}
}

func TestSafeConn_SuccessfullyWritesWhenConnectionOpen(t *testing.T) {
	c := &mockNetConn{
		setReadDeadline: func(t time.Time) error {
			return nil
		},
		read: func(b []byte) (int, error) {
			return 1, nil
		},
		write: func(b []byte) (int, error) {
			return len(b), nil
		},
	}

	s := SafeConn{
		netConn: c,
	}

	p := []byte("test_key:1|c\n")
	_, err := s.Write(p)
	if err != nil {
		t.Errorf("Error should have been nil, but instead it was: %v", err)
	}
}

func TestNewSafeConnWithDefaultTimeouts(t *testing.T) {
	for _, tc := range [...]struct {
		Name string
		Conn net.Conn
		Err  error
	}{
		{
			Name: `failure`,
			Err:  errors.New(`some error`),
		},
		{
			Name: `success`,
			Conn: new(mockNetConn),
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			type (
				DialIn struct {
					Network string
					Address string
					Timeout time.Duration
				}
				DialOut struct {
					Conn net.Conn
					Err  error
				}
			)
			var (
				dialIn  = make(chan DialIn)
				dialOut = make(chan DialOut)
			)
			defer close(dialIn)
			defer close(dialOut)
			defer func() func() {
				old := dialTimeout
				dialTimeout = func(network, address string, timeout time.Duration) (net.Conn, error) {
					dialIn <- DialIn{network, address, timeout}
					v := <-dialOut
					return v.Conn, v.Err
				}
				return func() { dialTimeout = old }
			}()()

			const (
				expectedNetwork = `tcp`
				expectedAddress = `127.0.0.1:21969`
			)

			done := make(chan struct{})
			go func() {
				defer close(done)
				if v := <-dialIn; v != (DialIn{Network: expectedNetwork, Address: expectedAddress, Timeout: defaultConnTimeout}) {
					t.Errorf("%+v", v)
				}
				dialOut <- DialOut{tc.Conn, tc.Err}
			}()

			conn, err := NewSafeConnWithDefaultTimeouts(expectedNetwork, expectedAddress)
			if err != tc.Err {
				t.Error(conn, err)
			}
			if (tc.Conn == nil) != (conn == nil) {
				t.Error(conn)
			} else if conn != nil {
				if conn.netConn != tc.Conn {
					t.Error(conn.netConn)
				}
				if conn.connTimeout != defaultConnTimeout {
					t.Error(conn.connTimeout)
				}
				if conn.readTimeout != defaultReadTimeout {
					t.Error(conn.readTimeout)
				}
			}

			<-done
		})
	}
}

func TestSafeConn_Close(t *testing.T) {
	a, b := net.Pipe()
	conn := SafeConn{
		netConn:     a,
		readTimeout: 1,
	}
	if conn.connIsClosed() {
		t.Error()
	}
	if err := (&SafeConn{netConn: b}).Close(); err != nil {
		t.Error(err)
	}
	if !conn.connIsClosed() {
		t.Error()
	}
	if err := conn.Close(); err != nil {
		t.Error(err)
	}
}
