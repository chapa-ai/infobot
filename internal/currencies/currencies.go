package currencies

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"infoBot/internal/models"
	"io"
	"net/http"
)

func GetCurrency(logger *logrus.Entry, currencyPair string) (*models.Currencies, error) {
	requestBody, err := json.Marshal("send request")
	if err != nil {
		logger.Errorf("failed json.Marshal: %s", err)
		return nil, err
	}
	pathUrl := fmt.Sprintf("https://api.kucoin.com/api/v1/market/stats?symbol=%s", currencyPair)
	req, err := http.NewRequest(http.MethodGet, pathUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		logger.Errorf("failed make get request: %s", err)
		return nil, err
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Errorf("failed http.DefaultClient.Do: %s", err)
		return nil, err
	}
	defer func() {
		err = response.Body.Close()
		if err != nil {
			logger.Errorf("closing response.Body failed: %s", err)
			return
		}
	}()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Errorf("failed io.ReadAll: %s", err)
		return nil, err
	}

	var currencies models.Currencies
	err = json.Unmarshal(body, &currencies)
	if err != nil {
		logger.Errorf("failed unmarshalling: %s", err)
		return nil, err
	}

	return &currencies, err
}
