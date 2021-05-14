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
	"fmt"
	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
	"strings"
)

var (
	camel          = gen.Funcs["camel"].(func(string) string)
	annotationName = entgql.Annotation{}.Name()
)

type EntGqlGenOption func(*Entgqlgen)

func WithSchemaHooks(hooks ...SchemaHook) EntGqlGenOption {
	return func(e *Entgqlgen) {
		e.hooks = hooks
	}
}

func WithDebug() EntGqlGenOption {
	return func(e *Entgqlgen) {
		e.debug = true
	}
}

type Entgqlgen struct {
	debug          bool
	genTypes       []*gen.Type
	scalarMappings map[string]string
	schema         *ast.Schema
	hooks          []SchemaHook
}

// SchemaHook hook to modify schema before printing
type SchemaHook func(schema *ast.Schema)

func (e *Entgqlgen) InjectSourceEarly() *ast.Source {
	e.builtIns()
	e.scalars()
	e.relayBuiltins()
	e.enums()
	err := e.types()
	if err != nil {
		panic(err)
	}
	for _, h := range e.hooks {
		h(e.schema)
	}
	input := e.print()
	if e.debug {
		fmt.Printf("Generated Graphql:\n%s", input)
	}
	return &ast.Source{
		Name:    "Entgqlgen.graphql",
		Input:   input,
		BuiltIn: false,
	}
}

func (e *Entgqlgen) print() string {
	sb := &strings.Builder{}
	printer := formatter.NewFormatter(sb)
	printer.FormatSchema(e.schema)
	return sb.String()
}

func getTypes(graph *gen.Graph) []*gen.Type {
	var types []*gen.Type
	for _, n := range graph.Nodes {
		ann := entgql.EntgqlAnnotate(n.Annotations)
		if ann != nil {
			if ann.GenType {
				types = append(types, n)
			}
		}
	}
	return types
}

func New(graph *gen.Graph, opts ...EntGqlGenOption) *Entgqlgen {
	types := getTypes(graph)
	// Include default mapping for time
	scalarMappings := map[string]string{
		"Time": "Time",
	}
	if graph.Annotations != nil {
		globalAnn := graph.Annotations[annotationName]
		// TODO: cleanup assertions
		if globalAnn != nil {
			if globalAnn.(entgql.Annotation).GqlScalarMappings != nil {
				scalarMappings = globalAnn.(entgql.Annotation).GqlScalarMappings
			}
		}
	}
	e := &Entgqlgen{
		genTypes:       types,
		scalarMappings: scalarMappings,
		schema: &ast.Schema{
			Types:         map[string]*ast.Definition{},
			Directives:    map[string]*ast.DirectiveDefinition{},
			PossibleTypes: map[string][]*ast.Definition{},
			Implements:    map[string][]*ast.Definition{},
		},
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func (e *Entgqlgen) Name() string {
	return "Entgqlgen"
}

func Generate(cfg *config.Config, graph *gen.Graph, opts ...EntGqlGenOption) error {
	MutateConfig(cfg, graph)
	return api.Generate(cfg,
		api.AddPlugin(New(graph, opts...)),
	)
}

//MutateConfig We do not use plugin.MutateConfig because it is called too late in the initialization
func MutateConfig(cfg *config.Config, graph *gen.Graph) {
	autobindPresent := false
	for _, ab := range cfg.AutoBind {
		if ab == graph.Package {
			autobindPresent = true
		}
	}
	if !autobindPresent {
		cfg.AutoBind = append(cfg.AutoBind, graph.Package)
	}
	if !cfg.Models.Exists("Node") {
		cfg.Models.Add("Node", fmt.Sprintf("%s.Noder", graph.Package))
	}
}

var _ plugin.EarlySourceInjector = &Entgqlgen{}
var _ plugin.Plugin = &Entgqlgen{}
