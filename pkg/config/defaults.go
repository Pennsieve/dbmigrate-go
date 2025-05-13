package config

type DefaultSettings map[string]string

func NewDefaultSettings() DefaultSettings {
	return make(DefaultSettings, 20)
}

func (s DefaultSettings) getWithFallback(key string, fallback string) string {
	if value, present := s[key]; present {
		return value
	}
	return fallback
}

func (s DefaultSettings) get(key string) string {
	return s.getWithFallback(key, "")
}
