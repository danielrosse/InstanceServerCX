package logs

import (
	"github.com/fatih/color"
	"log"
	"os"
)

func saveLog(message string) {
	file, err := os.OpenFile("logs.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	logger := log.New(file, "", log.LstdFlags)
	logger.Println(message)
}

func FatalLog(message string) {
	color.Set(color.FgRed)
	saveLog(" 💣 " + message)
	color.Unset()
	log.Fatal(" 💣  " + message)
}

func ErrorLog(message string) {
	color.Set(color.FgRed)
	log.Println(" 🚨 ", message)
	color.Unset()
	saveLog(" 🚨 " + message)
}

func InfoLog(message string) {
	color.Set(color.FgYellow)
	log.Println(" ⚠️ ", message)
	color.Unset()
	saveLog(" ⚠️  " + message)
}

func SuccessLog(message string) {
	color.Set(color.FgGreen)
	log.Println(" ✅ ", message)
	color.Unset()
	saveLog(" ✅ " + message)
}
