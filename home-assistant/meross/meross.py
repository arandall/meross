import logging
import json
import hashlib
import random
import string
import time

from homeassistant.const import STATE_ON, STATE_OFF
from .const import (
    DOMAIN,
    LISTEN_TOPIC,
    PUBLISH_TOPIC,
    DEVICE_ON,
    APP_CONTROL_TOGGLE,
    APP_CONTROL_ELEC,
    APP_METHOD_GET,
    APP_METHOD_PUSH,
    APP_METHOD_SET,
    APP_SYS_ALL,
    APP_SYS_CLOCK,
)
_LOGGER = logging.getLogger(__name__)

class MQTTDevice():
    """
       MQTTDevice represents a Meross MQTT device\.

       The intent is to use this for more devices as their use of the common protocol is discovered.
    """

    def __init__(self, id, key, validate, callback):
        self.id = id
        self.key = key
        self.validate = validate
        self.callback = callback

    def start(self, service):
        """ starts listening to topic for device requests/responses """
        self.service = service
        topic = LISTEN_TOPIC.format(self.id)
        _LOGGER.info("%s: device MQTT subscription to %s", DOMAIN, topic)
        self.service.subscribe(topic, self.message_received)
        return self

    def message_received(self, msg):
        """ handler for incoming messages from a device """
        p = Packet(json.loads(msg.payload))
        _LOGGER.debug("received %s", p)
        if not p.validSignature(self.key):
            if self.validate:
                _LOGGER.error("invalid signature %s: %s", self.id, p.header.namespace)
                return
            else:
                _LOGGER.info("ignoreing signature error: %s", self.id, p.header.namespace)

        if p.header.method == "ERROR":
            _LOGGER.error("error occured: %s", self.id, p.payload)
            return

        if p.header.method == APP_METHOD_PUSH and p.header.namespace == APP_CONTROL_TOGGLE:
            self.callback(
                ToggleState(
                    p.payload.get("channel", 0),
                    p.payload["toggle"]["onoff"],
                )
            )

        # Respond to clock events with the current time.
        if p.header.method == APP_METHOD_PUSH and p.header.namespace == APP_SYS_CLOCK:
            self.sendPacket(
                self.createPacket(
                    APP_METHOD_PUSH,
                    APP_SYS_CLOCK,
                    {
                        "clock": {
                            "timestamp": int(time.time()),
                        }
                    }
                )
            )

        if p.header.namespace == APP_CONTROL_ELEC:            
            self.callback(
                PowerUsage(p.payload["electricity"]["power"],p.payload["electricity"]["current"],p.payload["electricity"]["voltage"])
            )

        if p.header.namespace == APP_SYS_ALL:
            # TODO(arandall): check for channel (if applicable)
            self.callback(ToggleState(0, p.payload["all"]["control"]["toggle"]["onoff"]))
            self.callback(
                SystemState(
                    p.payload["all"]["system"]["hardware"]["macAddress"],
                    p.payload["all"]["system"]["firmware"]["innerIp"],
                    "{}-{} v{} - {} (fw v{})".format(
                        p.payload["all"]["system"]["hardware"]["type"],
                        p.payload["all"]["system"]["hardware"]["subType"],
                        p.payload["all"]["system"]["hardware"]["version"],
                        p.payload["all"]["system"]["hardware"]["chipType"],
                        p.payload["all"]["system"]["firmware"]["version"],
                    )
                )
            )

    def createPacket(self, method, namespace, payload):
        """ createPacket ready for transmission via MQTT or HTTP """
        p = Packet({
            "header": {
                "from": "homeassistant/meross/subscribe",
                "method": method,
                "namespace": namespace,
            },
            "payload": payload
        })
        p.sign(self.key)
        return p

    def sendPacket(self, p):
        """ sendPacket via MQTT """
        _LOGGER.debug("sending %s", p)
        self.service.publish(PUBLISH_TOPIC.format(self.id), json.dumps(p, default=serialize))

    def SetOnOff(self, channel, state):
        """ SetOnOff state for given channel """
        self.sendPacket(
            self.createPacket(
                APP_METHOD_SET,
                APP_CONTROL_TOGGLE,
                {
                    "channel": channel,
                    "toggle": {
                        "onoff": state,
                    }
                }
            )
        )

    def GetElectricityUsage(self):
        self.sendPacket(
            self.createPacket(
                APP_METHOD_GET,
                APP_CONTROL_ELEC,
                {}
            )
        )

    def SystemAll(self):
        """ Get all system parameters """
        self.sendPacket(
            self.createPacket(
                APP_METHOD_GET,
                APP_SYS_ALL,
                {}
            )
        )

class ToggleState():
    """ ToggleState represents the state of a switch """
    def __init__(self, channel, state):
        self.channel = channel
        if state == DEVICE_ON:
            self.state = STATE_ON
        else:
            self.state = STATE_OFF

    def __str__(self):
        return "ToggleState channel:{} state:{}".format(self.channel, self.state)

class SystemState():
    """ SystemState represents the state of a Meross Appliance """
    def __init__(self, mac, ip, version):
        self.mac = mac
        self.ip = ip
        self.version = version

    def __str__(self):
        return "SystemState mac:{} ip:{} version:{}".format(self.mac, self.ip, self.version)

class PowerUsage():
    """ PowerUsage represents the current power usage of a Meross Appliance """
    def __init__(self, power, current, voltage):
        self.power = power / 1000
        self.current = current / 1000
        self.voltage = voltage / 10

    def __str__(self):
        return "PowerUsage {}W".format(self.power/1000)

class Header():
    def __init__(self, pk):
        if pk == None:
            pk = {}
        self.from_ = pk.get("from")
        self.messageId = pk.get("messageId", ''.join(random.SystemRandom().choice(string.ascii_lowercase + string.digits) for _ in range(32)))
        self.method = pk.get("method")
        self.namespace = pk.get("namespace")
        self.payloadVersion = pk.get("payloadVersion", 1)
        self.sign = pk.get("sign")
        self.timestamp = pk.get("timestamp", int(time.time()))

class Packet():
    def __init__(self, pk):
        self.header = Header(pk.get("header"))
        self.payload = pk.get("payload")

    def __str__(self):
        return "meross-packet({}): {} - {} [{}] ".format(self.header.messageId, self.header.method, self.header.namespace, self.header.from_)

    def calcSignature(self, key):
        signatureString = ""
        for arg in [self.header.messageId, key, self.header.timestamp]:
            signatureString += str(arg)
        return hashlib.md5(signatureString.encode()).hexdigest()

    def sign(self, key):
        self.header.sign = self.calcSignature(key)

    def validSignature(self, key):
        return self.header.sign == self.calcSignature(key)

def serialize(obj):
    if obj is Header:
        dict = obj.__dict__
        dict["from"] = dict.pop("from_")
        return dict
    return obj.__dict__