package statsd

import (
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
