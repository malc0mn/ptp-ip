## Fujifilm X-T1 Device Proprty Codes list
### Supported from the PTP/IP standard
#### Property Code 0x5001 - Battery Level
```text
opreq 0x1015 0x5001


Received 13 bytes. HEX dump:
00000000  0d 00 00 00 02 00 15 10  01 52 00 00 00           |.........R...|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  01 52 00 00              |....... .R..|
```

#### Property Code 0x5003 - Image Size
```text
opreq 0x1015 0x5003

Received 165 bytes. HEX dump:
00000000  a5 00 00 00 02 00 15 10  03 52 00 00 4c 5f 34 2d  |.........R..L_4-|
00000010  33 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |3...............|
00000020  00 00 00 00 00 00 00 00  33 32 30 78 32 38 30 00  |........320x280.|
00000030  64 00 00 00 06 00 00 00  64 00 16 00 94 16 cb 01  |d.......d.......|
00000040  00 00 01 00 13 00 00 00  ec 12 cb 01 05 00 00 00  |................|
00000050  3c 19 cb 01 38 00 00 00  90 11 cb 01 03 00 00 00  |<...8...........|
00000060  68 1e cb 01 03 00 00 00  c8 18 cb 01 09 00 00 00  |h...............|
00000070  78 1c cb 01 05 00 00 00  34 82 8b 01 05 00 00 00  |x.......4.......|
00000080  49 82 8b 01 1f 00 00 00  4e 82 8b 01 03 00 02 00  |I.......N.......|
00000090  02 00 30 30 ff ff ff ff  ff ff ff ff ff ff ff ff  |..00............|
000000a0  ff ff ff ff ff                                    |.....|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  03 52 00 00              |....... .R..|
```

#### Property Code 0x5005 - White Balance
```text
opreq 0x1015 0x5005

Received 14 bytes. HEX dump:
00000000  0e 00 00 00 02 00 15 10  05 52 00 00 00 00        |.........R....|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  05 52 00 00              |....... .R..|
```

#### Property Code 0x5007 - F-Number
```text
opreq 0x1015 0x5007

Received 14 bytes. HEX dump:
00000000  0e 00 00 00 02 00 15 10  07 52 00 00 00 00        |.........R....|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  07 52 00 00              |....... .R..|
```

#### Property Code 0x500a - Focus Mode
```text
opreq 0x1015 0x500a

Received 14 bytes. HEX dump:
00000000  0e 00 00 00 02 00 15 10  0a 52 00 00 00 00        |.........R....|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  0a 52 00 00              |....... .R..|
```

#### Property Code 0x500c - Flash Mode
```text
opreq 0x1015 0x500c

Received 14 bytes. HEX dump:
00000000  0e 00 00 00 02 00 15 10  0c 52 00 00 02 00        |.........R....|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  0c 52 00 00              |....... .R..|
```

#### Property Code 0x500d - Exposure Time
```text
opreq 0x1015 0x500d

Received 16 bytes. HEX dump:
00000000  10 00 00 00 02 00 15 10  0d 52 00 00 00 00 00 00  |.........R......|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  0d 52 00 00              |....... .R..|
```

#### Property Code 0x500e - Exposure Program Mode
```text
opreq 0x1015 0x500e

Received 14 bytes. HEX dump:
00000000  0e 00 00 00 02 00 15 10  0e 52 00 00 00 00        |.........R....|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  0e 52 00 00              |....... .R..|
```

#### Property Code 0x5010 - Exposure Bias Compensation
```text
opreq 0x1015 0x5010

Received 14 bytes. HEX dump:
00000000  0e 00 00 00 02 00 15 10  10 52 00 00 00 00        |.........R....|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  10 52 00 00              |....... .R..|
```

#### Property Code 0x5012 - Capture Delay
```text
opreq 0x1015 0x5012

Received 14 bytes. HEX dump:
00000000  0e 00 00 00 02 00 15 10  12 52 00 00 00 00        |.........R....|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  12 52 00 00              |....... .R..|
```

### Fuji Extensions

