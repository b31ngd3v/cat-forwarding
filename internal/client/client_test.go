package client

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"strconv"
	"testing"
	"time"
)

type MockConn struct {
	io.Reader
	io.Writer
}

func (m *MockConn) Close() error                       { return nil }
func (m *MockConn) LocalAddr() net.Addr                { return &net.TCPAddr{} }
func (m *MockConn) RemoteAddr() net.Addr               { return &net.TCPAddr{} }
func (m *MockConn) SetDeadline(t time.Time) error      { return nil }
func (m *MockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *MockConn) SetWriteDeadline(t time.Time) error { return nil }

func TestPerformPawshake(t *testing.T) {
	s := &server{}

	tests := []struct {
		name         string
		serverReply  string
		expectedPort int
		expectedErr  error
	}{
		{
			name:         "Successful pawshake",
			serverReply:  pawshakeRequestSucceeded + "12345\n",
			expectedPort: 12345,
		},
		{
			name:        "Max connection limit exceeded",
			serverReply: maxConnLimitExceeded,
			expectedErr: errors.New("max connection limit exceeded"),
		},
		{
			name:        "Invalid response",
			expectedErr: errors.New("pawshake failed! please update the cat-forwading client"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			serverResponse := bytes.NewBufferString(tc.serverReply)
			clientBuffer := &bytes.Buffer{}

			mockConn := &MockConn{
				Reader: serverResponse,
				Writer: clientBuffer,
			}

			port, err := s.performPawshake(mockConn)

			if tc.expectedErr != nil {
				if err == nil || err.Error() != tc.expectedErr.Error() {
					t.Errorf("Expected error '%v', got '%v'", tc.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			} else {
				if port != tc.expectedPort {
					t.Errorf("Expected port %d, got %d", tc.expectedPort, port)
				}
			}

			sentData := clientBuffer.String()
			if sentData != pawshakeRequest {
				t.Errorf("Expected to send '%s', but sent '%s'", pawshakeRequest, sentData)
			}
		})
	}
}

func TestReceivePackets(t *testing.T) {
	s := &server{}

	inputData := []byte("test data")
	remoteConn := bytes.NewBuffer(inputData)
	localConn := &bytes.Buffer{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)

	go s.receivePackets(ctx, localConn, remoteConn, errCh)

	time.Sleep(100 * time.Millisecond)

	outputData := localConn.Bytes()
	if !bytes.Equal(inputData, outputData) {
		t.Errorf("Expected data '%s', got '%s'", inputData, outputData)
	}
}

func TestSendPackets(t *testing.T) {
	s := &server{}

	inputData := []byte("test data")
	localConn := bytes.NewBuffer(inputData)
	remoteConn := &bytes.Buffer{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)

	go s.sendPackets(ctx, localConn, remoteConn, errCh)

	time.Sleep(100 * time.Millisecond)

	outputData := remoteConn.Bytes()
	if !bytes.Equal(inputData, outputData) {
		t.Errorf("Expected data '%s', got '%s'", inputData, outputData)
	}
}

func TestHandleConn(t *testing.T) {
	s := &server{
		quitCh: make(chan struct{}),
	}
	defer close(s.quitCh)

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to start listener: %v", err)
	}
	defer listener.Close()

	_, portStr, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		t.Fatalf("Failed to get port: %v", err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("Failed to convert port: %v", err)
	}

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Errorf("Failed to accept connection: %v", err)
			return
		}
		defer conn.Close()
		io.Copy(conn, conn)
	}()

	remoteConn1, remoteConn2 := net.Pipe()

	deadline := time.Now().Add(100 * time.Millisecond)
	remoteConn2.SetReadDeadline(deadline)

	go s.handleConn(remoteConn1, port)

	testData := []byte("hello world")
	_, err = remoteConn2.Write(testData)
	if err != nil {
		t.Fatalf("Failed to write to remoteConn2: %v", err)
	}

	buffer := make([]byte, len(testData))
	_, err = io.ReadFull(remoteConn2, buffer)
	if err != nil {
		t.Fatalf("Failed to read from remoteConn2: %v", err)
	}

	if !bytes.Equal(testData, buffer) {
		t.Errorf("Expected '%s', got '%s'", testData, buffer)
	}
}
