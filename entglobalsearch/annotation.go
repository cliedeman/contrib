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

// Annotation annotates fields and edges with metadata for templates.
type Annotation struct {
	Filter string
}

// Name implements ent.Annotation interface.
func (Annotation) Name() string {
	return "EntGlobalSearch"
}

// OrderField returns an order field annotation.
func Filter(filter string) Annotation {
	return Annotation{Filter: filter}
}
