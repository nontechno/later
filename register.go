// Copyright 2024 The NonTechno Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package later

import (
	"fmt"
	"reflect"
	"runtime"
)

func Register(f interface{}, linkages ...string) {
	if f == nil || len(linkages) == 0 {
		onWarning("empty name and/or function")
		return
	}

	signature := getSignature(f)
	if fn, _ := isFunction(f); !fn {
		onError("supplied parameter is not a function (%v, %s)", f, signature)
	}

	guard.Lock()
	defer guard.Unlock()

	for _, linkage := range linkages {
		fullName := getFullname(linkage, signature)

		if entry, found := registry[fullName]; found && entry != nil {
			onWarning("entry (%s) already set/exists", linkage)
			return
		}

		msg := fmt.Sprintf("registered linkage (%s; %s)", linkage, signature)
		if pc, file, no, ok := runtime.Caller(1); ok {
			msg += fmt.Sprintf(". func (%s), file (%s), line (#%d)", runtime.FuncForPC(pc).Name(), file, no)
		}
		onReport(msg)

		registry[fullName] = f
	}
}

func resolve(fn reflect.Value, linkage string) interface{} {

	target, found := registry[linkage]
	if !found {
		onWarning("entry (%s) not found", linkage)
		return nil
	}

	arg := reflect.ValueOf(target)
	fn.Set(arg)

	return target
}
