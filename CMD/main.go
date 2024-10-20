package main

import (
	"RMS/Database"
	"RMS/Routes"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

const shutdownTimeout = 20 * time.Second

func main() {

	//channel to do the task
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Welcome to the app - RMS")

	connect := fmt.Sprintf("host %s port %s ", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))
	fmt.Println(connect)

	//server instance
	serv := Routes.SetupRoutes()
	if err := Database.ConnectAndMigrate(
		//dbHost, dbPort, dbName, dbUser, dbPassword,

		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		Database.SSLModeDisable); err != nil {
		logrus.Panicf("Failed to initilize and migrate database with err %+v", err)
	}
	logrus.Printf("migration sucessfull")

	go func() {
		if err := serv.Run(":8085"); err != nil {
			logrus.Panicf("Failed to run server with err %+v", err)
		}
	}()

	logrus.Printf("Server started at : 8085")

	<-done

	logrus.Printf("Server Shutdown")

	if err := Database.ShutdownDatabase(); err != nil {
		logrus.WithError(err).Error("failed to close database connection")
	}

	if err := serv.Shutdown(shutdownTimeout); err != nil {
		logrus.WithError(err).Panic("failed to gracefully shutdown server")
	}

}
