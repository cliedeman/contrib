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
//
// Code generated by entc, DO NOT EDIT.

package private

const (
	// Label holds the string label denoting the private type in the database.
	Label = "private"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// EdgeParent holds the string denoting the parent edge name in mutations.
	EdgeParent = "parent"
	// EdgeChildren holds the string denoting the children edge name in mutations.
	EdgeChildren = "children"
	// Table holds the table name of the private in the database.
	Table = "privates"
	// ParentTable is the table the holds the parent relation/edge.
	ParentTable = "privates"
	// ParentColumn is the table column denoting the parent relation/edge.
	ParentColumn = "private_children"
	// ChildrenTable is the table the holds the children relation/edge.
	ChildrenTable = "privates"
	// ChildrenColumn is the table column denoting the children relation/edge.
	ChildrenColumn = "private_children"
)

// Columns holds all SQL columns for private fields.
var Columns = []string{
	FieldID,
	FieldName,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "privates"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"private_children",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}
