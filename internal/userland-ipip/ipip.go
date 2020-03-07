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

func newIPIPConn(network string, local, remote *net.IPAddr) (ip4ip, ip6ip *net.IPConn, err error) {
	if network == "ip6" {
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

	if network == "ip4" {
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
		return
	}

	return nil, nil, fmt.Errorf("invalid address family: %v", unix.EAFNOSUPPORT)
}
