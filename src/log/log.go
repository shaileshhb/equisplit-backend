package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// init will create instance of logger
func InitializeLogger() zerolog.Logger {
	return logFormatter()
}

func logFormatter() zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	// output.FormatLevel = func(i interface{}) string {
	// 	return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	// }
	// output.FormatMessage = func(i interface{}) string {
	// 	if i != nil {
	// 		return fmt.Sprintf("***%s****", i)
	// 	}
	// 	return ""
	// }
	// output.FormatFieldName = func(i interface{}) string {
	// 	return fmt.Sprintf("%s:", i)
	// }
	// output.FormatFieldValue = func(i interface{}) string {
	// 	return strings.ToUpper(fmt.Sprintf("%s", i))
	// }

	log := zerolog.New(output).Level(zerolog.InfoLevel).With().Timestamp().Caller().Logger()
	return log
}
