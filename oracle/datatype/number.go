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

// Package datatype provides Go types for Oracle Database scalar datatypes
// that extend the value types supported by database/sql.
package datatype

import (
	"database/sql/driver"
	"fmt"
	"math"
	"strconv"
)

/*
Number represents an Oracle NUMBER value carried as its exact decimal text.
It preserves all significant digits and the declared scale without converting
through a binary floating-point type.

The driver produces Number when fetching NUMBER columns whose precision and
scale metadata require an exact decimal representation. Because Number is a
string-based type, scanning into a plain string destination keeps working,
and Number itself can be used directly as a scan destination or bind value
through its sql.Scanner and driver.Valuer implementations.

Use Int64 and Float64 to materialize a Number when a native Go numeric is
acceptable; both fail rather than silently truncate or round beyond what the
stored text supports.
*/
type Number string

// String returns the decimal text representation of the number. It
// implements fmt.Stringer.
func (n Number) String() string {
	return string(n)
}

/*
Int64 returns the number as an int64.

Errors:
  - Returns a *strconv.NumError when the stored text has a fractional part,
    exceeds the int64 range, or is not a valid base-10 integer.
*/
func (n Number) Int64() (int64, error) {
	return strconv.ParseInt(string(n), 10, 64)
}

/*
Float64 returns the number as a float64.

Errors:
  - Returns a *strconv.NumError when the stored text is not a valid
    floating-point number. Values with more significant digits than float64
    can represent are rounded to the nearest representable number, like
    strconv.ParseFloat.
*/
func (n Number) Float64() (float64, error) {
	return strconv.ParseFloat(string(n), 64)
}

/*
Scan implements sql.Scanner, storing the decimal text of src in n. It accepts
string, []byte, Number, int64, and finite float64 sources; []byte content is
copied.

Errors:
  - Returns an error when src is SQL NULL; use sql.Null[Number] for nullable
    columns.
  - Returns an error when src is a type without a decimal representation or a
    non-finite float64.
*/
func (n *Number) Scan(src any) error {
	switch value := src.(type) {
	case nil:
		return fmt.Errorf("datatype: cannot scan NULL into datatype.Number; use sql.Null[datatype.Number] for nullable columns")
	case Number:
		*n = value
	case string:
		*n = Number(value)
	case []byte:
		*n = Number(append([]byte(nil), value...))
	case int64:
		*n = Number(strconv.FormatInt(value, 10))
	case float64:
		if math.IsNaN(value) || math.IsInf(value, 0) {
			return fmt.Errorf("datatype: cannot scan non-finite float64 %v into datatype.Number", value)
		}
		*n = Number(strconv.FormatFloat(value, 'f', -1, 64))
	default:
		return fmt.Errorf("datatype: cannot scan %T into datatype.Number", src)
	}
	return nil
}

// Value implements driver.Valuer, binding the stored decimal text as a
// string driver.Value so the exact digits reach the server without
// float64 rounding.
func (n Number) Value() (driver.Value, error) {
	return string(n), nil
}
