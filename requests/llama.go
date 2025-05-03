package requests

import (
	"encoding/json"
	"fmt"
	"os"

	"bytes"

	"github.com/arcedo/financial-ai-backend/types"
	"github.com/arcedo/financial-ai-backend/utils"
)

type ProfileWrapper struct {
	Profile types.Insights `json:"profile"`
}

func RequestUpdateUserProfile(data types.Insights) (types.UserProfile, error) {
	url := fmt.Sprintf("%s/update?id=%s", os.Getenv("LLM_HOST"), data.User.ID.Hex())
	headers := map[string]string{
		"x-api-key":    os.Getenv("LLM_API_KEY"),
		"Content-Type": "application/json",
	}

	wrapped := ProfileWrapper{Profile: data}

	bodyBytes, err := json.Marshal(wrapped)
	if err != nil {
		return types.UserProfile{}, fmt.Errorf("failed to marshal data: %v", err)
	}

	// Wrap the byte slice into an io.Reader
	bodyReader := bytes.NewReader(bodyBytes)

	respBody, err := utils.MakeHTTPRequest("POST", url, headers, bodyReader)
	if err != nil {
		return types.UserProfile{}, fmt.Errorf("failed to make request: %v", err)
	}

	var result types.UserProfile
	if err := json.Unmarshal(respBody, &result); err != nil {
		return types.UserProfile{}, fmt.Errorf("failed to decode LLM response: %w", err)
	}

	return result, nil
}

func GetRecommendations(userID string) ([]types.Recommendation, error) {
	url := fmt.Sprintf("%s/get_recommendations?id=%s", os.Getenv("LLM_HOST"), userID)
	headers := map[string]string{
		"x-api-key":    os.Getenv("LLM_API_KEY"),
		"Content-Type": "application/json",
	}

	respBody, err := utils.MakeHTTPRequest("GET", url, headers, nil)
	if err != nil {
		return []types.Recommendation{}, fmt.Errorf("failed to make request: %v", err)
	}

	var result []types.Recommendation
	if err := json.Unmarshal(respBody, &result); err != nil {
		return []types.Recommendation{}, fmt.Errorf("failed to decode LLM response: %w", err)
	}

	return result, nil
}

func GetAdvice(userID string) (types.Advice, error) {
	url := fmt.Sprintf("%s/get_ai_assisted_investments?id=%s", os.Getenv("LLM_HOST"), userID)
	headers := map[string]string{
		"x-api-key":    os.Getenv("LLM_API_KEY"),
		"Content-Type": "application/json",
	}

	respBody, err := utils.MakeHTTPRequest("GET", url, headers, nil)
	if err != nil {
		return types.Advice{}, fmt.Errorf("failed to make request: %v", err)
	}

	var result types.Advice
	if err := json.Unmarshal(respBody, &result); err != nil {
		return types.Advice{}, fmt.Errorf("failed to decode LLM response: %w", err)
	}

	return result, nil
}

func GetAssetRecommendation(userID string, symbol string) (int, error) {
	url := fmt.Sprintf("%s/get_asset_recommendation?id=%s&symbol=%s", os.Getenv("LLM_HOST"), userID, symbol)
	headers := map[string]string{
		"x-api-key":    os.Getenv("LLM_API_KEY"),
		"Content-Type": "application/json",
	}

	respBody, err := utils.MakeHTTPRequest("GET", url, headers, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to make request: %v", err)
	}

	var result int
	if err := json.Unmarshal(respBody, &result); err != nil {
		return 0, fmt.Errorf("failed to decode LLM response: %w", err)
	}

	return result, nil
}
