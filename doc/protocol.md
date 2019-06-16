# Meross Device/Appliance Protocol

This document is to explain the protocol used between Meross IoT appliances and the Meross cloud service.

## Abilities

A appliance returns what abilities it is able to support. I have so far only had an `mss310` socket and assume there are
additional abilities.

I hope others will be able to contribute to this list.

| Ability                                                         | Description
|-----------------------------------------------------------------|------------
| [Appliance.Config.Key](#applianceconfigkey)                     | Used to configure MQTT servers
| [Appliance.Config.Trace](#applianceconfigtrace)                 | Returns WiFi and system details during setup
| [Appliance.Config.Wifi](#applianceconfigwifi)                   | Configures WiFi network to connect to during setup
| [Appliance.Config.WifiList](#applianceconfigwifiList)           | Lists aviable Wifi networks
| [Appliance.Control.Bind](#appliancecontrolbind)                 | Have not observed use.
| [Appliance.Control.ConsumptionX](#appliancecontrolconsumptionx) | Shows daily power consumption from the last 30 days
| [Appliance.Control.Electricity](#appliancecontrolelectricity)   | Returns present electricity usage
| [Appliance.Control.Timer](#appliancecontroltimer)               | Used to GET/SET Timer values
| [Appliance.Control.Toggle](#appliancecontroltoggle)             | Used Toggle switch on/off
| [Appliance.Control.Trigger](#appliancecontroltrigger)           | Used to GET/SET Trigger rules
| [Appliance.Control.Unbind](#appliancecontrolunbind)             | Have not observed use.
| [Appliance.Control.Upgrade](#appliancecontrolupgrade)           | Upgrade firmware from URL
| [Appliance.System.Ability](#appliancesystemability)             | List all abilities appliance supports
| [Appliance.System.All](#appliancesystemall)                     | List all system attributes (includes firmware, hardware, timer, trigger, toogle and online status
| [Appliance.System.Clock](#appliancesystemclock)                 | Used to request current time
| [Appliance.System.DNDMode](#appliancesystemdndmode)             | Used to set/unset DNDMode
| [Appliance.System.Debug](#appliancesystemdebug)                 | Get debug information about appliance OS
| [Appliance.System.Firmware](#appliancesystemfirmware)           | Get firmware information
| [Appliance.System.Hardware](#appliancesystemhardware)           | Get hardware information
| [Appliance.System.Online](#appliancesystemonline)               | Get online status
| [Appliance.System.Position](#appliancesystemposition)           | Get/Set Position (lat/lng) of appliance
| [Appliance.System.Report](#appliancesystemreport)               | PUSH data back to server (only seen time reporting)
| [Appliance.System.Runtime](#appliancesystemruntime)             | Get runtime, only observed WiFi "signal" strength
| [Appliance.System.Time](#appliancesystemtime)                   | Get/Set timezone and daylight savings rules

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
    "timestamp": 1557596606
  },
  "payload": {}
}
```

## Headers

Each header contains the following keys

| Field          | Description
|----------------|---
| from           | HTTP url of appliance config endpoint or MQTT topic of the endpoint that generated packet.
| messageId      | Arbitrary ID (32 characters of HEX) - responses will use the same `messageID`
| method         | used to describe request/response action `GET`/`GETACK`, `SET`/`SETACK`, `PUSH`, `ERROR`
| namespace      | [ability](#abilities) being used
| payloadVersion | I assume this is for future payload versions, only ever seen `1` being used
| sign           | Signing value equal to md5(`messageId` + `key` + `timestamp`) where `key` is defined in `Appliance.Config.Key`
| timestamp      | Time in seconds past Epoch

**Note:** When a appliance is waiting to be configured it does not validate the sign value.

## Errors

Errors are returned buy any `GET` or `SET` request that has a problem. Typically these are sign errors.

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

The payload of a packet is a JSON structure dependant on the ability being used. Most `GET` methods use a payload of
`{}` when requesting data however in some cases additional context is sent with the request.

### Config

The `namespaces` are used to configure an appliance. Typically these do not require a signing key as the appliance is yet to be
configured.

#### Appliance.Config.Key

Used to set the MQTT servers and signing key of a appliance.

Method: `SET`

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

| Field                   | Description
|-------------------------|---
| .key.gateway.host       | Primary MQTT host
| .key.gateway.port       | Primary MQTT port
| .key.gateway.secondHost | Secondary MQTT host
| .key.gateway.secondPort | Secondary MQTT port
| .key.key                | pre-shared key used for signing requests (usually assigned by iot.meross.com)
| .key.userId             | userId (usually assigned by iot.meross.com)

Method: `SETACK`

```json
{}
```

#### Appliance.Config.Trace

I have not found a use for this but it is used by the app when configuring a socket. Once a appliance is configured the
response contains empty values.

Method: `GET`

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

Used to set Wifi SSID/Password.

**Note:** Once appliance responds it will restart and connect to the network provided. If an error occurs the appliance will
reset and wait for configuration.

Method: `SET`

```json
{
  "wifi": {
    "bssid": "de-ad-00-00-be-ef",
    "channel": 3,
    "cipher": 3,
    "encryption": 6,
    "password": "cGFzc3dvcmQK",
    "ssid": "c3NpZAo="
  }
}
```

| Field            | Description
|------------------|---
| .wifi.bssid      | BSSID
| .wifi.channel    | Channel
| .wifi.cipher     | cipher used by Wifi network
| .wifi.encryption | encryption used by network
| .wifi.password   | base64 encoded string of password to connect to network
| .wifi.bssid      | base64 encoded string of SSID

Method: `SETACK`

```json
{}
```

#### Appliance.Config.WifiList

Used to list all Wifi networks the appliance can see.

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
      "bssid": "2c-6e-a4-55-c6-47",
      "signal": 100,
      "channel": 3,
      "encryption": 6,
      "cipher": 3
    },
    {
      "ssid": "dGVzdCBTU0lEMg==",
      "bssid": "4d-23-0e-3c-a2-22",
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

Gets supported abilities of an appliance.

Method: `GET`

```json
{}
```

Method: `GETACK`

```json
{
  "Appliance.Config.Key": {},
  "Appliance.Config.WifiList": {},
  "Appliance.Config.Wifi": {},
  "Appliance.Config.Trace": {},
  "Appliance.System.Online": {},
  "Appliance.System.All": {},
  "Appliance.System.Hardware": {},
  "Appliance.System.Firmware": {},
  "Appliance.System.Time": {},
  "Appliance.System.Clock": {},
  "Appliance.System.Debug": {},
  "Appliance.System.Ability": {},
  "Appliance.System.Runtime": {},
  "Appliance.System.Report": {},
  "Appliance.System.Position": {},
  "Appliance.System.DNDMode": {},
  "Appliance.Control.Toggle": {},
  "Appliance.Control.Timer": {},
  "Appliance.Control.Trigger": {},
  "Appliance.Control.Consumption": {},
  "Appliance.Control.ConsumptionX": {},
  "Appliance.Control.Electricity": {},
  "Appliance.Control.Upgrade": {},
  "Appliance.Control.Bind": {},
  "Appliance.Control.Unbind": {}
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

Method: `PUSH`

```json
{
  "clock": {
    "timestamp": 1560144444
  }
}
```

| Field            | Description
|------------------|---
| .clock.timestamp | Timestamp in seconds past epoch

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

| Field         | Description
|---------------|---
| .DNDMode.mode | 0 = off (status LED functional), 1 = on (status LED disabled)

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

| Field               | Description
|---------------------|---
| .position.longitude | Longitude multiplied by 1000000
| .position.latitude  | Latitude multiplied by 1000000

Method: `SETACK`

```json
{}
```

#### Appliance.System.Report

Used to report data back to Meross. I haven't seen this used by an Appliance with the exception of this example.

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

```json
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
      ...
    ]
  }
}
```

| Field             | Description
|-------------------|---
| .time.timestamp   | Current time of appliance
| .time.timezone    | Timezone of appliance
| .time.timeRule    | List of daylight savings rules
| .time.timeRule[0] | timestamp rule takes effect
| .time.timeRule[1] | UTC offset in seconds
| .time.timeRule[2] | Daylight savings active (0 = yes, 1 = no)

Method: `SETACK`

```json
{}
```

### Control

#### Appliance.Control.Bind

Not observed.

#### Appliance.Control.ConsumptionX

Method: `GET`

```json
{}
```

Method: `GETACK`

```json
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
    ...
  ]
}
```

| Field               | Description
|---------------------|---
| .consumptionx.date  | date in Y-m-d format
| .consumptionx.time  | start timestamp of period
| .consumptionx.value | Usage value in watt hours

#### Appliance.Control.Electricity

Method: `GET`

```json
{
  "electricity": {
    "channel": 0
  }
}
```

| Field                | Description
|----------------------|---
| .electricity.channel | I assume this is to support an appliance with multiple sockets

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

| Field                | Description
|----------------------|---
| .electricity.channel |
| .electricity.current | Current being consumed in milli-amps (mA)
| .electricity.voltage | Current voltage in deci-volts (dV) (/10 for Volts)
| .electricity.power   | Current power usage in Watts

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

| Field               | Description
|---------------------|---
| .timer[].id         | unique identifier of timer
| .timer[].type       | 1 = weekly, 2 = once
| .timer[].enable     | Enabled (0 = off, 1 = on)
| .timer[].alias      | User string to identify timer
| .timer[].time       | Minute of day to fire
| .timer[].week       | 8 bit bitset with MSB always on LSB -> MSB Sun, Mon, ... (eg. Monday 0b10000010)
| .timer[].duration   | set to 0 (not sure on usage)
| .timer[].createTime | Timestamp timer was created
| .timer[].extend     | [Appliance.Control.Toggle](#Appliance.Control.Toggle) object without channel

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

| Field         | Description
|---------------|---
| .channel      | I assume this is to support an appliance with multiple sockets
| .toggle.onoff | 0 = off, 1 = on, any other value locks current state

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

| Field                                 | Description
|---------------------------------------|---
| .trigger[].id                         | unique identifier of timer
| .trigger[].type                       | 1 = weekly, 2 = once
| .trigger[].enable                     | Enabled (0 = off, 1 = on)
| .trigger[].alias                      | User string to identify timer
| .trigger[].createTime                 | Timestamp timer was created
| .trigger[].rule._if_                  | [Appliance.Control.Toggle](#Appliance.Control.Toggle) object without channel
| .trigger[].rule._then_.delay.week     | 8 bit bitset with MSB always on LSB -> MSB Sun, Mon, ... (eg. Mon & Sun 0b10000011)
| .trigger[].rule._then_.delay.duration | Seconds to wait
| .trigger[].rule._do_                  | [Appliance.Control.Toggle](#Appliance.Control.Toggle) object without channel

Method: `SETACK`

```json
{}
```

#### Appliance.Control.Unbind

Not observed.

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

| Field        | Description
|--------------|---
| .upgrade.md5 | md5 hash of file
| .upgrade.url | URL to retrieve file

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
