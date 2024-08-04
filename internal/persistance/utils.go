package persistance

// usally this is another packae since it will likely be used in all the tests
func AddPointer[T any](data T) *T {
	return &data
}
