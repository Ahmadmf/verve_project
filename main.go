package main

import (
	"io"
	"net/http"
	"os"
	"time"
	"verve_project/server"

	"github.com/gin-gonic/gin"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		NoColor:    false,
	})

	// setup log file
	file, err := os.Create("logfile.log")
	if err != nil {
		log.Fatal().Err(err).Msg("error in create log file")
	}
	defer file.Close()

	fileWriter := io.MultiWriter(file, os.Stdout)
	log.Logger = zerolog.New(fileWriter).With().Caller().Timestamp().Logger()
	gin.DefaultWriter = fileWriter
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	engine := server.NewRouter()

	engine.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong!",
		})
	})

	log.Info().Msg("Server starting at localhost:8085")
	if err := engine.Run("localhost:8085"); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server at 8085")
	}
}
