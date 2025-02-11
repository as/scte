// Package scte implements a SCTE_35 packet reader based on the following standard:
//
// https://dutchguild.nl/event/13/attachments/82/203/SCTE_35_2023r1.pdf
//
// It decodes all non-deprecated commands and produces an informative data structure on the
// command and any optional segment descriptors attached to it. SCTE-35 packets are used
// primarily to signal AD-breaks and content decisioning. They are carried in MPEG transport
// streams as well as DASH and HLS media playlists in base64 or hexidecimal formats.
package scte

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

// Parse parses a SCTE-35 command and any descriptors. The input format can be a raw
// bitstream, hexidecimal, or base64 encoded data.
func Parse(data []byte) (p Packet, err error) {
	switch {
	case bytes.HasPrefix(data, []byte("/D")):
		data, err = base64.StdEncoding.DecodeString(strings.TrimSuffix(string(data), "\n"))
	case bytes.HasPrefix(data, []byte("0x")):
		fallthrough
	case bytes.HasPrefix(data, []byte("f")):
		fallthrough
	case bytes.HasPrefix(data, []byte("F")):
		data, err = hex.DecodeString(strings.TrimSuffix(string(data), "\n"))
	}
	if err != nil {
		return p, err
	}
	r := NewReader(data)
	p.Decode(r)
	return p, r.Err()
}

func printjson(v any) {
	p, _ := json.Marshal(v)
	fmt.Println(string(p))
}
