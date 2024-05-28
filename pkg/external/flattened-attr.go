package external

import "log/slog"

// FlattenedAttr is a struct that represents a flattened attribute.
//
// This is exactly the same than slog.Attr but:
//   - The Key is always a string with the format "group1.group2.group3.key" (if there are 3 encapsulated groups for this key)
//   - The Value can't be a group (this is the interest, attributes are flattened with group prefixes in the Key)
type FlattenedAttr struct {
	slog.Attr
}

// newFlattenedAttr creates a new FlattenedAttr from a slog.Attr and a currentGroup (can be empty).
//
// WARNING: attr must not be a group!
func newFlattenedAttr(attr slog.Attr, currentGroup string) FlattenedAttr {
	if currentGroup == "" {
		return FlattenedAttr{Attr: attr}
	}
	return FlattenedAttr{
		Attr: slog.Attr{
			Key:   currentGroup + "." + attr.Key,
			Value: attr.Value,
		},
	}
}

// newFlattenedAttrs creates a slice of FlattenedAttr from a slice of slog.Attr and a currentGroup (can be empty).
//
// note: groups in attrs are recursively flattened
func newFlattenedAttrs(attrs []slog.Attr, currentGroup string) []FlattenedAttr {
	res := []FlattenedAttr{}
	for _, attr := range attrs {
		if attr.Value.Kind() == slog.KindGroup {
			if attr.Key == "" {
				continue
			}
			var newCurrentGroup string
			if currentGroup == "" {
				newCurrentGroup = attr.Key
			} else {
				newCurrentGroup = currentGroup + "." + attr.Key
			}
			res = append(res, newFlattenedAttrs(attr.Value.Group(), newCurrentGroup)...)
		} else {
			res = append(res, newFlattenedAttr(attr, currentGroup))
		}
	}
	return res
}
