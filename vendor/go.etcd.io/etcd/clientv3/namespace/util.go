package namespace

func prefixInterval(pfx string, key, end []byte) (pfxKey []byte, pfxEnd []byte) {
	pfxKey = make([]byte, len(pfx)+len(key))
	copy(pfxKey[copy(pfxKey, pfx):], key)

	if len(end) == 1 && end[0] == 0 {
		// the edge of the keyspace
		pfxEnd = make([]byte, len(pfx))
		copy(pfxEnd, pfx)
		ok := false
		for i := len(pfxEnd) - 1; i >= 0; i-- {
			if pfxEnd[i]++; pfxEnd[i] != 0 {
				ok = true
				break
			}
		}
		if !ok {
			// 0xff..ff => 0x00
			pfxEnd = []byte{0}
		}
	} else if len(end) >= 1 {
		pfxEnd = make([]byte, len(pfx)+len(end))
		copy(pfxEnd[copy(pfxEnd, pfx):], end)
	}

	return pfxKey, pfxEnd
}
