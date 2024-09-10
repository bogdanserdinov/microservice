package main

import (
	"context"
	"database/sql"
	"github.com/caarlos0/env/v6"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
	"log"
	"microservice"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq" // using postgres driver.
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	onSigInt(func() {
		// starting graceful shutdown on context cancellation.
		cancel()
	})

	loggerCfg := zap.NewProductionConfig()
	logger, err := loggerCfg.Build()
	if err != nil {
		log.Fatalln("could not build logger from config", err)
		return
	}
	defer func() {
		err = logger.Sync()
		if err != nil {
			log.Fatalln("could not sync logger", err)
		}
	}()

	cfg, err := getConfigFromEnv()
	if err != nil {
		logger.Fatal("could not parse config", zap.Error(err))
	}
	logger.Info("cfg", zap.Any("cfg", cfg))

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		logger.Error("can't open connection to postgres", zap.Error(err))
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("can't close connection to postgres", zap.Error(err))
		}
	}()
	if err := db.Ping(); err != nil {
		logger.Error("can't ping database", zap.Error(err))
		return
	}

	app := microservice.New(logger, *cfg, db)
	if err != nil {
		logger.Error("can't initialize app", zap.Error(err))
		return
	}

	logger.Info("servers shutdown err", zap.Error(app.Run(ctx)))

}

// onSigInt fires on a SIGINT or SIGTERM event (usually CTRL+C).
func onSigInt(onSigInt func()) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-done
		onSigInt()
	}()
}

//func getConfigFromEnv() (*microservice.Config, error) {
//	cfg := new(microservice.Config)
//
//	viper.AutomaticEnv()
//	viper.SetConfigType("env")
//	//viper.AllowEmptyEnv(false)
//	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
//
//	// setting up filepath of env file to load, used only for local development.
//	//path := os.Getenv("ENV_PATH")
//	readFromFile := viper.GetBool("READ_FROM_FILE")
//	if readFromFile {
//		// looking for the .env file from the root.
//		viper.SetConfigName(".env")
//		viper.AddConfigPath("./")
//	}
//
//	err := viper.ReadInConfig()
//	if err != nil {
//		return cfg, err
//	}
//
//	err = viper.Unmarshal(cfg)
//	return cfg, nil
//}

func getConfigFromEnv() (*microservice.Config, error) {
	cfg := new(microservice.Config)
	err := env.Parse(cfg)
	return cfg, err
}
