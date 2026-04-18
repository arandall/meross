#pragma once

#include <Arduino.h>
#include <ArduinoJson.h>

// ----------------------------------------------------------------
// MDP (Meross Device Protocol) layer.
//
// All functions that call the network require the ESP8266 to be
// connected to the Meross device AP (10.10.10.1) before use.
// ----------------------------------------------------------------

// Generate a 32-character lowercase hex random message ID.
String mdp_random_id();

// Compute the MDP signing hash: md5(msgId + key + timestamp).
// key = "" is valid — unconfigured devices skip signature validation.
String mdp_sign(const String &msgId, const String &key, const String &ts);

// Build a complete MDP JSON packet string ready to POST.
String mdp_build_packet(const String &ns,
                        const String &method,
                        const String &payload,
                        const String &key);

// Base64 encode raw bytes / a String.
String mdp_base64_encode(const uint8_t *data, size_t len);
String mdp_base64_encode(const String &s);

// Base64 decode. Returns true on success.
bool mdp_base64_decode(const String &encoded, String &out);

// POST a packet to the Meross device and deserialise the JSON response.
// Returns false on HTTP error or JSON parse failure.
bool mdp_post(const String &packetJson, JsonDocument &doc);

// Extract the "method" field from a parsed response document.
String mdp_method(const JsonDocument &doc);

// Return true when the response is an ERROR packet.
bool mdp_is_error(const JsonDocument &doc);
