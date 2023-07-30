/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dth0/fda-scrape/internal/scrape"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts fda-scrape as a service",
	Run: func(_ *cobra.Command, _ []string) {

		cfg := scrape.NewConfig()
		cfg.CacheDir = viper.GetString("cacheDir")
		cfg.TargerAddr = viper.GetString("targetAddr")
		cfg.Address = viper.GetString("server.address")
		cfg.Port = viper.GetInt("server.port")

		log.Println("Starting Server at: ", cfg.Bind())

		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

		mux := http.NewServeMux()

		mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
		mux.HandleFunc("/api/fis", scrape.FisHandler(cfg))
		mux.HandleFunc("/api/acao", scrape.AcHandler(cfg))

		api := http.Server{
			Addr:    cfg.Bind(),
			Handler: mux,
		}

		serverError := make(chan error, 1)

		go func() {
			serverError <- api.ListenAndServe()
		}()

		select {
		case err := <-serverError:
			log.Fatal(err)
		case sig := <-shutdown:
			log.Printf("Received signal %v, terminating...", sig)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			if err := api.Shutdown(ctx); err != nil {
				log.Fatal(err)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
