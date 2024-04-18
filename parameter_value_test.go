// Copyright (c) Paul Schick
// SPDX-License-Identifier: MPL-2.0

package oddsapi

import (
	"reflect"
	"testing"
)

func TestParamPointer(t *testing.T) {
	american := AmericanOddsFormat
	if !american.Valid() {
		t.Errorf("expected AmericanOddsFormat.Valid() to be true, got false")
	}

	ptr := ParamPointer(american)

	if reflect.ValueOf(ptr).Kind() != reflect.Ptr {
		t.Errorf("expected ptr to be a pointer")
	}

	var american2 OddsFormat = "american"
	if !american2.Valid() {
		t.Errorf("expected 'american' to be true, got false")
	}

	expectedType := "oddsapi.OddsFormat"
	resultType := reflect.TypeOf(american2).String()
	if resultType != expectedType {
		t.Errorf("expected type '%s', got '%s'", expectedType, resultType)
	}

	var failing OddsFormat = "not-valid"
	if failing.Valid() {
		t.Errorf("expected %s to be invalid, it was valid", string(failing))
	}

	failingStr := failing.String()
	resultType2 := reflect.TypeOf(failingStr).String()
	if resultType2 != "string" {
		t.Errorf("invalid OddsFormat String() should still be a string, it was %s", resultType2)
	}
}
