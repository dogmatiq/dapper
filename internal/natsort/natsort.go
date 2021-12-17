package natsort

import (
	"regexp"
	"strconv"
)

// chunkre is the regular expression to break a string into the chunks with
// alpha- or numeric-only content.
var chunkre = regexp.MustCompile(`(\d+|\D+)`)

// chunks breaks the string into alpha and numeric chunks using chunkre regular
// expression.
func chunks(s string) []string {
	return chunkre.FindAllString(s, -1)
}

// Less reports whether string 'a' should sort before string 'b'. It uses the
// Alphanum algorithm described at this link:
//
// http://davekoelle.com/alphanum.html.
func Less(a, b string) bool {
	ac, bc := chunks(a), chunks(b)

	for i := range ac {
		// This is a safeguard in case the number of string 'b' chunks is
		// smaller than the number of string 'a' chunks.
		if i >= len(bc) {
			return false
		}

		ai64, aerr := strconv.ParseInt(ac[i], 10, 0)
		bi64, berr := strconv.ParseInt(bc[i], 10, 0)

		if aerr == nil && berr == nil {
			if ai64 == bi64 {
				continue
			}

			return ai64 < bi64
		}

		// At this point it is established that either or both of the chunks are
		// not numeric and they should be compared as strings.
		if ac[i] == bc[i] {
			continue
		}

		return ac[i] < bc[i]
	}

	// At this point all the chunks in string 'a' have been equal to the
	// corresponding chunks in string 'b'. However, if the number of chunks of
	// string "b" is greater, it is safe to assume that the string 'a' should
	// sort before string 'b'.
	//
	// Otherwise, it is established that the strings are equal and string 'a'
	// should not sort before string 'b'.
	return len(ac) < len(bc)
}
