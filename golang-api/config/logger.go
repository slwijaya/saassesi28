package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger adalah instance global yang bisa digunakan di seluruh aplikasi
var Logger = logrus.New()

func InitLogger() {
	// Set format log menjadi JSON agar lebih terstruktur
	Logger.SetFormatter(&logrus.JSONFormatter{})

	// Simpan log ke file
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Logger.SetOutput(file)
	} else {
		Logger.Info("Gagal menyimpan log ke file, menggunakan stdout")
	}

	// Atur level logging ke Info (bisa diganti menjadi Debug untuk debugging lebih detail)
	Logger.SetLevel(logrus.InfoLevel)
}
