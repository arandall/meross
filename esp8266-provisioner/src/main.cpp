#include <Arduino.h>
#include <ESP8266WiFi.h>
#include <ArduinoJson.h>
#include <time.h>

#include "config.h"
#include "mdp.h"

// ----------------------------------------------------------------
// State machine
// ----------------------------------------------------------------
enum State {
    NTP_SYNC,
    SCANNING,
    CONNECTING,
    TRACE,
    SYSTEM_ALL,
    WIFI_LIST_CHECK,
    WIFI_LIST_WAIT,
    PROVISION_KEY,
    PROVISION_WIFI,
    COOLDOWN
};

// ----------------------------------------------------------------
// Per-device attempt tracking (rate-limiting)
// ----------------------------------------------------------------
struct DeviceAttempt {
    char   bssid[18];   // "AA:BB:CC:DD:EE:FF" + '\0'
    unsigned long lastAttempt;
    bool   used;
};

static DeviceAttempt s_attempts[MAX_TRACKED_DEVICES];

static DeviceAttempt *find_attempt(const char *bssid) {
    for (int i = 0; i < MAX_TRACKED_DEVICES; i++) {
        if (s_attempts[i].used && strcmp(s_attempts[i].bssid, bssid) == 0)
            return &s_attempts[i];
    }
    return nullptr;
}

static DeviceAttempt *alloc_attempt(const char *bssid) {
    // Reuse existing slot
    DeviceAttempt *da = find_attempt(bssid);
    if (da) return da;
    // Find a free slot
    for (int i = 0; i < MAX_TRACKED_DEVICES; i++) {
        if (!s_attempts[i].used) {
            s_attempts[i].used = true;
            strncpy(s_attempts[i].bssid, bssid, sizeof(s_attempts[i].bssid) - 1);
            s_attempts[i].bssid[sizeof(s_attempts[i].bssid) - 1] = '\0';
            s_attempts[i].lastAttempt = 0;
            return &s_attempts[i];
        }
    }
    // No free slot — evict the oldest
    DeviceAttempt *oldest = &s_attempts[0];
    for (int i = 1; i < MAX_TRACKED_DEVICES; i++) {
        if (s_attempts[i].lastAttempt < oldest->lastAttempt)
            oldest = &s_attempts[i];
    }
    strncpy(oldest->bssid, bssid, sizeof(oldest->bssid) - 1);
    oldest->bssid[sizeof(oldest->bssid) - 1] = '\0';
    oldest->lastAttempt = 0;
    oldest->used = true;
    return oldest;
}

// ----------------------------------------------------------------
// Meross AP detection
// ----------------------------------------------------------------
static bool is_meross_ap(const String &ssid, const String &bssid) {
    if (ssid.startsWith("Meross_")) return true;
    // Match known Meross OUIs (first 8 chars of BSSID "XX:XX:XX")
    if (strncasecmp(bssid.c_str(), "48:e1:e9", 8) == 0) return true;
    if (strncasecmp(bssid.c_str(), "34:29:8f", 8) == 0) return true;
    return false;
}

// ----------------------------------------------------------------
// State variables
// ----------------------------------------------------------------
static State        s_state     = NTP_SYNC;
static unsigned long s_stateTs  = 0;       // millis() at last state entry
static String       s_targetBssid;
static String       s_targetSsid;
static DeviceAttempt *s_currentAttempt = nullptr;

// MAC address of the Meross device currently being provisioned.
static char s_deviceMAC[18];

// ----------------------------------------------------------------
// helpers
// ----------------------------------------------------------------
static void enter_state(State s) {
    s_state   = s;
    s_stateTs = millis();
}

static bool target_ssid_visible() {
    int n = WiFi.scanComplete();
    // A fresh scan is needed; check inline
    int count = WiFi.scanNetworks();
    for (int i = 0; i < count; i++) {
        if (WiFi.SSID(i) == TARGET_SSID) return true;
    }
    return false;
}

