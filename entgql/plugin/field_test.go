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
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2/ast"
	"testing"
)

func TestTypeFields(t *testing.T) {
	e := New(&gen.Graph{
		Config: &gen.Config{},
	})
	fields := e.typeFields(&gen.Type{
		ID: &gen.Field{
			Name: "Id",
			Type: &field.TypeInfo{
				Type: field.TypeInt,
			},
		},
		Fields: []*gen.Field{
			{
				Name: "name",
				Type: &field.TypeInfo{
					Type: field.TypeString,
				},
			},
		},
	})
	require.Equal(t, fields, ast.FieldList{
		{
			Name: "id",
			Type: namedType("ID", false),
		},
		{
			Name: "name",
			Type: namedType("String", false),
		},
	})
}

func TestFields(t *testing.T) {
	testCases := []struct {
		name        string
		fieldType   field.Type
		userDefined string
		expected    string
	}{
		{"firstname", field.TypeString, "", "String"},
		{"age", field.TypeInt, "", "Int"},
		{"f", field.TypeFloat64, "", "Float"},
		{"f", field.TypeFloat32, "", "Float"},
		{"status", field.TypeEnum, "", "Status"},
		{"status", field.TypeEnum, "StatusEnum", "StatusEnum"},
		{"status", field.TypeEnum, "StatusEnum", "StatusEnum"},
		{"timestamp", field.TypeTime, "", "Time"},
		{"active", field.TypeBool, "", "Boolean"},
		{"data", field.TypeBytes, "", "TODOBytes"},
		{"json", field.TypeJSON, "", "TODOJSON"},
		{"other", field.TypeOther, "", "Invalid"},
	}
	e := New(&gen.Graph{
		Config: &gen.Config{
			Annotations: map[string]interface{}{
				annotationName: entgql.Annotation{
					GqlScalarMappings: map[string]string{
						"Time": "Time",
					}},
			},
		},
	})
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s(%s)", tc.name, tc.fieldType.ConstName()), func(t *testing.T) {
			f := e.fieldType(&gen.Field{
				Name: tc.name,
				Type: &field.TypeInfo{
					Type: tc.fieldType,
				},
				Nillable: true,
				Annotations: map[string]interface{}{
					annotationName: map[string]interface{}{
						"GqlType": tc.userDefined,
					},
				},
			}, false)
			require.Equal(t, f.String(), tc.expected)
			f = e.fieldType(&gen.Field{
				Name: tc.name,
				Type: &field.TypeInfo{
					Type: tc.fieldType,
				},
				Nillable: false,
				Annotations: map[string]interface{}{
					annotationName: map[string]interface{}{
						"GqlType": tc.userDefined,
					},
				},
			}, false)
			require.Equal(t, f.String(), tc.expected+"!")
		})
	}
}

func TestIdField(t *testing.T) {
	e := New(&gen.Graph{
		Config: &gen.Config{},
	})
	f := e.fieldType(&gen.Field{
		Name: "id",
		Type: &field.TypeInfo{
			Type: field.TypeInt,
		},
	}, true)
	require.Equal(t, f.String(), "ID!")
}
