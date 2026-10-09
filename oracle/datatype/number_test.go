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

package datatype

import (
	"errors"
	"math"
	"strconv"
	"testing"
)

// TestNumber_String verifies that String returns the stored decimal text
// unchanged, including the empty Number.
func TestNumber_String(t *testing.T) {
	t.Parallel()

	cases := []struct {
		number Number
		want   string
	}{
		{number: "123.45", want: "123.45"},
		{number: "-0.25", want: "-0.25"},
		{number: "0", want: "0"},
		{number: "", want: ""},
	}

	for _, tc := range cases {
		if got := tc.number.String(); got != tc.want {
			t.Errorf("Number(%q).String() = %q, want %q", tc.number, got, tc.want)
		}
	}
}

// TestNumber_Int64 verifies integer materialization, including the failures
// for fractional, out-of-range, and malformed text.
func TestNumber_Int64(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		number  Number
		want    int64
		wantErr bool
	}{
		{name: "positive", number: "42", want: 42},
		{name: "negative", number: "-987654321", want: -987654321},
		{name: "zero", number: "0", want: 0},
		{name: "max int64", number: "9223372036854775807", want: math.MaxInt64},
		{name: "fractional", number: "123.45", wantErr: true},
		{name: "scaled zero fraction", number: "0.00", wantErr: true},
		{name: "overflow", number: "9223372036854775808", wantErr: true},
		{name: "malformed", number: "1e5", wantErr: true},
		{name: "empty", number: "", wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := tc.number.Int64()
			if tc.wantErr {
				if err == nil {
					t.Fatalf("Number(%q).Int64() = %d, want error", tc.number, got)
				}
				var numErr *strconv.NumError
				if !errors.As(err, &numErr) {
					t.Fatalf("Number(%q).Int64() error = %v, want *strconv.NumError", tc.number, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Number(%q).Int64() returned err = %v, want %d", tc.number, err, tc.want)
			}
			if got != tc.want {
				t.Errorf("Number(%q).Int64() = %d, want %d", tc.number, got, tc.want)
			}
		})
	}
}

// TestNumber_Float64 verifies floating-point materialization, including
// values whose precision exceeds float64 and malformed text.
func TestNumber_Float64(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		number  Number
		want    float64
		wantErr bool
	}{
		{name: "integer text", number: "42", want: 42},
		{name: "decimal text", number: "123.45", want: 123.45},
		{name: "negative scale", number: "-0.25", want: -0.25},
		{name: "beyond float64 precision", number: "123456789012345678901234567890.1", want: 1.2345678901234568e+29},
		{name: "malformed", number: "abc", wantErr: true},
		{name: "empty", number: "", wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := tc.number.Float64()
			if tc.wantErr {
				if err == nil {
					t.Fatalf("Number(%q).Float64() = %v, want error", tc.number, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("Number(%q).Float64() returned err = %v, want %v", tc.number, err, tc.want)
			}
			if got != tc.want {
				t.Errorf("Number(%q).Float64() = %v, want %v", tc.number, got, tc.want)
			}
		})
	}
}

// TestNumber_Scan verifies the sql.Scanner implementation across accepted
// source types and the NULL, unsupported-type, and non-finite rejections.
func TestNumber_Scan(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		src     any
		want    Number
		wantErr bool
	}{
		{name: "string", src: "123.45", want: "123.45"},
		{name: "empty string", src: "", want: ""},
		{name: "byte slice", src: []byte("-0.25"), want: "-0.25"},
		{name: "number", src: Number("42"), want: "42"},
		{name: "int64", src: int64(-987654321), want: "-987654321"},
		{name: "float64", src: 0.0000001, want: "0.0000001"},
		{name: "float64 large", src: 1e21, want: "1000000000000000000000"},
		{name: "null", src: nil, wantErr: true},
		{name: "nan", src: math.NaN(), wantErr: true},
		{name: "positive infinity", src: math.Inf(1), wantErr: true},
		{name: "negative infinity", src: math.Inf(-1), wantErr: true},
		{name: "unsupported bool", src: true, wantErr: true},
		{name: "unsupported int", src: 5, wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var got Number
			if err := (&got).Scan(tc.src); err != nil {
				if !tc.wantErr {
					t.Fatalf("Scan(%#v) returned err = %v, want %q", tc.src, err, tc.want)
				}
				return
			}
			if tc.wantErr {
				t.Fatalf("Scan(%#v) = %q, want error", tc.src, got)
			}
			if got != tc.want {
				t.Errorf("Scan(%#v) = %q, want %q", tc.src, got, tc.want)
			}
		})
	}

	// []byte sources must be copied, because drivers may reuse the buffer.
	buf := []byte("9.75")
	var copied Number
	if err := copied.Scan(buf); err != nil {
		t.Fatalf("Scan([]byte) returned err = %v", err)
	}
	for i := range buf {
		buf[i] = "0123456789"[i]
	}
	if copied != "9.75" {
		t.Fatalf("Scan([]byte) aliased the source buffer: got %q", copied)
	}
}

// TestNumber_Value verifies the driver.Valuer implementation binds the exact
// decimal text as a string driver.Value.
func TestNumber_Value(t *testing.T) {
	t.Parallel()

	value, err := Number("123.4500000000").Value()
	if err != nil {
		t.Fatalf("Number.Value() returned err = %v", err)
	}
	got, ok := value.(string)
	if !ok {
		t.Fatalf("Number.Value() = %#v (%T), want string", value, value)
	}
	if got != "123.4500000000" {
		t.Errorf("Number.Value() = %q, want %q", got, "123.4500000000")
	}
}

// TestNumber_PrecisionPreservation verifies that a 38-digit Oracle NUMBER
// survives as text and that conversions report rather than hide the loss of
// precision.
func TestNumber_PrecisionPreservation(t *testing.T) {
	t.Parallel()

	const exact = "99999999999999999999999999999999999999"
	n := Number(exact)

	if n.String() != exact {
		t.Fatalf("Number.String() = %q, want %q", n.String(), exact)
	}
	if _, err := n.Int64(); err == nil {
		t.Fatal("Number.Int64() succeeded for a value beyond int64 range")
	}
	f, err := n.Float64()
	if err != nil {
		t.Fatalf("Number.Float64() returned err = %v", err)
	}
	if back := strconv.FormatFloat(f, 'f', -1, 64); back == exact {
		t.Fatal("float64 unexpectedly round-tripped all 38 digits")
	}
}
