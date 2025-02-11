package scte

func splice_null(r *Reader, si *SpliceNull) {
}

func splice_schedule(r *Reader, s *SpliceSchedule) {
	r.Decode(&s.Count, 8)
	s.Splice = make([]SpliceInsert, s.Count)
	for i := 0; i < int(s.Count); i++ {
		splice_insert(r, &s.Splice[i])
	}
}

func splice_insert(r *Reader, si *SpliceInsert) {
	r.Decode(&si.EventID, 32)
	r.Decode(&si.Cancel, 1)
	r.Decode(&si.Res0, 7)

	r.Decode(&si.OutOfNetwork, 1)
	r.Decode(&si.HasSplice, 1)
	r.Decode(&si.HasDuration, 1)
	r.Decode(&si.Immediate, 1)
	r.Decode(&si.Compliant, 1)
	r.Decode(&si.Res1, 3)

	if si.HasSplice && !si.Immediate {
		splice_time(r, &si.Time)
	} else if !si.HasSplice {
		// this is deprecated as per SCTE_35_2023r1.pdf
		r.Decode(&si.CompLen, 8)
		si.Comp = make([]Comp, int(si.CompLen))
		for i := 0; i < len(si.Comp); i++ {
			r.Decode(&si.Comp[i].Tag, 8)
			if !si.Immediate {
				splice_time(r, &si.Comp[i].Time)
			}
		}
	}
	if si.HasDuration { // break_duration() ss9.8.2
		r.Decode(&si.AutoReturn, 1)
		r.Decode(&si.Res2, 6)
		r.Decode(&si.BreakDur, 33)
	}
	r.Decode(&si.ProgID, 16)
	r.Decode(&si.Avail, 8)
	r.Decode(&si.AvailExp, 8)
}

func time_signal(r *Reader, ts *TimeSignal) {
	splice_time(r, &ts.Time)
}

func bandwidth_res(r *Reader, bw *Bandwidth) {
}

func private(r *Reader, priv *Private) {

}

//
// Helper functions
//

// splice_time() ss9.8.1
func splice_time(r *Reader, t *Time) {
	r.Decode(&t.HasPTS, 1)
	if t.HasPTS {
		r.Decode(&t.Res, 6)
		r.Decode(&t.PTS, 33)
	} else {
		r.Decode(&t.Res, 7)
	}
}
