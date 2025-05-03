package requests

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	db "github.com/arcedo/financial-ai-backend/database"
	"github.com/arcedo/financial-ai-backend/types"
	"github.com/arcedo/financial-ai-backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SyncDailyStockData(ctx context.Context, store db.MongoStorage) error {
	stocksColl := store.Collection("stocks")
	productsColl := store.Collection("products")

	// Fetch all products
	cursor, err := productsColl.Find(ctx, bson.D{})
	if err != nil {
		return fmt.Errorf("failed to fetch products: %w", err)
	}
	defer cursor.Close(ctx)

	var products []types.Product
	if err := cursor.All(ctx, &products); err != nil {
		return fmt.Errorf("failed to decode products: %w", err)
	}

	// Loop over products to fetch stock data
	for _, product := range products {
		// Check for the latest stock data for the product symbol
		var latest types.Stock
		err := stocksColl.FindOne(ctx,
			bson.M{"symbol": product.Symbol},
			options.FindOne().SetSort(bson.D{{Key: "date", Value: -1}}), // Sort by latest date
		).Decode(&latest)

		var latestDate string
		if err == mongo.ErrNoDocuments {
			// No data yet for this stock, we'll fetch all data
			latestDate = ""
		} else if err != nil {
			log.Printf("error getting latest stock for %s: %v", product.Symbol, err)
			continue
		} else {
			latestDate = latest.Date
		}

		// Now we fetch the stock data, depending on whether we have an existing date or not
		var stockData map[string]types.NewStock
		if latestDate == "" {
			// Fetch all data if no date exists
			stockData, err = fetchStockData(product.Symbol)
		} else {
			// Fetch only missing data starting from the latestDate
			stockData, err = fetchStockDataAfter(product.Symbol, latestDate)
		}
		if err != nil {
			log.Printf("error fetching data for stock %s: %v", product.Symbol, err)
			continue
		}

		// Insert the new stock data into the database
		for date, stock := range stockData {
			// Convert the stock date string to the correct format
			parsedDate, err := time.Parse("2006-01-02", date)
			if err != nil {
				log.Printf("invalid date format for stock %s on %s: %v", product.Symbol, date, err)
				continue
			}

			// Insert stock data into the database
			_, err = stocksColl.InsertOne(ctx, types.NewStock{
				Symbol:     product.Symbol,
				Date:       parsedDate.Format("2006-01-02"), // Store the date in the standard format
				OpenPrice:  stock.OpenPrice,
				ClosePrice: stock.ClosePrice,
				HighPrice:  stock.HighPrice,
				LowPrice:   stock.LowPrice,
				Volume:     stock.Volume,
			})
			if err != nil {
				log.Printf("failed to insert stock data for %s on %s: %v", product.Symbol, date, err)
			}
		}
		time.Sleep(12 * time.Second)
	}
	return nil
}

// fetchStockData fetches the raw stock data from Alpha Vantage and converts it to a map[string]types.Stock.
func fetchStockData(symbol string) (map[string]types.NewStock, error) {
	// Construct the Alpha Vantage API URL
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=%s&apikey=%s", symbol, os.Getenv("ALPHA_VANTAGE_API_KEY"))
	// Make the API request using the MakeHTTPRequest function
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	respBody, err := utils.MakeHTTPRequest("GET", url, headers, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stock data from Alpha Vantage: %w", err)
	}

	// Decode the response into a map
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to decode Alpha Vantage response: %w", err)
	}

	// Handle Alpha Vantage errors
	if errMsg, ok := result["Error Message"]; ok {
		return nil, fmt.Errorf("Alpha Vantage error: %s", errMsg)
	}
	if note, ok := result["Note"]; ok {
		return nil, fmt.Errorf("Rate limit hit: %s", note)
	}

	// Check if the API returned data
	timeSeries, ok := result["Time Series (Daily)"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format from Alpha Vantage: missing 'Time Series (Daily)', full response: %+v", result)
	}

	// Convert the time series data into map[string]types.Stock
	stockData := make(map[string]types.NewStock)
	for date, values := range timeSeries {
		stockValues := values.(map[string]interface{})
		openStr := stockValues["1. open"].(string)
		open, _ := strconv.ParseFloat(openStr, 32)

		highStr := stockValues["2. high"].(string)
		high, _ := strconv.ParseFloat(highStr, 32)

		lowStr := stockValues["3. low"].(string)
		low, _ := strconv.ParseFloat(lowStr, 32)

		closeStr := stockValues["4. close"].(string)
		close, _ := strconv.ParseFloat(closeStr, 32)

		volumeStr := stockValues["5. volume"].(string)
		volume, _ := strconv.ParseFloat(volumeStr, 32)

		stockData[date] = types.NewStock{
			Symbol:     symbol,
			Date:       date,
			OpenPrice:  float32(open),
			ClosePrice: float32(close),
			HighPrice:  float32(high),
			LowPrice:   float32(low),
			Volume:     float32(volume),
		}
	}

	// Return the map of stocks
	return stockData, nil
}

func fetchStockDataAfter(symbol, latestDate string) (map[string]types.NewStock, error) {
	// Fetch all stock data and filter based on the latest date
	allStockData, err := fetchStockData(symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stock data for %s: %w", symbol, err)
	}

	// Filter the data to only include stocks after the latestDate
	newStockData := make(map[string]types.NewStock)
	for date, stock := range allStockData {
		if date > latestDate {
			newStockData[date] = stock
		}
	}

	return newStockData, nil
}
