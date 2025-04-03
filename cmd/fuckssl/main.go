package main

import (
	"fmt"
	"fuckssl/internal/config"
	"fuckssl/internal/deployer"
	"fuckssl/internal/provider"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "fuckssl",
		Short: "A tool for managing SSL certificates",
		Long:  `FuckSSL is a command line tool for managing SSL certificates from multiple providers and deploying them to multiple platforms.`,
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.fuckssl/config.yaml)")

	// 添加子命令
	rootCmd.AddCommand(applyAndDeployCmd())
	rootCmd.AddCommand(configCmd())
}

func initConfig() {
	if err := config.Init(); err != nil {
		fmt.Println("Error initializing config:", err)
		os.Exit(1)
	}

	// 初始化提供者
	cfg := config.GetConfig()
	for _, p := range cfg.Providers {
		switch p.Type {
		case "tencent":
			if p, err := provider.NewTencentProvider(
				p.Settings["secret_id"],
				p.Settings["secret_key"],
			); err == nil {
				provider.Register(p)
			}
		case "aliyun":
			if p, err := provider.NewAliyunProvider(
				p.Settings["access_key_id"],
				p.Settings["access_key_secret"],
			); err == nil {
				provider.Register(p)
			}
		}
	}

	// 初始化部署者
	for _, d := range cfg.Deployers {
		switch d.Type {
		case "qiniu":
			if d, err := deployer.NewQiniuDeployer(
				d.Settings["access_key"],
				d.Settings["secret_key"],
				d.Settings["bucket"],
				d.Settings["domain"],
			); err == nil {
				deployer.Register(d)
			}
		}
	}
}

func applyAndDeployCmd() *cobra.Command {
	var providerName, deployerName string

	cmd := &cobra.Command{
		Use:   "apply-deploy [domain]",
		Short: "Apply for a certificate and deploy it",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			domain := args[0]

			// 获取证书提供者
			p, err := provider.GetProvider(providerName)
			if err != nil {
				return fmt.Errorf("failed to get provider: %v", err)
			}

			// 申请证书
			cert, err := p.Apply(domain)
			if err != nil {
				return fmt.Errorf("failed to apply certificate: %v", err)
			}

			// 获取部署者
			d, err := deployer.GetDeployer(deployerName)
			if err != nil {
				return fmt.Errorf("failed to get deployer: %v", err)
			}

			// 部署证书
			if err := d.Deploy(cert); err != nil {
				return fmt.Errorf("failed to deploy certificate: %v", err)
			}

			fmt.Printf("Successfully applied and deployed certificate for %s\n", domain)
			return nil
		},
	}

	cmd.Flags().StringVarP(&providerName, "provider", "p", "", "certificate provider (required)")
	cmd.Flags().StringVarP(&deployerName, "deployer", "d", "", "certificate deployer (required)")
	cmd.MarkFlagRequired("provider")
	cmd.MarkFlagRequired("deployer")

	return cmd
}

func configCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
	}

	// 添加子命令
	cmd.AddCommand(configListCmd())
	cmd.AddCommand(configSetCmd())
	cmd.AddCommand(configGetCmd())

	return cmd
}

func configListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all configurations",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := config.GetConfig()
			fmt.Println("Providers:")
			for name, provider := range cfg.Providers {
				fmt.Printf("  %s:\n    Type: %s\n", name, provider.Type)
			}
			fmt.Println("\nDeployers:")
			for name, deployer := range cfg.Deployers {
				fmt.Printf("  %s:\n    Type: %s\n", name, deployer.Type)
			}
		},
	}
}

func configSetCmd() *cobra.Command {
	var providerType, deployerType string
	var settings map[string]string

	cmd := &cobra.Command{
		Use:   "set [provider|deployer] [name]",
		Short: "Set configuration for a provider or deployer",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.GetConfig()
			configType := args[0]
			name := args[1]

			switch configType {
			case "provider":
				if providerType == "" {
					return fmt.Errorf("provider type is required")
				}
				cfg.Providers[name] = config.ProviderConfig{
					Type:     providerType,
					Settings: settings,
				}
			case "deployer":
				if deployerType == "" {
					return fmt.Errorf("deployer type is required")
				}
				cfg.Deployers[name] = config.DeployerConfig{
					Type:     deployerType,
					Settings: settings,
				}
			default:
				return fmt.Errorf("invalid config type: %s", configType)
			}

			return config.SaveConfig()
		},
	}

	cmd.Flags().StringVar(&providerType, "provider-type", "", "type of the provider")
	cmd.Flags().StringVar(&deployerType, "deployer-type", "", "type of the deployer")
	cmd.Flags().StringToStringVar(&settings, "settings", nil, "key-value settings")

	return cmd
}

func configGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get [provider|deployer] [name]",
		Short: "Get configuration for a provider or deployer",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.GetConfig()
			configType := args[0]
			name := args[1]

			switch configType {
			case "provider":
				if provider, exists := cfg.Providers[name]; exists {
					fmt.Printf("Type: %s\nSettings:\n", provider.Type)
					for k, v := range provider.Settings {
						fmt.Printf("  %s: %s\n", k, v)
					}
				} else {
					return fmt.Errorf("provider %s not found", name)
				}
			case "deployer":
				if deployer, exists := cfg.Deployers[name]; exists {
					fmt.Printf("Type: %s\nSettings:\n", deployer.Type)
					for k, v := range deployer.Settings {
						fmt.Printf("  %s: %s\n", k, v)
					}
				} else {
					return fmt.Errorf("deployer %s not found", name)
				}
			default:
				return fmt.Errorf("invalid config type: %s", configType)
			}

			return nil
		},
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
