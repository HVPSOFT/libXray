package libXray

import (
	"encoding/json"

	"github.com/xtls/libxray/share"
)

// ConvertShareLinkToXrayJSON provides a simple String -> JSON API for Apple bindings.
func ConvertShareLinkToXrayJSON(link string) string {
	xrayJSON, err := share.ConvertShareLinkToXrayJSON(link)
	if err == nil {
		return xrayJSON
	}

	errorJSON, marshalErr := json.Marshal(map[string]string{
		"error": err.Error(),
	})
	if marshalErr != nil {
		return `{"error":"failed to encode error"}`
	}

	return string(errorJSON)
}

// ConvertShareLinkToXrayJSONWithXrayKnife provides a pure xray-knife String -> JSON API for Apple bindings.
func ConvertShareLinkToXrayJSONWithXrayKnife(link string) string {
	xrayJSON, err := share.ConvertShareLinkToXrayJSONWithXrayKnife(link)
	if err == nil {
		return xrayJSON
	}

	errorJSON, marshalErr := json.Marshal(map[string]string{
		"error": err.Error(),
	})
	if marshalErr != nil {
		return `{"error":"failed to encode error"}`
	}

	return string(errorJSON)
}
