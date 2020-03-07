# userland-ipip

Userland IPIP + IP6IP (6in4) / IPIP6 + IP6IP6 tunnel for Linux

## Features

userland-ipip reduces headache when you find that a `type ip6tnl mode any`
tunnel is not as reliable as you assume. Either IPv4 or IPv6 payload drops
silently at some magic time. You tried various methods, only to find that a
system reboot can solve the problem auto-magically.

userland-ipip calculates `local` address automatically, saving you time to write
scripts for an DHCP-assigned host.

userland-ipip also solves the problem when you want to fragment your tunnel.
(i.e. inner MTU larger than outer MTU.)

## Building

1. Download Go compiler. The newer version, the better.

2. Type
```sh
./build.sh
```

3. Pick your fruit at `./build/ipip`.

## Usage

Please change the names and the addresses below to suit your needs.

On the first machine (e.g. fox.localdomain)
```sh
sudo ip tuntap add mode tun name tun-rabbit
sudo ip address add 10.0.0.1 peer 10.0.0.2/32 dev tun-rabbit
sudo ip address add fd00:cafe::1 peer fd00:cafe::2/128 dev tun-rabbit
sudo ./build/ipip dev tun-rabbit remote rabbit.localdomain mtu 1460
```

On the second machine (e.g. rabbit.localdomain)
```sh
sudo ip tuntap add mode tun name tun-fox
sudo ip address add 10.0.0.2 peer 10.0.0.1/32 dev tun-fox
sudo ip address add fd00:cafe::2 peer fd00:cafe::1/128 dev tun-fox
sudo ./build/ipip dev tun-fox remote fox.localdomain mtu 1460
```

To stop the tunnel, press `Ctrl-C`, then type
```sh
sudo ip link delete tun-rabbit
```
or
```sh
sudo ip link delete tun-fox
```

## Preventing "connection refused"

You may find a lot of "connection refused" on the screen. They are caused by
the remote machine sending ICMP errors to us.

It is suggested to block these packets to save bandwidth. A dirty but effective
method is to use iptables on both sides running userland-ipip:
```sh
sudo iptables -A OUTPUT -d [PEER IPv4 ADDRESS] -p icmp --icmp-type 3/3 -j DROP
```

Luckily the problem does not happen over IPv6.

## License

This program is released under GNU General Public License version 3 or later.
I hope this program can be useful to you. But I provide **absolutely no
warranty**. In case the program causes any damage due to malfunctioning, I might
be willing to diagnose and fix the problem, but it is not my responsibility to
do so.
