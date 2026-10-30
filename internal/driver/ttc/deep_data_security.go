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

package ttc

import (
	"maps"
	"slices"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/oracle/go-oracledb/v26/internal/common"
	"github.com/oracle/go-oracledb/v26/internal/driver/ttc/oson"
	oracleErrors "github.com/oracle/go-oracledb/v26/oracle/errors"
	osql "github.com/oracle/go-oracledb/v26/oracle/sql"
)

// deepDataSecurity holds connection-wide end-user security context selection.
type deepDataSecurity struct {
	connection *connection
	mu         sync.Mutex
	context    osql.EndUserSecurityContext
}

// SetEndUserSecurityContext validates context and selects it for subsequent
// database calls on this connection.
func (c *connection) SetEndUserSecurityContext(context osql.EndUserSecurityContext) error {
	if c._isClosed || !c._isValid {
		return c.shelf.LocalizeError(common.NewOracleError(oracleErrors.InvalidConnection, nil))
	}
	if context == nil {
		common.Odl.Debug("Deep Data Security: nil context")
		return c.shelf.LocalizeError(common.NewOracleError(oracleErrors.InvalidEndUserSecurityContext, nil))
	}
	if err := c.checkDeepDataSecuritySupport(); err != nil {
		return c.shelf.LocalizeError(err)
	}

	// validation
	if _, err := encodeEndUserSecurityContext(context); err != nil {
		return c.shelf.LocalizeError(err)
	}
	state := c.shelf._deepDataSecurity
	state.mu.Lock()
	state.context = context
	state.mu.Unlock()
	common.Odl.Debug("Deep Data Security: explicit context selected")
	return nil
}

// ClearEndUserSecurityContext stops propagating the selected context.
// No database call is made.
func (c *connection) ClearEndUserSecurityContext() {
	state := c.shelf._deepDataSecurity
	if state == nil {
		return
	}
	state.mu.Lock()
	state.context = nil
	state.mu.Unlock()
	common.Odl.Debug("Deep Data Security: explicit context cleared")
}

// SetEndUserSecurityContext selects the connection-wide context for this statement.
func (s *Statement) SetEndUserSecurityContext(context osql.EndUserSecurityContext) error {
	return s.shelf._deepDataSecurity.connection.SetEndUserSecurityContext(context)
}

// ClearEndUserSecurityContext clears the connection-wide explicit context.
func (s *Statement) ClearEndUserSecurityContext() {
	s.shelf._deepDataSecurity.connection.ClearEndUserSecurityContext()
}

// SetEndUserSecurityContext selects the connection-wide context for this transaction.
func (t *transaction) SetEndUserSecurityContext(context osql.EndUserSecurityContext) error {
	return t._underlyingConnection.SetEndUserSecurityContext(context)
}

// ClearEndUserSecurityContext clears the connection-wide explicit context.
func (t *transaction) ClearEndUserSecurityContext() {
	t._underlyingConnection.ClearEndUserSecurityContext()
}

// checkDeepDataSecuritySupport rejects servers without context propagation support.
func (c *connection) checkDeepDataSecuritySupport() error {
	if !c.shelf.GetCapabilities()[kpccapCtbBackportFeature1Xsctx].IsSet {
		return common.NewOracleError(oracleErrors.DeepDataSecurityNotSupported, nil)
	}
	return nil
}

// prepare encodes the selected context before any bytes of a database call
// are written. Without a selection, the call has no security context.
func (state *deepDataSecurity) prepare() (*ttiOXSCTX, error) {
	state.mu.Lock()
	selected := state.context
	state.mu.Unlock()
	if selected == nil {
		return nil, nil
	}

	if err := state.connection.checkDeepDataSecuritySupport(); err != nil {
		return nil, err
	}
	payload, err := encodeEndUserSecurityContext(selected)
	if err != nil {
		return nil, err
	}
	common.Odl.Debug("Deep Data Security: full context prepared")
	message, err := state.connection.shelf.GetMessageFactory().GetMessageForFunction(TTIPFN, oXsctx)
	if err != nil {
		return nil, err
	}
	msg := message.(*ttiOXSCTX)
	msg.pairs = []*keywordValuePairWithName{
		newKeywordValuePairWithName("ORCL_XS_AUTHZ_CONTEXT", "", payload, 0),
	}
	return msg, nil
}

