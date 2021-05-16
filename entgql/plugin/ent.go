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

package plugin

import (
	"entgo.io/contrib/entgql"
	"fmt"
	"github.com/vektah/gqlparser/v2/ast"
)

func (e *Entgqlgen) entBuiltins() {
	e.insertDefinitions([]*ast.Definition{
		{
			Name: "OrderDirection",
			Kind: ast.Enum,
			EnumValues: ast.EnumValueList{
				&ast.EnumValueDefinition{
					Name: "ASC",
				},
				&ast.EnumValueDefinition{
					Name: "DESC",
				},
			},
		},
	})
}

func (e *Entgqlgen) entOrderBy() {
	for _, obj := range e.genTypes {
		ann := entgql.EntgqlAnnotate(obj.Annotations)
		if (ann == nil || !ann.Skip) && hasOrderBy(obj) {
			var enumValues ast.EnumValueList
			for _, f := range obj.Fields {
				fAnn := entgql.EntgqlAnnotate(f.Annotations)
				if fAnn != nil && fAnn.OrderField != "" {
					enumValues = append(enumValues, &ast.EnumValueDefinition{
						Name: fAnn.OrderField,
					})
				}
			}
			e.insertDefinitions([]*ast.Definition{
				{
					Name:       fmt.Sprintf("%sOrderField", obj.Name),
					Kind:       ast.Enum,
					EnumValues: enumValues,
				},
				{
					Name: fmt.Sprintf("%sOrder", obj.Name),
					Kind: ast.InputObject,
					Fields: ast.FieldList{
						{
							Name: "direction",
							Type: ast.NonNullNamedType("OrderDirection", nil),
						},
						{
							Name: "field",
							Type: ast.NonNullNamedType(fmt.Sprintf("%sOrderField", obj.Name), nil),
						},
					},
				},
			})
		}
	}
}
