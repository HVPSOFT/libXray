package share

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	xrayknife "github.com/xtls/libxray/third_party/xrayknife/pkg/core/xray"
	"github.com/xtls/xray-core/infra/conf"
)

// ConvertShareLinkToXrayJSON converts a single share link to a cleaned Xray JSON string.
func ConvertShareLinkToXrayJSON(link string) (string, error) {
	xrayConfig, err := ConvertShareLinksToXrayJson(link)
	if err != nil {
		return "", err
	}
	stripRunnableOutboundFields(xrayConfig)
	return marshalCleanJSON(xrayConfig)
}

// ConvertShareLinkToXrayJSONWithXrayKnife converts a single share link to Xray JSON
// using only the vendored xray-knife parser.
func ConvertShareLinkToXrayJSONWithXrayKnife(link string) (string, error) {
	xrayConfig, err := convertShareLinkWithXrayKnifeOnly(link)
	if err != nil {
		return "", err
	}
	return marshalCleanJSON(xrayConfig)
}

func convertSingleShareLinkWithXrayKnife(link string) (*conf.Config, error) {
	xrayConfig, err := convertShareLinkWithXrayKnifeOnly(link)
	if err != nil {
		return nil, err
	}

	if len(xrayConfig.OutboundConfigs) == 0 {
		return xrayConfig, nil
	}

	outbound := &xrayConfig.OutboundConfigs[0]
	if remark := extractXrayKnifeRemark(link); remark != "" {
		setOutboundName(outbound, remark)
	}

	return xrayConfig, nil
}

func convertShareLinkWithXrayKnifeOnly(link string) (*conf.Config, error) {
	xrayCore := xrayknife.NewXrayService(false, false)

	outboundProtocol, err := xrayCore.CreateProtocol(strings.TrimSpace(link))
	if err != nil {
		return nil, err
	}

	xrayOutbound, ok := outboundProtocol.(xrayknife.Protocol)
	if !ok {
		return nil, fmt.Errorf("xray-knife returned an unsupported protocol type")
	}

	if err := xrayOutbound.Parse(); err != nil {
		return nil, err
	}

	outbound, err := xrayOutbound.BuildOutboundDetourConfig(false)
	if err != nil {
		return nil, err
	}

	return &conf.Config{
		OutboundConfigs: []conf.OutboundDetourConfig{*outbound},
	}, nil
}

func extractXrayKnifeRemark(link string) string {
	xrayCore := xrayknife.NewXrayService(false, false)

	outboundProtocol, err := xrayCore.CreateProtocol(strings.TrimSpace(link))
	if err != nil {
		return ""
	}

	xrayOutbound, ok := outboundProtocol.(xrayknife.Protocol)
	if !ok {
		return ""
	}

	if err := xrayOutbound.Parse(); err != nil {
		return ""
	}

	return strings.TrimSpace(xrayOutbound.ConvertToGeneralConfig().Remark)
}

func marshalCleanJSON(v any) (string, error) {
	verboseBytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	var genericJSON any
	if err := json.Unmarshal(verboseBytes, &genericJSON); err != nil {
		return "", err
	}

	cleanedJSON := removeEmptyJSONValues(genericJSON)
	if cleanedJSON == nil {
		cleanedJSON = map[string]any{}
	}

	cleanBytes, err := json.MarshalIndent(cleanedJSON, "", "  ")
	if err != nil {
		return "", err
	}

	return string(cleanBytes), nil
}

func stripRunnableOutboundFields(xrayConfig *conf.Config) {
	if xrayConfig == nil {
		return
	}

	for i := range xrayConfig.OutboundConfigs {
		xrayConfig.OutboundConfigs[i].SendThrough = nil
	}
}

func removeEmptyJSONValues(data any) any {
	if data == nil {
		return nil
	}

	val := reflect.ValueOf(data)

	switch val.Kind() {
	case reflect.Map:
		cleanMap := make(map[string]any)
		for _, key := range val.MapKeys() {
			cleanedValue := removeEmptyJSONValues(val.MapIndex(key).Interface())
			if cleanedValue != nil {
				cleanMap[key.String()] = cleanedValue
			}
		}
		if len(cleanMap) == 0 {
			return nil
		}
		return cleanMap
	case reflect.Slice:
		if val.Len() == 0 {
			return nil
		}
		cleanSlice := make([]any, 0, val.Len())
		for i := range val.Len() {
			cleanedValue := removeEmptyJSONValues(val.Index(i).Interface())
			if cleanedValue != nil {
				cleanSlice = append(cleanSlice, cleanedValue)
			}
		}
		if len(cleanSlice) == 0 {
			return nil
		}
		return cleanSlice
	case reflect.Ptr, reflect.Interface:
		if val.IsNil() {
			return nil
		}
		return removeEmptyJSONValues(val.Elem().Interface())
	case reflect.String:
		if val.String() == "" {
			return nil
		}
	case reflect.Bool:
		if !val.Bool() {
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val.Int() == 0 {
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if val.Uint() == 0 {
			return nil
		}
	case reflect.Float32, reflect.Float64:
		if val.Float() == 0 {
			return nil
		}
	}

	return data
}