#### Property Code 0xd001 - Film Simulation
```text
opreq 0x1015 0xd001

Received 14 bytes. HEX dump:
00000000  0e 00 00 00 02 00 15 10  01 d2 00 00 00 00        |..............|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  01 d2 00 00              |....... ....|
```

#### Property Code 0xd018 - Image Quality
```text
opreq 0x1015 0xd018

Received 14 bytes. HEX dump:
00000000  0e 00 00 00 02 00 15 10  18 d2 00 00 01 00        |..............|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  18 d2 00 00              |....... ....|
```

#### Property Code 0xd019 - Rec Mode
```text
opreq 0x1015 0xd019

Received 14 bytes. HEX dump:
00000000  0e 00 00 00 02 00 15 10  19 d2 00 00 01 00        |..............|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  19 d2 00 00              |....... ....|
```

#### Property Code 0xd01d - [still unknown]
```text
opreq 0x1015 0xd01d

Received 14 bytes. HEX dump:
00000000  0e 00 00 00 02 00 15 10  1d d2 00 00 01 00        |..............|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  1d d2 00 00              |....... ....|
```

#### Property Code 0xd028 - Command Dial Mode
```text
opreq 0x1015 0xd028

Received 14 bytes. HEX dump:
00000000  0e 00 00 00 02 00 15 10  28 d2 00 00 00 00        |........(.....|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  28 d2 00 00              |....... (...|
```

#### Property Code 0xd02a - Exposure Index
```text
opreq 0x1015 0xd02a

Received 16 bytes. HEX dump:
00000000  10 00 00 00 02 00 15 10  2a d2 00 00 ff ff ff ff  |........*.......|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  2a d2 00 00              |....... *...|
```

#### Property Code 0xd170 - [still unknown]
```text
opreq 0x1015 0xd170

Received 14 bytes. HEX dump:
00000000  0e 00 00 00 02 00 15 10  70 d3 00 00 01 00        |........p.....|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  70 d3 00 00              |....... p...|
```

#### Property Code 0xd174 - Looks like the Fuji version of the PTP/IP Image Size
property.
```text
opreq 0x1015 0xd174

Received 115 bytes. HEX dump:
00000000  73 00 00 00 02 00 15 10  12 00 00 00 33 32 30 78  |s...........320x|
00000010  32 38 30 00 64 00 00 00  06 00 00 00 64 00 16 00  |280.d.......d...|
00000020  94 16 cb 01 00 00 01 00  13 00 00 00 ec 12 cb 01  |................|
00000030  05 00 00 00 3c 19 cb 01  38 00 00 00 90 11 cb 01  |....<...8.......|
00000040  03 00 00 00 68 1e cb 01  03 00 00 00 c8 18 cb 01  |....h...........|
00000050  09 00 00 00 78 1c cb 01  05 00 00 00 34 82 8b 01  |....x.......4...|
00000060  05 00 00 00 49 82 8b 01  1f 00 00 00 4e 82 8b 01  |....I.......N...|
00000070  03 00 02                                          |...|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  12 00 00 00              |....... ....|
```

#### Property Code 0xd17c - Focus Point
```text
opreq 0x1015 0xd17c

Received 16 bytes. HEX dump:
00000000  10 00 00 00 02 00 15 10  7c d3 00 00 00 00 00 00  |........|.......|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  7c d3 00 00              |....... |...|
```

#### Property Code 0xd212 - Current State
```text
opreq 0x1015 0xd212

Received 116 bytes. HEX dump:
00000000  74 00 00 00 02 00 15 10  12 d4 00 00 11 00 01 50  |t..............P|
00000010  02 00 00 00 41 d2 0a 00  00 00 05 50 02 00 00 00  |....A......P....|
00000020  0a 50 01 80 00 00 0c 50  0a 80 00 00 0e 50 02 00  |.P.....P.....P..|
00000030  00 00 10 50 b3 fe 00 00  12 50 00 00 00 00 01 d0  |...P.....P......|
00000040  02 00 00 00 18 d0 04 00  00 00 28 d0 00 00 00 00  |..........(.....|
00000050  2a d0 00 19 00 80 7c d1  02 07 02 03 09 d2 00 00  |*.....|.........|
00000060  00 00 1b d2 00 00 00 00  29 d2 d6 05 00 00 2a d2  |........).....*.|
00000070  8f 06 00 00                                       |....|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  12 d4 00 00              |....... ....|
```

