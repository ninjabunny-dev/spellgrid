package spellgrid

import (
	"errors"
	"testing"
)

func TestStore_AddingStrings(t *testing.T) {
	store := NewStore()

	store.AddField("strings", "")
	store.AddValues("strings", "hello", "world")

	info, err := store.GetFieldInfo("strings")

	if err != nil {
		t.Error(err)
	}

	if info.Amount != 2 {
		t.Error(ErrorStoreFieldAmountMismatch(2, info.Amount))
	}
}

func TestStore_AddingBools(t *testing.T) {
	store := NewStore()

	store.AddField("bools", false)
	store.AddValues("bools", true, false)
	info, err := store.GetFieldInfo("bools")

	if err != nil {
		t.Error(err)
	}

	if info.Amount != 2 {
		t.Error(ErrorStoreFieldAmountMismatch(2, info.Amount))
	}
}

func TestStore_AddingInts(t *testing.T) {
	store := NewStore()

	store.AddField("ints", 0)
	store.AddValues("ints", 1, 2)
	info, err := store.GetFieldInfo("ints")

	if err != nil {
		t.Error(err)
	}

	if info.Amount != 2 {
		t.Error(ErrorStoreFieldAmountMismatch(2, info.Amount))
	}
}

func TestStore_AddingFloats(t *testing.T) {
	store := NewStore()

	store.AddField("floats", 0.0)
	store.AddValues("floats", 1.2, 3.4)

	info, err := store.GetFieldInfo("floats")

	if err != nil {
		t.Error(err)
	}

	if info.Amount != 2 {
		t.Error(ErrorStoreFieldAmountMismatch(2, info.Amount))
	}
}

func TestStore_AddingCustomTypes(t *testing.T) {
	type customType struct {
		name string
	}

	store := NewStore()

	store.AddField("customType", customType{})
	store.AddValues("customType", customType{"foo"}, customType{"bar"})

	info, err := store.GetFieldInfo("customType")

	if err != nil {
		t.Error(err)
	}

	if info.Amount != 2 {
		t.Error(ErrorStoreFieldAmountMismatch(2, info.Amount))
	}
}

func TestStore_AddingValuesMultiableTimes(t *testing.T) {
	store := NewStore()

	store.AddField("ints", 0)
	store.AddValues("ints", 1, 2)
	store.AddValues("ints", 1, 2)
	info, err := store.GetFieldInfo("ints")

	if err != nil {
		t.Error(err)
	}

	if info.Amount != 4 {
		t.Error(ErrorStoreFieldAmountMismatch(2, info.Amount))
	}
}

func TestStore_AddingMixedTypes(t *testing.T) {
	store := NewStore()

	store.AddField("mixed", 0)
	err := store.AddValues("mixed", 1, "two")

	var typeError *ErrorStoreTypeMismatch

	if !errors.As(err, &typeError) {
		t.Error(err)
	}
}

func TestStore_GettingFieldInfo(t *testing.T) {
	testInfo := StoreFieldInfo{
		Name:   "ints",
		Type:   "int",
		Amount: 2,
	}

	store := NewStore()
	store.AddField(testInfo.Name, 1)
	store.AddValues(testInfo.Name, 1, 2)

	info, _ := store.GetFieldInfo(testInfo.Name)

	match, msg := ErrorStoreFieldMismatch(testInfo, info)

	if !match {
		t.Error(msg)
	}
}

func TestStore_GettingMultipleFieldsInfo(t *testing.T) {
	testInfo := []StoreFieldInfo{
		{Name: "ints", Type: "int", Amount: 2},
		{Name: "strings", Type: "string", Amount: 2},
		{Name: "floats", Type: "float64", Amount: 2},
	}

	store := NewStore()
	store.AddField(testInfo[0].Name, 1)
	store.AddField(testInfo[1].Name, "")

	store.AddValues(testInfo[0].Name, 1, 2)
	store.AddValues(testInfo[1].Name, "one", "two")

	info, _ := store.GetMultiFieldInfo(testInfo[0].Name, testInfo[1].Name)

	for index, field := range info {
		if ok, msg := ErrorStoreFieldMismatch(testInfo[index], field); !ok {
			t.Error(msg)
		}
	}
}

