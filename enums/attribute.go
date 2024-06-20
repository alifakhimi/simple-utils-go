package enums

// AttributeType ...
type AttributeType int

const (
	AttributeString AttributeType = iota + 1
	AttributeInteger
	AttributeFloat
	AttributeBool
	AttributeDateTime
	AttributeColor
)

var attributeTypeTextDic = map[AttributeType]string{
	AttributeString:   "String",
	AttributeInteger:  "Integer",
	AttributeFloat:    "Float",
	AttributeBool:     "Bool",
	AttributeDateTime: "Date Time",
	AttributeColor:    "Color",
}

func AttributeTypeText(qType AttributeType) string {
	return attributeTypeTextDic[qType]
}
