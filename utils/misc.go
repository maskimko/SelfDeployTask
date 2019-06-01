package utils

func Slice2String(arr []*string) string {
	buf := make([]byte, 0)
	var first bool = true
	for _, s := range arr {
		buf = append(buf, []byte(*s)...)
		if !first {
			buf = append(buf, []byte(", ")...)
		}
		first = false
	}
	return string(buf)
}
