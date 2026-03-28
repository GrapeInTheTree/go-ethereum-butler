package ethereum

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

// LookupSelector resolves a 4-byte function selector (e.g., "0xa9059cbb") to a
// human-readable function signature (e.g., "transfer(address,uint256)") using
// the OpenChain signature database. Returns empty string on lookup failure.
func LookupSelector(selector string) string {
	client := &http.Client{Timeout: 5 * time.Second}

	reqURL := "https://api.openchain.xyz/signature-database/v1/lookup?" +
		url.Values{"function": {selector}}.Encode()

	resp, err := client.Get(reqURL)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var result struct {
		OK     bool `json:"ok"`
		Result struct {
			Function map[string][]struct {
				Name string `json:"name"`
			} `json:"function"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || !result.OK {
		return ""
	}

	sigs, ok := result.Result.Function[selector]
	if !ok || len(sigs) == 0 {
		return ""
	}
	return sigs[0].Name
}

// LookupEventTopic resolves an event topic hash to event signature.
func LookupEventTopic(topic string) string {
	client := &http.Client{Timeout: 5 * time.Second}

	reqURL := "https://api.openchain.xyz/signature-database/v1/lookup?" +
		url.Values{"event": {topic}}.Encode()

	resp, err := client.Get(reqURL)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var result struct {
		OK     bool `json:"ok"`
		Result struct {
			Event map[string][]struct {
				Name string `json:"name"`
			} `json:"event"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || !result.OK {
		return ""
	}

	sigs, ok := result.Result.Event[topic]
	if !ok || len(sigs) == 0 {
		return ""
	}
	return sigs[0].Name
}

// Common selectors cached locally for zero-latency lookups
var commonSelectors = map[string]string{
	"0xa9059cbb": "transfer(address,uint256)",
	"0x095ea7b3": "approve(address,uint256)",
	"0x23b872dd": "transferFrom(address,address,uint256)",
	"0x70a08231": "balanceOf(address)",
	"0x18160ddd": "totalSupply()",
	"0x313ce567": "decimals()",
	"0x06fdde03": "name()",
	"0x95d89b41": "symbol()",
	"0xdd62ed3e": "allowance(address,address)",
}

// ResolveSelector tries local cache first, then OpenChain API.
func ResolveSelector(selector string) string {
	if name, ok := commonSelectors[selector]; ok {
		return name
	}
	return LookupSelector(selector)
}

// ResolveMethodName returns a display-friendly method name.
// If resolved, returns the function signature. Otherwise returns the raw selector.
func ResolveMethodName(selector string) string {
	if selector == "" || selector == "0x" {
		return "Transfer"
	}
	if name := ResolveSelector(selector); name != "" {
		// Extract just the function name for display
		return name
	}
	return selector
}

// ExtractFunctionName returns just the function name part from a full signature.
// e.g., "transfer(address,uint256)" → "transfer"
func ExtractFunctionName(sig string) string {
	for i, c := range sig {
		if c == '(' {
			return sig[:i]
		}
	}
	return sig
}
