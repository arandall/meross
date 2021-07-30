# Meross Home Assistant integration

A quick and dirty implmentation to get meross switches working in Home Assistant.

## Installation

Copy the contents of this directory to a directory `custom-components/meross` of you HA configuration.

## Attributes exposed

The meross switches monitor power usage and as a result I've exposed these values as attributes. You must use template
platform to expose these as sensors. (see below)

`current_power_w` power currently being used in Watts
`voltage` voltage measured in Volts
`current_a` current mesured in Amperes

## Configuration

```
switch:
  - platform: meross
    key: "00110011001100110011001100110011" # Key to use for switch if not overriden
    switches:
      - uuid: "11111111111111111111111111111111" # UUID of your switch
        name: "My Lamp" # Name to give switch
      - uuid: "22222222222222222222222222222222"
        name: "Heater"
        key: "11111111111111111111111111111111" # Override key if diffent from platform default
```

### Power usage graph

To enable power usage you need to create a sonsor from the switch attribute you desire

```
sensor:
  - platform: template
    sensors:
      my_lamp:
        friendly_name: "My Lamp"
        unit_of_measurement: 'W'
        device_class: power
        value_template: "{{ state_attr('switch.my_lamp','current_power_w') }}"
```