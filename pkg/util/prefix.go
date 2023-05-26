package util

var PM map[string]string

func StorePrefix(key, prefix string) {
	PM = map[string]string{
		"user": "",
		"task": "",
	}
	PM[key] = prefix
}
