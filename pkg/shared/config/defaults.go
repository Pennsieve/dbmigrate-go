package config

type DefaultSettings map[string]string

func NewDefaultSettings() DefaultSettings {
	return make(DefaultSettings, 20)
}
