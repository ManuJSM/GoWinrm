package utils

type Arguments = map[string]any

// optStringSlice es un helper para obtener []string desde map[string]any
func OptStringSlice(opts Arguments, key string, defaultVal []string) []string {
	if val, ok := opts[key]; ok {
		if slice, ok := val.([]string); ok {
			return slice
		}
	}
	return defaultVal
}

func OptOrDefault(options Arguments, key string, defaultValue any) any {
	if value, ok := options[key]; ok {
		return value
	}
	return defaultValue
}
