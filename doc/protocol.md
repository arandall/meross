# Meross Device/Appliance Protocol

This document is to explain the protocol used between Meross IoT appliances, the mobile app and the Meross cloud service.

## Abilities

An appliance returns what abilities it is able to support. I have so far only had an `mss310` socket and assume there are
additional abilities.

Note: I've only used a very old version `1.1.13` from 2019 and `6.1.10` which is equivalent to `6.1.8` from 2021.

| Ability                                                                   | First Seen | Description                                                                                       |
|---------------------------------------------------------------------------|------------|---------------------------------------------------------------------------------------------------|
| [Appliance.Config.Key](#applianceconfigkey)                               | 1.1.13     | Used to configure MQTT servers                                                                    |
| [Appliance.Config.Trace](#applianceconfigtrace)                           | 1.1.13     | Returns WiFi and system details during setup                                                      |
| [Appliance.Config.Wifi](#applianceconfigwifi)                             | 1.1.13     | Configures WiFi network to connect to during setup                                                |
| [Appliance.Config.WifiList](#applianceconfigwifiList)                     | 1.1.13     | Lists aviable Wifi networks                                                                       |
| [Appliance.Control.Bind](#appliancecontrolbind)                           | 1.1.13     | Sent by device after setup                                                                        |
| [Appliance.Control.ConsumptionConfig](#appliancecontrolconsumptionconfig) | 6.1.8      | Consumption ratio, unsure of purpose                                                              |
| [Appliance.Control.ConsumptionX](#appliancecontrolconsumptionx)           | 1.1.13     | Shows daily power consumption from the last 30 days                                               |
| [Appliance.Control.Electricity](#appliancecontrolelectricity)             | 1.1.13     | Returns present electricity usage                                                                 |
| [Appliance.Control.Multiple](#appliancecontrolmultiple)                   | 6.1.8      | Send multiple Appliance.Control requests in one                                                   |
| [Appliance.Control.Timer](#appliancecontroltimer)                         | 1.1.13     | GET/SET Timer values                                                                              |
| [Appliance.Control.TimerX](#appliancecontroltimerx)                       | 6.1.8      | GET/SET Timer values on newer firmware                                                            |
| [Appliance.Control.Toggle](#appliancecontroltoggle)                       | 1.1.13     | Toggle switch on/off                                                                              |
| [Appliance.Control.ToggleX](#appliancecontroltogglex)                     | 6.1.8      | Toggle switch on/off on newer firmware                                                            |
| [Appliance.Control.Trigger](#appliancecontroltrigger)                     | 1.1.13     | GET/SET Trigger rules                                                                             |
| [Appliance.Control.TriggerX](#appliancecontroltriggerx)                   | 6.1.8      | Trigger rules on newer firmware                                                                   |
| [Appliance.Control.Unbind](#appliancecontrolunbind)                       | 1.1.13     | Deactivate device                                                                                 |
| [Appliance.Control.Upgrade](#appliancecontrolupgrade)                     | 1.1.13     | Upgrade firmware from URL                                                                         |
| [Appliance.Digest.TimerX](#appliancedigesttimerx)                         | 6.1.8      | List configured timers                                                                            |
| [Appliance.Digest.TriggerX](#appliancedigesttriggerx)                     | 6.1.8      | List configured triggers                                                                          |
| [Appliance.System.Ability](#appliancesystemability)                       | 1.1.13     | List all abilities appliance supports                                                             |
| [Appliance.System.All](#appliancesystemall)                               | 1.1.13     | List all system attributes (includes firmware, hardware, timer, trigger, toogle and online status |
| [Appliance.System.Clock](#appliancesystemclock)                           | 1.1.13     | Used to request current time                                                                      |
| [Appliance.System.DNDMode](#appliancesystemdndmode)                       | 1.1.13     | Used to set/unset DNDMode                                                                         |
| [Appliance.System.Debug](#appliancesystemdebug)                           | 1.1.13     | Get debug information about appliance OS                                                          |
| [Appliance.System.Firmware](#appliancesystemfirmware)                     | 1.1.13     | Get firmware information                                                                          |
| [Appliance.System.Hardware](#appliancesystemhardware)                     | 1.1.13     | Get hardware information                                                                          |
| [Appliance.System.Online](#appliancesystemonline)                         | 1.1.13     | Get online status                                                                                 |
| [Appliance.System.Position](#appliancesystemposition)                     | 1.1.13     | Get/Set Position (lat/lng) of appliance                                                           |
| [Appliance.System.Report](#appliancesystemreport)                         | 1.1.13     | PUSH data back to server (only seen time reporting)                                               |
| [Appliance.System.Runtime](#appliancesystemruntime)                       | 1.1.13     | Get runtime, only observed WiFi "signal" strength                                                 |
| [Appliance.System.Time](#appliancesystemtime)                             | 1.1.13     | Get/Set timezone and daylight savings rules                                                       |

## Packets

All messages to/from a appliance use the following packet format.

```json
{
  "header": {
    "from": "http://10.10.10.1/config",
    "messageId": "0123456789abcdef01234567890abcde",
    "method": "GET",
    "namespace": "Appliance.System.All",
    "payloadVersion": 1,
    "sign": "0123456789abcdef01234567890abcde",
    "timestamp": 1557596606,
    "timestampMs": 100,
  },
  "payload": {}
}
```

## Headers

Each header contains the following keys

| Field          | Description                                                                                                    |
|----------------|----------------------------------------------------------------------------------------------------------------|
| from           | HTTP url of appliance config endpoint or MQTT topic of the endpoint that generated packet.                     |
| messageId      | Arbitrary ID (32 characters of HEX) - responses will use the same `messageID`                                  |
| method         | used to describe request/response action `GET`/`GETACK`, `SET`/`SETACK`, `PUSH`, `ERROR`                       |
| namespace      | [ability](#abilities) being used                                                                               |
| payloadVersion | I assume this is for future payload versions, only ever seen `1` being used                                    |
| sign           | Signing value equal to md5(`messageId` + `key` + `timestamp`) where `key` is defined in `Appliance.Config.Key` |
| timestamp      | Time in seconds past Epoch                                                                                     |
| timestampMs    | Time millisecond component (appeared in newer firmware versions `6.1.8+` probably sooner)                      |

**Note:** When a appliance is waiting to be configured it does not validate the sign value.

## Errors

Errors are returned by any `GET` or `SET` request that has a problem. Typically, these are sign errors.

Method: `ERROR`

```json
{
  "error": {
    "code": 5001,
    "detail": "sign error"
  }
}
```

## Payload

The payload of a packet is a JSON structure dependent on the ability being used. Most `GET` methods use a payload of
`{}` when requesting data however in some cases additional context is sent with the request.

### Config

The `namespaces` are used to configure an appliance. Typically, these do not require a signing key as the appliance is yet to be
configured.

#### Appliance.Config.Key

Used to set the MQTT servers and signing key of an appliance.

Method: `SET` / (`PUSH` with `6.1.8+`)

The Meross cloud service appears to be able to now `PUSH` key/config to a device. TODO: Confirm this can be used to
change a key.

```json
{
  "key": {
    "gateway": {
      "host": "iot.meross.com",
      "port": 2001,
      "secondHost": "smart.meross.com",
      "secondPort": 2001
    },
    "key": "0123456789abcdef01234567890abcde",
    "userId": "1234"
  }
}
```

| Field                   | Description                                                                   |
|-------------------------|-------------------------------------------------------------------------------|
| .key.gateway.host       | Primary MQTT host                                                             |
| .key.gateway.port       | Primary MQTT port                                                             |
| .key.gateway.secondHost | Secondary MQTT host                                                           |
| .key.gateway.secondPort | Secondary MQTT port                                                           |
| .key.key                | pre-shared key used for signing requests (usually assigned by iot.meross.com) |
| .key.userId             | userId (usually assigned by iot.meross.com)                                   |

Method: `SETACK`

```json
{}
```

#### Appliance.Config.Trace

I have not found a use for this, but it is used by the app when configuring. Once an appliance is configured the
response contains empty values.

Method: `GET`

Empty object `{}` appears to work too.

```json
{"trace":{}}
```

Method: `GETACK`

```json
{
  "trace": {
    "ssid": "My SSID",
    "bssid": "DE:AD:00:00:BE:EF",
    "rssi": "100",
    "code": 0,
    "info": "no error"
  }
}
```

#### Appliance.Config.Wifi

Used to set Wi-Fi SSID/Password.

**Note:** Once appliance responds it will restart and connect to the network provided. If an error occurs the appliance
will reset and wait for configuration.

Method: `SET`

```json
{
  "wifi": {
    "bssid": "de:ad:00:00:be:ef",
    "channel": 3,
    "cipher": 3,
    "encryption": 6,
    "password": "cGFzc3dvcmQK",
    "ssid": "c3NpZAo="
  }
}
```

| Field            | Description                                             |
|------------------|---------------------------------------------------------|
| .wifi.bssid      | BSSID                                                   |
| .wifi.channel    | Channel                                                 |
| .wifi.cipher     | cipher used by Wifi network                             |
| .wifi.encryption | encryption used by network                              |
| .wifi.password   | base64 encoded string of password to connect to network |
| .wifi.bssid      | base64 encoded string of SSID                           |

Method: `SETACK`

```json
{}
```

#### Appliance.Config.WifiList

Used to list all Wi-Fi networks the appliance can see.

At some point (observed from `6.1.8+`) the mac address changed from using `-` separators to `:`.

Method: `GET`

```json
{}
```

Method: `GETACK`

```json
{
  "wifiList": [
    {
      "ssid": "dGVzdCBTU0lE",
      "bssid": "2c:6e:a4:55:c6:47",
      "signal": 100,
      "channel": 3,
      "encryption": 6,
      "cipher": 3
    },
    {
      "ssid": "dGVzdCBTU0lEMg==",
      "bssid": "4d:23:0e:3c:a2:22",
      "signal": 20,
      "channel": 3,
      "encryption": 6,
      "cipher": 3
    }
  ]
}
```

### System

#### Appliance.System.Ability

Gets supported abilities of an appliance. Abilities may report additional limits.

Method: `GET`

```json
{}
```

Method: `GETACK`

```json
{
  "payloadVersion": 1,
  "ability": {
    "Appliance.Config.Key": {},
    "Appliance.Config.WifiList": {},
    "Appliance.Config.Wifi": {},
    "Appliance.Config.Trace": {},
    "Appliance.System.All": {},
    "Appliance.System.Hardware": {},
    "Appliance.System.Firmware": {},
    "Appliance.System.Debug": {},
    "Appliance.System.Online": {},
    "Appliance.System.Time": {},
    "Appliance.System.Clock": {},
    "Appliance.System.Ability": {},
    "Appliance.System.Runtime": {},
    "Appliance.System.Report": {},
    "Appliance.System.Position": {},
    "Appliance.System.DNDMode": {},
    "Appliance.Control.Multiple": {
      "maxCmdNum": 5
    },
    "Appliance.Control.ToggleX": {},
    "Appliance.Control.TimerX": {
      "sunOffsetSupport": 1
    },
    "Appliance.Control.TriggerX": {},
    "Appliance.Control.Bind": {},
    "Appliance.Control.Unbind": {},
    "Appliance.Control.Upgrade": {},
    "Appliance.Control.ConsumptionX": {},
    "Appliance.Control.Electricity": {},
    "Appliance.Control.ConsumptionConfig": {},
    "Appliance.Digest.TriggerX": {},
    "Appliance.Digest.TimerX": {}
  }
}
```

#### Appliance.System.All

Gets all system values.

Method: `GET`

```json
{}
```

Method: `GETACK`

```json
{
  "all": {
    "system": {
      "hardware": {
        "type": "mss310",
        "subType": "us",
        "version": "1.0.0",
        "chipType": "MT7688",
        "uuid": "11111111111111111111111111111111",
        "macAddress": "34:29:8f:ff:ff:ff"
      },
      "firmware": {
        "version": "1.1.13",
        "compileTime": "2018-06-01 09:56:24",
        "wifiMac": "ff:ff:ff:ff:ff:ff",
        "innerIp": "10.0.0.21",
        "server": "10.0.0.1",
        "port": 8883,
        "secondServer": "",
        "secondPort": 0,
        "userId": 12345
      },
      "time": {
        "timestamp": 1560144444,
        "timezone": "",
        "timeRule": []
      },
      "online": {
        "status": 1
      }
    },
    "control": {
      "toggle": {
        "onoff": 1,
        "lmTime": 1560124804
      },
      "trigger": [],
      "timer": []
    }
  }
}
```

#### Appliance.System.Clock

When an appliance first connects to MQTT it publishes a PUSH of the timestamp of the appliance. Every appliance push
expects a response back with the current timestamp. I assume this is to mimic NTP and setting of the appliance time.

When this message is returned to from the server to the client the `messageId`, `timestamp` and `timestampMs` headers
are identical resulting in the same key, and do not require recomputing. I assume this is to provide a quick response
and minimise time drift.

Method: `PUSH`

```json
{
  "clock": {
    "timestamp": 1560144444
  }
}
```

| Field            | Description                     |
|------------------|---------------------------------|
| .clock.timestamp | Timestamp in seconds past epoch |

#### Appliance.System.DNDMode

DNDMode turns off the status LED of the appliance.

Method: `SET`

```json
{
  "DNDMode": {
    "mode": 0
  }
}
```

| Field         | Description                                                   |
|---------------|---------------------------------------------------------------|
| .DNDMode.mode | 0 = off (status LED functional), 1 = on (status LED disabled) |

Method: `SETACK`

```json
{}
```

#### Appliance.System.Debug

Method: `GET`

```json
{}
```

Method: `GETACK`

```json
{
  "debug": {
    "system": {
      "version": "1.1.13",
      "compileTime": "2018-12-03 11:15:36",
      "sysUpTime": "2h39m32s",
      "memory": "Total 3319k, Free 1691k, Largest free block 1687k",
      "suncalc": "7:34;17:7",
      "localTime": "Sun Jun 16 17:21:45 2019",
      "localTimeOffset": 36000
    },
    "network": {
      "linkStatus": 0,
      "signal": 100,
      "wifiDisconnectCount": 0,
      "ssid": "SSID",
      "gatewayMac": "ff:ff:ff:ff:ff:ff",
      "innerIp": "192.168.0.1"
    },
    "cloud": {
      "activeServer": "iot.meross.com",
      "mainServer": "iot.meross.com",
      "mainPort": 2001,
      "secondServer": "smart.meross.com",
      "secondPort": 2001,
      "userId": 1234,
      "sysConnectTime": "Sun Jun 16 04:42:30 2019",
      "sysOnlineTime": "2h39m15s",
      "sysDisconnectCount": 0,
      "pingTrace": []
    }
  }
}
```

#### Appliance.System.Firmware

Method: `GET`

```json
{}
```

Method: `GETACK`

```json
{
  "firmware": {
    "version": "1.1.13",
    "compileTime": "2018-06-01 09:56:24",
    "wifiMac": "ff:ff:ff:ff:ff:ff",
    "innerIp": "10.0.0.21",
    "server": "10.0.0.1",
    "port": 8883,
    "secondServer": "",
    "secondPort": 0,
    "userId": 12345
  }
}
```

#### Appliance.System.Hardware

Method: `GET`

```json
{}
```

Method: `GETACK`

```json
{
  "hardware": {
    "type": "mss310",
    "subType": "us",
    "version": "1.0.0",
    "chipType": "MT7688",
    "uuid": "11111111111111111111111111111111",
    "macAddress": "34:29:8f:ff:ff:ff"
  }
}
```

#### Appliance.System.Online

Method: `GET`

```json
{}
```

Method: `GETACK`

```json
{
  "online": {
    "status": 1
  }
}
```

#### Appliance.System.Position

Used to set a lat/long on a device. I guess this is to support mapping.

Method: `GET`

```json
{}
```

Method: `GETACK`/`SET`

```json
{
  "position": {
    "longitude": 116363625,
    "latitude": 39913818
  }
}
```

| Field               | Description                     |
|---------------------|---------------------------------|
| .position.longitude | Longitude multiplied by 1000000 |
| .position.latitude  | Latitude multiplied by 1000000  |

Method: `SETACK`

```json
{}
```

#### Appliance.System.Report

Used to report data back to Meross. I haven't seen this used by an Appliance except for this `timestamp` example.

Method: `PUSH`

```json
{
  "report": [
    {
      "type": "1",
      "value": "0",
      "timestamp": 1560658220
    }
  ]
}
```

#### Appliance.System.Runtime

I think this is for Wifi signal strength.

Method: `GET`

```json
{}
```

Method: `GETACK`

```json
{
  "report": [
    {
      "type": "1",
      "value": "0",
      "timestamp": 1560658220
    }
  ]
}
```

#### Appliance.System.Time

Method: `GET`

```json
{}
```

Method: `GETACK` / `SET`

The `timeRule` key is essentially the output for the next 10 years. eg. `zdump -i -c 2019,2029 Australia/Sydney`

```yaml
{
  "time": {
    "timestamp": 1560670665,
    "timezone": "Australia/Sydney",
    "timeRule": [
      [
        1554566400,
        36000,
        0
      ],
      [
        1570291200,
        39600,
        1
      ],
      #...
    ]
  }
}
```

| Field             | Description                               |
|-------------------|-------------------------------------------|
| .time.timestamp   | Current time of appliance                 |
| .time.timezone    | Timezone of appliance                     |
| .time.timeRule    | List of daylight savings rules            |
| .time.timeRule[0] | timestamp rule takes effect               |
| .time.timeRule[1] | UTC offset in seconds                     |
| .time.timeRule[2] | Daylight savings active (0 = no, 1 = yes) |

Method: `SETACK`

```json
{}
```

### Control

#### Appliance.Control.Bind

Sent once over MQTT after successful configuration

Method: `PUSH`
```yaml
{
  "bind": {
    "bindTime": 1630498742,
    "time": {
      "timestamp": 1630498742,
      "timezone": "Australia/Sydney",
      "timeRule": [
        #timerules removed
      ]
    },
    "hardware": {
      "type": "mss310",
      "subType": "un",
      "version": "6.0.0",
      "chipType": "rtl8710cf",
      "uuid": "11111111111111111111111111111111",
      "macAddress": "48:e1:e9:ff:ff:ff"
    },
    "firmware": {
      "version": "6.1.8",
      "compileTime": "2021/04/07-16:08:36",
      "wifiMac": "ff:ff:ff:ff:ff:ff",
      "innerIp": "10.0.0.21",
      "server": "mqtt-ap-2.meross.com",
      "port": 443,
      "userId": 1234
    }
  }
}
```

#### Appliance.Control.ConsumptionConfig

On startup the device sends and the server returns this with `PUSH` not sure what significance this has. A `SET` results in an HTTP error.

Method: `GET` | `PUSH`

```json
{
  "config": {
    "voltageRatio": 188,
    "electricityRatio": 102,
    "maxElectricityCurrent":11000
  }
}
```

| Field                         | Description                      |
|-------------------------------|----------------------------------|
| .config.voltageRatio          | voltageRatio? (In AU set to 188) |
| .config.electricityRatio      | ? (In AU set to 102)             |
| .config.maxElectricityCurrent | current in mAh? only on `PUSH`   |

#### Appliance.Control.ConsumptionX

Method: `GET`

```json
{}
```

Method: `GETACK`

```yaml
{
  "consumptionx": [
    {
      "date": "2019-06-05",
      "time": 1559709311,
      "value": 0
    },
    {
      "date": "2019-06-06",
      "time": 1559795711,
      "value": 0
    },
    #...
  ]
}
```

| Field               | Description               |
|---------------------|---------------------------|
| .consumptionx.date  | date in Y-m-d format      |
| .consumptionx.time  | start timestamp of period |
| .consumptionx.value | Usage value in watt hours |

#### Appliance.Control.Electricity

Method: `GET`

```json
{
  "electricity": {
    "channel": 0
  }
}
```

| Field                | Description                                                    |
|----------------------|----------------------------------------------------------------|
| .electricity.channel | I assume this is to support an appliance with multiple sockets |

Method: `GETACK`

```json
{
  "electricity": {
    "channel": 0,
    "current": 36,
    "voltage": 2481,
    "power": 2782
  }
}
```

| Field                | Description                                        |
|----------------------|----------------------------------------------------|
| .electricity.channel |                                                    |
| .electricity.current | Current being consumed in milliamps (mA)           |
| .electricity.voltage | Current voltage in deci-volts (dV) (/10 for Volts) |
| .electricity.power   | Current power usage in milliwatts (mW)             |

#### Appliance.Control.Multiple

Send multiple control packets in one request.

Method is always `SET` even if requests within are `GET`s.

Method: `SET`

```json
{
  "multiple": [
      {
          "header": {
              "method": "SET",
              "namespace": "Appliance.Control.TimerX"
          },
          "payload": {
              "timerx": {
                  "alias": "test",
                  "channel": 0,
                  "count": 0,
                  "createTime": 1630825097,
                  "duration": 0,
                  "enable": 1,
                  "extend": {
                      "toggle": {
                          "channel": 0,
                          "onoff": 0
                      }
                  },
                  "id": "abc1234",
                  "sunOffset": 0,
                  "sunriseTime": 0,
                  "sunsetTime": 0,
                  "time": 963,
                  "type": 1,
                  "week": 192
              }
          }
      }
  ]
}
```

Method: `SETACK`

```json
{
  "multiple": [
      {
          "header": {
              "method": "SETACK",
              "namespace": "Appliance.Control.TimerX"
          },
          "payload": {}
      }
  ]
}
```

#### Appliance.Control.Timer

On or Off timer.

Method: `GET`

```json
{}
```

Method: `GETACK` / `SET`

```json
{
  "timer": [
    {
      "id": "abcdefghijklm123",
      "type": 1,
      "enable": 1,
      "alias": "on 20:52",
      "time": 1252,
      "week": 129,
      "duration": 0,
      "createTime": 1560513180,
      "extend": {
        "toggle": {
          "onoff": 1,
          "lmTime": 0
        }
      }
    }
  ]
}
```

| Field               | Description                                                                      |
|---------------------|----------------------------------------------------------------------------------|
| .timer[].id         | unique identifier of timer                                                       |
| .timer[].type       | 1 = weekly, 2 = once                                                             |
| .timer[].enable     | Enabled (0 = off, 1 = on)                                                        |
| .timer[].alias      | User string to identify timer                                                    |
| .timer[].time       | Minute of day to fire                                                            |
| .timer[].week       | 8 bit bitset with MSB always on LSB -> MSB Sun, Mon, ... (eg. Monday 0b10000010) |
| .timer[].duration   | set to 0 (not sure on usage)                                                     |
| .timer[].createTime | Timestamp timer was created                                                      |
| .timer[].extend     | [Appliance.Control.Toggle](#Appliance.Control.Toggle) object without channel     |

Method: `SETACK`

```json
{}
```

#### Appliance.Control.TimerX

On or Off timer.

To get a list of timers on the device use `Appliance.Digest.TimerX`

Method: `GET`

```json
{
  "timerx": {
    "id":"abc1234"
  }
}
```

| Field        | Description        |
|--------------|--------------------|
| .timerx[].id | id of timer to get |

Method: `SET` | `GETACK`

```json
{
  "digest": {
    "channel": 0,
    "id": "n1jxtruoknvotm8a",
    "count": 3
  },
  "timerx": {
    "id": "n1jxtruoknvotm8a",
    "alias": "test2",
    "type": 1,
    "enable": 1,
    "channel": 0,
    "createTime": 1630825097,
    "week": 192,
    "time": 963,
    "sunOffset": 0,
    "duration": 0,
    "extend": {
      "toggle": {
        "onoff": 0,
        "lmTime": 0
      }
    }
  }
}
```

| Field                | Description                                                                         |
|----------------------|-------------------------------------------------------------------------------------|
| .digest (`GET` only) | See Appliance.Digest.TriggerX                                                       |
| .timerx.id           | unique identifier of timer                                                          |
| .timerx.alias        | User string to identify timer                                                       |
| .timerx.type         | 1 = weekly, 2 = once                                                                |
| .timerx.enable       | Enabled (0 = off, 1 = on)                                                           |
| .timerx.channel      | I assume this is to support an appliance with multiple sockets (0 on single socket) |
| .timerx.createTime   | Timestamp timer was created                                                         |
| .timerx.week         | 8 bit bitset with MSB always on LSB -> MSB Sun, Mon, ... (eg. Monday 0b10000010)    |
| .timerx.time         | Minute of day to fire                                                               |
| .timerx.sunOffset    | set to 0 (not sure on usage)                                                        |
| .timerx.duration     | set to 0 (not sure on usage, possible time to stay on)                              |
| .timerx.extend       | [Appliance.Control.Toggle](#appliancecontroltoggle) object without channel          |

Method: `SETACK`

```json
{}
```

#### Appliance.Control.Toggle

Method: `SET`

```json
{
  "channel": 0,
  "toggle": {
    "onoff": 1
  }
}
```

| Field         | Description                                                    |
|---------------|----------------------------------------------------------------|
| .channel      | I assume this is to support an appliance with multiple sockets |
| .toggle.onoff | 0 = off, 1 = on, any other value locks current state           |

Method: `SETACK`

```json
{}
```

When the power state is toggled the following is sent from the appliance.

Method: `PUSH`

```json
{
  "toggle": {
    "onoff": 0,
    "lmTime":1559375884
  }
}
```

#### Appliance.Control.ToggleX

Method: `SET`

```json
{
  "togglex": {
    "channel": 0,
    "onoff": 1
  }
}
```

| Field            | Description                                                    |
|------------------|----------------------------------------------------------------|
| .togglex.channel | I assume this is to support an appliance with multiple sockets |
| .togglex.onoff   | 0 = off, 1 = on, any other value locks current state           |

Method: `SETACK`

```json
{}
```

When the power state is toggled the following is sent from the appliance.

Method: `PUSH`

```json
{
  "togglex": [
    {
      "channel": 0,
      "onoff": 1,
      "lmTime": 1559375884
    }
  ]
}
```


#### Appliance.Control.Trigger

Trigger a toggle on toggle.
> Eg. After appliance turned on wait 1hr then turn it off.

Method: `GET`

```json
{}
```

Method: `GETACK` / `SET`

```json
{
  "trigger": [
    {
      "id": "abcdefghijklm123",
      "type": 0,
      "enable": 1,
      "alias": "test auto off",
      "createTime": 1560513139,
      "rule": {
        "_if_": {
          "toggle": {
            "onoff": 1,
            "lmTime": 0
          }
        },
        "_then_": {
          "delay": {
            "week": 129,
            "duration": 69300
          }
        },
        "_do_": {
          "toggle": {
            "onoff": 0,
            "lmTime": 0
          }
        }
      }
    }
  ]
}
```

| Field                                 | Description                                                                         |
|---------------------------------------|-------------------------------------------------------------------------------------|
| .trigger[].id                         | unique identifier of timer                                                          |
| .trigger[].type                       | 1 = weekly, 2 = once                                                                |
| .trigger[].enable                     | Enabled (0 = off, 1 = on)                                                           |
| .trigger[].alias                      | User string to identify timer                                                       |
| .trigger[].createTime                 | Timestamp timer was created                                                         |
| .trigger[].rule._if_                  | [Appliance.Control.Toggle](#Appliance.Control.Toggle) object without channel        |
| .trigger[].rule._then_.delay.week     | 8 bit bitset with MSB always on LSB -> MSB Sun, Mon, ... (eg. Mon & Sun 0b10000011) |
| .trigger[].rule._then_.delay.duration | Seconds to wait                                                                     |
| .trigger[].rule._do_                  | [Appliance.Control.Toggle](#Appliance.Control.Toggle) object without channel        |

Method: `SETACK`

```json
{}
```

#### Appliance.Control.TriggerX

Use Appliance.Digest.TriggerX to get list of triggers

Trigger a toggle on toggle.
> Eg. After appliance turned on wait 1hr then turn it off.

Method: `GET`

```json
{
  "triggerx": {
    "id": "abc1234"
  }
}
```

| Field        | Description                  |
|--------------|------------------------------|
| .triggerx.id | unique identifier of trigger |

Method: `GETACK` / `SET`

```json
{
  "digest": {
    "channel": 0,
    "id": "abc1234",
    "count": 1
  },
  "triggerx": {
    "id": "abc1234",
    "type": 0,
    "enable": 1,
    "channel": 0,
    "alias": "tesy2",
    "createTime": 1630827106,
    "rule": {
      "week": 193,
      "duration": 900
    }
  }
}
```

| Field                  | Description                                                                         |
|------------------------|-------------------------------------------------------------------------------------|
| .digest (`GET` only)   | See Appliance.Digest.TriggerX                                                       |
| .trigger.id            | unique identifier of trigger                                                        |
| .trigger.type          | 1 = weekly, 2 = once                                                                |
| .trigger.enable        | Enabled (0 = off, 1 = on)                                                           |
| .trigger.channel       | I assume this is to support an appliance with multiple sockets                      |
| .trigger.alias         | User string to identify trigger                                                     |
| .trigger.createTime    | Timestamp timer was created                                                         |
| .trigger.rule.week     | 8 bit bitset with MSB always on LSB -> MSB Sun, Mon, ... (eg. Mon & Sun 0b10000011) |
| .trigger.rule.duration | Seconds to wait                                                                     |

Method: `SETACK`

```json
{}
```

#### Appliance.Control.Unbind

Sent to remove the device from the account (Meross account). Causes device to reset waiting to be configured.

Method: `PUSH`

```json
{}
```

#### Appliance.Digest.TimerX

Method: `GET`

```json
{}
```

Method: `GETACK`

```json
{
  "digest": [
    {
      "channel": 0,
      "id": "abc1234",
      "count": 6
    },
    {
      "channel": 0,
      "id": "def4321",
      "count": 3
    },
    {
      "channel": 0,
      "id": "ghi666",
      "count": 5
    }
  ]
}
```

#### Appliance.Digest.TriggerX

Method: `GET`

```json
{}
```

Method: `GETACK`

```json
{
  "digest": [
    {
      "channel": 0,
      "id": "abc1234",
      "count": 1
    }
  ]
}
```


#### Appliance.Control.Upgrade

Method: `SET`

```json
{
  "upgrade": {
    "md5": "60685b8a9dbc02f4fe18085f22d37c87",
    "url": "http://bucket-meross-static.meross.com/production/upload/2018/12/07/14/04/55/201812071404554673948.bin"
  }
}
```

| Field        | Description          |
|--------------|----------------------|
| .upgrade.md5 | md5 hash of file     |
| .upgrade.url | URL to retrieve file |

Method: `SETACK`

```json
{}
```

Once appliance has upgraded a status is reported to indicate success.

Method: `PUSH`

```json
{
  "upgradeInfo": {
    "status": 1
  }
}
```
