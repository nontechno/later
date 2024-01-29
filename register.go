// Copyright 2024 The NonTechno Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package later

import (
	"reflect"
	"sync"
)

var (
	guard         sync.Mutex
	registry      = map[string]interface{}{}
	onWarning     = localOnWarning
	onError       = localOnError
	entryCounter  int32
	remoteCounter int32
)

func Register(f interface{}, linkage string) {
	if f == nil || len(linkage) == 0 {
		onWarning("empty name and/or function")
		return
	}

	guard.Lock()
	defer guard.Unlock()

	if entry, found := registry[linkage]; found && entry != nil {
		onWarning("entry (%s) already set/exists", linkage)
		return
	}

	registry[linkage] = f
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