#### Property Code 0xd21b - Device Error
```text
opreq 0x1015 0xd21b

Received 14 bytes. HEX dump:
00000000  0e 00 00 00 02 00 15 10  1b d4 00 00 00 00        |..............|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  1b d4 00 00              |....... ....|
```

#### Property Code 0xd220 - [still unknown]
```text
opreq 0x1015 0xd220

Received 16 bytes. HEX dump:
00000000  10 00 00 00 02 00 15 10  20 d4 00 00 00 00 00 00  |........ .......|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  20 d4 00 00              |.......  ...|
```

#### Property Code 0xd222 - [still unknown]
```text
opreq 0x1015 0xd222

Received 16 bytes. HEX dump:
00000000  10 00 00 00 02 00 15 10  22 d4 00 00 00 00 00 00  |........".......|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  22 d4 00 00              |....... "...|
```

#### Property Code 0xd226 - [still unknown]
```text
opreq 0x1015 0xd226

Received 14 bytes. HEX dump:
00000000  0e 00 00 00 02 00 15 10  26 d4 00 00 00 00        |........&.....|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  26 d4 00 00              |....... &...|
```

#### Property Code 0xd227 - [still unknown]
```text
opreq 0x1015 0xd227

Received 14 bytes. HEX dump:
00000000  0e 00 00 00 02 00 15 10  27 d4 00 00 00 00        |........'.....|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  27 d4 00 00              |....... '...|
```

#### Property Code 0xd228 - Clearly shows device model and firmware version but
what else...?
```text
opreq 0x1015 0xd228

Received 77 bytes. HEX dump:
00000000  4d 00 00 00 02 00 15 10  28 d4 00 00 20 30 2f 30  |M.......(... 0/0|
00000010  00 00 00 00 8f 06 00 00  00 00 00 00 00 00 00 00  |................|
00000020  02 00 00 00 02 07 02 03  0a 00 02 00 ff ff ff ff  |................|
00000030  00 00 00 00 58 2d 54 31  00 00 00 00 00 00 00 00  |....X-T1........|
00000040  00 00 00 00 00 00 00 00  00 35 2e 35 31           |.........5.51|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  28 d4 00 00              |....... (...|
```

#### Property Code 0xd229 - Image Space SD
```text
opreq 0x1015 0xd229

Received 16 bytes. HEX dump:
00000000  10 00 00 00 02 00 15 10  29 d4 00 00 d6 05 00 00  |........).......|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  29 d4 00 00              |....... )...|
```

#### Property Code 0xd22a - Movie Remaining Time
```text
opreq 0x1015 0xd22a

Received 16 bytes. HEX dump:
00000000  10 00 00 00 02 00 15 10  2a d4 00 00 8f 06 00 00  |........*.......|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  2a d4 00 00              |....... *...|
```

#### Property Code 0xd240 - Shutter Speed
```text
opreq 0x1015 0xd240

Received 16 bytes. HEX dump:
00000000  10 00 00 00 02 00 15 10  40 d4 00 00 ff ff ff ff  |........@.......|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  40 d4 00 00              |....... @...|
```

#### Property Code 0xd241 - Image Aspect Ratio
```text
opreq 0x1015 0xd241

Received 14 bytes. HEX dump:
00000000  0e 00 00 00 02 00 15 10  41 d4 00 00 0a 00        |........A.....|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  41 d4 00 00              |....... A...|
```

#### Property Code 0xd400 - [still unknown]
```text
opreq 0x1015 0xd400

Received 16 bytes. HEX dump:
00000000  10 00 00 00 02 00 15 10  00 d6 00 00 00 00 00 00  |................|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  00 d6 00 00              |....... ....|
```

