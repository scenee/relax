package relfile

import (
	"fmt"
	"github.com/DHowett/go-plist"
	"os"
	"strings"
)

func prepareFile(out string) (f *os.File) {
	var err error
	// Remove an existing file
	if _, err = os.Stat(out); err == nil {
		os.Remove(out)
	}

	if f, err = os.OpenFile(out, os.O_RDWR|os.O_CREATE, 0600); err != nil {
		logger.Fatal(err)
	}

	return f
}

func getBundleID(infoPlist string) string {
	var (
		err     error
		decoder *plist.Decoder
		f       *os.File
		data    map[string]interface{}
	)

	f, err = os.Open(infoPlist)
	if err != nil {
		logger.Fatalf("open error: %v", err)
	} else {
		defer f.Close()
		decoder = plist.NewDecoder(f)
	}

	err = decoder.Decode(&data)
	if err != nil {
		logger.Fatalf("decode error: %v", err)
	}

	props, ok := data["ApplicationProperties"].(map[string]interface{})
	if ok {
		return props["CFBundleIdentifier"].(string)
	} else {
		return data["CFBundleIdentifier"].(string)
	}
}

func mergeMap(old map[string]interface{}, new map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range old {
		res[k] = v
	}
	for k, v := range new {
		if _v, ok := old[k]; ok {
			switch _v := _v.(type) {
			case map[string]interface{}:
				if v, ok := v.(map[string]interface{}); ok {
					v = mergeMap(_v, v)
					res[k] = v
					continue
				}
			}
		}
		res[k] = v
	}
	return res
}

func cleanupInterfaceArray(in []interface{}) []interface{} {
	res := make([]interface{}, len(in))
	for i, v := range in {
		res[i] = cleanupMapValue(v)
	}
	return res
}

func cleanupInterfaceMap(in map[interface{}]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range in {
		res[fmt.Sprintf("%v", k)] = cleanupMapValue(v)
	}
	return res
}

func cleanupMapValue(v interface{}) interface{} {
	switch v := v.(type) {
	case []interface{}:
		return cleanupInterfaceArray(v)
	case map[interface{}]interface{}:
		return cleanupInterfaceMap(v)
	default:
		return v
	}
}

func genSourceline(key, value string) string {
	k := strings.Join([]string{PREFIX, key}, "_")
	return fmt.Sprintf("export %v=\"%v\"\n", k, value)
}

func genSourceLine2(name string, key string, value interface{}) string {
	k := strings.Join([]string{PREFIX, name, key}, "_")
	return fmt.Sprintf("export %v=\"%v\"\n", k, value)
}
