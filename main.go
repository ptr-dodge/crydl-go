package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"crydl-go/downloader"
)

func main() {
	var (
		exchange, symbol, since, until, proxy, output  string
		compress, current                              bool
		findExchangesForSymbol, listSymbolsForExchange string
	)

	var rootCmd = &cobra.Command{
		Use:   "crydl-go",
		Short: "Crypto OHLCV CLI Downloader (Go)",
		Run: func(cmd *cobra.Command, args []string) {
			// Set default output filename if not provided
			if output == "" && exchange != "" && symbol != "" && since != "" {
				untilPart := until
				if current || until == "" {
					untilPart = "current"
				}
				safeSymbol := symbol
				safeSymbol = strings.ReplaceAll(safeSymbol, "/", "_")
				output = fmt.Sprintf("%s_%s_%s_%s.csv", exchange, safeSymbol, since, untilPart)
			}

			switch {
			case findExchangesForSymbol != "":
				exchanges, err := downloader.FindExchangesForSymbol(findExchangesForSymbol)
				if err != nil {
					fmt.Println("Error:", err)
					os.Exit(1)
				}
				fmt.Printf("Exchanges supporting %s:\n", findExchangesForSymbol)
				for _, ex := range exchanges {
					fmt.Println(" ", ex)
				}
			case listSymbolsForExchange != "":
				symbols, err := downloader.ListSymbolsForExchange(listSymbolsForExchange)
				if err != nil {
					fmt.Println("Error:", err)
					os.Exit(1)
				}
				fmt.Printf("Symbols on %s:\n", listSymbolsForExchange)
				for _, s := range symbols {
					fmt.Println(" ", s)
				}
			case exchange != "" && symbol != "" && since != "":
				if current {
					until = "current"
				}
				err := downloader.DownloadOHLCVData(exchange, symbol, since, until, proxy, compress, output)
				if err != nil {
					fmt.Println("Error:", err)
					os.Exit(1)
				}
			default:
				cmd.Help()
			}
		},
	}

	rootCmd.Flags().StringVar(&exchange, "exchange", "", "Exchange name (e.g., binanceus)")
	rootCmd.Flags().StringVar(&symbol, "symbol", "", "Trading pair symbol (e.g., BTC/USDT)")
	rootCmd.Flags().StringVar(&since, "since", "", "Start date (YYYY-MM-DD)")
	rootCmd.Flags().StringVar(&until, "until", "", "End date (YYYY-MM-DD), or use --current for now")
	rootCmd.Flags().BoolVar(&current, "current", false, "Set until to current UTC datetime")
	rootCmd.Flags().StringVar(&proxy, "proxy", "", "Proxy URL (optional)")
	rootCmd.Flags().BoolVar(&compress, "compress", false, "Enable ZIP compression of output CSV")
	rootCmd.Flags().StringVar(&output, "output", "", "Output CSV filename (default: exchange_pair_since_until.csv)")
	rootCmd.Flags().StringVar(&findExchangesForSymbol, "find-exchanges-for-symbol", "", "Show exchanges that support this symbol (e.g., BTC/USDT)")
	rootCmd.Flags().StringVar(&listSymbolsForExchange, "list-symbols-for-exchange", "", "List all symbols for this exchange (e.g., binanceus)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
