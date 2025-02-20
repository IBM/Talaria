package utils

// Keys returns the keys of the map m.
// The keys will be an indeterminate order.
//
// Copied over from https://cs.opensource.google/go/x/exp/+/39d4317d:maps/maps.go;l=10
// I don't want to include the whole exp package just for this function.
func MapKeys[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}
