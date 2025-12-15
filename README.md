# ⚡️ gonrg - Energy Smart Meter Reader

CLI tool, server, and Go library to connect to the optical D0 interfaces of
power meters and gather OBIS data (e.g. power consumption).

For this to work, you need a gadget to connect to the optical
D0 interface of your meter, which you preferably connect to your host via USB,
but any (tunneled) serial connection should also be fine.

Text-based D0 OBIS is currently supported (which most digital meters in Germany use),
**but no SML yet**, because I don't have a device at hand to test with.

The binary can be used as a standalone CLI tool, as a JSON/REST server serving
OBIS meter data, or as a client to connect with a remote gonrg server.

## Tested and compatible meters

- **eBZ DD3**, many variants
  - Power meter, consumption and production (two-direction).
  - Tested with and without extended data (detailed per-phase power).
- **Landis+Gyr Ultraheat T550**
  - District heating meter
  - Required params: `--baudrate 300`
  - Only poll daily or hourly because the device runs on battery.

gonrg is probably compatible with many other devices. If you can verify compatibility,
please open a pull request that extends the list.

## TODOs

- **CI/CD, autobuilds and binary releases on GitHub**
- Server OpenAPI specification
- Client code
- SML and HAN interface support (I don't have any devices available to test with)

## How to build

Building requires **go 1.25**.

Just hit `make` to obtain a statically linked binary in the project root.

## Set Me Up - CLI tool

```
john@doe:~/foo$ gonrg --help
⚡️ gonrg - a simple D0 OBIS energy meter CLI tool or server.

Usage:
  gonrg [flags]
  gonrg [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  server      run in server mode given a config

Flags:
  -b, --baudrate int    baud rate, 0 means choose best option
  -D, --debug           set debug log level
  -d, --device string   device to read from (default "/dev/ttyUSB0")
  -h, --help            help for gonrg
  -j, --json            output json instead of pretty table
  -S, --strict          strict mode for parsing - fail fast
  -v, --version         version for gonrg

Use "gonrg [command] --help" for more information about a command.
```

### Examples

Create a pretty tabular (but not machine-readable) overview over the OBIS
data of your meter: 

```
john@doe:~/foo$ gonrg --device /dev/ttyUSB1

⚡️gonrg version 0.1.0
Device ID: EBZ5DD12345ETA_104

Exact Key       Simple Key  Name         Value           Unit  
1-0:0.0.0*255   0.0.0                    1EBZ0102123456        
1-0:96.1.0*255  96.1.0                   1EBZ0102123456        
1-0:1.8.0*255   1.8.0       energy_cons  1924            kWh   
1-0:96.5.0*255  96.5.0                   001C0104              
0-0:96.8.0*255  96.8.0                   0086C9BB           
```

The same info can be obtained in JSON with the `--json` flag appended:

```
john@doe:~/foo$ gonrg --device /dev/ttyUSB1   
{
  "measurementTime": "2025-12-15T21:10:29.730101646+01:00",
  "deviceID": "EBZ5DD12345ETA_104",
  "list": [
    {
      "exactKey": "1-0:0.0.0*255",
      "simplifiedKey": "0.0.0",
      "name": "",
      "valueText": "1EBZ0102123456",
      "valueNum": 0,
      "valueScale": 0,
      "valueFloat": 0,
      "unit": ""
    },
    {
      "exactKey": "1-0:96.1.0*255",
      "simplifiedKey": "96.1.0",
      "name": "",
      "valueText": "1EBZ0102123456",
      "valueNum": 0,
      "valueScale": 0,
      "valueFloat": 0,
      "unit": ""
    },
    {
      "exactKey": "1-0:1.8.0*255",
      "simplifiedKey": "1.8.0",
      "name": "energy_cons",
      "valueText": "",
      "valueNum": 1924,
      "valueScale": 0,
      "valueFloat": 1924,
      "unit": "kWh"
    },
    {
      "exactKey": "1-0:96.5.0*255",
      "simplifiedKey": "96.5.0",
      "name": "",
      "valueText": "001C0104",
      "valueNum": 0,
      "valueScale": 0,
      "valueFloat": 0,
      "unit": ""
    },
    {
      "exactKey": "0-0:96.8.0*255",
      "simplifiedKey": "96.8.0",
      "name": "",
      "valueText": "0086C9BB",
      "valueNum": 0,
      "valueScale": 0,
      "valueFloat": 0,
      "unit": ""
    }
  ]
}
```

## Set Me Up - Meter Server

You can start a thread-safe server which listens to incoming JSON/REST connections
and serves OBIS data.

First, create a server configuration by using the [example_server.yaml](./example_server.yaml)
as a template. Then start the server with:

```
john@doe:~/foo$ gonrg server -C my_server_config.yaml
```

You can use curl to test fetching OBIS data:
```
john@doe:~/foo$ curl http://server:8080/meter/power/1.8.0
{
    "measurementTime": "2025-12-15T21:17:57.759573978+01:00",
    "deviceID": "EBZ5DD12345ETA_104",
    "list": {
        "exactKey": "1-0:1.8.0*255",
        "simplifiedKey": "1.8.0",
        "name": "energy_cons",
        "valueText": "",
        "valueNum": 234172889972,
        "valueScale": 8,
        "valueFloat": 2341.72889972,
        "unit": "kWh"
    }
}
```

## Miscellaneous

- If you encounter errors, try to append the `--debug` flag to see what is going on.
- Running `make gonrg-mock` renders a binary `gonrg-mock`, which does not connect to actual
  meters, but mocks some popular meters which can be used to evaluate gonrg without real
  hardware. If you specify a non-existing device name, it will give you a list with existing
  ones in the log.
