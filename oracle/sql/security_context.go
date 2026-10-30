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

package sql

import (
	"time"

	"github.com/oracle/go-oracledb/v26/internal/common"
	"github.com/oracle/go-oracledb/v26/oracle/datatype/json"
	oracleErrors "github.com/oracle/go-oracledb/v26/oracle/errors"
)

// SecurityContextConfig configures the credentials, data roles and attributes
// for NewSecurityContext. Only DatabaseAccessToken is required. EndUserToken and
// EndUserName cannot both be set; leaving both empty uses the database access token alone.
type SecurityContextConfig struct {
	// DatabaseAccessToken is the required, nonempty UTF-8 database access token.
	DatabaseAccessToken []byte
	// EndUserToken is an optional UTF-8 end-user token. Do not set with EndUserName.
	EndUserToken []byte
	// EndUserName is an optional UTF-8 end-user name. Do not set with EndUserToken.
	EndUserName string
	// LookupKey is an optional UTF-8 context identifier, used only with EndUserName.
	LookupKey []byte
	// DataRoles lists optional additional data roles. Nil or empty requests none.
	DataRoles []string
	// Attributes contains optional JSON objects keyed by fully qualified context
	// names (SCHEMA.CONTEXT_NAME). Nil or empty supplies none.
	Attributes map[string]json.JSON
}

// SecurityContext implements [EndUserSecurityContext].
type SecurityContext struct{}

// NewSecurityContext creates an immutable context from config. It requires a
// nonempty database access token and accepts at most one of EndUserToken and
// EndUserName. Leaving both empty selects a database-access-token-only context.
// Supplying both returns an oracle/errors.InvalidEndUserSecurityContext error.
// Other invalid credentials, roles or attributes return the same error code;
// expired credentials return oracle/errors.ExpiredToken.
//
// The configuration's mutable values are copied. Subsequent changes to config,
// its slices or its attribute map do not change the constructed context.
// The application is responsible for supplying and refreshing credentials.
//
// Implementation is pending. Currently every call returns nil and an
// oracle/errors.UnsupportedFeature error without validating config.
func NewSecurityContext(config SecurityContextConfig) (*SecurityContext, error) {
	// TODO
	return nil, common.NewOracleError(oracleErrors.UnsupportedFeature, nil, "NewSecurityContext")
}

// DatabaseAccessToken implements [EndUserSecurityContext.DatabaseAccessToken].
func (c *SecurityContext) DatabaseAccessToken() []byte {
	// TODO
	return nil
}

// DatabaseAccessTokenExpiration implements
// [EndUserSecurityContext.DatabaseAccessTokenExpiration].
func (c *SecurityContext) DatabaseAccessTokenExpiration() time.Time {
	// TODO
	return time.Time{}
}

// EndUserToken implements [EndUserSecurityContext.EndUserToken].
func (c *SecurityContext) EndUserToken() ([]byte, bool) {
	// TODO
	return nil, false
}

// EndUserTokenExpiration implements [EndUserSecurityContext.EndUserTokenExpiration].
func (c *SecurityContext) EndUserTokenExpiration() time.Time {
	// TODO
	return time.Time{}
}

// EndUserName implements [EndUserSecurityContext.EndUserName].
func (c *SecurityContext) EndUserName() (string, bool) {
	// TODO
	return "", false
}

// DataRoles implements [EndUserSecurityContext.DataRoles].
func (c *SecurityContext) DataRoles() []string {
	// TODO
	return nil
}

// Attributes implements [EndUserSecurityContext.Attributes].
func (c *SecurityContext) Attributes() map[string]json.JSON {
	// TODO
	return nil
}

// LookUpKey implements [EndUserSecurityContext.LookUpKey].
func (c *SecurityContext) LookUpKey() ([]byte, bool) {
	// TODO
	return nil, false
}
