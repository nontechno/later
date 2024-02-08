// Copyright 2024 The NonTechno Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package later

import (
	"reflect"
	"sync/atomic"
)

func Link(fptr interface{}, linkage string, fallback interface{}) {
	signature := ""
	if pointer, elem := isPointer(fptr); pointer {
		if yes, _ := isFunction(elem); !yes {
			onError("supplied parameter is not a pointer to a function (%v)", fptr)
		}
		signature = getSignature(elem)
	} else {
		onError("supplied parameter is not a pointer (%v)", fptr)
	}

	if fallback != nil {
		if fn, _ := isFunction(fallback); !fn {
			onError("supplied fallback is not a function (%v)", fptr)
		}
	}

	fn := reflect.ValueOf(fptr).Elem()

	universal := func(args []reflect.Value) []reflect.Value {

		// a simple guard against an endless recursion
		current := atomic.AddInt32(&entryCounter, 1)
		defer atomic.AddInt32(&entryCounter, -1)
		if current > maxDepth {
			onError("too many hops (%v) to resolve linkage (%s)", current, linkage)
			return []reflect.Value{}
		}

		if target := resolve(fn, linkage, signature); target != nil {
			onReport("resolved linkage (%s)", linkage)
			return reflect.ValueOf(target).Call(args) // this can be a chain call
		} else {
			if fallback != nil {
				arg := reflect.ValueOf(fallback)
				fn.Set(arg)

				onReport("unresolved linkage (%s), using fallback", linkage)
				return reflect.ValueOf(fallback).Call(args) // this can be a chain call
			}

			onError("failed to resolve linkage (%s)", linkage)
			return []reflect.Value{}
		}
	}

	v := reflect.MakeFunc(fn.Type(), universal)
	fn.Set(v)
}
