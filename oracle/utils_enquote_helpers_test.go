/*
** Copyright (c) 2026 Oracle and/or its affiliates.
**
** The Universal Permissive License (UPL), Version 1.0
 */

package oracle

import (
	"strings"
	"testing"

	oracleErrors "github.com/oracle/go-oracledb/v26/oracle/errors"
	oracleutils "github.com/oracle/go-oracledb/v26/oracle/utils"
)

func TestEnquoteLiteral(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "plain", in: "Hello", want: "'Hello'"},
		{name: "single quote", in: "G'Day", want: "'G''Day'"},
		{name: "already quoted", in: "'G''Day'", want: "'''G''''Day'''"},
		{name: "many quotes", in: "I'''M", want: "'I''''''M'"},
		{name: "empty", in: "", want: "''"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := oracleutils.EnquoteLiteral(tc.in)
			if got != tc.want {
				t.Fatalf("unexpected value: got %q want %q", got, tc.want)
			}
		})
	}
}

func TestEnquoteNCharLiteral(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "plain", in: "Hello", want: "N'Hello'"},
		{name: "single quote", in: "G'Day", want: "N'G''Day'"},
		{name: "already quoted", in: "'G''Day'", want: "N'''G''''Day'''"},
		{name: "many quotes", in: "I'''M", want: "N'I''''''M'"},
		{name: "already N prefixed text", in: "N'Hello'", want: "N'N''Hello'''"},
		{name: "empty", in: "", want: "N''"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := oracleutils.EnquoteNCharLiteral(tc.in)
			if got != tc.want {
				t.Fatalf("unexpected value: got %q want %q", got, tc.want)
			}
		})
	}
}

func TestIsSimpleIdentifier(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		in   string
		want bool
	}{
		{name: "empty", in: "", want: false},
		{name: "one char", in: "A", want: true},
		{name: "simple", in: "Hello", want: true},
		{name: "with underscore and digits", in: "A_12", want: true},
		{name: "starts with digit", in: "1abc", want: false},
		{name: "single quote", in: "G'Day", want: false},
		{name: "double quoted", in: `"Bruce Wayne"`, want: false},
		{name: "dollar sign", in: "GoodDay$", want: false},
		{name: "contains double quote", in: `Hello"World`, want: false},
		{name: "max length valid", in: "A" + strings.Repeat("b", 127), want: true},
		{name: "too long", in: "A" + strings.Repeat("b", 128), want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := oracleutils.IsSimpleIdentifier(tc.in)
			if got != tc.want {
				t.Fatalf("unexpected value: got %v want %v", got, tc.want)
			}
		})
	}
}

func TestEnquoteIdentifier(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		in      string
		want    string
		wantErr bool
	}{
		{name: "empty", in: "", wantErr: true},
		{name: "simple always quoted", in: "Hello", want: `"Hello"`, wantErr: false},
		{name: "not simple quoted", in: "G'Day", want: `"G'Day"`, wantErr: false},
		{name: "already quoted kept quoted", in: `"Bruce Wayne"`, want: `"Bruce Wayne"`, wantErr: false},
		{name: "special char quoted", in: "GoodDay$", want: `"GoodDay$"`, wantErr: false},
		{name: "contains double quote", in: `Hello"World`, wantErr: true},
		{name: "quoted with inner double quote", in: `"Hello"World"`, wantErr: true},
		{name: "contains null char", in: "ab\x00cd", wantErr: true},
		{name: "too long", in: "A" + strings.Repeat("b", 128), wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := oracleutils.EnquoteIdentifier(tc.in)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				serr, ok := err.(oracleErrors.SQLError)
				if !ok {
					t.Fatalf("expected SQLError, got %T: %v", err, err)
				}
				if serr.ErrorCode() != string(oracleErrors.InvalidIdentifier) {
					t.Fatalf("unexpected error code: got %s want %s", serr.ErrorCode(), oracleErrors.InvalidIdentifier)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("unexpected value: got %q want %q", got, tc.want)
			}
		})
	}
}