// encodeEndUserSecurityContext validates credentials and encodes the complete
// end-user context using the database's OSON schema.
func encodeEndUserSecurityContext(context osql.EndUserSecurityContext) ([]byte, error) {
	accessToken := context.DatabaseAccessToken()
	endUserToken, tokenBased := context.EndUserToken()
	endUserName, nameBased := context.EndUserName()
	lookupKey, hasKey := context.LookUpKey()

	// validate credentials
	if len(accessToken) == 0 || !utf8.Valid(accessToken) {
		common.Odl.Debug("Deep Data Security: invalid database access token")
		return nil, common.NewOracleError(oracleErrors.InvalidEndUserSecurityContext, nil)
	}
	if tokenBased == nameBased || (tokenBased && (len(endUserToken) == 0 || !utf8.Valid(endUserToken))) || (nameBased && (endUserName == "" || !utf8.ValidString(endUserName))) {
		common.Odl.Debug("Deep Data Security: exactly one nonempty principal is required")
		return nil, common.NewOracleError(oracleErrors.InvalidEndUserSecurityContext, nil)
	}

	// check expiration
	expires := context.DatabaseAccessTokenExpiration()
	if !time.Now().Before(expires) {
		common.Odl.Debug("Deep Data Security: database access token expired")
		return nil, common.NewOracleError(oracleErrors.ExpiredToken, nil)
	}
	if tokenBased {
		userExpires := context.EndUserTokenExpiration()
		if !time.Now().Before(userExpires) {
			common.Odl.Debug("Deep Data Security: end-user token expired")
			return nil, common.NewOracleError(oracleErrors.ExpiredToken, nil)
		}
	}

	// build payload
	payload := map[string]any{"ver": "1.0"}
	payload["database_access_token"] = string(accessToken)
	if tokenBased {
		payload["end_user_token"] = string(endUserToken)
	}
	if nameBased {
		payload["end_user_name"] = endUserName
		// add lookup key
		if hasKey {
			if len(lookupKey) == 0 || !utf8.Valid(lookupKey) {
				common.Odl.Debug("Deep Data Security: invalid lookup key")
				return nil, common.NewOracleError(oracleErrors.InvalidEndUserSecurityContext, nil)
			}
			payload["end_user_contextid"] = string(lookupKey)
		}
	}

	// add roles
	roles := slices.Clone(context.DataRoles())
	slices.Sort(roles)
	roles = slices.Compact(roles)
	if len(roles) > 0 {
		values := make([]any, len(roles))
		for index, role := range roles {
			if role == "" || !utf8.ValidString(role) {
				common.Odl.Debug("Deep Data Security: invalid data role")
				return nil, common.NewOracleError(oracleErrors.InvalidEndUserSecurityContext, nil)
			}
			values[index] = role
		}
		payload["data_roles"] = values
	}

	// add attributes
	attributes := context.Attributes()
	if len(attributes) > 0 {
		values := make([]any, 0, len(attributes))
		for _, name := range slices.Sorted(maps.Keys(attributes)) {
			if name == "" || !utf8.ValidString(name) {
				common.Odl.Debug("Deep Data Security: invalid attribute name")
				return nil, common.NewOracleError(oracleErrors.InvalidEndUserSecurityContext, nil)
			}
			object, err := attributes[name].AsJSONObject()
			if err != nil {
				common.Odl.Debug("Deep Data Security: attributes must be JSON objects")
				return nil, common.NewOracleError(oracleErrors.InvalidEndUserSecurityContext, err)
			}
			value, err := object.GetValue()
			if err != nil {
				common.Odl.Debug("Deep Data Security: cannot read attributes")
				return nil, common.NewOracleError(oracleErrors.InvalidEndUserSecurityContext, err)
			}
			values = append(values, map[string]any{"name": name, "values": value})
		}
		payload["attributes"] = values
	}

	// encode payload
	encoded, err := oson.Encode(payload)
	if err != nil {
		common.Odl.Debug("Deep Data Security: cannot encode context")
		return nil, common.NewOracleError(oracleErrors.InvalidEndUserSecurityContext, err)
	}

	if len(encoded) > maxXsctxValueLength {
		common.Odl.Debug("Deep Data Security: context exceeds 65536 bytes")
		return nil, common.NewOracleError(oracleErrors.InvalidEndUserSecurityContext, nil)
	}
	return encoded, nil
}
