package translator

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Translate function to perform translation
func Translate(text, from, to string) string {
	// Prepare the URL for the API request
	apiURL := fmt.Sprintf("https://translate.googleapis.com/translate_a/single?client=gtx&sl=%s&tl=%s&dt=t&q=%s",
		from, to, url.QueryEscape(text))

	// Make the HTTP GET request
	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Sprintf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response is OK
	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("translation failed: %s", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("failed to read response body: %v", err)
	}

	// Parse the JSON response
	var jsonResponse []interface{}
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		return fmt.Sprintf("unexpected response format: %v", err)
	}

	// Extract the translated text
	if len(jsonResponse) > 0 {
		if translations, ok := jsonResponse[0].([]interface{}); ok && len(translations) > 0 {
			if translation, ok := translations[0].([]interface{}); ok && len(translation) > 0 {
				if str, ok := translation[0].(string); ok {
					return str
				}
			}
		}
	}

	return "unexpected response structure"
}
