package source

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"time"
)

// SyslogSource reads log lines from a UDP syslog listener.
type SyslogSource struct {
	addr string
	conn *net.UDPConn
}

// NewSyslogSource creates a new SyslogSource listening on the given UDP address (e.g. ":514").
func NewSyslogSource(addr string) (*SyslogSource, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("syslog: resolve addr %q: %w", addr, err)
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, fmt.Errorf("syslog: listen %q: %w", addr, err)
	}
	return &SyslogSource{addr: addr, conn: conn}, nil
}

// Name returns a human-readable identifier for this source.
func (s *SyslogSource) Name() string {
	return fmt.Sprintf("syslog(%s)", s.addr)
}

// Lines streams log lines received over UDP until ctx is cancelled.
func (s *SyslogSource) Lines(ctx context.Context) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		defer s.conn.Close()

		buf := make([]byte, 65536)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			// Set a read deadline so we can periodically check ctx.Done()
			// instead of blocking forever on ReadFromUDP.
			s.conn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))

			n, _, err := s.conn.ReadFromUDP(buf)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				return
			}

			scanner := bufio.NewScanner(bytesReader(buf[:n]))
			for scanner.Scan() {
				line := scanner.Text()
				if line == "" {
					continue
				}
				select {
				case ch <- line:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return ch
}
