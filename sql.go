package simutils

import (
	valid "github.com/asaskevich/govalidator"
)

// Dict ...
var DictPostgre = map[rune]string{
	'آ': "(آ|ا|أ|إ|ع)",
	'ا': "(آ|ا|أ|إ|ع)",
	'أ': "(آ|ا|أ|إ|ع)",
	'إ': "(آ|ا|أ|إ|ع)",
	'ع': "(آ|ا|أ|إ|ع)",
	'ک': "(ک|ك)",
	'ك': "(ک|ك)",
	'و': "(و|ؤ)",
	'ؤ': "(و|ؤ)",
	'ه': "(ه|ة|ح)",
	'ة': "(ه|ة|ح)",
	'ح': "(ه|ة|ح)",
	'ی': "(ی|ي|ئ)",
	'ي': "(ی|ي|ئ)",
	'ئ': "(ی|ي|ئ)",
	'ت': "(ت|ط)",
	'ط': "(ت|ط)",
	'ق': "(ق|غ)",
	'غ': "(ق|غ)",
	'ز': "(ز|ذ|ظ|ض)",
	'ذ': "(ز|ذ|ظ|ض)",
	'ظ': "(ز|ذ|ظ|ض)",
	'ض': "(ز|ذ|ظ|ض)",
	'س': "(س|ث|ص)",
	'ث': "(س|ث|ص)",
	'ص': "(س|ث|ص)",
	'0': "(0|۰)",
	'۰': "(0|۰)",
	'1': "(1|۱)",
	'۱': "(1|۱)",
	'2': "(2|۲)",
	'۲': "(2|۲)",
	'3': "(3|۳)",
	'۳': "(3|۳)",
	'4': "(4|۴)",
	'۴': "(4|۴)",
	'5': "(5|۵)",
	'۵': "(5|۵)",
	'6': "(6|۶)",
	'۶': "(6|۶)",
	'7': "(7|۷)",
	'۷': "(7|۷)",
	'8': "(8|۸)",
	'۸': "(8|۸)",
	'9': "(9|۹)",
	'۹': "(9|۹)",
}

// Dict ...
var DictSql = map[rune]string{
	'آ': "[آاأإع]",
	'ا': "[آاأإع]",
	'أ': "[آاأإع]",
	'إ': "[آاأإع]",
	'ع': "[آاأإع]",
	'ک': "[کك]",
	'ك': "[کك]",
	'و': "[وؤ]",
	'ؤ': "[وؤ]",
	'ه': "[هةح]",
	'ة': "[هةح]",
	'ح': "[هةح]",
	'ی': "[یيئ]",
	'ي': "[یيئ]",
	'ئ': "[یيئ]",
	'ت': "[تط]",
	'ط': "[تط]",
	'ق': "[قغ]",
	'غ': "[قغ]",
	'ز': "[زذظض]",
	'ذ': "[زذظض]",
	'ظ': "[زذظض]",
	'ض': "[زذظض]",
	'س': "[سثص]",
	'ث': "[سثص]",
	'ص': "[سثص]",
	'0': "[0۰]",
	'۰': "[0۰]",
	'1': "[1۱]",
	'۱': "[1۱]",
	'2': "[2۲]",
	'۲': "[2۲]",
	'3': "[3۳]",
	'۳': "[3۳]",
	'4': "[4۴]",
	'۴': "[4۴]",
	'5': "[5۵]",
	'۵': "[5۵]",
	'6': "[6۶]",
	'۶': "[6۶]",
	'7': "[7۷]",
	'۷': "[7۷]",
	'8': "[8۸]",
	'۸': "[8۸]",
	'9': "[9۹]",
	'۹': "[9۹]",
}

// ArabicPersianAI ...
func ArabicPersianAI(driver DatabaseDriver, source string) (str string) {
	for _, ch := range source {
		switch driver {
		case SQLServer:
			if replaceCh, ok := DictSql[ch]; ok {
				str += replaceCh
			} else {
				str += string(ch)
			}
		default:
			if replaceCh, ok := DictPostgre[ch]; ok {
				str += replaceCh
			} else {
				str += string(ch)
			}
		}
	}

	return
}

func CorrectSimilarChars(driver DatabaseDriver, value interface{}) interface{} {
	switch v := value.(type) {
	case string:
		if !valid.IsASCII(v) {
			return ArabicPersianAI(driver, v)
		} else {
			return v
		}
	}

	return value
}

func NeedCorrectChar(value interface{}) bool {
	switch v := value.(type) {
	case string:
		// return !valid.IsFloat(v) && !valid.IsRFC3339(v) && !valid.IsRFC3339WithoutZone(v)
		return !valid.IsASCII(v)
	}

	return false
}
