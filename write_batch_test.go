package tls

import (
	"io"
	"net"
	"sync/atomic"
	"testing"
)

type countingWriteConn struct {
	net.Conn
	writes atomic.Int32
}

func (c *countingWriteConn) Write(b []byte) (int, error) {
	c.writes.Add(1)
	return c.Conn.Write(b)
}

// A Write larger than one record reaches the transport as a single write
// (miray T-031); the peer still reads the exact bytes.
func TestWriteSendsRecordsInOneWrite(t *testing.T) {
	c, s := localPipe(t)
	counted := &countingWriteConn{Conn: c}
	clientConfig := testConfig.Clone()
	client := Client(counted, clientConfig)
	server := Server(s, testConfig.Clone())
	defer client.Close()
	defer server.Close()

	payload := make([]byte, 5*maxPlaintext+123)
	for i := range payload {
		payload[i] = byte(i)
	}
	received := make(chan []byte, 1)
	go func() {
		buf := make([]byte, len(payload))
		if _, err := io.ReadFull(server, buf); err != nil {
			t.Error(err)
		}
		received <- buf
	}()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}

	before := counted.writes.Load()
	if _, err := client.Write(payload); err != nil {
		t.Fatal(err)
	}
	if writes := counted.writes.Load() - before; writes != 1 {
		t.Fatalf("Write of %d bytes made %d transport writes, want 1", len(payload), writes)
	}
	got := <-received
	for i := range payload {
		if got[i] != payload[i] {
			t.Fatalf("byte %d = %d, want %d", i, got[i], payload[i])
		}
	}
}
