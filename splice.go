package scte

// Cmd is a SCTE-35 Splice command, which can be any of the following:
// - SpliceNull (0x00)
// - TimeSignal (0x04)
// - SpliceSchedule (0x05)
// - SpliceInsert (0x06)
// - Bandwidth (0x07)
type Cmd interface {
	Name() string
	Type() int
}

// SpliceNull is an empty command. It is used as a placeholder to transmit
// descriptors without an explicit command
type SpliceNull struct {
}

// SpliceSchedule is a list of splice_insert commands signaled in advance
// for more information see the SpliceInsert documentation
type SpliceSchedule struct {
	Count  int // 8
	Splice []SpliceInsert
}

// TimeSignal is used to signal the future presence of SpliceInserts. It carries
// segmentation descriptors with ProviderPlacement opportunities or ProgramStarts
// for conditioning purposes
type TimeSignal struct {
	Time
}

// SpliceInsert signals an upcoming splice event. A splice event is an opportunity
// for downstream equipment to schedule an AD-break or AD-insertion. This can
// be a Cue-In where the new content starts or a Cue-Out where the content ends.
type SpliceInsert struct {
	Cmd          Cmd
	EventID      int  // 32
	Cancel       bool // 1
	Res0         int  // 7
	OutOfNetwork bool // 1
	HasSplice    bool // 1
	HasDuration  bool // 1
	Immediate    bool // 1
	Compliant    bool // 1
	Res1         int  // 6

	CompLen int    // 8
	Comp    []Comp `json:",omitempty"`

	Time       Time
	AutoReturn bool // 1
	Res2       int  // 6
	BreakDur   int

	ProgID   int
	Avail    int
	AvailExp int
}

// Bandwidth reservation is an empty command
// that can be dropped by networking equipment.
type Bandwidth struct {
}

// Private is a private command. This implementation currently ignores private
// commands at the time of writing.
type Private struct {
}

// Comp is a component. The use of this is deprecated by the standard
// and it only exists to parse the bitstream correctly
type Comp struct {
	Tag  int // 8
	Time Time
}

// Time is a splice_time, which is present in the SpliceInsert and TimeSignal
// commands, containing a PTS with a timebase of 90kHz. Combined with the
// PTS Adjustment in the Packet, represents the intended time of the splice point
type Time struct {
	HasPTS bool // 1
	Res    int  // 6 or 7
	PTS    int  // 3
}

func (c SpliceNull) Name() string     { return "splice_null" }
func (c SpliceSchedule) Name() string { return "splice_schedule" }
func (c SpliceInsert) Name() string   { return "splice_insert" }
func (c TimeSignal) Name() string     { return "time_signal" }
func (c Bandwidth) Name() string      { return "bandwidth_reservation" }

func (c SpliceNull) Type() int     { return 0x00 }
func (c SpliceSchedule) Type() int { return 0x04 }
func (c SpliceInsert) Type() int   { return 0x05 }
func (c TimeSignal) Type() int     { return 0x06 }
func (c Bandwidth) Type() int      { return 0x07 }
