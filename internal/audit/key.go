package audit

import (
	"encoding/hex"
	"fmt"
	"os"
)

func LoadSecretKey() ([]byte, error) {
	keyS := os.Getenv("AUDIT_HMAC_KEY")
	if keyS == "" {
		return nil, fmt.Errorf("missing audit secret key")
	}

	key, err := hex.DecodeString(keyS)
	if err != nil {
		return nil, fmt.Errorf("AUDIT_HMAC_KEY is not valid hex: %w", err)
	}

	if len(key) < 32 {
		return nil, fmt.Errorf("AUDIT_HMAC_KEY must be at least 32 bytes")
	}

	return key, nil
}
