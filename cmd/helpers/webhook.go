package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func PostToWebhook(ctx context.Context, url string, body map[string]interface{}) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read POST response body: %w", err)
	}
	if resp.StatusCode >= 400 {
		fmt.Println("POST request failed with status " + resp.Status + ": " + string(respBody))
		return fmt.Errorf("webhook request failed with status: %d", resp.StatusCode)
	}

	// if err := json.Unmarshal(respBody, &responseStruct); err != nil {
	// 	return fmt.Errorf("failed to unmarshal POST response: %w", err)
	// }

	return nil
}
