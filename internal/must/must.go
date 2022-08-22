package must

// NoError panics only if the error is not nil.
func NoError(err error) {
	if err != nil {
		panic(err)
	}
}
