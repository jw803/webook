package loggerx

func String(key, val string) Field {
	return Field{
		Key:   key,
		Value: val,
	}
}

func Int32(key string, val int32) Field {
	return Field{
		Key:   key,
		Value: val,
	}
}

func Int64(key string, val int64) Field {
	return Field{
		Key:   key,
		Value: val,
	}
}

func Map[T any](key string, val map[string]T) Field {
	return Field{
		Key:   key,
		Value: val,
	}
}

func Slice[T any](key string, val []T) Field {
	return Field{
		Key:   key,
		Value: val,
	}
}

func Error(err error) Field {
	return Field{
		Key:   "error",
		Value: err,
	}
}

func Any(key string, val any) Field {
	return Field{
		Key:   key,
		Value: val,
	}
}
