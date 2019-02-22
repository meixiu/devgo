package helper

func ArgsInt(d int, args ...int) int {
	if len(args) > 0 {
		return args[0]
	}
	return d
}

func ArgsString(d string, args ...string) string {
	if len(args) > 0 {
		return args[0]
	}
	return d
}
