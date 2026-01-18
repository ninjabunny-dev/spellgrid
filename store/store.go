package spellgrid

import (
	"reflect"
	"sync"
)

type StoreFieldInfo struct {
	Name   string
	Type   string
	Amount int
}

type Store struct {
	data map[string]reflect.Value
	mu   sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]reflect.Value),
		mu:   sync.RWMutex{},
	}
}

func (s *Store) boundsCheck(fieldName string, index int) (bool, error) {
	if s.data[fieldName].Len() <= index || index < 0 {
		return false, &ErrorStoreOutOfBonds{Expected: s.data[fieldName].Len(), Actual: index}
	}
	return true, nil
}

func (s *Store) AddField(name string, typ interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[name]; ok {
		return ErrorStoreTypeAlreadyAdded
	}

	sliceType := reflect.SliceOf(reflect.TypeOf(typ))
	slice := reflect.MakeSlice(sliceType, 0, 1)
	s.data[name] = slice

	return nil
}

func (s *Store) AddValues(fieldname string, values ...interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, ok := s.data[fieldname]

	if !ok {
		return ErrorStoreTypeNotAdded
	}

	storedType := data.Type().Elem()

	for _, value := range values {
		if reflect.TypeOf(value) != storedType {
			return &ErrorStoreTypeMismatch{storedType, reflect.TypeOf(value)}
		}
	}

	for _, value := range values {
		data = reflect.Append(data, reflect.ValueOf(value))
	}
	s.data[fieldname] = data

	return nil
}

func (s *Store) GetFieldInfo(fieldname string) (StoreFieldInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, ok := s.data[fieldname]

	if !ok {
		return StoreFieldInfo{}, ErrorStoreTypeNotAdded
	}

	storedType := data.Type().Elem()

	return StoreFieldInfo{
		Name:   fieldname,
		Type:   storedType.String(),
		Amount: data.Len(),
	}, nil
}

func (s *Store) GetMultiFieldInfo(fieldnames ...string) ([]StoreFieldInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	fieldData := make([]StoreFieldInfo, 0, 1)

	for _, fieldname := range fieldnames {
		info, err := s.GetFieldInfo(fieldname)
		if err != nil {
			return nil, err
		}
		fieldData = append(fieldData, info)
	}

	return fieldData, nil
}

func (s *Store) GetAllFieldInfo() ([]StoreFieldInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	fieldData := make([]StoreFieldInfo, 0, 1)

	for fieldname, _ := range s.data {
		info, err := s.GetFieldInfo(fieldname)
		if err != nil {
			return nil, err
		}
		fieldData = append(fieldData, info)
	}

	return fieldData, nil
}

func (s *Store) GetAtIndex(fieldName string, index int) (interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, ok := s.data[fieldName]; !ok {
		return nil, ErrorStoreTypeNotAdded
	}

	if ok, err := s.boundsCheck(fieldName, index); !ok {
		return nil, err
	}

	return s.data[fieldName].Index(index).Interface(), nil
}

func (s *Store) EditAtIndex(fieldName string, index int, value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[fieldName]; !ok {
		return ErrorStoreTypeNotAdded
	}

	if ok, err := s.boundsCheck(fieldName, index); !ok {
		return err
	}

	storedType := s.data[fieldName].Type().Elem()

	if storedType != reflect.TypeOf(value) {
		return &ErrorStoreTypeMismatch{storedType, reflect.TypeOf(value)}
	}

	s.data[fieldName].Index(index).Set(reflect.ValueOf(value))

	return nil
}

func (s *Store) RemoveAtIndex(fieldName string, index int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[fieldName]; !ok {
		return ErrorStoreTypeNotAdded
	}

	if ok, err := s.boundsCheck(fieldName, index); !ok {
		return err
	}

	slice := s.data[fieldName]

	for i := index; i < slice.Len()-1; i++ {
		slice.Index(i).Set(slice.Index(i + 1))
	}

	s.data[fieldName] = slice.Slice(0, slice.Len()-1)

	return nil
}

func (s *Store) RemoveField(fieldName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[fieldName]; !ok {
		return ErrorStoreTypeNotAdded
	}

	delete(s.data, fieldName)

	return nil
}

func (s *Store) GetFirst(fieldName string) (interface{}, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[fieldName]; !ok {
		return nil, ErrorStoreTypeNotAdded
	}

	if s.data[fieldName].Len() == 0 {
		return nil, ErrorStoreNoDataAddedToField
	}

	return s.data[fieldName].Index(0).Interface(), nil
}

func (s *Store) GetLast(fieldName string) (interface{}, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[fieldName]; !ok {
		return nil, ErrorStoreTypeNotAdded
	}

	if s.data[fieldName].Len() == 0 {
		return nil, ErrorStoreNoDataAddedToField
	}
	return s.data[fieldName].Index(s.data[fieldName].Len() - 1).Interface(), nil
}
