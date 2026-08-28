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
	"sync"

	oracleErrors "github.com/oracle/go-oracledb/v26/oracle/errors"
	"github.com/oracle/go-oracledb/v26/oracle/providers"
)

// maxProviders maximum number of provider allowed
const maxProviders = 50

// ProviderRegistry stores connector runtime providers and exposes safe access
// to registered providers.
type ProviderRegistry interface {
	// RegisterProvider appends provider to the registry.
	// When the registry already contains maxProviders entries, the oldest
	// registered provider is removed before provider is appended.
	//
	// Parameters:
	//   - provider: the provider to register.
	RegisterProvider(provider providers.Provider)
	// Provider gets the first provider found in the registry that implements
	// the desired type
	//
	// Parameters:
	//   - providerType: the interface or concrete provider type to resolve from
	//     the registry.
	//
	// Returns:
	//   - the first registered provider matching providerType, or nil if no
	//     matching provider was found
	//   - an error if providerType is invalid
	Provider(providerType reflect.Type) (providers.Provider, error)
}

// providerRegistry is the default thread-safe ProviderRegistry implementation.
type providerRegistry struct {
	providers             []providers.Provider
	providerRegistryMutex sync.RWMutex
}

// NewProviderRegistry creates an empty provider registry.
//
// Returns:
//   - the initialized provider registry.
func NewProviderRegistry() *providerRegistry {
	return &providerRegistry{}
}

// RegisterProvider appends provider to the registry.
// When the registry already contains maxProviders entries, the oldest
// registered provider is removed before provider is appended.
//
// Parameters:
//   - provider: the provider to register.
func (providerRegistry *providerRegistry) RegisterProvider(provider providers.Provider) {
	providerRegistry.providerRegistryMutex.Lock()
	defer providerRegistry.providerRegistryMutex.Unlock()
	if len(providerRegistry.providers) >= maxProviders {
		Odl.Warn("Number of provider exceeded the maximum allowed, removing the oldest provider")
		providerRegistry.providers = providerRegistry.providers[1:]
	}
	providerRegistry.providers = append(providerRegistry.providers, provider)
}

// Provider gets the first provider found in the registry that implements
// the desired type
//
// Parameters:
//   - providerType: the interface or concrete provider type to resolve from
//     the registry.
//
// Returns:
//   - the first registered provider matching providerType, or nil if no
//     matching provider was found
//   - an error if providerType is invalid
func (providerRegistry *providerRegistry) Provider(providerType reflect.Type) (providers.Provider, error) {
	providerRegistry.providerRegistryMutex.RLock()
	defer providerRegistry.providerRegistryMutex.RUnlock()

	if providerType == nil {
		return nil, NewOracleError(oracleErrors.ProviderNotFound, nil)
	}
	for _, item := range providerRegistry.providers {
		itemType := reflect.TypeOf(item)
		if itemType.AssignableTo(providerType) || (providerType.Kind() == reflect.Interface && itemType.Implements(providerType)) {
			return item, nil
		}
	}
	return nil, nil
}
