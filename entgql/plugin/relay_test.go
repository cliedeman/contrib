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
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRelayBuiltins(t *testing.T) {
	e := New(&gen.Graph{
		Config: &gen.Config{},
	})
	e.relayBuiltins()
	require.Equal(t, e.print(), `scalar Cursor
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

func TestCreateRelayConnection(t *testing.T) {
	require.True(t, createRelayConnection(&gen.Type{
		Annotations: map[string]interface{}{
			annotationName: map[string]interface{}{
				"RelayConnection": true,
			},
		},
	}))
	require.False(t, createRelayConnection(&gen.Type{}))
}

func TestRelayConnection(t *testing.T) {
	e := New(&gen.Graph{
		Config: &gen.Config{},
	})
	e.relayConnection(&gen.Type{
		Name: "Todo",
	})
	require.Equal(t, e.print(), `type TodoConnection {
	edges: [TodoEdge]
	pageInfo: PageInfo!
	totalCount: Int!
}
type TodoEdge {
	node: Todo
	cursor: Cursor
}
`)
}
