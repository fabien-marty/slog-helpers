package slogh

type StringifiedAttr struct {
	Key   string
	Value string
}

func NewStringifiedAttr(attr FlattenedAttr) StringifiedAttr {
	return StringifiedAttr{
		Key:   attr.Key,
		Value: attr.Value.Resolve().String(),
	}
}

func NewStringifiedAttrs(attrs []FlattenedAttr) []StringifiedAttr {
	res := make([]StringifiedAttr, len(attrs))
	for i, attr := range attrs {
		res[i] = NewStringifiedAttr(attr)
	}
	return res
}

func (sa StringifiedAttr) String() string {
	return sa.Key + "=" + sa.Value
}
