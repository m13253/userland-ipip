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

import "golang.org/x/sys/unix"

const (
	_IFF_TUN         = 0x0001
	_IFF_TAP         = 0x0002
	_IFF_MULTI_QUEUE = 0x0100
)

type (
	ifreq_flags struct {
		ifr_name  [unix.IFNAMSIZ]byte
		ifr_flags int16
		_         int16
		_         [20]byte
	}
	ifreq_mtu struct {
		ifr_name [unix.IFNAMSIZ]byte
		ifr_mtu  int32
		_        [20]byte
	}
)
