package bot

/*
#include "../../build/Kada//lib/c/collections/string_view.h"

bool string_is_empty(string_t string)
{
	return string_equals(string, string_null);
}
*/
import "C"
import "strings"

type SizedString C.string_t // Representation of a sized string.

// Represent the given sized c-string as a [SizedString].
// Return the given sized c-string as a [SizedString].
func SizedStringFrom(sizedString C.string_t) SizedString {
	return SizedString(sizedString)
}

// Represent a [SizedString] as a go string.
// Returns a given sized string as a go string.
func (sizedString SizedString) String() string {
	return C.GoStringN(sizedString.data, C.int(sizedString.count))
}

// Return a given sized string as a go string without whitespace.
// Returns a sized string as a go string without whitespace.
func (sizedString SizedString) Trim() string {
	return strings.TrimSpace(sizedString.String())
}

// Determine if the sized string is empty.
// Returns true if the length of the string is equal to zero and the internal string data is equal to nil, else false.
func (sizedString SizedString) IsEmpty() bool {
	return bool(C.string_is_empty(C.string_t(sizedString)))
}
