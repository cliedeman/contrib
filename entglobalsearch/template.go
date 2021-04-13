// Copyright 2019-present Facebook
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package entglobalsearch

import (
	"entgo.io/contrib/entglobalsearch/internal"
	"entgo.io/ent/entc/gen"
	_ "github.com/go-bindata/go-bindata"
)

var (
	// TODO
	GlobalSearchTemplate = parse("template/globalsearch.tmpl")

	// AllTemplates holds all templates for extending ent to support GraphQL.
	AllTemplates = []*gen.Template{
		GlobalSearchTemplate,
	}
)

//go:generate go run github.com/go-bindata/go-bindata/go-bindata -o=internal/bindata.go -pkg=internal -modtime=1 ./template

func parse(path string) *gen.Template {
	text := string(internal.MustAsset(path))
	return gen.MustParse(gen.NewTemplate(path).
		Funcs(gen.Funcs).
		Parse(text))
}
