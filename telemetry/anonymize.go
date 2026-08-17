package telemetry

import (
	"crypto/sha256"
	"encoding/hex"
)

func Anonymize(value string) string {
	hash := sha256.Sum256([]byte(value))
	return hex.EncodeToString(hash[:])
}

func AnonymizeWiFi(info WiFiInfo) WiFiInfo {
	info.BSSID = Anonymize(info.BSSID)
	return info
}