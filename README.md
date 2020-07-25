# PTP/IP protocol implementation in Go

This project started out as an implementation of
- the Picture Transfer Protocol (PTP) ISO-15740
- the PTP over IP protocol (PTP-IP) DC-X005-2005

while the goals were
- get to learn the Go programming language better; all code improvement
suggestions are thus always welcomed (there will undoubtedly be lots of those
;-) ) so feel free to create a PR
- stick to the Go standard packages and use as little external dependencies as
possible
- have a working PTP/IP implementation for a Fuji X-T1 camera on firmware 5.51
as a nice reverse engineering exercise :-p
- make the implementation such that the camera can be *programmed* using any
scripting language to get features such as: focus stacking, time-lapse
photography (and ultimately even triggers based on what the camera is actually
*seeing* using https://gocv.io)

## What it has become

### The `ptp` package
This package holds all of the PTP protocol related stuff. It's pretty basic for
now and needs some work to make a lot more usable as a stand alone package.
As the Fuji implementation deviates quite a bit from the PTP/IP standard, the
work on this package is somewhat limited because of the custom stuff needed to
talk to a Fuji device.

Whenever possible, this package is used to stick to the standard as much as
possible.

### The `ip` package
This one holds the IP transport layer implementation of the PTP protocol. As
with the `ptp` package, it is not fully developed because of the same reason:
the Fuji PTP/IP implementation takes lots of bits and pieces from the standard
but extends and drops just as much.

The Fuji bits are in `_fuji` files and any other future vendor that gets added
should use the same approach.

### The `fmt` package
All things related to formatting that are not at all part of the PTP nor PTP/IP
protocols are in here. The `ptp` and `ip` packages are meant to be usable
without the need for the `fmt` package.

### The `cmd` package
A command line interface implementation of the PTP/IP protocol that uses the
`ptp`, `ip` and `fmt` packages. See *CLI command* for further info.

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

### Usage
Executing the `ptpip` command without arguments or with the `-?` flag will
print its usage:
```text
Usage of ptpip:
  -?	Display usage information.
  -c string
        The command to send to the responder.
  -f string
        Read all settings from a config file. The config file will override any command line flags present.
  -g string
        A custom GUID to use for the initiator. (default random)
  -h string
        The responder host to connect to. (default "192.168.0.1")
  -n string
        A custom friendly name to use for the initiator.
  -p value
        The responder port to connect to. Use this flag when the responder has only ONE port for all channels! (default 15740)
  -pc value
        The responder port used for the Command/Data connection.
  -pe value
        The responder port used for the Event connection.
  -ps value
        The responder port used for the streamer or 'live view' connection.
  -s    This will run the ptpip command as a server
  -sa string
        To be used in combination with '-s': this defines the server address to listen on. (default "127.0.0.1")
  -sp value
        To be used in combination with '-s': this defines the server port to listen on. (default 15740)
  -t string
        The vendor of the responder that will be connected to. (default "generic")
  -v    Display version info.
```

### Config file
The config file is in the classic INI file format. Some examples:
```ini
; This is us
[initiator]
friendly_name = "Golang PTP/IP generic client"
; Generate a new random one using uuidgen or some other tool!
guid = "cca455de-79ac-4b12-9731-91e433a899cf"

; The target we will be connecting to
[responder]
host = "192.168.0.1"
port = 15740

; Config when running as a server
[server]
; Setting this to true will enable server mode
enabled = true
address = "127.0.0.1"
port = 15740
```
```ini
; This is us
[initiator]
friendly_name = "Golang PTP/IP Fuji client"
; Generate a new random one using uuidgen or some other tool!
guid = "9fe5160c-4951-404d-9505-10baaf725606"

; The target we will be connecting to
[responder]
vendor = "fuji"
cmd_data_port = 55740
event_port = 55741
stream_port = 55742

; Config when running as a daemon
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

### Supported commands

#### `capture`
This command will make the responder capture a single image. Some devices will
return a preview of the captured image. To save this preview to disk, you can
pass a path to write the preview to as the first parameter. E.g.:
```text
capture /tmp/my-preview.jpg
```
There are three aliases for this command: `shoot`, `shutter` and `snap`.

#### `info`
The info command will display the current info about the camera. The output
will vary from vendor to vendor.
There is one additional parameter for this command: `json`. It is no doubt
clear what it does: it will print the data as parsable JSON output, but again
it will differ from vendor to vendor!
Finally the `json` parameter itself has the option `pretty` to print indented
JSON output, e.g.:
```text
info json pretty
```

##### `get`
This command will request a property from the camera and return its current
value. The parameter defining the property can be a hexadecimal property code,
like `0x5005`, or a unified property name. The currently supported names are:
1. `delay`: delay before releasing shutter
2. `effect`: like sepia or other vendor specific effects or film simulations
3. `exposure`: exposure time
4. `exp-bias`: exposure bias compensation
5. `flashmode`
6. `focusmtr`: focus metering mode, or focus point
7. `iso`
8. `whitebalance`

TODO: add `recmode`!

#### `opreq`
This command is intended for reverse engineering and/or debugging purposes. It
takes two parameters in hexadecimal form: the first one is the operation code
to execute, and the second one is a parameter for the operation. Whether or not
this second parameter is mandatory depends on the operation being executed.
An example would be to describe (`0x1014`) a responder's image size property
(`0x5003`) by calling:
```text
opreq 0x1015 0x5003
```
The output will always be a **hexadecimal dump** of the packets received from the
responder.

See *server mode* below for example output.

#### `set`
This command will set a property on the camera to the requested value. The
first parameter indicating the property to be set, can be a hexadecimal
property code, like `0x5005`, or a unified property name. The currently
supported names are:
1. `delay`
2. `effect`
3. `exposure`
4. `exp-bias`
5. `flashmode`
6. `focusmtr`
7. `iso`
8. `whitebalance`

TODO: add `recmode`!

The second parameter is the value to set the property to. E.g.:
```text
set iso 800
```

#### `state`
This command is, for now, only supported by Fuji cameras and will display the
current state of a fixed list of camera dependent properties.
Like the `info` command, `state` also has the `json` parameter to output the
data in JSON parsable format with the additional `pretty` for indented JSON
output:
```text
state json pretty
```

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

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  06 00 00 00              |....... ....|
```
Take note that the `0x902B` code is Fuji specific and not part of the PTP/IP
standard!

As you can see the `opreq` command requires at least one parameter: the
operation code to perform which must be in hexadecimal notation.

It also supports an additional parameter, again in hex, to pass along with the
operation request. An example of executing `GetDevicePropValue` from the PTP
specification would be:
```text
$ nc 127.0.0.1 15740
opreq 0x1015 0xD212

Received 116 bytes. HEX dump:
00000000  74 00 00 00 02 00 15 10  06 00 00 00 11 00 01 50  |t..............P|
00000010  03 00 00 00 41 d2 0a 00  00 00 05 50 02 00 00 00  |....A......P....|
00000020  0a 50 01 80 00 00 0c 50  0a 80 00 00 0e 50 02 00  |.P.....P.....P..|
00000030  00 00 10 50 00 00 00 00  12 50 00 00 00 00 01 d0  |...P.....P......|
00000040  02 00 00 00 18 d0 04 00  00 00 28 d0 00 00 00 00  |..........(.....|
00000050  2a d0 00 19 00 80 7c d1  04 04 02 03 09 d2 00 00  |*.....|.........|
00000060  00 00 1b d2 00 00 00 00  29 d2 97 05 00 00 2a d2  |........).....*.|
00000070  8f 06 00 00                                       |....|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  06 00 00 00              |....... ....|
```
Again: the `0xD212` code is Fuji specific and not part of the PTP/IP standard!

The output depends on the command executed and can be one
single packet or, depending on the data phase, an *end of data* packet as well.

## Library
### Usage examples
Start by creating a new PTP IP client:
```go
package main

import(
    "github.com/malc0mn/ptp-ip/ip"
)

c := ip.NewClient(ip.DefaultVendor, "192.168.0.1", ip.DefaultPort, "MyClient", "", ip.LevelDebug)
```

### Credits

Projects that were used to realise this library:
- https://github.com/atotto/ptpip
- https://github.com/hkr/fuji-cam-wifi-tool