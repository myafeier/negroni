package negroni

import (
	"log"
	"net/http"
	"os"
	"time"
	"fmt"
)

// ALogger interface
type ALogger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

// Logger is a middleware handler that logs the request as it goes in and the response as it goes out.
type Logger struct {
	// ALogger implements just enough log.Logger interface to be compatible with other implementations
	ALogger
}

// NewLogger returns a new Logger instance
func NewLogger() *Logger {
	return &Logger{log.New(os.Stdout, "[negroni] ", 0)}
}

func (l *Logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	log.Printf("Started %s %s for %s\n", r.Method, r.RequestURI, r.RemoteAddr)

	next(rw, r)


	res := rw.(ResponseWriter)
	content := fmt.Sprintf("Completed %s %s in %v\n", r.RequestURI, http.StatusText(res.Status()), time.Since(start))

	switch res.Status() {
	case 200, 201, 202:
		content = fmt.Sprintf("\033[1;32m%s\033[0m", content)
	case 301, 302:
		content = fmt.Sprintf("\033[1;37m%s\033[0m", content)
	case 304:
		content = fmt.Sprintf("\033[1;33m%s\033[0m", content)
	case 401, 403:
		content = fmt.Sprintf("\033[4;31m%s\033[0m", content)
	case 404:
		content = fmt.Sprintf("\033[1;31m%s\033[0m", content)
	case 500:
		content = fmt.Sprintf("\033[1;36m%s\033[0m", content)
	}

	log.Println(content)
}
