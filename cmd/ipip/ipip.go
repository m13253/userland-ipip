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

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	textwidth "github.com/m13253/go-textwidth"
	userland_ipip "github.com/m13253/userland-ipip/internal/userland-ipip"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	var (
		useIPv4       bool
		useIPv6       bool
		deviceName    *string
		localAddress  *string
		remoteAddress *string
		mtu           *uint16
	)
	state := 0

	for i, arg := range os.Args[1:] {
		switch {
		case (state == 0 || state == 1) && arg == "--":
			state = 1
		case state == 0 && arg == "--help":
			printHelp()
			return
		case state == 0 && arg == "-4":
			if useIPv6 {
				reportArguentError(i, "you cannot specify -4 and -6 at the same time")
			}
			useIPv4 = true
		case state == 0 && arg == "-6":
			if useIPv4 {
				reportArguentError(i, "you cannot specify -4 and -6 at the same time")
			}
			useIPv6 = true
		case (state == 0 || state == 1) && arg == "dev":
			if deviceName != nil {
				reportArguentError(i, "you already specified the TUN device name")
			}
			state = 2 | state
		case (state == 0 || state == 1) && arg == "local":
			if localAddress != nil {
				reportArguentError(i, "you already specified the local address")
			}
			state = 4 | state
		case (state == 0 || state == 1) && arg == "remote":
			if remoteAddress != nil {
				reportArguentError(i, "you already specified the remote address")
			}
			state = 6 | state
		case (state == 0 || state == 1) && arg == "mtu":
			if mtu != nil {
				reportArguentError(i, "you already specified the MTU")
			}
			state = 8 | state
		case state == 2 || state == 3:
			deviceName = &os.Args[i]
			state &= 1
		case state == 4 || state == 5:
			localAddress = &os.Args[i]
			state &= 1
		case state == 6 || state == 7:
			remoteAddress = &os.Args[i]
			state &= 1
		case state == 8 || state == 9:
			mtu64, err := strconv.ParseUint(arg, 0, 16)
			if err != nil {
				reportArguentError(i, "invalid MTU value")
			}
			mtu16 := uint16(mtu64)
			mtu = &mtu16
			state &= 1
		default:
			reportArguentError(i, "unknown option")
		}
	}
	if deviceName == nil {
		reportArguentError(0, "TUN device name not specified")
	}
	if localAddress == nil {
		localAddress = new(string)
	}
	if remoteAddress == nil {
		remoteAddress = new(string)
	}
	if !useIPv4 && !useIPv6 {
		useIPv4, useIPv6 = true, true
	}
	if mtu == nil {
		mtu = new(uint16)
	}

	err := userland_ipip.RunIPIP(*deviceName, *localAddress, *remoteAddress, useIPv4, useIPv6, *mtu)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Printf("Usage: %s [-4 | -6] dev DEVICE [local ADDRESS] [remote ADDRESS] [mtu MTU]\n", os.Args[0])
	fmt.Print("Userland IPIP / IPIP6 / IP6IP (6in4) / IP6IP6 tunnel for Linux.\n\n")
	fmt.Print("This program establishes IPIP and IP6IP (6in4) tunnel, or IPIP6 and IP6IP6 \n")
	fmt.Print("tunnel on a TUN device.\n\n")
	fmt.Print("Options:\n")
	fmt.Print("  -4                         use IPv4 to resolve addresses.\n")
	fmt.Print("  -6                         use IPv6 to resolve addresses.\n")
	fmt.Print("                               otherwise, IPv6 will be tried first, then IPv4.\n\n")
}

func reportArguentError(index int, reason string) {
	fmt.Fprintf(os.Stderr, "Command line error: %s\n\n", reason)

	if index == 0 {
		os.Exit(1)
	}
	startCol, endCol := 0, 0

	var b strings.Builder
	b.WriteByte('>')
	for i := 0; i < len(os.Args); i++ {
		b.WriteByte(' ')
		if i == index {
			_, startCol = textwidth.GetTextOffset(b.String(), 0, 0)
		}

		arg := os.Args[i]
		if strings.IndexAny(arg, "\x00\x01\x02\x03\x04\x05\x06\x07\x08\t\n\x0b\x0c\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e !\"#$&'()*,;<>?[\\]^`{|}~\x7f") == -1 {
			b.WriteString(arg)
		} else {
			b.WriteByte('\'')
			b.WriteString(strings.ReplaceAll(arg, "'", "'\\''"))
			b.WriteByte('\'')
		}

		if i == index {
			_, endCol = textwidth.GetTextOffset(b.String(), 0, 0)
		}
	}

	b.WriteByte('\n')
	for i := 0; i < startCol; i++ {
		b.WriteByte(' ')
	}
	for i := startCol; i < endCol; i++ {
		b.WriteByte('~')
	}
	b.WriteString("\n\n")
	fmt.Print(b.String())

	os.Exit(1)
}
