package simutils

import "github.com/gosimple/slug"

type Slug string

func (s Slug) IsValid() bool {
	return slug.IsSlug(string(s))
}

func MakeSlug(s string) Slug {
	return Slug(slug.Make(s))
}

func (s Slug) ToString() string {
	return string(s)
}
