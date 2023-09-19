package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

func ComputeChecksum(data []byte, removeWhiteSpaces bool) string {
	var dataTrimmed []byte
	if removeWhiteSpaces {
		dataString := strings.Join(strings.Fields(string(data)), "")
		dataTrimmed = []byte(dataString)
	} else {
		dataTrimmed = data
	}
	dataChecksumBytes := md5.Sum(dataTrimmed)
	return hex.EncodeToString(dataChecksumBytes[:])
}
