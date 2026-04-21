package features

var FEATURE_SEND_EMAIL = "SEND_EMAIL"

// TODO: use config file
var FEATURE_MAP = map[string]bool{
	FEATURE_SEND_EMAIL: false,
}

func HasFeature(feature string) bool {
	return FEATURE_MAP[feature]
}
