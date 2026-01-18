package spellgrid

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrorStoreTypeAlreadyAdded   = errors.New("Already added type")
	ErrorStoreTypeNotAdded       = errors.New("Type not added")
	ErrorStoreNoDataAddedToField = errors.New("No data added to field")
)

type ErrorStoreTypeMismatch struct {
	Expected reflect.Type
	Actual   reflect.Type
}

func (e *ErrorStoreTypeMismatch) Error() string {
	return fmt.Sprintf("Type mismatch: expected %v, got %v", e.Expected, e.Actual)
}

type ErrorStoreOutOfBonds struct {
	Expected int
	Actual   int
}

func (e *ErrorStoreOutOfBonds) Error() string {
	return fmt.Sprintf("Index out of bonds: expected %d, got %d", e.Expected, e.Actual)
}

func ErrorStoreFieldAmountMismatch(expected int, actual int) string {
	return fmt.Sprintf("Field length mismatch: expected %d, got %d", expected, actual)
}

func ErrorStoreFieldMismatch(expected StoreFieldInfo, actual StoreFieldInfo) (bool, string) {
	errorString := ""
	match := true
	if expected.Name != actual.Name {
		errorString = fmt.Sprintf("Name mismatch: expected %s, got %s", expected.Name, actual.Name)
		match = false
	}
	if expected.Type != actual.Type {
		errorString = fmt.Sprintf("Type mismatch: expected %s, got %s", expected.Type, actual.Type)
		match = false
	}
	if expected.Amount != actual.Amount {
		errorString = fmt.Sprintf("Amount mismatch: expected %d, got %d", expected.Amount, actual.Amount)
		match = false
	}
	return match, errorString
}
