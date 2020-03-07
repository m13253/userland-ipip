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
	"fmt"
	"net"

	"golang.org/x/sys/unix"
)

func newIPIPConn(local, remote *net.IPAddr) (ip4ip, ip6ip *net.IPConn, err error) {
	bothIPv4, bothIPv6 := true, true
	if local != nil {
		if ipv4 := local.IP.To4(); ipv4 != nil {
			local = &net.IPAddr{
				IP:   ipv4,
				Zone: local.Zone,
			}
			bothIPv6 = false
		} else {
			bothIPv4 = false
		}
	}
	if remote != nil {
		if ipv4 := remote.IP.To4(); ipv4 != nil {
			remote = &net.IPAddr{
				IP:   ipv4,
				Zone: remote.Zone,
			}
			bothIPv6 = false
		} else {
			bothIPv4 = false
		}
	}

	if bothIPv6 {
		if remote != nil {
			ip6ip, err = net.DialIP("ip6:41", local, remote)
			if err != nil {
				return
			}
			ip4ip, err = net.DialIP("ip6:4", local, remote)
			if err != nil {
				ip6ip.Close()
				return
			}
		} else {
			ip6ip, err = net.ListenIP("ip6:41", local)
			if err != nil {
				return
			}
			ip4ip, err = net.ListenIP("ip6:4", local)
			if err != nil {
				ip6ip.Close()
				return
			}
		}
		return
	}

	if bothIPv4 {
		if remote != nil {
			ip6ip, err = net.DialIP("ip4:41", local, remote)
			if err != nil {
				return
			}
			ip4ip, err = net.DialIP("ip4:94", local, remote)
			if err != nil {
				ip6ip.Close()
				return
			}
		} else {
			ip6ip, err = net.ListenIP("ip4:41", local)
			if err != nil {
				return
			}
			ip4ip, err = net.ListenIP("ip4:94", local)
			if err != nil {
				ip6ip.Close()
				return
			}
		}
	}

	return nil, nil, fmt.Errorf("local and remote addresses are not from the same address family: %v", unix.EAFNOSUPPORT)
}
