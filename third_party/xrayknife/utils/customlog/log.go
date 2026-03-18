package customlog

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Type is the log level/category.
type Type uint8

// Defines the available log types.
var (
	Success    Type = 0x00
	Failure    Type = 0x01
	Processing Type = 0x02
	Finished   Type = 0x03
	Info       Type = 0x04
	Warning    Type = 0x05
	// None is for un-styled text, providing a neutral default.
	None Type = 0x06
)

// TypesDetails holds the visual properties for each log type.
type TypesDetails struct {
	symbol string
}

// logTypeMap maps a log Type to its visual details (symbol and color).
var logTypeMap = map[Type]TypesDetails{
	Success:    {symbol: "✅"},
	Failure:    {symbol: "❌"},
	Processing: {symbol: "⚙️ "},
	Finished:   {symbol: "🎉"},
	Info:       {symbol: "ℹ️ "},
	Warning:    {symbol: "⚠️ "},
	None:       {symbol: ""},
}

var (
	// Default output is os.Stderr
	output io.Writer = os.Stderr
	mu     sync.Mutex
)

// SetOutput redirects log output (e.g. to a file or websocket).
func SetOutput(w io.Writer) {
	mu.Lock()
	defer mu.Unlock()
	output = w
}

func GetOutput() io.Writer {
	mu.Lock()
	defer mu.Unlock()
	return output
}

// Printf prints a formatted, timestamped, and colored log message.
// It prepends the corresponding symbol and current time to the message.
func Printf(logType Type, format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()

	t, ok := logTypeMap[logType]
	if !ok {
		// Fallback for an undefined type to prevent a panic.
		t = logTypeMap[None]
	}

	// Prepare the prefix with a symbol (if it exists) and a timestamp.
	prefix := ""
	if t.symbol != "" {
		prefix = t.symbol + " "
	}
	currentTime := time.Now()
	fullFormat := prefix + currentTime.Format("15:04:05") + " " + format

	fmt.Fprintf(output, fullFormat, v...)
}

// Println prints the given arguments to the designated output, followed by a newline.
func Println(v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Fprintln(output, v...)
}

// GetColor wraps text in the ANSI color for the given log type.
func GetColor(logType Type, text string) string {
	return text
}
