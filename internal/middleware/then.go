// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package middleware

import (
	"net/http"
	"slices"
)

type Chain []func(http.Handler) http.Handler

func (c Chain) ThenFunc(h http.HandlerFunc) http.Handler {
	return c.Then(h)
}

func (c Chain) Then(h http.Handler) http.Handler {
	for _, mw := range slices.Backward(c) {
		h = mw(h)
	}
	return h
}
