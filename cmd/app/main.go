package main

import (
	"zatrasz75/go_test/configs"
	"zatrasz75/go_test/internal/app"
	"zatrasz75/go_test/pkg/logger"
)

func main() {
	l := logger.NewLogger()

	// Configuration
	cfg, err := configs.NewConfig(l)
	if err != nil {
		l.Fatal("ошибка при разборе конфигурационного файла", err)
	}
	// Run
	app.Run(cfg, l)
}
