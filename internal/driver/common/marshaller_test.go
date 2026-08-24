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

import "testing"

func TestB1Array_String(t *testing.T) {
	t.Parallel()

	if got := B1Array(nil).String(); got != "" {
		t.Fatalf("expected nil B1Array to stringify to empty string, got %q", got)
	}

	input := B1Array{'H', 'e', 'l', 'l', 'o'}
	if got := input.String(); got != "Hello" {
		t.Fatalf("expected %q, got %q", "Hello", got)
	}
}

func TestB1Array_Equals(t *testing.T) {
	t.Parallel()

	left := B1Array{1, 2, 3}
	same := B1Array{1, 2, 3}
	different := B1Array{1, 2, 4}

	if !left.Equals(same) {
		t.Fatal("expected byte arrays with identical contents to be equal")
	}

	if left.Equals(different) {
		t.Fatal("expected byte arrays with different contents to be unequal")
	}
}

func TestKeyValue_String(t *testing.T) {
	t.Parallel()

	kv := &KeyValue{
		Key:   StringToB1Array("user"),
		Value: StringToB1Array("scott"),
		Flag:  SB4(-7),
	}

	if got := kv.String(); got != "[user=scott,-7]" {
		t.Fatalf("unexpected KeyValue string representation: %q", got)
	}
}

func TestKeyValue_Equals(t *testing.T) {
	t.Parallel()

	base := &KeyValue{
		Key:   StringToB1Array("alpha"),
		Value: StringToB1Array("beta"),
		Flag:  SB4(42),
	}

	if !base.Equals(base) {
		t.Fatal("expected KeyValue to equal itself")
	}

	if !base.Equals(&KeyValue{
		Key:   StringToB1Array("alpha"),
		Value: StringToB1Array("beta"),
		Flag:  SB4(42),
	}) {
		t.Fatal("expected identical KeyValue contents to be equal")
	}

	if base.Equals(&KeyValue{
		Key:   StringToB1Array("alpha"),
		Value: StringToB1Array("beta"),
		Flag:  SB4(7),
	}) {
		t.Fatal("expected different flags to make KeyValue unequal")
	}

	if base.Equals(&KeyValue{
		Key:   StringToB1Array("gamma"),
		Value: StringToB1Array("beta"),
		Flag:  SB4(42),
	}) {
		t.Fatal("expected different keys to make KeyValue unequal")
	}

	if base.Equals(&KeyValue{
		Key:   StringToB1Array("alpha"),
		Value: StringToB1Array("delta"),
		Flag:  SB4(42),
	}) {
		t.Fatal("expected different values to make KeyValue unequal")
	}
}
