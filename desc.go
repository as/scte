package scte

// Desc is a splice descriptor. A splice descriptor is an extension
// to splice commands which allows them to transmit additional
// data along with their original command messages in the Packet
type Desc interface {
	Name() string
	Kind() int
}

// DescAny is a generic splice descriptor with unparsed bytes
type DescAny struct {
	Tag     byte   // 8
	Len     byte   // 8
	ID      int    // 32
	Data    []byte `json:",omitempty"`
	TagName string `json:",omitempty"`
	Error   string `json:",omitempty"`
}

// DescSegment is a segmentation descriptor that extends the time_signal
// and splice_insert commands. It is only valid for splice_insert, time_signal,
// and splice_null commands, and should be transmitted at least 4 seconds in
// advance of the signaled splice_time so that a device can interpolate the
// Packet (splice_info) section correctly
type DescSegment struct {
	DescAny
	EventID   int
	Cancel    bool
	Compliant bool
	Res0      int

	Segmented            bool
	HasDuration          bool
	DeliveryUnrestricted bool
	WebDelivery          bool
	NoBlackout           bool
	CanArchive           bool
	Restrictions         int

	Res1 int // 5

	Duration    int // 40
	UPIDType    int // 8
	UPIDLen     int // 8
	UPID        []byte
	SegType     byte // 8
	SegTypeName string
	SegNum      int // 8
	SegExp      int // 8
	SubSegNum   int // 8
	SubSegExp   int // 8
}

// DescAvail is an extension to splice_insert that allows authorization
// identifier transmission. Its purpose is to replicate the CUE tone used
// in analog systems for AD insertions and is only valid in the context of
// a splice_insert command.
type DescAvail struct {
	DescAny
	ProviderID int // 32
}

// DescDTMF (Dual-Tone Multi Frequency) descriptor is another extension
// to splice_insert that allows the reciever to generate an analog sequence
// based on the Packet (splice_info)
type DescDTMF struct {
	DescAny
	Preroll int // 8
	DTMFLen int // 3
	Res     int // 5
	DTMF    []byte
}

// DescTime specifies a time descriptor for the Precision Time Protocol (PTP)
// which uses a time format similar to UTC but without the addition of leap
// seconds. The descriptor stores the difference between the PTP TAI standard
// and the UTC standard so timestamps can be converted between the two
type DescTime struct {
	DescAny
	TAIseconds int // 48
	TAIns      int // 32
	UTCOffset  int // 16
}

// DescAudio is an audio descriptor for multi-channel video programming
// descributors (MPVDs) that can't signal dynamic audio language changes
// due to their audio formats. The descriptor is used to signal such changes
// instead and is only valid with the time_signal command and segmentation
// descriptors ProgramStart or ProgramOverlapStart
type DescAudio struct {
	DescAny
	Count int // 4
	Res   int // 4
	Audio []Audio
}

// Audio describes the structure of the audio tracks in a DescAudio descriptor
type Audio struct {
	Tag         int  // 8
	ISO         int  // 24
	Mode        int  // 3
	Channels    int  // 4
	FullService bool // 1
}

func (c DescAvail) Kind() int   { return 0x00 }
func (c DescDTMF) Kind() int    { return 0x01 }
func (c DescSegment) Kind() int { return 0x02 }
func (c DescTime) Kind() int    { return 0x03 }
func (c DescAudio) Kind() int   { return 0x04 }
func (c DescAny) Kind() int     { return int(c.Tag) }

func (c DescAvail) Name() string   { return "avail_descriptor" }
func (c DescDTMF) Name() string    { return "DTMF_descriptor" }
func (c DescSegment) Name() string { return "segmentation_descriptor" }
func (c DescTime) Name() string    { return "time_descriptor" }
func (c DescAudio) Name() string   { return "splice_null" }
func (c DescAny) Name() string     { return c.TagName }

var segtype2name = map[byte]string{
	0x00: "NotIndicated",
	0x01: "ContentIdentification",
	0x02: "Private",
	0x10: "ProgramStart",
	0x11: "ProgramEnd",
	0x12: "ProgramEarlyTermination",
	0x13: "ProgramBreakaway",
	0x14: "ProgramResumption",
	0x15: "ProgramRunoverPlanned",
	0x16: "ProgramRunoverUnplanned",
	0x17: "ProgramOverlapStart",
	0x18: "ProgramBlackoutOverride",
	0x19: "ProgramJoin",
	0x20: "ChapterStart",
	0x21: "ChapterEnd",
	0x22: "BreakStart",
	0x23: "BreakEnd",
	0x24: "OpeningCreditStart_deprecated",
	0x25: "OpeningCreditEnd_deprecated",
	0x26: "ClosingCreditStart_deprecated",
	0x27: "ClosingCreditEnd_deprecated",
	0x30: "ProviderAdvertisementStart",
	0x31: "ProviderAdvertisementEnd",
	0x32: "DistributorAdvertisementStart",
	0x33: "DistributorAdvertisementEnd",
	0x34: "ProviderPlacementOpportunityStart",
	0x35: "ProviderPlacementOpportunityEnd",
	0x36: "DistributorPlacementOpportunityStart",
	0x37: "DistributorPlacementOpportunityEnd",
	0x38: "ProviderOverlayPlacementOpportunityStart",
	0x39: "ProviderOverlayPlacementOpportunityEnd",
	0x3A: "DistributorOverlayPlacementOpportunityStart",
	0x3B: "DistributorOverlayPlacementOpportunityEnd",
	0x3C: "ProviderPromoStart",
	0x3D: "ProviderPromoEnd",
	0x3E: "DistributorPromoStart",
	0x3F: "DistributorPromoEnd",
	0x40: "UnscheduledEventStart",
	0x41: "UnscheduledEventEnd",
	0x42: "AlternateContentOpportunityStart",
	0x43: "AlternateContentOpportunityEnd",
	0x44: "ProviderAdBlockStart",
	0x45: "ProviderAdBlockEnd",
	0x46: "DistributorAdBlockStart",
	0x47: "DistributorAdBlockEnd",
	0x50: "NetworkStart",
	0x51: "NetworkEnd",
}
