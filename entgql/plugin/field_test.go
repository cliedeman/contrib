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

func TestIncludeField(t *testing.T) {
	require.True(t, includeField(&gen.Field{}))
	require.False(t, includeField(&gen.Field{
		Annotations: map[string]interface{}{
			annotationName: map[string]interface{}{
				"Skip": true,
			},
		},
	}))
	require.True(t, includeField(&gen.Field{
		Annotations: map[string]interface{}{
			annotationName: map[string]interface{}{
				"Skip": false,
			},
		},
	}))
}

func TestTypeFields(t *testing.T) {
	e := New(&gen.Graph{
		Config: &gen.Config{},
	}, nil)
	fields, err := e.typeFields(&gen.Type{
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
	require.NoError(t, err)
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
		name         string
		fieldType    field.Type
		userDefined  string
		expectedType string
		err          error
	}{
		{"firstname", field.TypeString, "", "String", nil},
		{"age", field.TypeInt, "", "Int", nil},
		{"f", field.TypeFloat64, "", "Float", nil},
		{"f", field.TypeFloat32, "", "Float", nil},
		{"status", field.TypeEnum, "", "Status", nil},
		{"status", field.TypeEnum, "StatusEnum", "StatusEnum", nil},
		{"status", field.TypeEnum, "StatusEnum", "StatusEnum", nil},
		{"timestamp", field.TypeTime, "", "Time", nil},
		{"active", field.TypeBool, "", "Boolean", nil},
		{"data", field.TypeBytes, "", "", fmt.Errorf("bytes type not implemented")},
		{"json", field.TypeJSON, "", "", fmt.Errorf("json type not implemented")},
		{"other", field.TypeOther, "", "Invalid", fmt.Errorf("other type must have typed defined")},
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
	}, nil)
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s(%s)", tc.name, tc.fieldType.ConstName()), func(t *testing.T) {
			f, err := e.fieldType(&gen.Field{
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
			require.Equal(t, err, tc.err)
			if tc.err == nil {
				require.Equal(t, f.String(), tc.expectedType)
			}
			f, err = e.fieldType(&gen.Field{
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
			require.Equal(t, err, tc.err)
			if tc.err == nil {
				require.Equal(t, f.String(), tc.expectedType+"!")
			}
		})
	}
}

func TestIdField(t *testing.T) {
	e := New(&gen.Graph{
		Config: &gen.Config{},
	}, nil)
	f, err := e.fieldType(&gen.Field{
		Name: "id",
		Type: &field.TypeInfo{
			Type: field.TypeInt,
		},
	}, true)
	require.NoError(t, err)
	require.Equal(t, f.String(), "ID!")
}