func TestStore_GettingAllFieldsInfo(t *testing.T) {
	testInfo := []StoreFieldInfo{
		{Name: "ints", Type: "int", Amount: 2},
		{Name: "strings", Type: "string", Amount: 2},
		{Name: "floats", Type: "float64", Amount: 2},
	}

	store := NewStore()
	store.AddField(testInfo[0].Name, 1)
	store.AddField(testInfo[1].Name, "")
	store.AddField(testInfo[2].Name, 1.0)

	store.AddValues(testInfo[0].Name, 1, 2)
	store.AddValues(testInfo[1].Name, "one", "two")
	store.AddValues(testInfo[2].Name, 1.0, 2.0)

	info, err := store.GetMultiFieldInfo()

	if err != nil {
		t.Errorf("%s", err.Error())
	}

	for index, field := range info {
		if ok, msg := ErrorStoreFieldMismatch(testInfo[index], field); !ok {
			t.Error(msg)
		}
	}
}

func TestStore_GettingAllFieldsInfoWithNoAdded(t *testing.T) {
	testInfo := []StoreFieldInfo{}

	store := NewStore()

	info, err := store.GetMultiFieldInfo()

	if err != nil {
		t.Errorf("%s", err.Error())
	}

	for index, field := range info {
		if ok, msg := ErrorStoreFieldMismatch(testInfo[index], field); !ok {
			t.Error(msg)
		}
	}
}

func TestStore_GettingNotAddedFieldInfo(t *testing.T) {
	store := NewStore()
	_, err := store.GetFieldInfo("strings")

	if !errors.Is(err, ErrorStoreTypeNotAdded) {
		t.Error(err)
	}
}

func TestStore_GettingNotAddedMultipleFieldsInfo(t *testing.T) {
	store := NewStore()
	_, err := store.GetMultiFieldInfo("strings", "ints")

	if !errors.Is(err, ErrorStoreTypeNotAdded) {
		t.Error(err)
	}
}

func TestStore_GettingSomeNotAddedMultipleFieldsInfo(t *testing.T) {
	store := NewStore()
	store.AddField("strings", "")

	_, err := store.GetMultiFieldInfo("strings", "ints")

	if !errors.Is(err, ErrorStoreTypeNotAdded) {
		t.Error(err)
	}
}

func TestStore_AddingDuplicateFields(t *testing.T) {
	store := NewStore()
	store.AddField("strings", "")
	err := store.AddField("strings", "")

	if !errors.Is(err, ErrorStoreTypeAlreadyAdded) {
		t.Error(err)
	}
}

func TestStore_GettingFieldItemThatIsOutBounds(t *testing.T) {
	store := NewStore()
	store.AddField("strings", "")

	_, err := store.GetAtIndex("strings", 1)

	var testError *ErrorStoreOutOfBonds

	if errors.Is(err, testError) {
		t.Error(err)
	}
}

func TestStore_GettingFieldItemThatIsOutBoundsNegative(t *testing.T) {
	store := NewStore()
	store.AddField("strings", "")

	_, err := store.GetAtIndex("strings", -1)

	var testError *ErrorStoreOutOfBonds

	if errors.Is(err, testError) {
		t.Error(err)
	}
}

func TestStore_GettingFieldItemThatFieldIsNotAdded(t *testing.T) {
	store := NewStore()

	_, err := store.GetAtIndex("strings", 0)

	if !errors.Is(err, ErrorStoreTypeNotAdded) {
		t.Error(err)
	}
}

func TestStore_GettingAtIndex(t *testing.T) {
	store := NewStore()

	testValue := 1

	store.AddField("ints", 1)
	store.AddValues("ints", testValue)

	value, err := store.GetAtIndex("ints", 0)

	if err != nil {
		t.Error(err)
	}

	intValue, ok := value.(int)

	if ok && intValue != testValue {
		t.Errorf("Expected %d, got %d", testValue, value)
	}
}

