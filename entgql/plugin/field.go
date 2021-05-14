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
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
	"fmt"
	"github.com/vektah/gqlparser/v2/ast"
	"strings"
)

func includeField(f *gen.Field) bool {
	ann := entgql.EntgqlAnnotate(f.Annotations)
	if ann != nil && ann.Skip {
		return false
	}
	return true
}

func (e *entgqlgen) typeFields(t *gen.Type) (ast.FieldList, error) {
	var fields ast.FieldList
	if t.ID != nil && includeField(t.ID) {
		ft, err := e.fieldType(t.ID, true)
		if err != nil {
			return nil, fmt.Errorf("field(%s): %w", t.ID.Name, err)
		}
		fields = append(fields, &ast.FieldDefinition{
			Name:       camel(t.ID.Name),
			Type:       ft,
			Directives: e.fieldDirectives(t.ID),
		})
	}
	for _, f := range t.Fields {
		if !includeField(f) {
			continue
		}
		ft, err := e.fieldType(f, false)
		if err != nil {
			return nil, fmt.Errorf("field(%s): %w", t.ID.Name, err)
		}
		fields = append(fields, &ast.FieldDefinition{
			Name:       camel(f.Name),
			Type:       ft,
			Directives: e.fieldDirectives(f),
		})
	}
	return fields, nil
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

func (e *entgqlgen) fieldType(f *gen.Field, idField bool) (*ast.Type, error) {
	// TODO: handle array
	userDefinedType := e.fieldUserDefinedType(f)
	if userDefinedType != nil {
		return userDefinedType, nil
	}
	nillable := f.Nillable
	typ := f.Type.Type
	typeName := strings.TrimPrefix(typ.ConstName(), "Type")

	switch {
	case e.scalarMappings[typeName] != "":
		return namedType(e.scalarMappings[typeName], nillable), nil
	case idField:
		// Id cannot be null for node interface
		return namedType("ID", false), nil
	case f.IsEnum():
		// Guess enum type
		return namedType(strings.Title(f.Name), nillable), nil
	case typ.Float():
		return namedType("Float", nillable), nil
	case typ.Integer():
		return namedType("Int", nillable), nil
	case typ == field.TypeString:
		return namedType("String", nillable), nil
	case typ == field.TypeBool:
		return namedType("Boolean", nillable), nil
	case typ == field.TypeBytes:
		return nil, fmt.Errorf("bytes type not implemented")
	case typ == field.TypeJSON:
		return nil, fmt.Errorf("json type not implemented")
	case typ == field.TypeOther:
		return nil, fmt.Errorf("other type must have typed defined")
	default:
		return nil, fmt.Errorf("unexpected type: %s", typ.String())
	}
}

func (e *entgqlgen) fieldUserDefinedType(f *gen.Field) *ast.Type {
	ann := entgql.EntgqlAnnotate(f.Annotations)
	if ann != nil && ann.GqlType != "" {
		return namedType(ann.GqlType, f.Nillable)
	}
	return nil
}
