// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package util

import "reflect"

// ReflectValue returns reflect value underlying given value, unwrapping pointer and slice
func ReflectValue(v any) reflect.Value {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Slice {
		rv = rv.Elem()
	}
	return rv
}

// ReflectType returns reflect type underlying given value, unwrapping pointer and slice
func ReflectType(v any) reflect.Type {
	var rt reflect.Type
	if rvt, ok := v.(reflect.Type); ok {
		rt = rvt
	} else {
		rv, ok := v.(reflect.Value)
		if !ok {
			rv = reflect.ValueOf(v)
		}
		rt = rv.Type()
	}
	for rt.Kind() == reflect.Ptr || rt.Kind() == reflect.Slice {
		rt = rt.Elem()
	}
	return rt
}