#### Property Code 0xd401 - [still unknown]
```text
opreq 0x1015 0xd401

Received 16 bytes. HEX dump:
00000000  10 00 00 00 02 00 15 10  01 d6 00 00 00 00 00 00  |................|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  01 d6 00 00              |....... ....|
```

#### Property Code 0xd500 - Some form of GEO coordinates? But 'K' and 'M' are
odd though...
```text
opreq 0x1015 0xd500

Received 109 bytes. HEX dump:
00000000  6d 00 00 00 02 00 15 10  00 d7 00 00 30 30 30 30  |m...........0000|
00000010  2e 30 30 30 30 30 30 2c  4e 30 30 30 30 30 2e 30  |.000000,N00000.0|
00000020  30 30 30 30 30 2c 45 30  30 30 30 30 2e 30 30 2c  |00000,E00000.00,|
00000030  4d 20 30 30 30 2e 30 2c  4b 30 30 30 30 3a 30 30  |M 000.0,K0000:00|
00000040  3a 30 30 30 30 3a 30 30  3a 30 30 2e 30 30 30 00  |:0000:00:00.000.|
00000050  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
00000060  00 00 00 00 00 00 00 00  00 00 00 00 00           |.............|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  00 d7 00 00              |....... ....|
```

#### Property Code 0xd52f - [still unknown]
```text
opreq 0x1015 0xd52f

Received 16 bytes. HEX dump:
00000000  10 00 00 00 02 00 15 10  2f d7 00 00 00 00 00 00  |......../.......|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  2f d7 00 00              |....... /...|
```

#### Property Code 0xdf00 - [still unknown]
```text
opreq 0x1015 0xdf00

Received 14 bytes. HEX dump:
00000000  0e 00 00 00 02 00 15 10  00 e1 00 00 00 00        |..............|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  00 e1 00 00              |....... ....|
```

#### Property Code 0xdf01 - The init sequence to use during connection
establishment.
```text
opreq 0x1015 0xdf01

Received 14 bytes. HEX dump:
00000000  0e 00 00 00 02 00 15 10  01 e1 00 00 00 00        |..............|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  01 e1 00 00              |....... ....|
```

#### Property Code 0xdf21 - [still unknown]
```text
opreq 0x1015 0xdf21

Received 16 bytes. HEX dump:
00000000  10 00 00 00 02 00 15 10  21 e1 00 00 03 00 00 00  |........!.......|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  21 e1 00 00              |....... !...|
```

#### Property Code 0xdf22 - [still unknown]
```text
opreq 0x1015 0xdf22

Received 16 bytes. HEX dump:
00000000  10 00 00 00 02 00 15 10  22 e1 00 00 03 00 00 00  |........".......|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  22 e1 00 00              |....... "...|
```

#### Property Code 0xdf23 - [still unknown]
```text
opreq 0x1015 0xdf23

Received 16 bytes. HEX dump:
00000000  10 00 00 00 02 00 15 10  23 e1 00 00 01 00 00 00  |........#.......|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  23 e1 00 00              |....... #...|
```

#### Property Code 0xdf24 - The minimal supported Application Version
```text
opreq 0x1015 0xdf24

Received 16 bytes. HEX dump:
00000000  10 00 00 00 02 00 15 10  24 e1 00 00 01 00 02 00  |........$.......|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  24 e1 00 00              |....... $...|
```

#### Property Code 0xdf25 - [still unknown]
```text
opreq 0x1015 0xdf25

Received 16 bytes. HEX dump:
00000000  10 00 00 00 02 00 15 10  25 e1 00 00 01 00 00 00  |........%.......|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  25 e1 00 00              |....... %...|
```

#### Property Code 0xdf31 - [still unknown]
```text
opreq 0x1015 0xdf31

Received 16 bytes. HEX dump:
00000000  10 00 00 00 02 00 15 10  31 e1 00 00 02 00 00 00  |........1.......|

Received 12 bytes. HEX dump:
00000000  0c 00 00 00 03 00 01 20  31 e1 00 00              |....... 1...|
```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```

#### Property Code 0x - 
```text

```
