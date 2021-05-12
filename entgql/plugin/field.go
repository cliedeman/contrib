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
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
	"github.com/vektah/gqlparser/v2/ast"
	"strings"
)

func (e *entgqlgen) typeFields(t *gen.Type) ast.FieldList {
	var fields ast.FieldList
	if t.ID != nil {
		fields = append(fields, &ast.FieldDefinition{
			Name:       camel(t.ID.Name),
			Type:       e.fieldType(t.ID, true),
			Directives: e.fieldDirectives(t.ID),
		})
	}
	for _, f := range t.Fields {
		fields = append(fields, &ast.FieldDefinition{
			Name:       camel(f.Name),
			Type:       e.fieldType(f, false),
			Directives: e.fieldDirectives(f),
		})
	}
	return fields
}

func (e *entgqlgen) fieldDirectives(*gen.Field) ast.DirectiveList {
	// TODO
	return nil
}

func namedType(name string, nillable bool) *ast.Type {
	if !nillable {
		return ast.NonNullNamedType(name, nil)
	}
	return ast.NamedType(name, nil)
}

func (e *entgqlgen) fieldType(f *gen.Field, idField bool) *ast.Type {
	// TODO: handle array
	userDefinedType := e.fieldUserDefinedType(f)
	if userDefinedType != nil {
		return userDefinedType
	}
	nillable := f.Nillable
	typ := f.Type.Type
	typeName := strings.TrimPrefix(typ.ConstName(), "Type")

	switch {
	case e.scalarMappings[typeName] != "":
		return namedType(e.scalarMappings[typeName], nillable)
	case idField:
		// Id cannot be null for node interface
		return namedType("ID", false)
	case f.IsEnum():
		// Guess enum type
		return namedType(strings.Title(f.Name), nillable)
	case typ.Float():
		return namedType("Float", nillable)
	case typ.Integer():
		return namedType("Int", nillable)
	case typ == field.TypeString:
		return namedType("String", nillable)
	case typ == field.TypeBool:
		return namedType("Boolean", nillable)
	case typ == field.TypeBytes:
		return namedType("TODOBytes", nillable)
	case typ == field.TypeJSON:
		return namedType("TODOJSON", nillable)
	default:
		// TODO: error
		return namedType("Invalid", nillable)
	}
}

func (e *entgqlgen) fieldUserDefinedType(f *gen.Field) *ast.Type {
	ann := entgqlAnnotate(f.Annotations)
	if ann != nil && ann.GqlType != "" {
		return namedType(ann.GqlType, f.Nillable)
	}
	return nil
}
