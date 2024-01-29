// Copyright 2024 The NonTechno Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package later

import (
	"reflect"
	"sync/atomic"
)

func Link(fptr interface{}, linkage string) {
	fn := reflect.ValueOf(fptr).Elem()

	universal := func(args []reflect.Value) []reflect.Value {

		current := atomic.AddInt32(&entryCounter, 1)
		defer atomic.AddInt32(&entryCounter, -1)

		if current > 12 {
			onError("too many hops (%v) to resolve linkage (%s)", current, linkage)
			return []reflect.Value{}
		}

		if target := resolve(fn, linkage); target != nil {
			return reflect.ValueOf(target).Call(args) // this can be a chain call
		} else {
			onError("failed to resolve linkage (%s)", linkage)
			return []reflect.Value{}
		}
	}

	v := reflect.MakeFunc(fn.Type(), universal)
	fn.Set(v)
}