// ----------------------------------------------------------------
// setup / loop
// ----------------------------------------------------------------
void setup() {
    Serial.begin(115200);
    delay(100);
    Serial.println("\n[boot] Meross provisioner starting");

    memset(s_attempts, 0, sizeof(s_attempts));

    // Connect to home WiFi to get NTP sync
    Serial.printf("[wifi] Connecting to %s\n", HOME_SSID);
    WiFi.mode(WIFI_STA);
    Serial.printf("[boot] Provisioner MAC: %s\n", WiFi.macAddress().c_str());
    WiFi.begin(HOME_SSID, HOME_PASSWORD);
    while (WiFi.status() != WL_CONNECTED) {
        delay(500);
        Serial.print(".");
    }
    Serial.printf("\n[wifi] Connected, IP: %s\n", WiFi.localIP().toString().c_str());

    configTime(0, 0, NTP_SERVER_1, NTP_SERVER_2);
    enter_state(NTP_SYNC);
}

void loop() {
    switch (s_state) {

    // ---- Wait for NTP clock ------------------------------------
    case NTP_SYNC: {
        if (time(nullptr) > 1700000000L) {
            Serial.println("[ntp] Time synced");
            enter_state(SCANNING);
        }
        break;
    }

    // ---- Scan for Meross APs -----------------------------------
    case SCANNING: {
        if (millis() - s_stateTs < SCAN_INTERVAL_MS && s_stateTs != 0) break;

        Serial.println("[scan] Scanning for Meross APs...");
        int n = WiFi.scanNetworks();
        for (int i = 0; i < n; i++) {
            String ssid  = WiFi.SSID(i);
            String bssid = WiFi.BSSIDstr(i);
            if (!is_meross_ap(ssid, bssid)) continue;

            DeviceAttempt *da = find_attempt(bssid.c_str());
            if (da && (millis() - da->lastAttempt < ATTEMPT_COOLDOWN_MS)) {
                Serial.printf("[scan] Skipping %s (cooldown)\n", bssid.c_str());
                continue;
            }

            Serial.printf("[scan] Found Meross AP: %s (%s)\n",
                          ssid.c_str(), bssid.c_str());
            s_targetSsid  = ssid;
            s_targetBssid = bssid;
            s_currentAttempt = alloc_attempt(bssid.c_str());
            enter_state(CONNECTING);
            return;
        }
        s_stateTs = millis();  // reset scan timer
        break;
    }

    // ---- Connect to the Meross device AP -----------------------
    case CONNECTING: {
        Serial.printf("[wifi] Connecting to Meross AP: %s\n", s_targetSsid.c_str());
        WiFi.disconnect();
        WiFi.begin(s_targetSsid.c_str());  // open AP, no password
        unsigned long deadline = millis() + 15000UL;
        while (WiFi.status() != WL_CONNECTED && millis() < deadline) {
            delay(200);
            Serial.print(".");
        }
        Serial.println();

        if (WiFi.status() != WL_CONNECTED) {
            Serial.println("[wifi] Connection failed, back to SCANNING");
            enter_state(SCANNING);
            break;
        }
        Serial.printf("[wifi] Connected to Meross AP, IP: %s\n",
                      WiFi.localIP().toString().c_str());

        // Decide whether to run TRACE first
        if (s_currentAttempt && s_currentAttempt->lastAttempt != 0) {
            enter_state(TRACE);
        } else {
            enter_state(SYSTEM_ALL);
        }
        break;
    }

    // ---- Appliance.Config.Trace --------------------------------
    case TRACE: {
        Serial.println("[mdp] Sending Appliance.Config.Trace");
        String pkt = mdp_build_packet(
            "Appliance.Config.Trace", "GET", "{}", "");
        JsonDocument doc;
        if (mdp_post(pkt, doc)) {
            Serial.print("[mdp] Trace response: ");
            serializeJson(doc, Serial);
            Serial.println();
        } else {
            Serial.println("[mdp] Trace failed");
        }
        enter_state(SYSTEM_ALL);
        break;
    }

    // ---- Appliance.System.All ----------------------------------
    case SYSTEM_ALL: {
        Serial.println("[mdp] Sending Appliance.System.All GET");
        String pkt = mdp_build_packet(
            "Appliance.System.All", "GET", "{}", "");
        JsonDocument doc;
        if (!mdp_post(pkt, doc) || mdp_is_error(doc)) {
            Serial.println("[mdp] System.All failed, back to SCANNING");
            enter_state(SCANNING);
            break;
        }
        strncpy(s_deviceMAC,
                doc["payload"]["all"]["system"]["hardware"]["macAddress"] | "",
                sizeof(s_deviceMAC) - 1);
        s_deviceMAC[sizeof(s_deviceMAC) - 1] = '\0';
        enter_state(WIFI_LIST_CHECK);
        break;
    }

    // ---- Check Appliance.Config.WifiList -----------------------
    case WIFI_LIST_CHECK: {
        Serial.println("[mdp] Checking Appliance.Config.WifiList");
        String pkt = mdp_build_packet(
            "Appliance.Config.WifiList", "GET", "{}", "");
        JsonDocument doc;
        if (!mdp_post(pkt, doc) || mdp_is_error(doc)) {
            Serial.println("[mdp] WifiList GET failed");
            enter_state(SCANNING);
            break;
        }

        // Check whether the target SSID is visible to the device
        JsonArray list = doc["payload"]["wifiList"].as<JsonArray>();
        bool found = false;
        for (JsonObject entry : list) {
            String enc = entry["ssid"].as<String>();
            String dec;
            if (mdp_base64_decode(enc, dec) && dec == TARGET_SSID) {
                found = true;
                break;
            }
        }

        if (found) {
            Serial.printf("[mdp] Target SSID '%s' visible, provisioning\n", TARGET_SSID);
            enter_state(PROVISION_KEY);
        } else {
            Serial.printf("[mdp] Target SSID '%s' not visible, waiting\n", TARGET_SSID);
            enter_state(WIFI_LIST_WAIT);
        }
        break;
    }

    // ---- Wait and retry WifiList -------------------------------
    case WIFI_LIST_WAIT: {
        if (millis() - s_stateTs < WIFI_LIST_RETRY_MS) break;
        enter_state(WIFI_LIST_CHECK);
        break;
    }

    // ---- Appliance.Config.Key ----------------------------------
    case PROVISION_KEY: {
        Serial.println("[mdp] Sending Appliance.Config.Key");

        String payload;
        payload  = "{\"key\":{";
        payload += "\"gateway\":{";
        payload += "\"host\":\"";       payload += MQTT_HOST;        payload += "\",";
        payload += "\"port\":";         payload += MQTT_PORT;         payload += ",";
        payload += "\"secondHost\":\""; payload += MQTT_SECOND_HOST; payload += "\",";
        payload += "\"secondPort\":";   payload += MQTT_SECOND_PORT;  payload += "},";
        payload += "\"key\":\"";        payload += MQTT_KEY;          payload += "\",";
        payload += "\"userId\":\"";     payload += MQTT_USER_ID;      payload += "\"}}";

        String pkt = mdp_build_packet("Appliance.Config.Key", "SET", payload, "");
        JsonDocument doc;
        if (!mdp_post(pkt, doc) || mdp_is_error(doc)) {
            Serial.println("[mdp] Config.Key failed");
            enter_state(SCANNING);
            break;
        }
        enter_state(PROVISION_WIFI);
        break;
    }

    // ---- Appliance.Config.Wifi ---------------------------------
    case PROVISION_WIFI: {
        Serial.println("[mdp] Sending Appliance.Config.Wifi");

        String encSsid = mdp_base64_encode(TARGET_SSID);
        String encPass = mdp_base64_encode(TARGET_PASSWORD);

        String payload;
        payload  = "{\"wifi\":{";
        payload += "\"ssid\":\"";     payload += encSsid; payload += "\",";
        payload += "\"password\":\""; payload += encPass; payload += "\",";
        payload += "\"channel\":0,\"encryption\":6,\"cipher\":2}}";

        String pkt = mdp_build_packet("Appliance.Config.Wifi", "SET", payload, "");
        JsonDocument doc;
        if (!mdp_post(pkt, doc) || mdp_is_error(doc)) {
            Serial.println("[mdp] Config.Wifi failed");
        } else {
            Serial.println("[mdp] Provisioning complete!");
            String sig = mdp_sign(String(s_deviceMAC), MQTT_KEY, "");
            Serial.println("[mdp] If using auth in your MQTT server use these credentials");
            Serial.printf("[mdp] Username: %s\n", s_deviceMAC);
            Serial.printf("[mdp] Password: %s_%s\n", MQTT_USER_ID, sig.c_str());
        }

        // Record the attempt time and enter cooldown
        if (s_currentAttempt) {
            s_currentAttempt->lastAttempt = millis();
        }
        enter_state(SCANNING);
        break;
    }
    }
}
