```text
info
```
Result:
```text
#TODO: add correct table output

capture delay: off || 0 - [0 0] - 0x0000 - 0x0
flash mode: disabled || 32777 - [9 128] - 0x0980 - 0x8009
white balance: automatic || 2 - [2 0] - 0x0200 - 0x2
exposure bias compensation: off || 0 - [0 0] - 0x0000 - 0x0
film simulation: Velvia || 2 - [2 0] - 0x0200 - 0x2
ISO: S 6400 || 6400 - [0 25 0 128] - 0x00190080 - 0x1900
rec mode:  || 1 - [1 0] - 0x0100 - 0x1
focus point:  || 1028 - [4 4 2 3] - 0x04040203 - 0x404
```

```text
info json
```
Result:
```json
[{"DevicePropertyCode":{"code":"0x5012","label":"capture delay"},"dataType":"uint16","readOnly":false,"FactoryDefaultValue":{"value":"0x0","label":"off"},"CurrentValue":{"value":"0x0","label":"off"},"formType":"enum","form":{"values":[{"value":"0x0","label":"off"},{"value":"0x2","label":"2 seconds"},{"value":"0x4","label":"10 seconds"}]}},{"DevicePropertyCode":{"code":"0x500c","label":"flash mode"},"dataType":"uint16","readOnly":false,"FactoryDefaultValue":{"value":"0x2","label":"off"},"CurrentValue":{"value":"0x8009","label":"disabled"},"formType":"enum","form":{"values":[{"value":"0x8009","label":"disabled"},{"value":"0x800a","label":"enabled"}]}},{"DevicePropertyCode":{"code":"0x5005","label":"white balance"},"dataType":"uint16","readOnly":false,"FactoryDefaultValue":{"value":"0x2","label":"automatic"},"CurrentValue":{"value":"0x2","label":"automatic"},"formType":"enum","form":{"values":[{"value":"0x2","label":"automatic"},{"value":"0x4","label":"daylight"},{"value":"0x8006","label":"shade"},{"value":"0x8001","label":"fluorescent 1"},{"value":"0x8002","label":"fluorescent 2"},{"value":"0x8003","label":"fluorescent 3"},{"value":"0x6","label":"tungsten"},{"value":"0x800a","label":"underwater"},{"value":"0x800b","label":"temprerature"},{"value":"0x800c","label":"custom"}]}},{"DevicePropertyCode":{"code":"0x5010","label":"exposure bias compensation"},"dataType":"int16","readOnly":false,"FactoryDefaultValue":{"value":"0x0","label":"off"},"CurrentValue":{"value":"0x0","label":"off"},"formType":"enum","form":{"values":[{"value":"0xf448","label":"-3"},{"value":"0xf595","label":"-2 2/3"},{"value":"0xf6e3","label":"-2 1/3"},{"value":"0xf830","label":"-2"},{"value":"0xf97d","label":"-1 2/3"},{"value":"0xfacb","label":"-1 1/3"},{"value":"0xfc18","label":"-1"},{"value":"0xfd65","label":"-2/3"},{"value":"0xfeb3","label":"-1/3"},{"value":"0x0","label":"off"},{"value":"0x14d","label":"+1/3"},{"value":"0x29b","label":"+2/3"},{"value":"0x3e8","label":"+1"},{"value":"0x535","label":"+1 1/3"},{"value":"0x683","label":"+1 2/3"},{"value":"0x7d0","label":"+2"},{"value":"0x91d","label":"+2 1/3"},{"value":"0xa6b","label":"+2 2/3"},{"value":"0xbb8","label":"+3"}]}},{"DevicePropertyCode":{"code":"0xd001","label":"film simulation"},"dataType":"uint16","readOnly":false,"FactoryDefaultValue":{"value":"0x1","label":"Provia"},"CurrentValue":{"value":"0x2","label":"Velvia"},"formType":"enum","form":{"values":[{"value":"0x1","label":"Provia"},{"value":"0x2","label":"Velvia"},{"value":"0x3","label":"Astia"},{"value":"0x4","label":"Monochrome"},{"value":"0x5","label":"Sepia"},{"value":"0x6","label":"Pro. Neg. Hi"},{"value":"0x7","label":"Pro Neg. Std."},{"value":"0x8","label":"Monochrome + Ye Filter"},{"value":"0x9","label":"Monochrome + R Filter"},{"value":"0xa","label":"Monochrome + G Filter"},{"value":"0xb","label":"Classic Chrome"}]}},{"DevicePropertyCode":{"code":"0xd02a","label":"ISO"},"dataType":"uint32","readOnly":false,"FactoryDefaultValue":{"value":"0xffffffff","label":"auto"},"CurrentValue":{"value":"0x80001900","label":"S 6400"},"formType":"enum","form":{"values":[{"value":"0x80000190","label":"S 400"},{"value":"0x80000320","label":"S 800"},{"value":"0x80000640","label":"S 1600"},{"value":"0x80000c80","label":"S 3200"},{"value":"0x80001900","label":"S 6400"},{"value":"0x40000064","label":"L 100"},{"value":"0xc8","label":"200"},{"value":"0xfa","label":"250"},{"value":"0x140","label":"320"},{"value":"0x190","label":"400"},{"value":"0x1f4","label":"500"},{"value":"0x280","label":"640"},{"value":"0x320","label":"800"},{"value":"0x3e8","label":"1000"},{"value":"0x4e2","label":"1250"},{"value":"0x640","label":"1600"},{"value":"0x7d0","label":"2000"},{"value":"0x9c4","label":"2500"},{"value":"0xc80","label":"3200"},{"value":"0xfa0","label":"4000"},{"value":"0x1388","label":"5000"},{"value":"0x1900","label":"6400"},{"value":"0x40003200","label":"H 12800"},{"value":"0x40006400","label":"H 25600"},{"value":"0x4000c800","label":"H 51200"}]}},{"DevicePropertyCode":{"code":"0xd019","label":"rec mode"},"dataType":"uint16","readOnly":false,"FactoryDefaultValue":{"value":"0x1","label":""},"CurrentValue":{"value":"0x1","label":""},"formType":"enum","form":{"values":[{"value":"0x0","label":""},{"value":"0x1","label":""}]}},{"DevicePropertyCode":{"code":"0xd17c","label":"focus point"},"dataType":"uint32","readOnly":false,"FactoryDefaultValue":{"value":"0x0","label":""},"CurrentValue":{"value":"0x3020404","label":""},"formType":"range","form":{"min":"0x0","max":"0x10090707","step":"0x1"}}]
```

