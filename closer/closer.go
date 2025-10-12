package closer

func RecoverAndCleanup(cleanup func()) {
	if r := recover(); r != nil {
		if cleanup != nil {
			cleanup()
		}
		panic(r)
	}
}
