#include "mdp.h"
#include "config.h"

#include <Arduino.h>
#include <ESP8266HTTPClient.h>
#include <WiFiClient.h>
#include <MD5Builder.h>
#include <base64.hpp>
#include <time.h>

// ----------------------------------------------------------------
// Random message ID: 32 lowercase hex chars from hardware RNG.
// ----------------------------------------------------------------
String mdp_random_id() {
    String id;
    id.reserve(32);
    for (int i = 0; i < 8; i++) {
        uint32_t r = RANDOM_REG32;  // ESP8266 hardware RNG register
        char buf[9];
        snprintf(buf, sizeof(buf), "%08x", r);
        id += buf;
    }
    return id;
}

// ----------------------------------------------------------------
// Signing: md5(msgId + key + timestamp)
// Matches the Go mdp package: GenerateSignature(msgId, key, ts)
// ----------------------------------------------------------------
String mdp_sign(const String &msgId, const String &key, const String &ts) {
    MD5Builder md5;
    md5.begin();
    md5.add(msgId);
    md5.add(key);
    md5.add(ts);
    md5.calculate();
    return md5.toString();
}

// ----------------------------------------------------------------
// Packet builder.
// Constructs the full JSON packet string, signs the header.
// key = "" for unconfigured devices (signature not validated).
// ----------------------------------------------------------------
String mdp_build_packet(const String &ns,
                        const String &method,
                        const String &payload,
                        const String &key) {
    String msgId = mdp_random_id();
    long ts = (long)time(nullptr);
    String tsStr = String(ts);
    String sign  = mdp_sign(msgId, key, tsStr);

    String pkt;
    pkt.reserve(256 + payload.length());
    pkt  = "{\"header\":{";
    pkt += "\"messageId\":\"";  pkt += msgId;   pkt += "\",";
    pkt += "\"method\":\"";     pkt += method;  pkt += "\",";
    pkt += "\"namespace\":\"";  pkt += ns;      pkt += "\",";
    pkt += "\"payloadVersion\":1,";
    pkt += "\"sign\":\"";       pkt += sign;    pkt += "\",";
    pkt += "\"timestamp\":";    pkt += tsStr;   pkt += "},";
    pkt += "\"payload\":";      pkt += payload;
    pkt += "}";
    return pkt;
}

// ----------------------------------------------------------------
// Base64 — wrappers around Densaugeo/base64 (Densaugeo/base64 @ ^1.2.1).
// encode_base64 / decode_base64 are pure C++, no Arduino deps.
// ----------------------------------------------------------------
String mdp_base64_encode(const uint8_t *data, size_t len) {
    unsigned int outLen = encode_base64_length(len);
    uint8_t *buf = new uint8_t[outLen + 1];
    encode_base64(const_cast<uint8_t *>(data), len, buf);
    buf[outLen] = '\0';
    String result(reinterpret_cast<char *>(buf));
    delete[] buf;
    return result;
}

String mdp_base64_encode(const String &s) {
    return mdp_base64_encode(
        reinterpret_cast<const uint8_t *>(s.c_str()), s.length());
}

bool mdp_base64_decode(const String &encoded, String &out) {
    const uint8_t *src    = reinterpret_cast<const uint8_t *>(encoded.c_str());
    unsigned int   outLen = decode_base64_length(const_cast<uint8_t *>(src));
    uint8_t *buf = new uint8_t[outLen + 1];
    decode_base64(const_cast<uint8_t *>(src), buf);
    buf[outLen] = '\0';
    out = reinterpret_cast<char *>(buf);
    delete[] buf;
    return true;
}

// ----------------------------------------------------------------
// HTTP POST to the Meross device and parse the JSON response.
// ----------------------------------------------------------------
bool mdp_post(const String &packetJson, JsonDocument &doc) {
    WiFiClient wifiClient;
    HTTPClient http;

    http.begin(wifiClient, MEROSS_DEVICE_URL);
    http.setTimeout(HTTP_TIMEOUT_MS);
    http.addHeader("Content-Type", "application/json; charset=UTF-8");

#if MDP_DEBUG
    Serial.printf("[mdp] --> POST %s\n", MEROSS_DEVICE_URL);
    Serial.printf("[mdp] --> %s\n", packetJson.c_str());
#endif

    int code = http.POST(packetJson);
    if (code != 200) {
        Serial.printf("[mdp] HTTP error %d\n", code);
        http.end();
        return false;
    }

    String body = http.getString();
    http.end();

#if MDP_DEBUG
    Serial.printf("[mdp] <-- %d\n", code);
    Serial.printf("[mdp] <-- %s\n", body.c_str());
#endif

    DeserializationError err = deserializeJson(doc, body);
    if (err) {
        Serial.printf("[mdp] JSON parse error: %s\n", err.c_str());
        return false;
    }
    return true;
}

// ----------------------------------------------------------------
// Helpers for inspecting response packets.
// ----------------------------------------------------------------
String mdp_method(const JsonDocument &doc) {
    const char *m = doc["header"]["method"] | "";
    return String(m);
}

bool mdp_is_error(const JsonDocument &doc) {
    return mdp_method(doc) == "ERROR";
}
