#pragma once

// ============================================================
// WiFi the ESP8266 uses on boot to sync NTP time.
// May be the same as TARGET_SSID / TARGET_PASSWORD.
// ============================================================
#define HOME_SSID     "YourWiFiSSID"
#define HOME_PASSWORD "YourWiFiPassword"

// ============================================================
// Target WiFi network that Meross devices will be provisioned
// to connect to after configuration.
// ============================================================
#define TARGET_SSID     "YourTargetSSID"
#define TARGET_PASSWORD "YourTargetPassword"

// ============================================================
// MQTT broker details used in Appliance.Config.Key.
// Set MQTT_SECOND_HOST / MQTT_SECOND_PORT to empty/0 if you
// only have one broker.
// ============================================================
#define MQTT_HOST        "your.mqtt.broker.com"
#define MQTT_PORT        8883
#define MQTT_SECOND_HOST ""
#define MQTT_SECOND_PORT 0

// Pre-shared signing key assigned to your Meross account.
#define MQTT_KEY     "0123456789abcdef0123456789abcdef"

// User ID assigned to your Meross account (sent as a string).
#define MQTT_USER_ID "9876543210"

// ============================================================
// NTP servers used for time sync on boot.
// ============================================================
#define NTP_SERVER_1 "pool.ntp.org"
#define NTP_SERVER_2 "time.nist.gov"

// ============================================================
// Scanning & timing
// ============================================================

// How often (ms) to scan for Meross APs when idle.
#define SCAN_INTERVAL_MS 10000UL

// If a device was last attempted within this window, run
// Appliance.Config.Trace before continuing (5 minutes).
#define ATTEMPT_COOLDOWN_MS (5UL * 60UL * 1000UL)

// How long (ms) to wait between WifiList retries when the
// target SSID is not visible to the Meross device.
#define WIFI_LIST_RETRY_MS 20000UL

// How many distinct Meross BSSID → last-attempt entries to
// track simultaneously.
#define MAX_TRACKED_DEVICES 10

// Fixed IP / URL of every unconfigured Meross device AP.
#define MEROSS_DEVICE_URL "http://10.10.10.1/config"

// ============================================================
// Debug logging.
// Set to 1 to print every HTTP request/response body over Serial.
// ============================================================
#define MDP_DEBUG 0

// HTTP client timeout (ms) for requests to the Meross device.
#define HTTP_TIMEOUT_MS       30000
