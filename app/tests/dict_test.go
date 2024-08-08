package ordereddict_test

import (
	"github.com/utkarsh5026/justdoit/app/ordereddict"
	"reflect"
	"testing"
)

func TestNewOrderedDict(t *testing.T) {
	od := ordereddict.New()
	if od == nil {
		t.Fatal("New() returned nil")
	}
	if od.Len() != 0 {
		t.Errorf("New() dict not empty, got %d items", od.Len())
	}
}

func TestOrderedDict_Set(t *testing.T) {
	od := ordereddict.New()
	od.Set("key1", "value1")
	od.Set("key2", "value2")
	od.Set("key1", "newvalue1")

	if od.Len() != 2 {
		t.Errorf("Expected 2 items, got %d", od.Len())
	}
	value, exists := od.Get("key1")
	if !exists || value != "newvalue1" {
		t.Errorf("Expected 'newvalue1', got %v", value)
	}
}

func TestOrderedDict_Get(t *testing.T) {
	od := ordereddict.New()
	od.Set("key1", "value1")

	value, exists := od.Get("key1")
	if !exists {
		t.Error("Get() returned false for existing key")
	}
	if value != "value1" {
		t.Errorf("Expected 'value1', got %v", value)
	}

	_, exists = od.Get("nonexistent")
	if exists {
		t.Error("Get() returned true for non-existent key")
	}
}

func TestOrderedDict_Delete(t *testing.T) {
	od := ordereddict.New()
	od.Set("key1", "value1")
	od.Set("key2", "value2")
	od.Delete("key1")

	if od.Len() != 1 {
		t.Errorf("Expected 1 item, got %d", od.Len())
	}
	keys := od.Keys()
	if len(keys) != 1 || keys[0] != "key2" {
		t.Errorf("Expected ['key2'], got %v", keys)
	}
}

func TestOrderedDict_Keys(t *testing.T) {
	od := ordereddict.New()
	od.Set("key1", "value1")
	od.Set("key2", "value2")

	keys := od.Keys()
	expectedKeys := []string{"key1", "key2"}
	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Errorf("Expected %v, got %v", expectedKeys, keys)
	}
}

func TestOrderedDict_Len(t *testing.T) {
	od := ordereddict.New()
	od.Set("key1", "value1")
	od.Set("key2", "value2")

	if od.Len() != 2 {
		t.Errorf("Expected length 2, got %d", od.Len())
	}
}

func TestOrderedDict_Range(t *testing.T) {
	od := ordereddict.New()
	od.Set("key1", "value1")
	od.Set("key2", "value2")
	od.Set("key3", "value3")

	expectedKeys := []string{"key1", "key2", "key3"}
	expectedValues := []string{"value1", "value2", "value3"}
	index := 0

	od.Range(func(key string, value interface{}) bool {
		if key != expectedKeys[index] {
			t.Errorf("Expected key %s, got %s", expectedKeys[index], key)
		}
		if value != expectedValues[index] {
			t.Errorf("Expected value %s, got %v", expectedValues[index], value)
		}
		index++
		return true
	})

	if index != 3 {
		t.Errorf("Range didn't iterate over all elements, got %d iterations", index)
	}
}

func TestOrderedDict_RangeBreak(t *testing.T) {
	od := ordereddict.New()
	od.Set("key1", "value1")
	od.Set("key2", "value2")
	od.Set("key3", "value3")

	count := 0
	od.Range(func(key string, value interface{}) bool {
		count++
		return count < 2
	})

	if count != 2 {
		t.Errorf("Range didn't break as expected, got %d iterations", count)
	}
}