func TestStore_EditAtFieldIndex(t *testing.T) {
	store := NewStore()

	initValue := 1
	testValue := 2
	testField := "ints"

	store.AddField("ints", 1)
	store.AddValues("ints", initValue)

	store.EditAtIndex(testField, 0, testValue)

	value, err := store.GetAtIndex("ints", 0)

	if err != nil {
		t.Error(err)
	}

	intValue, ok := value.(int)

	if ok && intValue != testValue {
		t.Errorf("Expected %d, got %d", testValue, value)
	}
}

func TestStore_EditAtFieldIndexWithWrongType(t *testing.T) {
	store := NewStore()
	store.AddField("ints", 1)
	store.AddValues("ints", 1)

	err := store.EditAtIndex("ints", 0, "test")

	var typeError *ErrorStoreTypeMismatch

	if !errors.As(err, &typeError) {
		t.Error(err)
	}
}

func TestStore_EditAtFieldIndexWithNotAddedField(t *testing.T) {
	store := NewStore()

	err := store.EditAtIndex("strings", 0, "test")

	if !errors.Is(err, ErrorStoreTypeNotAdded) {
		t.Error(err)
	}
}

func TestStore_EditFieldItemThatIsOutBounds(t *testing.T) {
	store := NewStore()
	store.AddField("strings", "")

	err := store.EditAtIndex("strings", -1, "test")

	var testError *ErrorStoreOutOfBonds

	if errors.Is(err, testError) {
		t.Error(err)
	}
}

func TestStore_EditFieldItemThatIsOutBoundsNegative(t *testing.T) {
	store := NewStore()
	store.AddField("strings", "")

	err := store.EditAtIndex("strings", 1, "test")

	var testError *ErrorStoreOutOfBonds

	if errors.Is(err, testError) {
		t.Error(err)
	}
}

func TestStore_RemoveAtField(t *testing.T) {
	store := NewStore()
	store.AddField("strings", "")
	store.AddValues("strings", "one")

	value, err := store.GetFieldInfo("strings")

	if err != nil {
		t.Error(err)
	}

	if value.Name != "strings" {
		t.Errorf("Field not added")
	}

	store.RemoveField("strings")

	values, err := store.GetAllFieldInfo()

	if err != nil {
		t.Error(err)
	}

	if len(values) != 0 {
		t.Errorf("Field not removed")
	}
}

func TestStore_RemoveAtFieldStartIndex(t *testing.T) {
	testArray := []int{1, 2, 3, 4, 5}

	store := NewStore()
	store.AddField("ints", 0)

	for _, value := range testArray {
		store.AddValues("ints", value)
	}

	err := store.RemoveAtIndex("ints", 0)

	if err != nil {
		t.Error(err)
	}

	testArray = testArray[1:]

	for index, _ := range testArray {
		value, err := store.GetAtIndex("ints", index)
		if err != nil {
			t.Errorf("%v: %s", index, err)
		} else if value != testArray[index] {
			t.Errorf("Value expected %d, got %d", testArray[index], value)
		}
	}

	info, err := store.GetFieldInfo("ints")

	if err != nil {
		t.Error(err)
	}

	if info.Amount != len(testArray) {
		t.Errorf("Lenght of array is not correct, Expected %d, got %d", len(testArray), info.Amount)
	}
}

func TestStore_RemoveAtFieldIndexMiddle(t *testing.T) {
	testArray := []int{1, 2, 3, 4, 5}

	store := NewStore()
	store.AddField("ints", 0)

	for _, value := range testArray {
		store.AddValues("ints", value)
	}

	err := store.RemoveAtIndex("ints", 2)

	if err != nil {
		t.Error(err)
	}

	newArray := testArray[0:2]
	newArray = append(newArray, testArray[3:]...)

	for index, _ := range newArray {
		value, err := store.GetAtIndex("ints", index)

		if err != nil {
			t.Errorf("%v: %s", index, err)
		} else if value != newArray[index] {
			t.Errorf("Value expected %d, got %d", testArray[index], value)
		}
	}

	info, err := store.GetFieldInfo("ints")

	if err != nil {
		t.Error(err)
	}

	if info.Amount != len(newArray) {
		t.Errorf("Length of array is not correct, Expected %d, got %d", len(newArray), info.Amount)
	}
}

