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
	"errors"
	"math"
	"testing"

	oracleErrors "github.com/oracle/go-oracledb/v26/oracle/errors"
)

// TestBufferReadWrite verifies typed reads and writes for both byte orders.
func TestBufferReadWrite(t *testing.T) {
	for _, order := range []ByteOrder{BIG_ENDIAN, LITTLE_ENDIAN} {
		t.Run(order.name, func(t *testing.T) {
			writer := NewWriteBuffer(order)
			if got := writer.Reserve(0); got != 0 {
				t.Fatalf("Reserve(0) = %d, want 0", got)
			}
			if got := writer.Reserve(22); got != 0 {
				t.Fatalf("Reserve(22) = %d, want 0", got)
			}
			if writer.Position() != 22 || writer.Size() != 22 {
				t.Fatalf("reserved buffer position/size = %d/%d, want 22/22", writer.Position(), writer.Size())
			}
			if err := writer.WriteUB1At(0, 0x12); err != nil {
				t.Fatal(err)
			}
			if err := writer.WriteUB2At(1, 0x3456); err != nil {
				t.Fatal(err)
			}
			if err := writer.WriteUB4At(3, 0x789abcde); err != nil {
				t.Fatal(err)
			}
			if err := writer.WriteUB8At(7, 0x0123456789abcdef); err != nil {
				t.Fatal(err)
			}
			if err := writer.WriteSB1At(15, -1); err != nil {
				t.Fatal(err)
			}
			if err := writer.WriteSB2At(16, -2); err != nil {
				t.Fatal(err)
			}
			if err := writer.WriteSB4At(18, -3); err != nil {
				t.Fatal(err)
			}

			reader := NewReadBuffer(writer.Bytes(), order)
			if got, err := reader.ReadUB1(); err != nil || got != 0x12 {
				t.Fatalf("ReadUB1() = %#x, %v", got, err)
			}
			if got, err := reader.ReadUB2(); err != nil || got != 0x3456 {
				t.Fatalf("ReadUB2() = %#x, %v", got, err)
			}
			if got, err := reader.ReadUB4(); err != nil || got != 0x789abcde {
				t.Fatalf("ReadUB4() = %#x, %v", got, err)
			}
			if got, err := reader.ReadUB8(); err != nil || got != 0x0123456789abcdef {
				t.Fatalf("ReadUB8() = %#x, %v", got, err)
			}
			if got, err := reader.ReadSB1(); err != nil || got != -1 {
				t.Fatalf("ReadSB1() = %d, %v", got, err)
			}
			if got, err := reader.ReadSB2(); err != nil || got != -2 {
				t.Fatalf("ReadSB2() = %d, %v", got, err)
			}
			if got, err := reader.ReadSB4(); err != nil || got != -3 {
				t.Fatalf("ReadSB4() = %d, %v", got, err)
			}
			if reader.Position() != reader.Size() {
				t.Fatalf("reader position = %d, want %d", reader.Position(), reader.Size())
			}

			if got, err := reader.ReadUB4At(3); err != nil || got != 0x789abcde {
				t.Fatalf("ReadUB4At() = %#x, %v", got, err)
			}
			if got, err := reader.ReadUB1At(0); err != nil || got != 0x12 {
				t.Fatalf("ReadUB1At() = %#x, %v", got, err)
			}
			if got, err := reader.ReadUB2At(1); err != nil || got != 0x3456 {
				t.Fatalf("ReadUB2At() = %#x, %v", got, err)
			}
			if got, err := reader.ReadUB8At(7); err != nil || got != 0x0123456789abcdef {
				t.Fatalf("ReadUB8At() = %#x, %v", got, err)
			}
			if got, err := reader.ReadSB1At(15); err != nil || got != -1 {
				t.Fatalf("ReadSB1At() = %d, %v", got, err)
			}
			if got, err := reader.ReadSB2At(16); err != nil || got != -2 {
				t.Fatalf("ReadSB2At() = %d, %v", got, err)
			}
			if got, err := reader.ReadSB4At(18); err != nil || got != -3 {
				t.Fatalf("ReadSB4At() = %d, %v", got, err)
			}
		})
	}

	writer := NewWriteBuffer(BIG_ENDIAN)
	writer.WriteUB1(1)
	writer.WriteUB2(2)
	writer.WriteUB4(3)
	writer.WriteUB8(4)
	writer.WriteSB1(-1)
	writer.WriteSB2(-2)
	writer.WriteSB4(-3)
	if writer.Size() != 22 {
		t.Fatalf("append writer size = %d, want 22", writer.Size())
	}
	reader := NewReadBuffer(writer.Bytes(), BIG_ENDIAN)
	if got, err := reader.ReadUB1(); err != nil || got != 1 {
		t.Fatalf("append ReadUB1() = %d, %v", got, err)
	}
	if got, err := reader.ReadUB2(); err != nil || got != 2 {
		t.Fatalf("append ReadUB2() = %d, %v", got, err)
	}
	if got, err := reader.ReadUB4(); err != nil || got != 3 {
		t.Fatalf("append ReadUB4() = %d, %v", got, err)
	}
	if got, err := reader.ReadUB8(); err != nil || got != 4 {
		t.Fatalf("append ReadUB8() = %d, %v", got, err)
	}
}

