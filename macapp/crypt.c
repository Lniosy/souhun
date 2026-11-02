#include <CommonCrypto/CommonCryptor.h>
#include <CommonCrypto/CommonKeyDerivation.h>
#include <stdlib.h>
#include <string.h>

int grok_decrypt_v10(const uint8_t *cipher, int cipher_len, const char *password, uint8_t **out, int *out_len) {
    uint8_t key[16];
    uint8_t iv[16];
    uint8_t *buf;
    size_t moved = 0;
    if (cipher_len < 0) return -1;
    CCCryptorStatus status;

    if (cipher == NULL || password == NULL || out == NULL || out_len == NULL) return -1;
    if (CCKeyDerivationPBKDF(kCCPBKDF2, password, strlen(password), (const uint8_t *)"saltysalt", 9,
                             kCCPRFHmacAlgSHA1, 1003, key, sizeof(key)) != kCCSuccess) {
        return -2;
    }
    memset(iv, ' ', sizeof(iv));
    buf = malloc(cipher_len + kCCBlockSizeAES128);
    if (buf == NULL) return -3;
    status = CCCrypt(kCCDecrypt, kCCAlgorithmAES, kCCOptionPKCS7Padding, key, sizeof(key), iv, cipher, cipher_len, buf,
                     cipher_len + kCCBlockSizeAES128, &moved);
    if (status != kCCSuccess) {
        free(buf);
        return -4;
    }
    *out = buf;
    *out_len = (int)moved;
    return 0;
}
