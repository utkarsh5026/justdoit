package ordereddict

// OrderedDict represents a dictionary that maintains the order of keys.
// It is implemented using a map and a slice of keys.
type OrderedDict struct {
	dict map[string]interface{}
	keys []string
}

func New() *OrderedDict {
	return &OrderedDict{
		dict: make(map[string]interface{}),
		keys: []string{},
	}
}

// Set adds a key-value pair to the OrderedDict. If the key already exists, its value is updated.
// If the key does not exist, it is added to the keys slice to maintain order.
//
// Parameters:
// - key: The key to be added or updated in the dictionary.
// - value: The value associated with the key.
func (od *OrderedDict) Set(key string, value interface{}) {
	if _, exists := od.dict[key]; !exists {
		od.keys = append(od.keys, key)
	}
	od.dict[key] = value
}

// Get retrieves the value associated with the given key from the OrderedDict.
// It returns the value and a boolean indicating whether the key exists.
//
// Parameters:
// - key: The key to look up in the dictionary.
//
// Returns:
// - interface{}: The value associated with the key.
// - bool: True if the key exists, false otherwise.
func (od *OrderedDict) Get(key string) (interface{}, bool) {
	value, exists := od.dict[key]
	return value, exists
}

// Delete removes the key-value pair associated with the given key from the OrderedDict.
// If the key exists, it is removed from both the dictionary and the keys slice.
//
// Parameters:
// - key: The key to be removed from the dictionary.
func (od *OrderedDict) Delete(key string) {
	if _, exists := od.dict[key]; exists {
		delete(od.dict, key)
		for i, k := range od.keys {
			if k == key {
				od.keys = append(od.keys[:i], od.keys[i+1:]...)
				break
			}
		}
	}
}

// Keys returns a slice of keys in the OrderedDict.
func (od *OrderedDict) Keys() []string {
	return od.keys
}

// Len Values returns a slice of values in the OrderedDict.
func (od *OrderedDict) Len() int {
	return len(od.keys)
}

// Range iterates over the key-value pairs in the OrderedDict in the order of keys.
// The provided function f is called for each key-value pair. If f returns false, the iteration stops.
//
// Parameters:
//   - f: A function that takes a key and a value as arguments and returns a boolean.
//     If the function returns false, the iteration stops.
func (od *OrderedDict) Range(f func(key string, value interface{}) bool) {
	for _, key := range od.keys {
		if !f(key, od.dict[key]) {
			break
		}
	}
}
