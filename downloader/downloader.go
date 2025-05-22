package downloader

import (
	"archive/zip"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// DownloadOHLCVData fetches OHLCV data and writes to CSV (optionally zipped)
func DownloadOHLCVData(exchange, symbol, since, until, proxy string, compress bool, output string) error {
	// For demo: Use Binance US public REST API for 1m klines
	if strings.ToLower(exchange) != "binanceus" {
		return fmt.Errorf("only binanceus supported in Go demo")
	}
	baseURL := "https://api.binance.us"
	endpoint := "/api/v3/klines"
	symbolAPI := strings.ReplaceAll(symbol, "/", "")
	interval := "1m"

	startTime, err := time.Parse("2006-01-02", since)
	if err != nil {
		return err
	}
	var endTime time.Time
	if until == "current" || until == "" {
		endTime = time.Now().UTC()
	} else {
		endTime, err = time.Parse("2006-01-02", until)
		if err != nil {
			return err
		}
	}

	outFile := output
	var writer *csv.Writer
	var file *os.File

	if compress {
		zipFile, err := os.Create(strings.TrimSuffix(output, ".csv") + ".zip")
		if err != nil {
			return err
		}
		defer zipFile.Close()
		zipWriter := zip.NewWriter(zipFile)
		defer zipWriter.Close()
		csvWriter, err := zipWriter.Create(output)
		if err != nil {
			return err
		}
		writer = csv.NewWriter(csvWriter)
	} else {
		file, err = os.Create(outFile)
		if err != nil {
			return err
		}
		defer file.Close()
		writer = csv.NewWriter(file)
	}

	// Write header
	writer.Write([]string{"timestamp", "open", "high", "low", "close", "volume"})

	// Binance API limit: 1000 per request
	limit := 1000
	start := startTime
	for start.Before(endTime) {
		url := fmt.Sprintf("%s%s?symbol=%s&interval=%s&startTime=%d&endTime=%d&limit=%d",
			baseURL, endpoint, symbolAPI, interval,
			start.UnixMilli(),
			endTime.UnixMilli(),
			limit,
		)
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return err
		}
		var klines [][]interface{}
		if err := json.Unmarshal(body, &klines); err != nil {
			return err
		}
		if len(klines) == 0 {
			break
		}
		for _, k := range klines {
			ts := time.UnixMilli(int64(k[0].(float64))).UTC().Format("2006-01-02 15:04:05")
			writer.Write([]string{
				ts,
				k[1].(string),
				k[2].(string),
				k[3].(string),
				k[4].(string),
				k[5].(string),
			})
			start = time.UnixMilli(int64(k[0].(float64))).Add(time.Minute)
		}
		writer.Flush()
		if len(klines) < limit {
			break
		}
	}
	fmt.Println("Saved:", outFile)
	return nil
}

// FindExchangesForSymbol returns a list of exchanges supporting the symbol (demo: only binanceus)
func FindExchangesForSymbol(symbol string) ([]string, error) {
	// For demo, only binanceus is checked
	syms, err := ListSymbolsForExchange("binanceus")
	if err != nil {
		return nil, err
	}
	symbol = strings.ToUpper(symbol)
	for _, s := range syms {
		if s == symbol {
			return []string{"binanceus"}, nil
		}
	}
	return []string{}, nil
}

// ListSymbolsForExchange lists all symbols for an exchange (demo: only binanceus)
func ListSymbolsForExchange(exchange string) ([]string, error) {
	if strings.ToLower(exchange) != "binanceus" {
		return nil, fmt.Errorf("only binanceus supported in Go demo")
	}
	url := "https://api.binance.us/api/v3/exchangeInfo"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var data struct {
		Symbols []struct {
			Symbol string `json:"symbol"`
			Status string `json:"status"`
			Base   string `json:"baseAsset"`
			Quote  string `json:"quoteAsset"`
		} `json:"symbols"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	var pairs []string
	for _, s := range data.Symbols {
		if s.Status == "TRADING" {
			pairs = append(pairs, fmt.Sprintf("%s/%s", s.Base, s.Quote))
		}
	}
	return pairs, nil
}
