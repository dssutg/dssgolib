package utils

// NaturalCmp compares two strings using natural ordering.
// For example, "file2" < "file10", unlike lexicographical ordering
// that seems weird in this particular case and should not be used.
// This routine is very useful for sorting lists with numeric suffixes.
// List can be sorted like this: slices.SortFunc(v, NaturalCmp).
// Returns -1 if a < b, 1 if a > b, 0 if equal.
// This function is optimized for speed and does not allocate anything.
// Only ASCII digits are considered as digits.
// The comparison is done in linear time.
func NaturalCmp(a, b string) int {
	ai, bi := 0, 0
	an, bn := len(a), len(b)

	for ai < an || bi < bn { // while either of the strings can be iterated
		if ai >= an || bi >= bn { // if one of the strings exhausted
			if ai >= an { // a ended earlier, b has remaining, so a < b
				return -1
			}
			if bi >= bn { // b ended earlier, a has remaining, so a > b
				return 1
			}
		}

		// Get next segment type: digit run or non-digit run.
		aIsDigit := a[ai] >= '0' && a[ai] <= '9'
		bIsDigit := b[bi] >= '0' && b[bi] <= '9'

		if aIsDigit && bIsDigit {
			// Compare numeric segments.
			aiStart, biStart := ai, bi

			// Skip leading zeros but count them.
			for ai < an && a[ai] == '0' {
				ai++
			}
			for bi < bn && b[bi] == '0' {
				bi++
			}

			// Find end of number (remaining digits after leading zeros).
			aiDigitsStart := ai
			for ai < an && a[ai] >= '0' && a[ai] <= '9' {
				ai++
			}
			biDigitsStart := bi
			for bi < bn && b[bi] >= '0' && b[bi] <= '9' {
				bi++
			}

			aDigitsLen := ai - aiDigitsStart
			bDigitsLen := bi - biDigitsStart

			// If both numbers are all zeros (e.g., "000"), digitsLen == 0 => treat as zero.
			if aDigitsLen == 0 && bDigitsLen == 0 {
				// Tie on numeric value; break tie by shorter leading-zero sequence first.
				aLeading := aiDigitsStart - aiStart
				bLeading := biDigitsStart - biStart
				if aLeading != bLeading {
					if aLeading < bLeading {
						return -1
					}
					return 1
				}
				continue
			}

			// Compare by length of significant digits.
			if aDigitsLen != bDigitsLen {
				if aDigitsLen < bDigitsLen {
					return -1
				}
				return 1
			}

			// Same length -> lexicographic compare of digit runs.
			for i := 0; i < aDigitsLen; i++ {
				ac := a[aiDigitsStart+i]
				bc := b[biDigitsStart+i]
				if ac != bc {
					if ac < bc {
						return -1
					}
					return 1
				}
			}

			// Numeric values equal; tie-breaker: fewer leading zeros sorts first.
			aLeading := aiDigitsStart - aiStart
			bLeading := biDigitsStart - biStart
			if aLeading != bLeading {
				if aLeading < bLeading {
					return -1
				}
				return 1
			}

			// Continue from current ai, bi (they already point after digit runs).
			continue
		}

		if aIsDigit != bIsDigit {
			// One digit-run vs non-digit-run: decide by comparing the first character types.
			// We consider digit < non-digit (matches common natural sort).
			if aIsDigit {
				return -1
			}
			return 1
		}

		// Both non-digits: compare run of non-digits.

		// Find the end of both non-digit runs.
		aiStart, biStart := ai, bi
		for ai < an && (a[ai] < '0' || a[ai] > '9') {
			ai++
		}
		for bi < bn && (b[bi] < '0' || b[bi] > '9') {
			bi++
		}

		// Compare lexicographically the bytes that are in both runs.
		aLen := ai - aiStart
		bLen := bi - biStart
		for i := range min(aLen, bLen) {
			ac := a[aiStart+i]
			bc := b[biStart+i]
			if ac != bc {
				if ac < bc {
					return -1
				}
				return 1
			}
		}

		// If one run shorter, shorter run sorts first.
		if aLen != bLen {
			if aLen < bLen {
				return -1
			}
			return 1
		}
	}

	// All inequality checks are done, both strings exhausted,
	// so they are equal.
	return 0
}
