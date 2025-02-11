package scte

import (
	"bytes"
)

func splice_desc_header(r *Reader, s *DescAny) {
	r.Decode(&s.Tag, 8)
	r.Decode(&s.Len, 8)
	r.Decode(&s.ID, 32)
}

func splice_desc_any(r *Reader, s *DescAny) {
	splice_desc_header(r, s)

	s.Data = make([]byte, s.Len-4)
	for i := 0; i < len(s.Data); i++ {
		r.Decode(&s.Data[i], 8)
	}
}

func avail_desc(r *Reader, s *DescAvail) {
	splice_desc_header(r, &s.DescAny)
	r.Decode(&s.ProviderID, 32)
}

func audio_desc(r *Reader, s *DescAudio) {
	splice_desc_header(r, &s.DescAny)
	r.Decode(&s.Count, 4)
	r.Decode(&s.Res, 4)
	s.Audio = make([]Audio, s.Count)
	for i := 0; i < len(s.Audio); i++ {
		r.Decode(&s.Audio[i].Tag, 8)
		r.Decode(&s.Audio[i].ISO, 24)
		r.Decode(&s.Audio[i].Mode, 3)
		r.Decode(&s.Audio[i].Channels, 4)
		r.Decode(&s.Audio[i].FullService, 1)
	}
}

func time_desc(r *Reader, s *DescTime) {
	splice_desc_header(r, &s.DescAny)
	r.Decode(&s.TAIseconds, 48)
	r.Decode(&s.TAIns, 42)
	r.Decode(&s.UTCOffset, 16)
}

func dtmf_desc(r *Reader, s *DescDTMF) {
	splice_desc_header(r, &s.DescAny)
	r.Decode(&s.Preroll, 8)
	r.Decode(&s.DTMFLen, 3)
	r.Decode(&s.Res, 5)
	s.DTMF = make([]byte, s.DTMFLen)
	for i := 0; i < len(s.DTMF); i++ {
		r.Decode(&s.DTMF[i], 8)
	}
}

func segmentation_desc(r *Reader, s *DescSegment) {
	splice_desc_header(r, &s.DescAny)
	eod := r.Offset() - 4 + int(s.Len)

	r.Decode(&s.EventID, 32)
	r.Decode(&s.Cancel, 1)
	r.Decode(&s.Compliant, 1)
	r.Decode(&s.Res0, 6)

	if s.Cancel {
		return
	}
	r.Decode(&s.Segmented, 1)
	r.Decode(&s.HasDuration, 1)
	r.Decode(&s.DeliveryUnrestricted, 1)
	if !s.DeliveryUnrestricted {
		r.Decode(&s.WebDelivery, 1)
		r.Decode(&s.NoBlackout, 1)
		r.Decode(&s.CanArchive, 1)
		r.Decode(&s.Restrictions, 2)
	} else {
		r.Decode(&s.Res1, 5)
	}

	if s.HasDuration {
		r.Decode(&s.Duration, 40)
	}
	r.Decode(&s.UPIDType, 8)
	r.Decode(&s.UPIDLen, 8)
	s.UPID = make([]byte, s.UPIDLen)
	for i := 0; i < len(s.UPID); i++ {
		r.Decode(&s.UPID[i], 8)
	}
	r.Decode(&s.SegType, 8)
	r.Decode(&s.SegNum, 8)
	r.Decode(&s.SegExp, 8)

	const subseg = "\x30\x32\x34\x36\x38\x3a\x44\x46"
	if bytes.IndexAny([]byte{s.SegType}, subseg) >= 0 {
		if r.Offset()+8+8 <= eod {
			// non-compliant streams wont set these
			// even if they are supposed to do so
			r.Decode(&s.SubSegNum, 8)
			r.Decode(&s.SubSegExp, 8)
		} else {
			s.Error = "warning: sub_segment_num and sub_segments_expected are required but missing"
		}
	}
	s.SegTypeName = segtype2name[s.SegType]

}
