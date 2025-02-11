package scte

// Packet is the start of a splice_info_section, containing command metadata
// the splice command, and a trailer possibly containing additional descriptors.
type Packet struct {
	Header // 112

	// Cmd is the actual splice command. It may be a null command
	// that only carries descriptors.
	Cmd Cmd

	// Trailer carries a variable length section of descriptors as well
	// as a checksum for encrypted and plaintext streams.
	Trailer // ?
}

type Header struct {
	Table   int   // 8
	SSI     bool  // 1
	Priv    bool  // 1
	SAP     int   // 2
	Len     int   // 12
	Ver     int   // 8
	Enc     bool  // 1
	EncAlg  int   // 6
	PTSA    int64 // 33
	CWI     int   // 8
	Tier    int   // 12
	CmdLen  int   // 12
	CmdType int   // 8

	CmdName string // human-readable command type; not encoded
}

type Trailer struct {
	DescLen  int // 16
	Desc     []Desc
	Stuffing []byte
	ECRC32   int // 32
	CRC32    int // 32
}

// Decode decodes the packet from the binary reader
func (c *Packet) Decode(r *Reader) error {
	r.Decode(&c.Table, 8)
	r.Decode(&c.SSI, 1)
	r.Decode(&c.Priv, 1)
	r.Decode(&c.SAP, 2)
	r.Decode(&c.Len, 12)
	r.Decode(&c.Ver, 8) // 4 bytes total
	r.Decode(&c.Enc, 1)
	r.Decode(&c.EncAlg, 6)
	r.Decode(&c.PTSA, 33)
	r.Decode(&c.CWI, 8) // 6 bytes total
	r.Decode(&c.Tier, 12)
	r.Decode(&c.CmdLen, 12)
	r.Decode(&c.CmdType, 8)

	switch c.CmdType {
	case 0x00:
		s := SpliceNull{}
		splice_null(r, &s)
		c.Cmd = s
	case 0x04:
		s := SpliceSchedule{}
		splice_schedule(r, &s)
		c.Cmd = s
	case 0x05:
		s := SpliceInsert{}
		splice_insert(r, &s)
		c.Cmd = s
	case 0x06:
		s := TimeSignal{}
		time_signal(r, &s)
		c.Cmd = s
	case 0x07:
		s := Bandwidth{}
		bandwidth_res(r, &s)
		c.Cmd = s
	case 0x08:
	}
	if c.Cmd != nil {
		c.CmdName = c.Cmd.Name()
	}

	r.Decode(&c.DescLen, 16)
	len := r.Offset() + (c.DescLen * 8)
	for r.Offset() < len {
		switch r.Peek(8) {
		case 0x02:
			desc := DescSegment{}
			desc.TagName = desc.Name()
			segmentation_desc(r, &desc)
			c.Desc = append(c.Desc, desc)
		default:
			desc := DescAny{}
			splice_desc_any(r, &desc)
			c.Desc = append(c.Desc, desc)
		}
	}

	// TODO: Stuffing?
	if c.Enc {
		r.Decode(&c.ECRC32, 32)
	}
	r.Decode(&c.CRC32, 32)
	return r.Err()
}