// TestBufferSlicesAndCopies verifies slice operations, automatic growth, and output ownership.
func TestBufferSlicesAndCopies(t *testing.T) {
	empty := NewWriteBuffer(BIG_ENDIAN).Bytes()
	if empty == nil || len(empty) != 0 {
		t.Fatalf("empty Bytes() = %#v, want non-nil empty slice", empty)
	}

	writer := NewWriteBufferWithCapacity(1, BIG_ENDIAN)
	writer.WriteBytes(B1Array{1, 2, 3})
	if err := writer.WriteBytesAt(writer.Size(), B1Array{}); err != nil {
		t.Fatalf("WriteBytesAt(size, empty) error = %v", err)
	}
	if err := writer.WriteBytesAt(1, B1Array{9, 8}); err != nil {
		t.Fatal(err)
	}
	encoded := writer.Bytes()
	encoded[0] = 7
	if got := writer.Bytes()[0]; got != 1 {
		t.Fatalf("Bytes() exposed writer data: got %d", got)
	}

	reader := NewReadBuffer(writer.Bytes(), BIG_ENDIAN)
	view, err := reader.ReadBytes(2)
	if err != nil || string(view) != string(B1Array{1, 9}) {
		t.Fatalf("ReadBytes() = %v, %v", view, err)
	}
	view[0] = 6
	if got, err := reader.ReadUB1At(0); err != nil || got != 6 {
		t.Fatalf("ReadBytes() did not return a view: %d, %v", got, err)
	}
	if err := reader.SetPosition(reader.Size()); err != nil {
		t.Fatal(err)
	}
	if err := reader.SetPosition(0); err != nil {
		t.Fatal(err)
	}
	if got, err := reader.ReadUB1(); err != nil || got != 6 {
		t.Fatalf("ReadUB1() after backward seek = %d, %v", got, err)
	}
}

// TestBufferErrors verifies every buffer error code and overflow-safe range checks.
func TestBufferErrors(t *testing.T) {
	reader := NewReadBuffer(nil, BIG_ENDIAN)
	_, err := reader.ReadUB1()
	assertBufferErrorCode(t, err, oracleErrors.BufferUnderflow)
	_, err = reader.ReadBytes(1)
	assertBufferErrorCode(t, err, oracleErrors.BufferUnderflow)
	if err := reader.SetPosition(1); err == nil {
		t.Fatal("SetPosition(1) error = nil")
	} else {
		assertBufferErrorCode(t, err, oracleErrors.BufferInvalidPosition)
	}
	if _, err := reader.ReadBytesAt(math.MaxUint, 1); err == nil {
		t.Fatal("ReadBytesAt(max, 1) error = nil")
	} else {
		assertBufferErrorCode(t, err, oracleErrors.BufferInvalidRange)
	}
	if _, err := reader.ReadBytesAt(0, 0); err != nil {
		t.Fatalf("ReadBytesAt(0, 0) error = %v", err)
	}
	for _, read := range []func() error{
		func() error { _, err := reader.ReadUB2(); return err },
		func() error { _, err := reader.ReadUB4(); return err },
		func() error { _, err := reader.ReadUB8(); return err },
	} {
		assertBufferErrorCode(t, read(), oracleErrors.BufferUnderflow)
	}
	for _, read := range []func() error{
		func() error { _, err := reader.ReadUB1At(0); return err },
		func() error { _, err := reader.ReadUB2At(0); return err },
		func() error { _, err := reader.ReadUB4At(0); return err },
		func() error { _, err := reader.ReadUB8At(0); return err },
		func() error { _, err := reader.ReadSB1At(0); return err },
		func() error { _, err := reader.ReadSB2At(0); return err },
		func() error { _, err := reader.ReadSB4At(0); return err },
	} {
		assertBufferErrorCode(t, read(), oracleErrors.BufferInvalidRange)
	}

	writer := NewWriteBuffer(BIG_ENDIAN)
	writer.WriteUB1(1)
	if err := writer.WriteUB2At(0, 2); err == nil {
		t.Fatal("WriteUB2At() error = nil")
	} else {
		assertBufferErrorCode(t, err, oracleErrors.BufferInvalidRange)
	}
	if err := writer.WriteUB1At(math.MaxUint, 2); err == nil {
		t.Fatal("WriteUB1At(max) error = nil")
	} else {
		assertBufferErrorCode(t, err, oracleErrors.BufferInvalidRange)
	}
	for _, write := range []func() error{
		func() error { return writer.WriteUB4At(0, 0) },
		func() error { return writer.WriteUB8At(0, 0) },
		func() error { return writer.WriteSB1At(1, 0) },
		func() error { return writer.WriteSB2At(0, 0) },
		func() error { return writer.WriteSB4At(0, 0) },
		func() error { return writer.WriteBytesAt(1, B1Array{0}) },
	} {
		assertBufferErrorCode(t, write(), oracleErrors.BufferInvalidRange)
	}
}

// assertBufferErrorCode verifies that err has the expected driver error code.
func assertBufferErrorCode(t *testing.T, err error, want oracleErrors.ErrorCode) {
	t.Helper()
	var sqlErr oracleErrors.SQLError
	if !errors.As(err, &sqlErr) || sqlErr.ErrorCode() != string(want) {
		t.Fatalf("error = %v, want code %s", err, want)
	}
}
