package libXray

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertShareLinkToXrayJSONVmess(t *testing.T) {
	vmessPayload := `{"v":"2","ps":"Demo","add":"example.com","port":"443","id":"123e4567-e89b-12d3-a456-426614174000","aid":"0","net":"ws","type":"none","host":"example.com","path":"/ws","tls":"tls","sni":"example.com"}`
	link := "vmess://" + base64.StdEncoding.EncodeToString([]byte(vmessPayload))

	output := ConvertShareLinkToXrayJSON(link)

	var result map[string]any
	require.NoError(t, json.Unmarshal([]byte(output), &result))

	assert.NotContains(t, result, "error")

	outbounds, ok := result["outbounds"].([]any)
	require.True(t, ok)
	require.Len(t, outbounds, 1)

	outbound, ok := outbounds[0].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "vmess", outbound["protocol"])
	assert.NotContains(t, outbound, "sendThrough")
}

func TestConvertShareLinkToXrayJSONFallsBackForHy2(t *testing.T) {
	link := "hy2://password@example.com:443?sni=example.com#Fallback"

	output := ConvertShareLinkToXrayJSON(link)

	var result map[string]any
	require.NoError(t, json.Unmarshal([]byte(output), &result))

	assert.NotContains(t, result, "error")

	outbounds, ok := result["outbounds"].([]any)
	require.True(t, ok)
	require.Len(t, outbounds, 1)

	outbound, ok := outbounds[0].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "hysteria", outbound["protocol"])
	assert.NotContains(t, outbound, "sendThrough")
}

func TestConvertShareLinkToXrayJSONReturnsErrorJSON(t *testing.T) {
	output := ConvertShareLinkToXrayJSON("not-a-share-link")

	var result map[string]string
	require.NoError(t, json.Unmarshal([]byte(output), &result))
	assert.NotEmpty(t, result["error"])
}

func TestConvertShareLinkToXrayJSONWithXrayKnifeVmess(t *testing.T) {
	vmessPayload := `{"v":"2","ps":"KnifeOnly","add":"example.com","port":"443","id":"123e4567-e89b-12d3-a456-426614174000","aid":"0","net":"ws","type":"none","host":"example.com","path":"/ws","tls":"tls","sni":"example.com"}`
	link := "vmess://" + base64.StdEncoding.EncodeToString([]byte(vmessPayload))

	output := ConvertShareLinkToXrayJSONWithXrayKnife(link)

	var result map[string]any
	require.NoError(t, json.Unmarshal([]byte(output), &result))
	assert.NotContains(t, result, "error")

	outbounds, ok := result["outbounds"].([]any)
	require.True(t, ok)
	require.Len(t, outbounds, 1)

	outbound, ok := outbounds[0].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "vmess", outbound["protocol"])
	assert.NotContains(t, outbound, "sendThrough")
}

func TestConvertShareLinkToXrayJSONWithXrayKnifeDoesNotFallbackForHy2(t *testing.T) {
	output := ConvertShareLinkToXrayJSONWithXrayKnife("hy2://password@example.com:443?sni=example.com#NoFallback")

	var result map[string]string
	require.NoError(t, json.Unmarshal([]byte(output), &result))
	assert.NotEmpty(t, result["error"])
}
