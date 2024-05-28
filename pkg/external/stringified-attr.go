package external

// StringifiedAttr is a struct that represents a stringified slog.Attr attribute.
//
// This is exactly the same than FlattenedAttr but the Value is resolved to a string.
type StringifiedAttr struct {
	Key   string
	Value string
}

func newStringifiedAttr(attr FlattenedAttr) StringifiedAttr {
	return StringifiedAttr{
		Key:   attr.Key,
		Value: attr.Value.Resolve().String(),
	}
}

func newStringifiedAttrs(attrs []FlattenedAttr) []StringifiedAttr {
	res := make([]StringifiedAttr, len(attrs))
	for i, attr := range attrs {
		res[i] = newStringifiedAttr(attr)
	}
	return res
}

// String returns the string representation of the StringifiedAttr as "key=value".
func (sa StringifiedAttr) String() string {
	return sa.Key + "=" + sa.Value
}
