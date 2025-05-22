# crydl-go

A Go CLI for downloading crypto OHLCV data and listing supported exchanges/symbols.

## Usage

- Download OHLCV data (Binance only):

`crydl --help`
```
Crypto OHLCV CLI Downloader (Go)

Usage:
  crydl-go [flags]

Flags:
      --compress                           Enable ZIP compression of output CSV
      --current                            Set until to current UTC datetime
      --exchange string                    Exchange name (e.g., binanceus)
      --find-exchanges-for-symbol string   Show exchanges that support this symbol (e.g., BTC/USDT)
  -h, --help                               help for crydl-go
      --list-symbols-for-exchange string   List all symbols for this exchange (e.g., binanceus)
      --output string                      Output CSV filename (default: exchange_pair_since_until.csv)
      --proxy string                       Proxy URL (optional)
      --since string                       Start date (YYYY-MM-DD)
      --symbol string                      Trading pair symbol (e.g., BTC/USDT)
      --until string                       End date (YYYY-MM-DD), or use --current for now
```