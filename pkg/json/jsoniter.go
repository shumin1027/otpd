// Copyright 2017 Bo-Yi Wu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

//go:build jsoniter
// +build jsoniter

package json

import (
	"github.com/iancoleman/strcase"
	jsoniter "github.com/json-iterator/go"
	"strings"
	"unicode"
	"unsafe"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
	// Marshal is exported by gin/json package.
	Marshal = json.Marshal
	// Unmarshal is exported by gin/json package.
	Unmarshal = json.Unmarshal
	// MarshalIndent is exported by gin/json package.
	MarshalIndent = json.MarshalIndent
	// NewDecoder is exported by gin/json package.
	NewDecoder = json.NewDecoder
	// NewEncoder is exported by gin/json package.
	NewEncoder = json.NewEncoder
)

func init() {
	// 配置json风格
	jsoniter.RegisterExtension(NewJSONStyleExtension(true, SnakeCase))

	//jsoniter.RegisterTypeEncoder("map[string]interface {}", &MapNamingStrategyEncoder{SnakeCase})
	//jsoniter.RegisterTypeEncoder("map[string]string", &MapNamingStrategyEncoder{SnakeCase})
}

/*
命名规则:

CamelCase: 	"persionId"

PascalCase: "PersonId"

SnakeCase:  "person_id"

KebabCase: 	"KebabCase"
*/
type NamingStrategy string

const (
	CamelCase  NamingStrategy = "CamelCase"  // persionId
	PascalCase NamingStrategy = "PascalCase" // PersonId
	SnakeCase  NamingStrategy = "SnakeCase"  // person_id
	KebabCase  NamingStrategy = "KebabCase"  // person-id
)

func NewJSONStyleExtension(override bool, namingStrategy NamingStrategy) *JSONStyleExtension {
	ext := new(JSONStyleExtension)
	ext.Override = override
	ext.NamingStrategy = namingStrategy
	return ext
}

/*
json-iterator JSONStyleExtension
*/
type JSONStyleExtension struct {
	jsoniter.DummyExtension
	NamingStrategy NamingStrategy // 命名规则
	Override       bool           // 是否覆盖已明确指定命名的 json key
}

func (extension *JSONStyleExtension) UpdateStructDescriptor(structDescriptor *jsoniter.StructDescriptor) {

	fields := structDescriptor.Fields

	for _, binding := range fields {
		if unicode.IsLower(rune(binding.Field.Name()[0])) || binding.Field.Name()[0] == '_' {
			continue
		}
		tag, hastag := binding.Field.Tag().Lookup("json")
		if hastag {
			tagParts := strings.Split(tag, ",")
			if tagParts[0] == "-" {
				continue // hidden field
			}
			if tagParts[0] != "" && !extension.Override {
				continue // field explicitly named
			}
		}
		binding.ToNames = []string{extension.translate(binding.Field.Name())}
		binding.FromNames = []string{extension.translate(binding.Field.Name())}
	}

}

func (extension *JSONStyleExtension) translate(str string) string {
	namingStrategy := extension.NamingStrategy
	switch namingStrategy {
	case PascalCase:
		{
			return strcase.ToCamel(str)
		}
	case CamelCase:
		{
			return strcase.ToLowerCamel(str)
		}
	case SnakeCase:
		{
			return strcase.ToSnake(str)
		}
	case KebabCase:
		{
			return strcase.ToKebab(str)
		}
	}
	return str
}

type MapNamingStrategyEncoder struct {
	NamingStrategy NamingStrategy
}

func (codec *MapNamingStrategyEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return len(*((*map[string]interface{})(ptr))) == 0
}

func (codec *MapNamingStrategyEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	namingStrategy := codec.NamingStrategy
	m := *((*map[string]interface{})(ptr))
	for k, v := range m {
		switch namingStrategy {
		case PascalCase:
			{
				k = strcase.ToCamel(k)
				stream.WriteObjectField(k)
				stream.WriteVal(v)
			}

		case CamelCase:
			{
				k = strcase.ToLowerCamel(k)
				stream.WriteObjectField(k)
				stream.WriteVal(v)
			}
		case SnakeCase:
			{
				k = strcase.ToSnake(k)
				stream.WriteObjectField(k)
				stream.WriteVal(v)
			}
		case KebabCase:
			{
				k = strcase.ToKebab(k)
				stream.WriteObjectField(k)
				stream.WriteVal(v)
			}
		}
	}
}
