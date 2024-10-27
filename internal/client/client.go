package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/inconshreveable/muxado"
)

const bufferSize = 1024

var (
	c                        = GetConfig()
	pawshakeTimeout          = 5 * time.Second
	versionString            = fmt.Sprintf("cat_forwarding_v%s", c.Version)
	maxConnLimitExceeded     = fmt.Sprintf("%s\nmax_conn_limit_exceeded", versionString)
	pawshakeRequest          = fmt.Sprintf("%s\npawshake_request\n", versionString)
	pawshakeRequestSucceeded = fmt.Sprintf("%s\npawshake_successful\nport_", versionString)
)

func Run(port int) {
	s := newServerConnection(c.ServerAddr)
	s.connectAndForwardPackets(port)
}

type server struct {
	addr   string
	quitCh chan struct{}
}

func newServerConnection(addr string) *server {
	return &server{
		addr:   addr,
		quitCh: make(chan struct{}),
	}
}

func (s *server) connectAndForwardPackets(port int) {
	conn, err := net.Dial("tcp", s.addr)
	if err != nil {
		fmt.Println("failed to connect to the server")
		os.Exit(1)
	}
	defer conn.Close()

	remotePort, err := s.performPawshake(conn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	host := strings.Split(s.addr, ":")[0]
	fmt.Printf("🐾 Cat Forwarding activated on port %v to %s:%v! Your cat is now purring across the globe! 🐾\n", port, host, remotePort)

	sess := muxado.Client(conn, nil)
	defer sess.Close()

	go func() {
		for {
			select {
			case <-s.quitCh:
				return
			default:
				stream, err := sess.Accept()
				if err != nil {
					continue
				}
				go s.handleConn(stream, port)
			}
		}
	}()

	sess.Wait()
	close(s.quitCh)
}

func (s *server) performPawshake(conn net.Conn) (int, error) {
	conn.SetDeadline(time.Now().Add(pawshakeTimeout))
	defer conn.SetDeadline(time.Time{})

	defaultErr := errors.New("pawshake failed! please update the cat-forwading client")

	conn.Write([]byte(pawshakeRequest))

	res := make([]byte, bufferSize)
	bc, err := conn.Read(res)
	if err != nil {
		return 0, defaultErr
	}

	resStr := string(res[:bc])

	if resStr == maxConnLimitExceeded {
		return 0, errors.New("max connection limit exceeded")
	}

	portStr, success := strings.CutPrefix(resStr, pawshakeRequestSucceeded)
	if !success {
		return 0, defaultErr
	}
	portStr = strings.TrimSuffix(portStr, "\n")

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, defaultErr
	}

	return port, nil
}

func (s *server) handleConn(remoteConn net.Conn, localPort int) {
	defer remoteConn.Close()

	localAddress := fmt.Sprintf(":%d", localPort)
	localConn, err := net.Dial("tcp", localAddress)
	if err != nil {
		return
	}
	defer localConn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 2)

	go s.receivePackets(ctx, localConn, remoteConn, errCh)
	go s.sendPackets(ctx, localConn, remoteConn, errCh)

	select {
	case <-s.quitCh:
	case <-errCh:
		cancel()
	}
}

func (s *server) receivePackets(ctx context.Context, localConn io.Writer, remoteConn io.Reader, errCh chan error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			d := make([]byte, bufferSize)
			bc, err := remoteConn.Read(d)
			if err == io.EOF {
				errCh <- io.EOF
				return
			} else if err != nil {
				continue
			}
			_, err = localConn.Write(d[:bc])
			if err != nil {
				continue
			}
		}
	}
}

func (s *server) sendPackets(ctx context.Context, localConn io.Reader, remoteConn io.Writer, errCh chan error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			d := make([]byte, bufferSize)
			bc, err := localConn.Read(d)
			if err != nil {
				continue
			}
			_, err = remoteConn.Write(d[:bc])
			if err == io.EOF {
				errCh <- io.EOF
				return
			} else if err != nil {
				continue
			}
		}
	}
}
