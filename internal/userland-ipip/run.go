// Userland-IPIP
// Copyright (C) 2020  StarBrilliant
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package userland_ipip

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

func RunIPIP(devName, local, remote string, useIPv4, useIPv6 bool, mtu uint16) error {
	var ip4ip, ip6ip *net.IPConn
	var err, err1, err2, err3, err4 error

	if useIPv6 {
		var localAddr, remoteAddr *net.IPAddr
		if len(local) != 0 {
			localAddr, err1 = net.ResolveIPAddr("ip6", local)
		}
		if err1 == nil && len(remote) != 0 {
			remoteAddr, err2 = net.ResolveIPAddr("ip6", remote)
		}

		if err1 == nil && err2 == nil {
			ip4ip, ip6ip, err = newIPIPConn("ip6", localAddr, remoteAddr)
			if err != nil {
				return fmt.Errorf("failed to start IPIP6 and IP6IP6: %v", err)
			}
			defer ip4ip.Close()
			defer ip6ip.Close()
		}
	}

	if ip4ip == nil && ip6ip == nil && useIPv4 {
		var localAddr, remoteAddr *net.IPAddr
		if len(local) != 0 {
			localAddr, err3 = net.ResolveIPAddr("ip4", local)
		}
		if err3 == nil && len(remote) != 0 {
			remoteAddr, err4 = net.ResolveIPAddr("ip4", remote)
		}

		if err3 == nil && err4 == nil {
			ip4ip, ip6ip, err = newIPIPConn("ip4", localAddr, remoteAddr)
			if err != nil {
				return fmt.Errorf("failed to start IPIP and IP6IP (6in4): %v", err)
			}
			defer ip4ip.Close()
			defer ip6ip.Close()
		}
	}

	if ip4ip == nil && ip6ip == nil {
		if err1 != nil {
			err = err1
		} else if err2 != nil {
			err = err2
		} else if err3 != nil {
			err = err3
		} else if err4 != nil {
			err = err4
		}
		if err != nil {
			return fmt.Errorf("failed to resolve address: %v", err)
		}
	}

	tun, err := newTunDevice(devName, mtu)
	if err != nil {
		return fmt.Errorf("failed to create TUN device: %v", err)
	}
	defer tun.Close()

	fmt.Printf("Tunnel created, local %v, remote %v\n", ip6ip.LocalAddr(), ip6ip.RemoteAddr())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errChan := make(chan error)
	go forwardIP6ToTun(ctx, tun, ip6ip, errChan)
	go forwardIP4ToTun(ctx, tun, ip4ip, errChan)
	go forwardTunToIP(ctx, ip4ip, ip6ip, tun, errChan)

	err = <-errChan
	return err
}

func forwardTunToIP(ctx context.Context, ip4ip, ip6ip *net.IPConn, tun *os.File, errChan chan<- error) {
	var buf [65540]byte

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		n, err := tun.Read(buf[:])
		if err != nil {
			errChan <- fmt.Errorf("failed to read from TUN device: %v", err)
			return
		}
		if n == 0 {
			errChan <- nil
			return
		}
		if n < 4 {
			continue
		}

		ethertype := binary.BigEndian.Uint16(buf[2:4])
		packet := buf[4:n]

		switch ethertype {
		case etherTypeIPv4:
			_, err = ip4ip.Write(packet)
			if err != nil {
				errChan <- fmt.Errorf("failed to send IPv4 tunneled data: %v", err)
				return
			}
		case etherTypeIPv6:
			_, err = ip6ip.Write(packet)
			if err != nil {
				errChan <- fmt.Errorf("failed to send IPv6 tunneled data: %v", err)
				return
			}
		}
	}
}

func forwardIP4ToTun(ctx context.Context, tun *os.File, ip4ip *net.IPConn, errChan chan<- error) {
	var buf [65540]byte
	binary.BigEndian.PutUint16(buf[2:4], etherTypeIPv4)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		n, _, err := ip4ip.ReadFromIP(buf[4:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "IPv4 tunnel returns: %v (maybe caused by ICMP)\n", err)
			continue
		}
		if n == 0 {
			errChan <- nil
			return
		}

		_, err = tun.Write(buf[:n+4])
		if err != nil {
			errChan <- fmt.Errorf("failed to write to TUN device: %v", err)
			return
		}
	}
}

func forwardIP6ToTun(ctx context.Context, tun *os.File, ip6ip *net.IPConn, errChan chan<- error) {
	var buf [65540]byte
	binary.BigEndian.PutUint16(buf[2:4], etherTypeIPv6)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		n, _, err := ip6ip.ReadFromIP(buf[4:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "IPv6 tunnel returns: %v (maybe caused by ICMP)\n", err)
			continue
		}
		if n == 0 {
			errChan <- nil
			return
		}

		_, err = tun.Write(buf[:n+4])
		if err != nil {
			errChan <- fmt.Errorf("failed to write to TUN device: %v", err)
			return
		}
	}
}
