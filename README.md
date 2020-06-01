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
related things, as the camera might be overwhelmed and simply has no time or
resources to hand out an IP address.

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