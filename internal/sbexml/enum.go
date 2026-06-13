package sbexml

func buildEnum(enum indexedEnum) enumType {
	result := enumType{
		Name:         enum.name,
		EncodingType: "uint8",
	}

	for _, value := range enum.descriptor.Value {
		result.Values = append(result.Values, validValue{
			Name:  value.GetName(),
			Value: value.GetNumber(),
		})
	}

	return result
}
