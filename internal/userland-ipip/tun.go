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
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

func newTunDevice(name string, mtu uint16) (*os.File, error) {
	ifreqSock, err := syscall.Socket(syscall.AF_UNIX, syscall.SOCK_DGRAM|syscall.SOCK_CLOEXEC, 0)
	if err != nil {
		return nil, os.NewSyscallError("socket", err)
	}

	f, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0666)
	if err != nil {
		syscall.Close(ifreqSock)
		return nil, err
	}

	ifreq_flags := &ifreq_flags{}
	copy(ifreq_flags.ifr_name[:], name)
	ifreq_flags.ifr_flags = unix.IFF_TUN | unix.IFF_MULTI_QUEUE

	r1, _, err := syscall.Syscall(unix.SYS_IOCTL, f.Fd(), unix.TUNSETIFF, uintptr(unsafe.Pointer(ifreq_flags)))
	if r1 != 0 {
		f.Close()
		syscall.Close(ifreqSock)
		return nil, os.NewSyscallError("ioctl (TUNSETIFF)", err)
	}

	if mtu != 0 {
		ifreq_mtu := &ifreq_mtu{}
		copy(ifreq_mtu.ifr_name[:], ifreq_flags.ifr_name[:])
		ifreq_mtu.ifr_mtu = int32(mtu)

		r1, _, err = syscall.Syscall(unix.SYS_IOCTL, uintptr(ifreqSock), unix.SIOCSIFMTU, uintptr(unsafe.Pointer(ifreq_mtu)))
		if r1 != 0 {
			f.Close()
			syscall.Close(ifreqSock)
			return nil, os.NewSyscallError("ioctl (SIOCSIFMTU)", err)
		}
	}

	r1, _, err = syscall.Syscall(unix.SYS_IOCTL, uintptr(ifreqSock), unix.SIOCGIFFLAGS, uintptr(unsafe.Pointer(ifreq_flags)))
	if r1 != 0 {
		f.Close()
		syscall.Close(ifreqSock)
		return nil, os.NewSyscallError("ioctl (SIOCGIFFLAGS)", err)
	}
	ifreq_flags.ifr_flags |= unix.IFF_UP
	r1, _, err = syscall.Syscall(unix.SYS_IOCTL, uintptr(ifreqSock), unix.SIOCSIFFLAGS, uintptr(unsafe.Pointer(ifreq_flags)))
	if r1 != 0 {
		f.Close()
		syscall.Close(ifreqSock)
		return nil, os.NewSyscallError("ioctl (SIOCSIFFLAGS)", err)
	}

	syscall.Close(ifreqSock)
	return f, nil
}
