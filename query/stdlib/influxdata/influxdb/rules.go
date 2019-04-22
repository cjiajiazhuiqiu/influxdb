package influxdb

import (
	"github.com/influxdata/flux"
	"github.com/influxdata/flux/execute"
	"github.com/influxdata/flux/plan"
	"github.com/influxdata/flux/semantic"
	"github.com/influxdata/flux/stdlib/universe"
)

// func init() {
// 	plan.RegisterPhysicalRules(
// 		PushDownRangeRule{},
// 		PushDownFilterRule{},
// 		PushDownReadTagKeysRule{},
// 		PushDownReadTagValuesRule{},
// 	)
// }

// PushDownRangeRule pushes down a range filter to storage
type PushDownRangeRule struct{}

func (rule PushDownRangeRule) Name() string {
	return "PushDownRangeRule"
}

// Pattern matches 'from |> range'
func (rule PushDownRangeRule) Pattern() plan.Pattern {
	return plan.Pat(universe.RangeKind, plan.Pat(FromKind))
}

// Rewrite converts 'from |> range' into 'ReadRange'
func (rule PushDownRangeRule) Rewrite(node plan.Node) (plan.Node, bool, error) {
	fromNode := node.Predecessors()[0]
	fromSpec := fromNode.ProcedureSpec().(*FromProcedureSpec)

	rangeSpec := node.ProcedureSpec().(*universe.RangeProcedureSpec)
	return plan.CreatePhysicalNode("ReadRange", &ReadRangePhysSpec{
		Bucket:   fromSpec.Bucket,
		BucketID: fromSpec.BucketID,
		Bounds:   rangeSpec.Bounds,
	}), true, nil
}

// PushDownFilterRule is a rule that pushes filters into from procedures to be evaluated in the storage layer.
// This rule is likely to be replaced by a more generic rule when we have a better
// framework for pushing filters, etc into sources.
type PushDownFilterRule struct{}

func (PushDownFilterRule) Name() string {
	return "PushDownFilterRule"
}

func (PushDownFilterRule) Pattern() plan.Pattern {
	return plan.Pat(universe.FilterKind, plan.Pat(ReadRangePhysKind))
}

func (PushDownFilterRule) Rewrite(pn plan.Node) (plan.Node, bool, error) {
	filterSpec := pn.ProcedureSpec().(*universe.FilterProcedureSpec)
	fromNode := pn.Predecessors()[0]
	fromSpec := fromNode.ProcedureSpec().(*ReadRangePhysSpec)

	bodyExpr, ok := filterSpec.Fn.Block.Body.(semantic.Expression)
	if !ok {
		return pn, false, nil
	}

	if len(filterSpec.Fn.Block.Parameters.List) != 1 {
		// I would expect that type checking would catch this, but just to be safe...
		return pn, false, nil
	}

	paramName := filterSpec.Fn.Block.Parameters.List[0].Key.Name

	pushable, notPushable, err := semantic.PartitionPredicates(bodyExpr, func(e semantic.Expression) (bool, error) {
		return isPushableExpr(paramName, e)
	})
	if err != nil {
		return nil, false, err
	}

	if pushable == nil {
		// Nothing could be pushed down, no rewrite can happen
		return pn, false, nil
	}

	newFromSpec := fromSpec.Copy().(*ReadRangePhysSpec)
	if newFromSpec.FilterSet {
		newBody := semantic.ExprsToConjunction(newFromSpec.Filter.Block.Body.(semantic.Expression), pushable)
		newFromSpec.Filter.Block.Body = newBody
	} else {
		newFromSpec.FilterSet = true
		newFromSpec.Filter = filterSpec.Fn.Copy().(*semantic.FunctionExpression)
		newFromSpec.Filter.Block.Body = pushable
	}

	if notPushable == nil {
		// All predicates could be pushed down, so eliminate the filter
		mergedNode, err := plan.MergeToPhysicalNode(pn, fromNode, newFromSpec)
		if err != nil {
			return nil, false, err
		}
		return mergedNode, true, nil
	}

	err = fromNode.ReplaceSpec(newFromSpec)
	if err != nil {
		return nil, false, err
	}

	newFilterSpec := filterSpec.Copy().(*universe.FilterProcedureSpec)
	newFilterSpec.Fn.Block.Body = notPushable
	if err := pn.ReplaceSpec(newFilterSpec); err != nil {
		return nil, false, err
	}

	return pn, true, nil
}

// PushDownReadTagKeysRule matches 'ReadRange |> keys() |> keep() |> distinct()'.
// The 'from()' must have already been merged with 'range' and, optionally,
// may have been merged with 'filter'.
// If any other properties have been set on the from procedure,
// this rule will not rewrite anything.
type PushDownReadTagKeysRule struct{}

func (rule PushDownReadTagKeysRule) Name() string {
	return "PushDownReadTagKeysRule"
}

func (rule PushDownReadTagKeysRule) Pattern() plan.Pattern {
	return plan.Pat(universe.DistinctKind,
		plan.Pat(universe.SchemaMutationKind,
			plan.Pat(universe.KeysKind,
				plan.Pat(ReadRangePhysKind))))
}

