import homeassistant.helpers.config_validation as cv
import voluptuous as vol
from typing import Any
from collections.abc import Mapping
from homeassistant.components.switch import SwitchEntity, PLATFORM_SCHEMA, ATTR_CURRENT_POWER_W
from homeassistant.const import STATE_ON, STATE_UNKNOWN, CONF_SWITCHES, CONF_NAME, ATTR_VOLTAGE
from .meross import MQTTDevice, SystemState, ToggleState, PowerUsage

from .const import (
    CONF_KEY,
    CONF_UUID,
    CONF_VALIDATE,
    
    ATTR_VERSION,
    ATTR_MAC,
    ATTR_IP,
    ATTR_CURRENT_A,
)

'''
    uuid: UUID as reported by meross firmware (required)
    name: Name to identify appliance/device eg. Living Room Lamp (required)
    key:  Auth key, if different from platform key (optional)
'''
SWITCH_SCHEMA = vol.Schema({
    vol.Required(CONF_UUID): cv.string,
    vol.Required(CONF_NAME): cv.string,
    vol.Optional(CONF_KEY): cv.string,
})

'''
    switches: [] one or more SWITCH_SCHEMA (required)
    key:  Auth key to use for all appliances/devices (optional)
    validate: validate key for incomming requests (default: True)
'''
PLATFORM_SCHEMA = PLATFORM_SCHEMA.extend({
    vol.Required(CONF_SWITCHES): cv.ensure_list(SWITCH_SCHEMA),
    vol.Optional(CONF_KEY): cv.string,
    vol.Optional(CONF_VALIDATE, default=True): cv.boolean,
})

def setup_platform(hass, config, add_entities, discovery_info=None):
    switches = []
    for swConfig in config[CONF_SWITCHES]:
        switches.append(
            MerossSwitch(
                hass.components.mqtt,
                swConfig[CONF_UUID],
                next(filter(None, [swConfig.get(CONF_KEY), config.get(CONF_KEY)])),
                swConfig[CONF_NAME],
                config.get(CONF_VALIDATE),
            )
        )
    return add_entities(switches, update_before_add=True)

class MerossSwitch(SwitchEntity):
    def __init__(self, mqtt, id, key, name, validate):
        self.device = MQTTDevice(id, key, validate, self.mqttUpdate).start(mqtt)

        # preconfigured attributes
        self._name = name
        self._unique_id = id
        self.key = key

        # discovered later
        self._state = STATE_UNKNOWN
        self.mac = ""
        self.ip = ""
        self.version = ""
        
        self._emeter_params = {}

    def turn_on(self, **kwargs) -> None:
        self.device.SetOnOff(0, 1)

    def turn_off(self, **kwargs):
        self.device.SetOnOff(0, 0)

    @property
    def is_on(self):
        return self._state == STATE_ON

    @property
    def name(self):
        """Return the display name of this switch."""
        return self._name

    @property
    def should_poll(self) -> bool:
        return True

    @property
    def extra_state_attributes(self) -> Mapping[str, Any]:
        """Return the state attributes of the device."""
        return self._emeter_params

    @property
    def state_attributes(self):
        state_attr = super().state_attributes
        state_attr[ATTR_VERSION] = self.version
        state_attr[ATTR_MAC] = self.mac
        state_attr[ATTR_IP] = self.ip
        return state_attr

    @property
    def unique_id(self) -> str:
        """Return a unique ID."""
        return "meross_{}".format(self._unique_id)

    def update(self):
        self.device.SystemAll()
        self.device.GetElectricityUsage()

    def mqttUpdate(self, state):
        """ Callback from MQTT device to report on state changes """
        if isinstance(state, ToggleState):
            self._state = state.state

        if isinstance(state, SystemState):
            self.mac = state.mac
            self.ip = state.ip
            self.version = state.version

        if isinstance(state, PowerUsage):
            # Inspired from components/tplink/switch.py
            self._emeter_params[ATTR_CURRENT_POWER_W] = round(float(state.power), 2)
            self._emeter_params[ATTR_VOLTAGE]         = round(float(state.voltage), 1)
            self._emeter_params[ATTR_CURRENT_A]       = round(float(state.current), 1)

        self.async_write_ha_state()
