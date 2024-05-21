package slogh

import "log/slog"

type FlattenedAttr struct {
	slog.Attr
}

func NewFlattenedAttr(attr slog.Attr, currentGroup string) FlattenedAttr {
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

func NewFlattenedAttrs(attrs []slog.Attr, currentGroup string) []FlattenedAttr {
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
			res = append(res, NewFlattenedAttrs(attr.Value.Group(), newCurrentGroup)...)
		} else {
			res = append(res, NewFlattenedAttr(attr, currentGroup))
		}
	}
	return res
}
