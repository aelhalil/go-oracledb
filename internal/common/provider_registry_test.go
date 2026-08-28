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
	"reflect"
	"testing"

	oracleProviders "github.com/oracle/go-oracledb/v26/oracle/providers"
)

type namedProvider interface {
	oracleProviders.Provider
	Name() string
}

type mockProviderRegistryProvider struct {
	name string
}

func (m mockProviderRegistryProvider) Name() string {
	return m.name
}

// TestProviderRegistryGetProviderReturnsFirstMatch verifies that GetProvider
// returns the first registered provider matching the requested type.
func TestProviderRegistryGetProviderReturnsFirstMatch(t *testing.T) {
	t.Parallel()

	registry := NewProviderRegistry()
	first := mockProviderRegistryProvider{name: "first"}
	second := mockProviderRegistryProvider{name: "second"}

	registry.RegisterProvider(first)
	registry.RegisterProvider(second)

	gotProvider, err := registry.Provider(reflect.TypeOf((*namedProvider)(nil)).Elem())
	if err != nil {
		t.Fatalf("GetProvider returned error: %v", err)
	}
	got := gotProvider.(namedProvider)
	if got.Name() != "first" {
		t.Fatalf("expected first provider, got %q", got.Name())
	}
}

// TestProviderRegistryRegisterProviderEvictsOldestWhenCapacityExceeded
// verifies that registering more than maxProviders providers removes the
// oldest provider, so the next matching provider becomes visible.
func TestProviderRegistryRegisterProviderEvictsOldestWhenCapacityExceeded(t *testing.T) {
	t.Parallel()

	registry := NewProviderRegistry()
	for i := 0; i < maxProviders; i++ {
		registry.RegisterProvider(mockProviderRegistryProvider{name: string(rune('a' + i))})
	}
	registry.RegisterProvider(mockProviderRegistryProvider{name: "overflow"})

	gotProvider, err := registry.Provider(reflect.TypeOf((*namedProvider)(nil)).Elem())
	if err != nil {
		t.Fatalf("GetProvider returned error: %v", err)
	}
	got := gotProvider.(namedProvider)
	if got.Name() != "b" {
		t.Fatalf("expected oldest provider to be evicted, got %q", got.Name())
	}
}

// TestProviderRegistryGetProviderReturnsRequestedInterface verifies that
// GetProvider can populate a requested provider interface from the registry.
func TestProviderRegistryGetProviderReturnsRequestedInterface(t *testing.T) {
	t.Parallel()

	registry := NewProviderRegistry()
	original := mockProviderRegistryProvider{name: "original"}
	registry.RegisterProvider(original)

	gotProvider, err := registry.Provider(reflect.TypeOf((*namedProvider)(nil)).Elem())
	if err != nil {
		t.Fatalf("GetProvider returned error: %v", err)
	}
	got := gotProvider.(namedProvider)
	if got.Name() != "original" {
		t.Fatalf("expected original provider, got %q", got.Name())
	}
}

// TestProviderRegistryGetProviderReturnsErrorWhenUninitialized verifies that a
// newly created registry returns nil when no provider is found.
func TestProviderRegistryGetProviderReturnsErrorWhenUninitialized(t *testing.T) {
	t.Parallel()

	registry := NewProviderRegistry()

	provider, _ := registry.Provider(reflect.TypeOf((*namedProvider)(nil)).Elem())
	if provider != nil {
		t.Fatal("expected GetProvider to fail for empty registry")
	}
}

var _ oracleProviders.Provider = mockProviderRegistryProvider{}
