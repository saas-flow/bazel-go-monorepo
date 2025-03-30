package logger

import (
	"bufio"
	"os"
	"sync"

	_ "github.com/natefinch/lumberjack"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var once sync.Once

var Module = fx.Module("logger",
	fx.Provide(New),
)

type Config struct {
	LogPath      string
	LogFilename  string
	LogMaxSize   int
	LogMaxBackup int
	LogMaxAge    int
	LogCompress  bool
}

func New() (log *zap.Logger) {

	// var filepath string
	// if config.LogPath != "" {
	// 	filepath = fmt.Sprintf("%s", config.LogPath)
	// }

	// var filename string = "combine.log"
	// if config.LogFilename != "" {
	// 	filename = config.LogFilename
	// }

	// rollingFile := &lumberjack.Logger{
	// 	Filename:   fmt.Sprintf("%s/%s", filepath, filename),
	// 	MaxSize:    config.LogMaxSize,
	// 	MaxBackups: config.LogMaxBackup,
	// 	MaxAge:     config.LogMaxAge,
	// 	Compress:   config.LogCompress,
	// }

	fields := zap.Fields(
		zap.String("service_name", os.Getenv("SERVICE_NAME")),
		zap.String("service_version", os.Getenv("SERVICE_VERSION")),
	)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	bufferedWriter := zapcore.AddSync(&BufferedWriteSyncer{
		writer: bufio.NewWriterSize(os.Stdout, 1024), // Buffer 1KB
	})

	// fileInfoCore := &MaskingCore{
	// 	zapcore.NewCore(
	// 		zapcore.NewJSONEncoder(encoderConfig),
	// 		zapcore.AddSync(rollingFile),
	// 		zap.InfoLevel,
	// 	),
	// }

	// fileErrorCore := &MaskingCore{
	// 	zapcore.NewCore(
	// 		zapcore.NewJSONEncoder(encoderConfig),
	// 		zapcore.AddSync(rollingFile),
	// 		zap.ErrorLevel,
	// 	),
	// }

	consoleInfoCore := &MaskingCore{
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(bufferedWriter),
			zap.InfoLevel,
		),
	}

	consoleErrorCore := &MaskingCore{
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(bufferedWriter),
			zap.ErrorLevel,
		),
	}

	core := zapcore.NewTee(
		consoleInfoCore,
		consoleErrorCore,
	)
	log = zap.New(core, fields)
	defer log.Sync()

	zap.ReplaceGlobals(log)

	return
}