func (rule PushDownReadTagKeysRule) Rewrite(pn plan.Node) (plan.Node, bool, error) {
	// Retrieve the nodes and specs for all of the predecessors.
	distinctSpec := pn.ProcedureSpec().(*universe.DistinctProcedureSpec)
	keepNode := pn.Predecessors()[0]
	keepSpec := keepNode.ProcedureSpec().(*universe.SchemaMutationProcedureSpec)
	keysNode := keepNode.Predecessors()[0]
	keysSpec := keysNode.ProcedureSpec().(*universe.KeysProcedureSpec)
	fromNode := keysNode.Predecessors()[0]
	fromSpec := fromNode.ProcedureSpec().(*ReadRangePhysSpec)

	// A filter spec would have already been merged into the
	// from spec if it existed so we will take that one when
	// constructing our own replacement. We do not care about it
	// at the moment though which is why it is not in the pattern.

	// The schema mutator needs to correspond to a keep call
	// on the column specified by the keys procedure.
	if len(keepSpec.Mutations) != 1 {
		return nil, false, nil
	} else if m, ok := keepSpec.Mutations[0].(*universe.KeepOpSpec); !ok {
		return nil, false, nil
	} else if m.Predicate != nil || len(m.Columns) != 1 {
		// We have a keep mutator, but it uses a function or
		// it retains more than one column so it does not match
		// what we want.
		return nil, false, nil
	} else if m.Columns[0] != keysSpec.Column {
		// We are not keeping the value column so this optimization
		// will not work.
		return nil, false, nil
	}

	// The distinct spec should keep only the value column.
	if distinctSpec.Column != keysSpec.Column {
		return nil, false, nil
	}

	// We have passed all of the necessary prerequisites
	// so construct the procedure spec.
	return plan.CreatePhysicalNode("ReadTagKeys", &ReadTagKeysPhysSpec{
		ReadRangePhysSpec: *fromSpec.Copy().(*ReadRangePhysSpec),
		ValueColumnName:   keysSpec.Column,
	}), true, nil
}

// PushDownReadTagValuesRule matches 'ReadRange |> group(columns: [tag]) |> distinct(column: tag) |> keep()'.
// The 'from()' must have already been merged with 'range' and, optionally,
// may have been merged with 'filter'.
// If any other properties have been set on the from procedure,
// this rule will not rewrite anything.
type PushDownReadTagValuesRule struct{}

func (rule PushDownReadTagValuesRule) Name() string {
	return "PushDownReadTagValuesRule"
}

func (rule PushDownReadTagValuesRule) Pattern() plan.Pattern {
	return plan.Pat(universe.DistinctKind,
		plan.Pat(universe.GroupKind,
			plan.Pat(universe.SchemaMutationKind,
				plan.Pat(ReadRangePhysKind))))
}

func (rule PushDownReadTagValuesRule) Rewrite(pn plan.Node) (plan.Node, bool, error) {
	// Retrieve the nodes and specs for all of the predecessors.
	distinctNode := pn
	distinctSpec := distinctNode.ProcedureSpec().(*universe.DistinctProcedureSpec)
	groupNode := distinctNode.Predecessors()[0]
	groupSpec := groupNode.ProcedureSpec().(*universe.GroupProcedureSpec)
	schemaNode := groupNode.Predecessors()[0]
	schemaSpec := schemaNode.ProcedureSpec().(*universe.SchemaMutationProcedureSpec)
	fromNode := schemaNode.Predecessors()[0]
	fromSpec := fromNode.ProcedureSpec().(*ReadRangePhysSpec)

	// A filter spec would have already been merged into the
	// from spec if it existed so we will take that one when
	// constructing our own replacement. We do not care about it
	// at the moment though which is why it is not in the pattern.

	// All of the values need to be grouped into the same table.
	if groupSpec.GroupMode != flux.GroupModeBy {
		return nil, false, nil
	} else if len(groupSpec.GroupKeys) > 0 {
		return nil, false, nil
	}

	// The column that distinct is for will be the tag key.
	tagKey := distinctSpec.Column
	if !isValidTagKeyForTagValues(tagKey) {
		return nil, false, nil
	}

	// If the schema mutators do not drop the tag key column,
	// then we can perform the optimization.
	hasTagKey, hasValue := true, true
	for _, m := range schemaSpec.Mutations {
		if !hasTagKey {
			// We dropped the tag key and cannot
			// get it back so we cannot perform this
			// optimization.
			return nil, false, nil
		}

		switch m := m.(type) {
		case *universe.KeepOpSpec:
			if m.Predicate != nil {
				return nil, false, nil
			}

			var foundTagKey, foundValue bool
			for _, name := range m.Columns {
				if name == tagKey {
					foundTagKey = true
				} else if name == execute.DefaultValueColLabel {
					foundValue = true
				}
			}

			hasTagKey = foundTagKey
			if hasValue && !foundValue {
				// We did not keep the value column.
				hasValue = false
			}
		case *universe.DropOpSpec:
			if m.Predicate != nil {
				return nil, false, nil
			}

			for _, name := range m.Columns {
				if name == tagKey {
					hasTagKey = false
				} else if name == execute.DefaultValueColLabel {
					hasValue = false
				}
			}
		default:
			// We do not optimize if a different mutator is used.
			return nil, false, nil
		}
	}

	// We have passed all of the necessary prerequisites
	// so construct the procedure spec.
	return plan.CreatePhysicalNode("ReadTagValues", &ReadTagValuesPhysSpec{
		ReadRangePhysSpec: *fromSpec.Copy().(*ReadRangePhysSpec),
		TagKey:            tagKey,
	}), true, nil
}

var invalidTagKeysForTagValues = []string{
	execute.DefaultTimeColLabel,
	execute.DefaultValueColLabel,
	execute.DefaultStartColLabel,
	execute.DefaultStopColLabel,
}

// isValidTagKeyForTagValues returns true if the given key can
// be used in a tag values call.
func isValidTagKeyForTagValues(key string) bool {
	for _, k := range invalidTagKeysForTagValues {
		if k == key {
			return false
		}
	}
	return true
}
