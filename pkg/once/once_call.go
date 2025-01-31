package once

import (
	"reflect"
	"sync"
)

var (
	instanced = make(map[string]any)
	locks     = sync.Map{}
)

// Call returns a function that calls the given function only once. */
func Call[T any](singletonFunc func() (T, error)) (T, error) {
	var (
		instanceKey = reflect.TypeOf((*T)(nil)).Elem().String()
		lock        = getLockForKey(instanceKey)
	)

	lock.Lock()
	defer lock.Unlock()
	return getInstanceOrCreate[T](instanceKey, singletonFunc)
}

func CallFlush() {
	instanced = make(map[string]interface{})
}

// getLockForKey returns the lock for the given key.
func getLockForKey(key string) *sync.Mutex {
	v, _ := locks.LoadOrStore(key, &sync.Mutex{})
	return v.(*sync.Mutex)
}

// getInstanceOrCreate returns an instance of the given type. If the instance
// does not exist, it is created using the given function.
func getInstanceOrCreate[T any](instanceKey string, singletonFunc func() (T, error)) (T, error) {
	if instance, ok := instanced[instanceKey]; ok {
		if i, ok := instance.(T); ok {
			return i, nil
		}
	}

	uc, err := singletonFunc()
	if err != nil {
		return genericEmptyInstance[T](), err
	}

	instanced[instanceKey] = uc
	return uc, nil
}

// genericEmptyInstance returns an empty instance of the given type.
func genericEmptyInstance[T any]() T {
	return reflect.Zero(reflect.TypeOf((*T)(nil)).Elem()).Interface().(T)
}