func TestStore_RemoveAtFieldIndexEnd(t *testing.T) {
	testArray := []int{1, 2, 3, 4, 5}

	store := NewStore()
	store.AddField("ints", 0)

	for _, value := range testArray {
		store.AddValues("ints", value)
	}

	err := store.RemoveAtIndex("ints", 4)

	if err != nil {
		t.Error(err)
	}

	testArray = testArray[:4]

	for index, _ := range testArray {
		value, err := store.GetAtIndex("ints", index)

		if err != nil {
			t.Errorf("%v: %s", index, err)
		} else if value != testArray[index] {
			t.Errorf("Value expected %d, got %d", testArray[index], value)
		}
	}

	info, err := store.GetFieldInfo("ints")

	if err != nil {
		t.Error(err)
	}

	if info.Amount != len(testArray) {
		t.Errorf("Length of array is not correct, Expected %d, got %d", len(testArray), info.Amount)
	}
}

func TestStore_RemoveAtFieldIndexWithNotAddedField(t *testing.T) {
	store := NewStore()

	err := store.RemoveAtIndex("strings", 0)

	if !errors.Is(err, ErrorStoreTypeNotAdded) {
		t.Error(err)
	}
}

func TestStore_RemoveFieldItemThatIsOutBounds(t *testing.T) {
	store := NewStore()
	store.AddField("strings", "")

	err := store.RemoveAtIndex("strings", 1)

	var testError *ErrorStoreOutOfBonds

	if errors.Is(err, testError) {
		t.Error(err)
	}
}

func TestStore_RemoveFieldItemThatIsOutBoundsNegative(t *testing.T) {
	store := NewStore()
	store.AddField("strings", "")

	err := store.RemoveAtIndex("strings", -1)

	var testError *ErrorStoreOutOfBonds

	if errors.Is(err, testError) {
		t.Error(err)
	}
}

func TestStore_RemoveFieldNotAdded(t *testing.T) {
	store := NewStore()

	err := store.RemoveField("strings")

	if !errors.Is(err, ErrorStoreTypeNotAdded) {
		t.Error(err)
	}
}

func TestStore_GettingFirstFieldItem(t *testing.T) {
	firstValue := 1
	secondValue := 2

	store := NewStore()
	store.AddField("ints", 0)

	store.AddValues("ints", firstValue, secondValue)

	value, err := store.GetFirst("ints")

	if err != nil {
		t.Error(err)
	}

	if value != firstValue {
		t.Errorf("Expected %d, got %d", firstValue, value)
	}
}

func TestStore_GettingLastFieldItem(t *testing.T) {
	firstValue := 1
	secondValue := 2

	store := NewStore()
	store.AddField("ints", 0)

	store.AddValues("ints", firstValue, secondValue)

	value, err := store.GetLast("ints")

	if err != nil {
		t.Error(err)
	}

	if value != secondValue {
		t.Errorf("Expected %d, got %d", firstValue, value)
	}
}

func TestStore_GettingFirstItemFieldNotAdded(t *testing.T) {
	store := NewStore()

	_, err := store.GetFirst("strings")

	if !errors.Is(err, ErrorStoreTypeNotAdded) {
		t.Error(err)
	}
}

func TestStore_GettingLastItemFieldNotAdded(t *testing.T) {
	store := NewStore()

	_, err := store.GetFirst("strings")

	if !errors.Is(err, ErrorStoreTypeNotAdded) {
		t.Error(err)
	}
}

func TestStore_GettingFirstFieldWithNoItems(t *testing.T) {
	store := NewStore()
	store.AddField("strings", "")

	_, err := store.GetFirst("strings")

	if !errors.Is(err, ErrorStoreNoDataAddedToField) {
		t.Error(err)
	}
}

func TestStore_GettingLastFieldWithNoItems(t *testing.T) {
	store := NewStore()
	store.AddField("strings", "")

	_, err := store.GetFirst("strings")

	if !errors.Is(err, ErrorStoreNoDataAddedToField) {
		t.Error(err)
	}
}
