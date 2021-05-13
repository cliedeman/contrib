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
	"encoding/json"
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

type entgqlgen struct {
	genTypes       []*gen.Type
	scalarMappings map[string]string
	schema         *ast.Schema
	hooks          []SchemaHook
}

// SchemaHook hook to modify schema before printing
type SchemaHook func(schema *ast.Schema)

func (e *entgqlgen) InjectSourceEarly() *ast.Source {
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
	return &ast.Source{
		Name:    "entgqlgen.graphql",
		Input:   e.print(),
		BuiltIn: false,
	}
}

func (e *entgqlgen) print() string {
	sb := &strings.Builder{}
	printer := formatter.NewFormatter(sb)
	printer.FormatSchema(e.schema)
	return sb.String()
}

func getTypes(graph *gen.Graph) []*gen.Type {
	var types []*gen.Type
	for _, n := range graph.Nodes {
		ann := entgqlAnnotate(n.Annotations)
		if ann != nil {
			if ann.GenType {
				types = append(types, n)
			}
		}
	}
	return types
}

func New(graph *gen.Graph, hooks []SchemaHook) *entgqlgen {
	types := getTypes(graph)
	var scalarMappings map[string]string
	if graph.Annotations != nil {
		globalAnn := graph.Annotations[annotationName]
		// TODO: cleanup assertions
		if globalAnn != nil {
			if globalAnn.(entgql.Annotation).GqlScalarMappings != nil {
				scalarMappings = globalAnn.(entgql.Annotation).GqlScalarMappings
			}
		}
	}
	return &entgqlgen{
		genTypes:       types,
		scalarMappings: scalarMappings,
		hooks:          hooks,
		schema: &ast.Schema{
			Types:         map[string]*ast.Definition{},
			Directives:    map[string]*ast.DirectiveDefinition{},
			PossibleTypes: map[string][]*ast.Definition{},
			Implements:    map[string][]*ast.Definition{},
		},
	}
}

func (e *entgqlgen) Name() string {
	return "entgqlgen"
}

func Generate(cfg *config.Config, graph *gen.Graph, hooks ...SchemaHook) error {
	modifyConfig(cfg, graph)
	return api.Generate(cfg,
		api.AddPlugin(New(graph, hooks)),
	)
}

func modifyConfig(cfg *config.Config, graph *gen.Graph) {
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

func entgqlAnnotate(annotation map[string]interface{}) *entgql.Annotation {
	annotate := &entgql.Annotation{}
	if annotation == nil || annotation[annotate.Name()] == nil {
		return nil
	}
	if buf, err := json.Marshal(annotation[annotate.Name()]); err == nil {
		_ = json.Unmarshal(buf, &annotate)
	}
	return annotate
}

var _ plugin.EarlySourceInjector = &entgqlgen{}
var _ plugin.Plugin = &entgqlgen{}
