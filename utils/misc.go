package utils

func Slice2String(arr []*string) string {
	buf := make([]byte, 0)
	var first bool = true
	for _, s := range arr {
		if !first {
			buf = append(buf, []byte(", ")...)
		}
		buf = append(buf, []byte(*s)...)

		first = false
	}
	return string(buf)
}
