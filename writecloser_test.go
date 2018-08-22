//
// writecloser_test.go
// Copyright (C) 2018 Grigorii Sokolik <g.sokol99@g-sokol.info>
//
// Distributed under terms of the MIT license.
//

package statsd

import (
	"bytes"
	"testing"
)

func TestWriteCloser(t *testing.T) {
	connMock := newConnMockExample()
	client, err := New(CustomConn(connMock))
	if err != nil {
		t.Fatalf("Failed to initialize client %v", err)
	}
	client.Count(testKey, 5)
	client.Close()

	if len(*connMock) != 2 {
		t.Errorf(
			"Infalid connection calls number, got:\n%d\nwant:\n%d",
			len(*connMock),
			2,
		)
		return
	}
	bts, ok := (*connMock)[0].([]byte)
	if !ok || !bytes.Equal(bts, []byte("test_key:5|c")) {
		t.Errorf(
			"Invalid output, got:\n%+v\nwant:\n%+v",
			(*connMock)[0],
			[]byte("test_key:5|c"),
		)
		return
	}
	if (*connMock)[1] != closeCall {
		t.Errorf(
			"Invalid call, got\n%+v\nwant:\n%+v",
			(*connMock)[1],
			closeCall,
		)
	}
}

type connMockExample []call

func newConnMockExample() *connMockExample {
	return new(connMockExample)
}

func (c *connMockExample) Write(p []byte) (int, error) {
	c.append(call(p))
	return len(p), nil
}

func (c *connMockExample) Close() error {
	c.append(closeCall)
	return nil
}

func (cm *connMockExample) append(c ...call) {
	cmBuf := connMockExample(append(([]call)(*cm), c...))
	*cm = cmBuf
}

type call interface{}

var closeCall = call(nil)
