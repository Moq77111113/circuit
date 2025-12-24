package auth

import (
	"crypto/subtle"
	"encoding/base64"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

// isArgon2 checks if a password string is in argon2 PHC format.
func isArgon2(s string) bool {
	return strings.HasPrefix(s, "$argon2id$") || strings.HasPrefix(s, "$argon2i$")
}

// verifyArgon2 verifies a password against an argon2 PHC encoded hash.
func verifyArgon2(encoded, password string) bool {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 {
		return false
	}

	variant := parts[1]
	if variant != "argon2id" && variant != "argon2i" {
		return false
	}

	params := parseParams(parts[3])
	memory, ok := params["m"]
	if !ok {
		return false
	}
	time, ok := params["t"]
	if !ok {
		return false
	}
	threads, ok := params["p"]
	if !ok {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}

	var computedHash []byte
	if variant == "argon2id" {
		computedHash = argon2.IDKey([]byte(password), salt, time, memory, uint8(threads), uint32(len(expectedHash)))
	} else {
		computedHash = argon2.Key([]byte(password), salt, time, memory, uint8(threads), uint32(len(expectedHash)))
	}

	return subtle.ConstantTimeCompare(computedHash, expectedHash) == 1
}

// parseParams parses PHC parameter string like "m=65536,t=3,p=4".
func parseParams(s string) map[string]uint32 {
	params := make(map[string]uint32)
	for _, pair := range strings.Split(s, ",") {
		kv := strings.Split(pair, "=")
		if len(kv) != 2 {
			continue
		}
		val, err := strconv.ParseUint(kv[1], 10, 32)
		if err != nil {
			continue
		}
		params[kv[0]] = uint32(val)
	}
	return params
}
