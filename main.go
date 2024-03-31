package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
	"log"
	"os"
)

func printResponse(resp *genai.GenerateContentResponse) {
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				fmt.Println(part)
			}
		}
	}
	fmt.Println("---")
}
func main() {

	var doubt string
	var username string
	var key string
	var getConfig bool
	type Config struct {
		Username string `json:"username"`
		Key      string `json:"key"`
	}

	configDir := "/home/vinico/.gg/"

	var rootCommand = &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			// validations
			if doubt == "" {
				fmt.Println("You must supply a doubt.")
				return
			}

			data, err := os.ReadFile(configDir + "config.json")
			if err != nil {
				fmt.Println("Error reading config file:", err)
				return
			}

			var config Config
			if err := json.Unmarshal(data, &config); err != nil {
				fmt.Println("Error parsing config file:", err)
				return
			}

			ctx := context.Background()
			client, err := genai.NewClient(ctx, option.WithAPIKey(config.Key))
			if err != nil {
				log.Fatal(err)
			}
			defer client.Close()

			model := client.GenerativeModel("gemini-pro")

			resp, err := model.GenerateContent(ctx, genai.Text(doubt))
			if err != nil {
				log.Fatal(err)
			}

			printResponse(resp)
		},
	}

	var configCommand = &cobra.Command{
		Use:   "configure",
		Short: "Configure user information and API key",
		Run: func(cmd *cobra.Command, args []string) {

			if getConfig {
				data, err := os.ReadFile(configDir + "config.json")
				if err != nil {
					fmt.Println("Error reading config file:", err)
					return
				}
				fmt.Println("Current configuration:", string(data))
			} else {
				config := Config{
					Username: username,
					Key:      key,
				}

				if _, err := os.Stat(configDir + "config.json"); err == nil {
					fmt.Println("Config file already exist.")

					data, err := os.ReadFile(configDir + "config.json")
					if err != nil {
						fmt.Println("Error reading config file:", err)
						return
					}

					var existingConfig Config
					if err := json.Unmarshal(data, &existingConfig); err != nil {
						fmt.Println("Error parsing config file:", err)
						return
					}

					existingConfig.Username = config.Username
					existingConfig.Key = config.Key

					data, err = json.Marshal(existingConfig)
					if err != nil {
						fmt.Println("Error marshalling config:", err)
						return
					}
					if err := os.WriteFile(configDir+"config.json", data, 0644); err != nil {
						fmt.Println("Error writing config file:", err)
						return
					}
				}
				if _, err := os.Stat(configDir); err != nil {
					if err := os.Mkdir(configDir, os.ModePerm); err != nil {
						log.Fatal(err)
					}
					return
				}
				data, err := json.Marshal(config)
				if err != nil {
					fmt.Println("Error marshalling config:", err)
					return
				}
				if err := os.WriteFile(configDir+"config.json", data, 0644); err != nil {
					fmt.Println("Error writing config file:", err)
					return
				}
			}
		},
	}

	rootCommand.Flags().StringVarP(&doubt, "query", "q", "", "The question to be answred by AI")
	configCommand.Flags().StringVarP(&username, "username", "u", "", "Name or nick for user.")
	configCommand.Flags().StringVarP(&key, "key", "k", "", "API Key")
	configCommand.Flags().BoolVarP(&getConfig, "get", "g", false, "Get the current configuration")
	rootCommand.AddCommand(configCommand)
	rootCommand.Execute()
}
