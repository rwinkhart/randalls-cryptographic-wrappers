# RCW (Randall's Cryptographic Wrappers)
RCW is a cascading symmetric cryptography agent meant to be embedded within Go programs.

It encrypts all data with both AES256-GCM and ChaCha20-Poly1305.

Passphrases are securely cached for three minutes and RPC authentication is used to
ensure that only the binary+user responsible for caching the passphrase can utilize it.
This feature is supported on Linux, FreeBSD, MacOS, and Windows.

RCW also features a sanity check to ensure no data loss occurs due to a user entering the
incorrect passphrase during encryption.

Please note that RCW is a work-in-progress and breaking changes should be expected.
Future versions may not be capable of decrypting the output of the current version.

> [!WARNING]
>It is your responsibility to assess the security and stability of RCW and to ensure it meets your needs before using it.
>I am not responsible for any data loss or breaches of your information resulting from the use of RCW.
>RCW is a new project that is constantly being updated, and though safety and security are priorities, they cannot be guaranteed.

# Usage
For now, please reference [example.go](https://github.com/rwinkhart/randalls-cryptographic-wrappers/blob/main/example.go).
