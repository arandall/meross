#include <unity.h>
#include <base64.hpp>
#include <string.h>
#include <stdint.h>

// ----------------------------------------------------------------
// Helpers: encode/decode to/from C strings.
// ----------------------------------------------------------------
static void b64_enc(const char *input, char *output) {
    encode_base64(
        reinterpret_cast<uint8_t *>(const_cast<char *>(input)),
        strlen(input),
        reinterpret_cast<uint8_t *>(output));
}

static void b64_dec(const char *input, char *output) {
    decode_base64(
        reinterpret_cast<uint8_t *>(const_cast<char *>(input)),
        reinterpret_cast<uint8_t *>(output));
}

// ----------------------------------------------------------------
// Test vectors sourced from:
//   - mdp/testdata/GETACK-wifi-scan.json  (real device responses)
//   - doc/protocol.md                     (protocol examples)
//   - RFC 4648 §10 test vectors
// ----------------------------------------------------------------

void setUp()    {}
void tearDown() {}

// -- Encode: RFC 4648 canonical vectors ----------------------------

void test_encode_empty() {
    char out[4] = {};
    encode_base64((uint8_t *)"", 0, (uint8_t *)out);
    TEST_ASSERT_EQUAL_STRING("", out);
}

void test_encode_one_byte() {
    char out[8] = {};
    b64_enc("a", out);
    TEST_ASSERT_EQUAL_STRING("YQ==", out);
}

void test_encode_two_bytes() {
    char out[8] = {};
    b64_enc("ab", out);
    TEST_ASSERT_EQUAL_STRING("YWI=", out);
}

void test_encode_three_bytes_no_padding() {
    char out[8] = {};
    b64_enc("ABC", out);
    TEST_ASSERT_EQUAL_STRING("QUJD", out);
}

void test_encode_rfc_man() {
    char out[8] = {};
    b64_enc("Man", out);
    TEST_ASSERT_EQUAL_STRING("TWFu", out);
}

// -- Encode: real Meross testdata vectors --------------------------

// SSID values from mdp/testdata/GETACK-wifi-scan.json
void test_encode_test_ssid() {
    char out[20] = {};
    b64_enc("test SSID", out);
    TEST_ASSERT_EQUAL_STRING("dGVzdCBTU0lE", out);
}

void test_encode_test_ssid2() {
    char out[24] = {};
    b64_enc("test SSID2", out);
    TEST_ASSERT_EQUAL_STRING("dGVzdCBTU0lEMg==", out);
}

// Password example from doc/protocol.md ("password\n")
void test_encode_password_with_newline() {
    char out[20] = {};
    b64_enc("password\n", out);
    TEST_ASSERT_EQUAL_STRING("cGFzc3dvcmQK", out);
}

// SSID example from doc/protocol.md ("ssid\n")
void test_encode_ssid_with_newline() {
    char out[12] = {};
    b64_enc("ssid\n", out);
    TEST_ASSERT_EQUAL_STRING("c3NpZAo=", out);
}

// -- Decode: real Meross testdata vectors --------------------------

void test_decode_test_ssid() {
    char out[16] = {};
    b64_dec("dGVzdCBTU0lE", out);
    TEST_ASSERT_EQUAL_STRING("test SSID", out);
}

void test_decode_test_ssid2() {
    char out[16] = {};
    b64_dec("dGVzdCBTU0lEMg==", out);
    TEST_ASSERT_EQUAL_STRING("test SSID2", out);
}

void test_decode_password_with_newline() {
    char out[16] = {};
    b64_dec("cGFzc3dvcmQK", out);
    TEST_ASSERT_EQUAL_STRING("password\n", out);
}

void test_decode_ssid_with_newline() {
    char out[12] = {};
    b64_dec("c3NpZAo=", out);
    TEST_ASSERT_EQUAL_STRING("ssid\n", out);
}

void test_decode_empty() {
    char out[4] = {};
    unsigned int len = decode_base64((uint8_t *)"", (uint8_t *)out);
    TEST_ASSERT_EQUAL_UINT(0, len);
}

// -- Round-trip ----------------------------------------------------

void test_roundtrip_short_string() {
    const char *original = "hello";
    char enc[16] = {};
    char dec[16] = {};
    encode_base64((uint8_t *)original, strlen(original), (uint8_t *)enc);
    decode_base64((uint8_t *)enc, (uint8_t *)dec);
    TEST_ASSERT_EQUAL_STRING(original, dec);
}

void test_roundtrip_binary_bytes() {
    const uint8_t data[] = {0x00, 0xFF, 0x80, 0x7F, 0x01, 0xFE};
    char enc[16] = {};
    uint8_t dec[8] = {};
    encode_base64(const_cast<uint8_t *>(data), sizeof(data), (uint8_t *)enc);
    unsigned int decLen = decode_base64((uint8_t *)enc, dec);
    TEST_ASSERT_EQUAL_INT(sizeof(data), decLen);
    TEST_ASSERT_EQUAL_MEMORY(data, dec, sizeof(data));
}

void test_roundtrip_length_not_multiple_of_3() {
    // 10 bytes exercises all three padding cases
    const char *original = "1234567890";
    char enc[20] = {};
    char dec[16] = {};
    encode_base64((uint8_t *)original, strlen(original), (uint8_t *)enc);
    decode_base64((uint8_t *)enc, (uint8_t *)dec);
    TEST_ASSERT_EQUAL_STRING(original, dec);
}

// ----------------------------------------------------------------

int main(int argc, char **argv) {
    UNITY_BEGIN();

    RUN_TEST(test_encode_empty);
    RUN_TEST(test_encode_one_byte);
    RUN_TEST(test_encode_two_bytes);
    RUN_TEST(test_encode_three_bytes_no_padding);
    RUN_TEST(test_encode_rfc_man);
    RUN_TEST(test_encode_test_ssid);
    RUN_TEST(test_encode_test_ssid2);
    RUN_TEST(test_encode_password_with_newline);
    RUN_TEST(test_encode_ssid_with_newline);

    RUN_TEST(test_decode_test_ssid);
    RUN_TEST(test_decode_test_ssid2);
    RUN_TEST(test_decode_password_with_newline);
    RUN_TEST(test_decode_ssid_with_newline);
    RUN_TEST(test_decode_empty);

    RUN_TEST(test_roundtrip_short_string);
    RUN_TEST(test_roundtrip_binary_bytes);
    RUN_TEST(test_roundtrip_length_not_multiple_of_3);

    return UNITY_END();
}
