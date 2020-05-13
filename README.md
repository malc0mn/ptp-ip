The Picture Transfer Protocol (PTP) ISO-15740
The PTP over IP protocol (PTP-IP) DC-X005-2005

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
host = "192.168.0.1"
port = 15740

; Config when running as a server
[server]
; Setting this to true will enable server mode
enabled = true
address = "127.0.0.1"
port = 15740
```

### Usage examples
Start by creating a new PTP IP client:
```go
package main

import(
    "github.com/google/uuid"
    "github.com/malc0mn/ptp-ip/ip"
)

c := NewClient("192.168.0.1", ip.DefaultPort, "MyClient", uuid.Nil)
```

### Credits

Projects that were used to realise this library:
- https://github.com/atotto/ptpip
- https://github.com/hkr/fuji-cam-wifi-tool