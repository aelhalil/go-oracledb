/*
** Copyright (c) 2026 Oracle and/or its affiliates.
**
** The Universal Permissive License (UPL), Version 1.0
**
** Subject to the condition set forth below, permission is hereby granted to any
** person obtaining a copy of this software, associated documentation and/or data
** (collectively the "Software"), free of charge and under any and all copyright
** rights in the Software, and any and all patent rights owned or freely
** licensable by each licensor hereunder covering either (i) the unmodified
** Software as contributed to or provided by such licensor, or (ii) the Larger
** Works (as defined below), to deal in both
**
** (a) the Software, and
** (b) any piece of software and/or hardware listed in the lrgrwrks.txt file if
** one is included with the Software (each a "Larger Work" to which the Software
** is contributed by such licensors),
**
** without restriction, including without limitation the rights to copy, create
** derivative works of, display, perform, and distribute the Software and make,
** use, sell, offer for sale, import, export, have made, and have sold the
** Software and the Larger Work(s), and to sublicense the foregoing rights on
** either these or other terms.
**
** This license is subject to the following condition:
** The above copyright notice and either this complete permission notice or at
** a minimum a reference to the UPL must be included in all copies or
** substantial portions of the Software.
**
** THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
** IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
** FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
** AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
** LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
** OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
** SOFTWARE.
 */

package common

import (
	"flag"
	"os"
	"strings"
	"testing"
)

// TestCategory category of tests to be un
var TestCategory string

func TestMain(m *testing.M) {
	flag.StringVar(&TestCategory, "test.category", "", "testing category, can be unitary, functional, performance, robustness")
	os.Exit(m.Run())
}

var testCases = []struct {
	name       string
	categories string
	exclusive  bool
	f          func(t *testing.T)
}{
	{"TestError3113default", "unitary", false, TestError3113default},
	{"TestErrorUnknownCode", "unitary", false, TestErrorUnknownCode},
	{"TestErrorBasic3113", "unitary", false, TestErrorBasic3113},
	{"TestErrorBasic3113Fr", "unitary", false, TestErrorBasic3113Fr},
	{"TestError3113", "unitary", false, TestError3113},
	{"TestError3113Fr", "unitary", false, TestError3113Fr},
	{"TestError3113pt", "unitary", false, TestError3113pt},
	{"TestError3113NoLanguage", "unitary", false, TestError3113NoLanguage},
	{"TestErrorUnwrap", "unitary", false, TestErrorUnwrap},
	{"TestError3113InvalidLanguage", "unitary", false, TestError3113InvalidLanguage},
	{"TestNewOERMessageError", "unitary", false, TestNewOERMessageError},
	{"TestConstants_GetLogonModeFromString", "unitary", false, TestConstants_GetLogonModeFromString},
	{"TestConstants_LogonModeEnabled", "unitary", false, TestConstants_LogonModeEnabled},
	{"TestConstants_LogonModeString", "unitary", false, TestConstants_LogonModeString},
	{"TestProviderRegistryGetProviderReturnsFirstMatch", "unitary", false, TestProviderRegistryGetProviderReturnsFirstMatch},
	{"TestProviderRegistryRegisterProviderEvictsOldestWhenCapacityExceeded", "unitary", false, TestProviderRegistryRegisterProviderEvictsOldestWhenCapacityExceeded},
	{"TestProviderRegistryGetProviderReturnsRequestedInterface", "unitary", false, TestProviderRegistryGetProviderReturnsRequestedInterface},
	{"TestProviderRegistryGetProviderReturnsErrorWhenUninitialized", "unitary", false, TestProviderRegistryGetProviderReturnsErrorWhenUninitialized},
}

func TestCategoryExecutor(t *testing.T) {
	var regularCases, exclusiveCases []struct {
		name       string
		categories string
		exclusive  bool
		f          func(t *testing.T)
	}

	for _, c := range testCases {
		cats := strings.Split(c.categories, ",")
		for _, p := range cats {
			if strings.Compare(strings.TrimSpace(p), TestCategory) == 0 {
				if c.exclusive {
					exclusiveCases = append(exclusiveCases, c)
				} else {
					regularCases = append(regularCases, c)
				}
				break
			}
		}
	}

	if len(regularCases) > 0 {
		t.Run("parallel", func(t *testing.T) {
			t.Parallel()
			for _, c := range regularCases {
				t.Run(c.name, c.f)
			}
		})
	}
	for _, c := range exclusiveCases {
		t.Run(c.name, c.f)
	}
}
