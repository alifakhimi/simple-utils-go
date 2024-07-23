package scheme

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type Label struct {
	RawKey Key    `json:"key"`
	Scope  string `json:"scope"`
	Keys   []Key  `json:"keys"`
	Parent *Label `json:"-"`
}

func (l *Label) UnmarshalJSON(data []byte) (err error) {
	var lbl *Label
	if lbl, err = ParseLabel(string(data)); err != nil {
		return err
	}
	*l = *lbl
	return nil
}

func (l *Label) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.GetKey())
}

func (l *Label) GetKey() Key {
	if l.RawKey == "" {
		l.RawKey = l.makeKey()
	}

	return l.RawKey
}

func (l *Label) makeKey() Key {
	lKeys := make([]string, len(l.Keys))
	for i, k := range l.Keys {
		lKeys[i] = string(k)
	}

	lkey := fmt.Sprintf("(%s:[%s])", l.Scope, strings.Join(lKeys, keySeparator))

	if l.Parent != nil {
		lkey = fmt.Sprintf("%s%s%v", l.Parent, scopeSeparator, lkey)
	}

	return Key(lkey)
}

func (l *Label) invalidate() *Label {
	l.RawKey = l.makeKey()
	return l
}

func (l Label) String() string {
	return string(l.GetKey())
}

func ParseLabelWithTag(t string) (*Label, error) {
	re := regexp.MustCompile(`(?i).*\((.+):\[(.*)\]\)`)
	sub := re.FindStringSubmatch(t)

	if len(sub) != 3 {
		return nil, ErrLabelParse
	}

	scope := sub[1]
	keys := strings.Split(sub[2], keySeparator)

	keyArr := make([]Key, len(keys))
	for i, k := range keys {
		keyArr[i] = Key(k)
	}

	return (&Label{
		Scope: scope,
		Keys:  keyArr,
	}).invalidate(), nil
}

func ParseRelationLabelWithTag(t string) (from, to *Label, err error) {
	re := regexp.MustCompile(`(?i).*\((.+)` + relSeparator + `(.+)\)$`)
	sub := re.FindStringSubmatch(t)

	if len(sub) != 3 {
		return nil, nil, ErrLabelParse
	}

	from, err = ParseLabelWithTag(sub[1])
	if err != nil {
		return nil, nil, err
	}

	to, err = ParseLabelWithTag(sub[2])
	if err != nil {
		return nil, nil, err
	}

	return from, to, nil
}

func ParseLabel(t string) (l *Label, err error) {
	var (
		currentLabel   *Label
		parentLabel    *Label
		globalRegx     = regexp.MustCompile(`(?iUm)\(.+\)`)
		submatchGlobal = globalRegx.FindAllStringSubmatch(t, -1)
	)

	if len(submatchGlobal) >= 1 {
		currentLabel, err = ParseLabelWithTag(submatchGlobal[len(submatchGlobal)-1][0])
		if err != nil {
			return nil, err
		}
	}

	if len(submatchGlobal) >= 2 {
		parentLabel, err = ParseLabelWithTag(submatchGlobal[len(submatchGlobal)-2][0])
		if err != nil {
			return nil, err
		}

		doc := GetDocumentByLabel(parentLabel)
		if doc != nil {
			parentLabel = &doc.Label
		}
	}

	if currentLabel == nil {
		return nil, ErrLabelParse
	}

	currentLabel.Parent = parentLabel

	return currentLabel.invalidate(), nil
}

// func ParseRelLabel(t string) (from, to *Label, err error) {
// 	var (
// 		relFromLabel   *Label
// 		relToLabel     *Label
// 		relRegx        = regexp.MustCompile(`(?iUm)\(rel:\[.+\]\)$`)
// 		submatchRel    = relRegx.FindAllStringSubmatch(t, -1)
// 	)

// 	if len(submatchRel) >= 1 {
// 		relFromLabel, relToLabel, err = ParseRelationLabelWithTag(submatchRel[len(submatchRel)-1][0])
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	if len(submatchGlobal) >= 1 {
// 		currentLabel, err = ParseLabelWithTag(submatchGlobal[len(submatchGlobal)-1][0])
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	if len(submatchGlobal) >= 2 {
// 		parentLabel, err = ParseLabelWithTag(submatchGlobal[len(submatchGlobal)-2][0])
// 		if err != nil {
// 			return nil, err
// 		}

// 		doc := GetDocumentByLabel(parentLabel)
// 		if doc != nil {
// 			parentLabel = &doc.Label
// 		}
// 	}

// 	if currentLabel == nil {
// 		return nil, ErrLabelParse
// 	}

// 	currentLabel.Parent = parentLabel

// 	return currentLabel.invalidate(), nil
// }

func BuildLabel(scope string, keys ...Key) *Label {
	return (&Label{
		Scope: scope,
		Keys:  keys,
	}).invalidate()
}

func (l *Label) Append(scope string, keys ...Key) *Label {
	label := BuildLabel(scope, keys...)
	label.Parent = l
	return label.invalidate()
}
