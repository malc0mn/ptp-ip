The Picture Transfer Protocol (PTP) ISO-15740
The PTP over IP protocol (PTP-IP) DC-X005-2005

## Connecting to your camera
The first and obvious step is to enable the camera's wifi. Have your network
manager scan for new SSIDs and connect to the one from your camera. It will most
likely have an obvious name. A Fujifilm X-T1 SSID, for example, starts with
`FUJIFILM-X-T1` followed by four more characters.
## Linux `NetworkManager` troubleshooting
If you have trouble establishing a WiFi connection to your camera, start off by
tailing the logs: `sudo journalctl -f`. When those are open, connect to your
camera's SSID and look closely at what it spews out.
### IPv6 errors
If you see IPv6 related errors in the logs, make sure you disable IPv6 for the
connection to your camera's SSID. You can do this using the UI: under the IPv6
settings select the method `disabled`. Or you can edit the config file
directly:
`sudo vi /etc/NetworkManager/system-connections/[SSID].nmconnection`.

Look for the `[ipv6]` section or add it if it's not there and make sure that
this line is present: `method=disabled`.
### Cannot get an IP from the camera's DHCP server
If you are using `NetworkManager` with its built-in, and rather buggy, DHCP
client, you might have trouble getting a DHCP address from the camera.
In this case, you could try using `dhclient` as follows:
1. Make sure `dhclient` is installed: simply run `dhclient --version` from the
CLI and if you see output in the sense of `isc-dhclient-4.4.2` then you you're
good to go :-). If not, install it first.
2. Now let's tell `NetworkManager` to use it by adding some config:
`sudo vi /etc/NetworkManager/conf.d/dhcp-client.conf`
3. Paste the following config and save it:
```text
[main]
dhcp=dhclient
```
4. Finish by restarting the `NetworkManager` service:
`sudo systemctl restart NetworkManager`
5. Do another connection attempt to your camera's SSID, which should now
complete as expected.
6. If you're still having trouble connecting, stop as many applications and/or
services as possible that might be fighting for DNS requests or other network
related things, as the camera could simply be overwhelmed and has no time or
resources to hand out an IP address.
The first shutdown candidates here are any web browser (Chrome, Firefox, Edge)
or chat applications such as Slack, WhatsApp etc.

## CLI command
### Config file
The config file is in the classic INI file format.
```ini
; This is us
[initiator]
friendly_name = "Golang PTP/IP client"
; Generate a new random one using uuidgen or some other tool!
guid = "cca455de-79ac-4b12-9731-91e433a899cf"

; The target we will be connecting to
[responder]
vendor = "fuji"
host = "192.168.0.1"
port = 15740

; Config when running as a server
[server]
; Setting this to true will enable server mode
enabled = true
address = "127.0.0.1"
port = 15740
```

### Exit codes
Depending on the error, the exit code of the `ptpip` command will differ:
1. Unspecified: `1`
2. Error opening config file: `102`
3. Error creating client: `104`
4. Error connecting to responder: `105`

### Server mode
When executing the command with the `-s` flag, it will first connect to your
specified camera and when that succeeds a socket is opened on `127.0.0.1`
port `15740`, unless you specified a custom listen address and/or port using
the `-sa` and `-sp` flags.
As soon as you see output along the lines of:
```text
[Local server] listening on 127.0.0.1:15740...
[Local server] awaiting messages... (CTRL+C to quit)
```
you can start sending messages. You can use any language you like to
communicate with the socket. From a linux command line interface such as bash,
you can simply use `nc` to connect and send a message.
An example of using the `opreq` command for debugging or reverse engineering
purposes would be:
```text
$ nc 127.0.0.1 15740
opreq 0x902B
Received 356 bytes. HEX dump:
00000000  64 01 00 00 02 00 2b 90  07 00 00 00 08 00 00 00  |d.....+.........|
00000010  16 00 00 00 12 50 04 00  01 00 00 00 00 02 03 00  |.....P..........|
00000020  00 00 02 00 04 00 14 00  00 00 0c 50 04 00 01 02  |...........P....|
00000030  00 09 80 02 02 00 09 80  0a 80 24 00 00 00 05 50  |..........$....P|
00000040  04 00 01 02 00 02 00 02  0a 00 02 00 04 00 06 80  |................|
00000050  01 80 02 80 03 80 06 00  0a 80 0b 80 0c 80 36 00  |..............6.|
00000060  00 00 10 50 03 00 01 00  00 00 00 02 13 00 48 f4  |...P..........H.|
00000070  95 f5 e3 f6 30 f8 7d f9  cb fa 18 fc 65 fd b3 fe  |....0.}.....e...|
00000080  00 00 4d 01 9b 02 e8 03  35 05 83 06 d0 07 1d 09  |..M.....5.......|
00000090  6b 0a b8 0b 26 00 00 00  01 d0 04 00 01 01 00 02  |k...&...........|
000000a0  00 02 0b 00 01 00 02 00  03 00 04 00 05 00 06 00  |................|
000000b0  07 00 08 00 09 00 0a 00  0b 00 78 00 00 00 2a d0  |..........x...*.|
000000c0  06 00 01 ff ff ff ff 00  19 00 80 02 19 00 90 01  |................|
000000d0  00 80 20 03 00 80 40 06  00 80 80 0c 00 80 00 19  |.. ...@.........|
000000e0  00 80 64 00 00 40 c8 00  00 00 fa 00 00 00 40 01  |..d..@........@.|
000000f0  00 00 90 01 00 00 f4 01  00 00 80 02 00 00 20 03  |.............. .|
00000100  00 00 e8 03 00 00 e2 04  00 00 40 06 00 00 d0 07  |..........@.....|
00000110  00 00 c4 09 00 00 80 0c  00 00 a0 0f 00 00 88 13  |................|
00000120  00 00 00 19 00 00 00 32  00 40 00 64 00 40 00 c8  |.......2.@.d.@..|
00000130  00 40 14 00 00 00 19 d0  04 00 01 01 00 01 00 02  |.@..............|
00000140  02 00 00 00 01 00 1e 00  00 00 7c d1 06 00 01 00  |..........|.....|
00000150  00 00 00 04 04 02 03 01  00 00 00 00 07 07 09 10  |................|
00000160  01 00 00 00                                       |....|
```
As you can see the `opreq` command requires at least one parameter: the
operation code to perform which must be in hexadecimal notation. It also
supports an additional parameter, again in hex, to pass along with the
operation request.

## Library
### Usage examples
Start by creating a new PTP IP client:
```go
package main

import(
    "github.com/malc0mn/ptp-ip/ip"
)

c := NewClient("192.168.0.1", ip.DefaultPort, "MyClient", "")
```

### Credits

Projects that were used to realise this library:
- https://github.com/atotto/ptpip
- https://github.com/hkr/fuji-cam-wifi-tool