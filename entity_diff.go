package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/google/go-cmp/cmp"
)

type DiffReporter struct {
	path cmp.Path
	modelName string
	diffs []string
}

func (r *DiffReporter) PushStep(ps cmp.PathStep) {
	r.path = append(r.path, ps)
}

func (r *DiffReporter) Report(rs cmp.Result) {
	if !rs.Equal() {
		vx, vy := r.path.Last().Values()

		if !vx.IsValid() && vy.IsValid() && vy.Kind() == reflect.Struct {
			r.reportAddedStruct(vy.Interface())
			return
		}

		oldVal := r.formatValue(vx)
		newVal := r.formatValue(vy)

		customPath := r.formatPath()

		r.diffs = append(r.diffs, fmt.Sprintf("%#v:\n\t-: %+v\n\t+: %v\n", customPath, oldVal, newVal))
	}
}

func (r *DiffReporter) reportAddedStruct(addedStruct any) {
	addedType := reflect.TypeOf(addedStruct)
	zeroValue := reflect.New(addedType).Elem().Interface()

	tempReporter := &DiffReporter{
		path: r.path,
		modelName: r.modelName,
		diffs: []string{},
	}

	cmp.Equal(zeroValue, addedStruct, cmp.Reporter(tempReporter))

	for _, diff := range tempReporter.diffs {
		r.diffs = append(r.diffs, "+ "+diff)
	}
}

func (r *DiffReporter) formatPath() string {
	parts := []string{r.modelName}
	for _, step := range r.path[1:] {
		parts = append(parts, fmt.Sprintf("%v", step))
	}
	return strings.Join(parts, "")
}

func (r *DiffReporter) formatValue(v reflect.Value) string {
	if !v.IsValid() {
		return "(new)"
	}

	valStr := fmt.Sprintf("%#v", v.Interface())
	return strings.TrimPrefix(valStr, "main.")
}

func (r *DiffReporter) PopStep() {
	r.path = r.path[:len(r.path)-1]
}

func (r *DiffReporter) String() string {
	return strings.Join(r.diffs, "\n")
}


func EntityDiff(entity, updatedEntity FoundationalModel) string {
	var r = DiffReporter{
		modelName: entity.Name,
		diffs: []string{},
	}
	cmp.Equal(entity, updatedEntity, cmp.Reporter(&r))

	return r.String()
}
