package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/kkgo-software-engineering/workshop/database"
	"github.com/kkgo-software-engineering/workshop/router"
	"go.uber.org/zap"

	_ "github.com/lib/pq"
)

func main() {

	cfg := config.New().All()
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	sql, err := sql.Open("postgres", cfg.DBConnection)
	if err != nil {
		logger.Fatal("unable to configure database", zap.Error(err))
	}

	createTBAccount := `CREATE TABLE IF NOT EXISTS accounts (id SERIAL PRIMARY KEY, balance FLOAT)`
	createTBPocket := `CREATE TABLE IF NOT EXISTS pockets (id SERIAL PRIMARY KEY, account_id INT, name TEXT, category TEXT, currency TEXT, balance FLOAT, CONSTRAINT fk_account_id FOREIGN KEY(account_id) REFERENCES accounts(id))`

	_, err = sql.Exec(createTBAccount)
	if err != nil {
		log.Fatal("can't create table accounts", err)
	}

	_, err = sql.Exec(createTBPocket)
	if err != nil {
		log.Fatal("can't create table pockets", err)
	}

	e := router.RegRoute(cfg, logger, sql)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Hostname, cfg.Server.Port)

	go func() {
		err := e.Start(addr)
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal("unexpected shutdown the server", zap.Error(err))
		}
		logger.Info("gracefully shutdown the server")
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	gCtx := context.Background()
	ctx, cancel := context.WithTimeout(gCtx, 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logger.Fatal("unexpected shutdown the server", zap.Error(err))
	}
}