```text
info json pretty
```
Result:
```json
[
    {
        "DevicePropertyCode": {
            "code": "0x5012",
            "label": "capture delay"
        },
        "dataType": "uint16",
        "readOnly": false,
        "FactoryDefaultValue": {
            "value": "0x0",
            "label": "off"
        },
        "CurrentValue": {
            "value": "0x0",
            "label": "off"
        },
        "formType": "enum",
        "form": {
            "values": [
                {
                    "value": "0x0",
                    "label": "off"
                },
                {
                    "value": "0x2",
                    "label": "2 seconds"
                },
                {
                    "value": "0x4",
                    "label": "10 seconds"
                }
            ]
        }
    },
    {
        "DevicePropertyCode": {
            "code": "0x500c",
            "label": "flash mode"
        },
        "dataType": "uint16",
        "readOnly": false,
        "FactoryDefaultValue": {
            "value": "0x2",
            "label": "off"
        },
        "CurrentValue": {
            "value": "0x8009",
            "label": "disabled"
        },
        "formType": "enum",
        "form": {
            "values": [
                {
                    "value": "0x8009",
                    "label": "disabled"
                },
                {
                    "value": "0x800a",
                    "label": "enabled"
                }
            ]
        }
    },
    {
        "DevicePropertyCode": {
            "code": "0x5005",
            "label": "white balance"
        },
        "dataType": "uint16",
        "readOnly": false,
        "FactoryDefaultValue": {
            "value": "0x2",
            "label": "automatic"
        },
        "CurrentValue": {
            "value": "0x2",
            "label": "automatic"
        },
        "formType": "enum",
        "form": {
            "values": [
                {
                    "value": "0x2",
                    "label": "automatic"
                },
                {
                    "value": "0x4",
                    "label": "daylight"
                },
                {
                    "value": "0x8006",
                    "label": "shade"
                },
                {
                    "value": "0x8001",
                    "label": "fluorescent 1"
                },
                {
                    "value": "0x8002",
                    "label": "fluorescent 2"
                },
                {
                    "value": "0x8003",
                    "label": "fluorescent 3"
                },
                {
                    "value": "0x6",
                    "label": "tungsten"
                },
                {
                    "value": "0x800a",
                    "label": "underwater"
                },
                {
                    "value": "0x800b",
                    "label": "temprerature"
                },
                {
                    "value": "0x800c",
                    "label": "custom"
                }
            ]
        }
    },
    {
        "DevicePropertyCode": {
            "code": "0x5010",
            "label": "exposure bias compensation"
        },
        "dataType": "int16",
        "readOnly": false,
        "FactoryDefaultValue": {
            "value": "0x0",
            "label": "off"
        },
        "CurrentValue": {
            "value": "0x0",
            "label": "off"
        },
        "formType": "enum",
        "form": {
            "values": [
                {
                    "value": "0xf448",
                    "label": "-3"
                },
                {
                    "value": "0xf595",
                    "label": "-2 2/3"
                },
                {
                    "value": "0xf6e3",
                    "label": "-2 1/3"
                },
                {
                    "value": "0xf830",
                    "label": "-2"
                },
                {
                    "value": "0xf97d",
                    "label": "-1 2/3"
                },
                {
                    "value": "0xfacb",
                    "label": "-1 1/3"
                },
                {
                    "value": "0xfc18",
                    "label": "-1"
                },
                {
                    "value": "0xfd65",
                    "label": "-2/3"
                },
                {
                    "value": "0xfeb3",
                    "label": "-1/3"
                },
                {
                    "value": "0x0",
                    "label": "off"
                },
                {
                    "value": "0x14d",
                    "label": "+1/3"
                },
                {
                    "value": "0x29b",
                    "label": "+2/3"
                },
                {
                    "value": "0x3e8",
                    "label": "+1"
                },
                {
                    "value": "0x535",
                    "label": "+1 1/3"
                },
                {
                    "value": "0x683",
                    "label": "+1 2/3"
                },
                {
                    "value": "0x7d0",
                    "label": "+2"
                },
                {
                    "value": "0x91d",
                    "label": "+2 1/3"
                },
                {
                    "value": "0xa6b",
                    "label": "+2 2/3"
                },
                {
                    "value": "0xbb8",
                    "label": "+3"
                }
            ]
        }
    },
    {
        "DevicePropertyCode": {
            "code": "0xd001",
            "label": "film simulation"
        },
        "dataType": "uint16",
        "readOnly": false,
        "FactoryDefaultValue": {
            "value": "0x1",
            "label": "Provia"
        },
        "CurrentValue": {
            "value": "0x2",
            "label": "Velvia"
        },
        "formType": "enum",
        "form": {
            "values": [
                {
                    "value": "0x1",
                    "label": "Provia"
                },
                {
                    "value": "0x2",
                    "label": "Velvia"
                },
                {
                    "value": "0x3",
                    "label": "Astia"
                },
                {
                    "value": "0x4",
                    "label": "Monochrome"
                },
                {
                    "value": "0x5",
                    "label": "Sepia"
                },
                {
                    "value": "0x6",
                    "label": "Pro. Neg. Hi"
                },
                {
                    "value": "0x7",
                    "label": "Pro Neg. Std."
                },
                {
                    "value": "0x8",
                    "label": "Monochrome + Ye Filter"
                },
                {
                    "value": "0x9",
                    "label": "Monochrome + R Filter"
                },
                {
                    "value": "0xa",
                    "label": "Monochrome + G Filter"
                },
                {
                    "value": "0xb",
                    "label": "Classic Chrome"
                }
            ]
        }
    },
    {
        "DevicePropertyCode": {
            "code": "0xd02a",
            "label": "ISO"
        },
        "dataType": "uint32",
        "readOnly": false,
        "FactoryDefaultValue": {
            "value": "0xffffffff",
            "label": "auto"
        },
        "CurrentValue": {
            "value": "0x80001900",
            "label": "S 6400"
        },
        "formType": "enum",
        "form": {
            "values": [
                {
                    "value": "0x80000190",
                    "label": "S 400"
                },
                {
                    "value": "0x80000320",
                    "label": "S 800"
                },
                {
                    "value": "0x80000640",
                    "label": "S 1600"
                },
                {
                    "value": "0x80000c80",
                    "label": "S 3200"
                },
                {
                    "value": "0x80001900",
                    "label": "S 6400"
                },
                {
                    "value": "0x40000064",
                    "label": "L 100"
                },
                {
                    "value": "0xc8",
                    "label": "200"
                },
                {
                    "value": "0xfa",
                    "label": "250"
                },
                {
                    "value": "0x140",
                    "label": "320"
                },
                {
                    "value": "0x190",
                    "label": "400"
                },
                {
                    "value": "0x1f4",
                    "label": "500"
                },
                {
                    "value": "0x280",
                    "label": "640"
                },
                {
                    "value": "0x320",
                    "label": "800"
                },
                {
                    "value": "0x3e8",
                    "label": "1000"
                },
                {
                    "value": "0x4e2",
                    "label": "1250"
                },
                {
                    "value": "0x640",
                    "label": "1600"
                },
                {
                    "value": "0x7d0",
                    "label": "2000"
                },
                {
                    "value": "0x9c4",
                    "label": "2500"
                },
                {
                    "value": "0xc80",
                    "label": "3200"
                },
                {
                    "value": "0xfa0",
                    "label": "4000"
                },
                {
                    "value": "0x1388",
                    "label": "5000"
                },
                {
                    "value": "0x1900",
                    "label": "6400"
                },
                {
                    "value": "0x40003200",
                    "label": "H 12800"
                },
                {
                    "value": "0x40006400",
                    "label": "H 25600"
                },
                {
                    "value": "0x4000c800",
                    "label": "H 51200"
                }
            ]
        }
    },
    {
        "DevicePropertyCode": {
            "code": "0xd019",
            "label": "rec mode"
        },
        "dataType": "uint16",
        "readOnly": false,
        "FactoryDefaultValue": {
            "value": "0x1",
            "label": ""
        },
        "CurrentValue": {
            "value": "0x1",
            "label": ""
        },
        "formType": "enum",
        "form": {
            "values": [
                {
                    "value": "0x0",
                    "label": ""
                },
                {
                    "value": "0x1",
                    "label": ""
                }
            ]
        }
    },
    {
        "DevicePropertyCode": {
            "code": "0xd17c",
            "label": "focus point"
        },
        "dataType": "uint32",
        "readOnly": false,
        "FactoryDefaultValue": {
            "value": "0x0",
            "label": ""
        },
        "CurrentValue": {
            "value": "0x3020404",
            "label": ""
        },
        "formType": "range",
        "form": {
            "min": "0x0",
            "max": "0x10090707",
            "step": "0x1"
        }
    }
]
```