# userland-ipip

Userland IPIP + IP6IP (6in4) / IPIP6 + IP6IP6 tunnel for Linux

userland-ipip sets up either an IPIP + IP6IP (6in4) tunnel, or an IPIP6 + IP6IP6
tunnel between two hosts.

## Features

- userland-ipip reduces headache when you find that a `type ip6tnl mode any`
  tunnel is not as reliable as you assume. Either IPv4 or IPv6 payload drops
  silently at some magic time. You tried various methods, only to find that a
  system reboot can solve the problem auto-magically.

- userland-ipip calculates `local` address automatically, saving you time to
  write scripts for an DHCP-assigned host.

- userland-ipip also solves the problem when you want to fragment your tunnel.
  (i.e. inner MTU larger than outer MTU.)

## Building

1. Download Go compiler. The newer version, the better.

2. Type
```bash
./build.sh
```

3. Pick your fruit at `./build/ipip`.

## Usage

```
Usage: ipip [-4 | -6] dev DEVICE [local ADDRESS] remote ADDRESS [mtu MTU]
Userland IPIP + IP6IP (6in4) / IPIP6 + IP6IP6 tunnel for Linux.

This program establishes IPIP and IP6IP (6in4) tunnel, or IPIP6 and IP6IP6
tunnel on a TUN device.

Options:
  -4            use IPv4 to resolve addresses.
  -6            use IPv6 to resolve addresses.
                  otherwise, IPv6 will be tried first, then IPv4.

Project web page: https://github.com/m13253/userland-ipip
```

## Example

Please change the names and the addresses below to suit your needs.

On the first machine (e.g. fox.localdomain)
```bash
sudo ip tuntap add mode tun name tun-rabbit
sudo ip address add 10.0.0.1 peer 10.0.0.2/32 dev tun-rabbit
sudo ip address add fd00:cafe::1 peer fd00:cafe::2/128 dev tun-rabbit
sudo ./build/ipip dev tun-rabbit remote rabbit.localdomain mtu 1460
```

On the second machine (e.g. rabbit.localdomain)
```bash
sudo ip tuntap add mode tun name tun-fox
sudo ip address add 10.0.0.2 peer 10.0.0.1/32 dev tun-fox
sudo ip address add fd00:cafe::2 peer fd00:cafe::1/128 dev tun-fox
sudo ./build/ipip dev tun-fox remote fox.localdomain mtu 1460
```

To stop the tunnel, press `Ctrl-C`, then type
```bash
sudo ip link delete tun-rabbit
```
or
```bash
sudo ip link delete tun-fox
```

## Preventing “connection refused”

You may find a lot of “connection refused” on the screen. They are caused by
the remote machine sending ICMP errors to us.

It is suggested to block these packets to save bandwidth. A dirty but effective
method is to use iptables on both sides running userland-ipip:
```bash
sudo iptables -A OUTPUT -d [PEER IPv4 ADDRESS] -p icmp --icmp-type 3/3 -j DROP
sudo ip6tables -A OUTPUT -d [PEER IPv6 ADDRESS] -p icmpv6 --icmpv6-type 1/4 -j DROP
```

## Use userland-ipip with systemd

I don't provide a systemd service file out-of-the-box, since you may want to
write one systemd service for each tunnel you want to create.

Here is a template that you can modify based on:
```systemd
[Unit]
Description=Userland IPIP for rabbit.localdomain
Documentation=https://github.com/m13253/userland-ipip
After=network.target

[Service]
ExecStartPre=-/usr/bin/env ip tunnel delete tun-rabbit
ExecStartPre=/usr/bin/env ip tuntap add mode tun name tun-rabbit
ExecStartPre=/usr/bin/env ip address add 10.0.0.1 peer 10.0.0.2/32 dev tun-rabbit
ExecStartPre=/usr/bin/env ip address add fd00:cafe::1 peer fd00:cafe::2/128 dev tun-rabbit
ExecStart=/path/to/ipip dev tun-rabbit local fox.localdomain remote rabbit.localdomain mtu 1460
ExecStopPost=/usr/bin/env ip tunnel delete tun-rabbit
Restart=always
RestartSec=3
Type=simple

[Install]
WantedBy=multi-user.target
```

## Use userland-ipip with `/etc/network/interfaces`

```conf
auto tun-rabbit
iface tun-rabbit inet static
    address 10.0.0.1
    pointopoint 10.0.0.2
    pre-up ip tuntap add mode tun name $IFACE
    up /path/to/ipip dev $IFACE local fox.localdomain remote rabbit.localdomain mtu 1460 &
    post-down ip link del $IFACE
iface tun-rabbit inet6 static
    address fd00:cafe::1/128
    up ip route add fd00:cafe::2 dev $IFACE metric 256
```

## License

This program is released under GNU General Public License version 3 or later.
I hope this program can be useful to you. But I provide **absolutely no
warranty**. In case the program causes any damage due to malfunctioning, I might
be willing to diagnose and fix the problem, but it is not my obligation to
do so.
