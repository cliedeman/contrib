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
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEmpty(t *testing.T) {
	e := New(&gen.Graph{
		Config: &gen.Config{},
	}, nil)
	require.Equal(t, e.print(), ``)
}

func TestModifyConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	modifyConfig(cfg, &gen.Graph{
		Config: &gen.Config{
			Package: "example.com",
		},
	})
	expected := config.DefaultConfig()
	expected.AutoBind = append(expected.AutoBind, "example.com")
	expected.Models["Node"] = config.TypeMapEntry{
		Model: []string{"example.com.Noder"},
	}
	require.Equal(t, cfg, expected)
}

func TestModifyConfig_autoBindPresent(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.AutoBind = append(cfg.AutoBind, "example.com")
	modifyConfig(cfg, &gen.Graph{
		Config: &gen.Config{
			Package: "example.com",
		},
	})
	expected := config.DefaultConfig()
	expected.AutoBind = append(expected.AutoBind, "example.com")
	expected.Models["Node"] = config.TypeMapEntry{
		Model: []string{"example.com.Noder"},
	}
	require.Equal(t, cfg, expected)
}

func TestModifyConfig_noderPresent(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Models["Node"] = config.TypeMapEntry{
		Model: []string{"example.com.CustomNoder"},
	}
	modifyConfig(cfg, &gen.Graph{
		Config: &gen.Config{
			Package: "example.com",
		},
	})
	expected := config.DefaultConfig()
	expected.AutoBind = append(expected.AutoBind, "example.com")
	expected.Models["Node"] = config.TypeMapEntry{
		Model: []string{"example.com.CustomNoder"},
	}
	require.Equal(t, cfg, expected)
}

func TestGetTypes(t *testing.T) {
	t1 := &gen.Type{
		Name: "T1",
		Annotations: map[string]interface{}{
			annotationName: map[string]interface{}{
				"GenType": true,
			},
		},
	}
	t2 := &gen.Type{
		Name: "T1",
	}
	require.Equal(t, []*gen.Type{t1}, getTypes(&gen.Graph{
		Nodes: []*gen.Type{t1, t2},
	}))
}

func TestInjectSourceEarlyEmpty(t *testing.T) {
	e := New(&gen.Graph{
		Config: &gen.Config{},
	}, nil)
	s := e.InjectSourceEarly()
	require.False(t, s.BuiltIn)
	require.Equal(t, s.Input, `scalar Cursor
interface Node {
	id: ID!
}
type PageInfo {
	hasNextPage: Boolean!
	hasPreviousPage: Boolean!
	startCursor: Cursor
	endCursor: Cursor
}
`)
}

func TestInjectSourceEarly(t *testing.T) {
	graph, err := entc.LoadGraph("../internal/todoplugin/ent/schema", &gen.Config{})
	require.NoError(t, err)
	plugin := New(graph, nil)
	s := plugin.InjectSourceEarly()
	require.Equal(t, s.Input, expected)
}

var expected = `scalar Cursor
interface Node {
	id: ID!
}
type PageInfo {
	hasNextPage: Boolean!
	hasPreviousPage: Boolean!
	startCursor: Cursor
	endCursor: Cursor
}
enum Role {
	ADMIN
	USER
	UNKNOWN
}
enum Status {
	IN_PROGRESS
	COMPLETED
}
type Todo implements Node {
	id: ID!
	createdAt: Invalid!
	status: Status!
	priority: Int!
	text: String!
}
type TodoConnection {
	edges: [TodoEdge]
	pageInfo: PageInfo!
	totalCount: Int!
}
type TodoEdge {
	node: Todo
	cursor: Cursor
}
type User implements Node {
	id: ID!
	username: String!
	age: Float!
	amount: Float!
	role: Role!
}
`
