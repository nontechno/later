// Copyright 2024 The NonTechno Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package later

import "sync/atomic"

type reportFunc func(format string, args ...interface{})

func callRemote(remote, format string, args ...interface{}) {

	current := atomic.AddInt32(&remoteCounter, 1)
	defer atomic.AddInt32(&remoteCounter, -1)

	if current > 12 {
		// this is too deep - most likely an endless recursion
		return
	}

	if target, found := registry[remote]; found && target != nil {
		if operation, compliant := target.(reportFunc); compliant && operation != nil {
			operation(format, args...)
		}
	}
}

func localOnWarning(format string, args ...interface{}) {
	callRemote("warning", format, args...)
}

func localOnError(format string, args ...interface{}) {
	callRemote("error", format, args...)
}
