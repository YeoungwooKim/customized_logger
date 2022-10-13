package main

import (
	"encoding/json"
	"fmt"
	"log_test/colorLog"
	"log_test/zapLog"
	"net/http"
	"time"
)

func main() {
	// colorLog := &colorLog.Logger{}
	// colorLog = colorLog.NewLogger()

	// TRACE, DEBUG, INFO, WARN, ERROR, FATAL.
	// colorLog.Trace("hello world %v", 1234-234)
	// colorLog.Debug("hello world %v", 1234-234)
	// colorLog.Info("hello world %v", 1234-234)
	// colorLog.Warn("hello world %v", 1234-234)
	// colorLog.Error("hello world %v", 1234-234)
	// colorLog.Fatal("hello world %v", 1234-234)

	http.Handle("/healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"time":    string(time.Now().Format("2006-01-02 15:04:05.000")),
			"content": "hello world",
		}

		jData, err := json.Marshal(data)
		if err != nil {
			colorLog.Error("err %v", err)
		}

		// TRACE, DEBUG, INFO, WARN, ERROR, FATAL.
		// colorLog.SetLogLevel(colorLog.WARN)
		zapLog.Error("hello world ~~~` ")
		// colorLog.Trace("%v", data)
		// colorLog.Debug("%v", data)
		// colorLog.Info("%v", data)
		// colorLog.Warn("%v", data)
		// colorLog.Error("%v", data)
		// colorLog.Fatal("%v", data)
		zapLog.Info("dddddㅇㅇㅇㅇ22222222222")

		w.Header().Set("Content-Type", "application/json")
		w.Write(jData)
	}))
	// colorLog.Info("hello")
	http.ListenAndServe(":8080", nil)

}

func getPanic() {
	colorLog.Info("i am inside")
	names := []string{"lobster", "sea urchin", "sea cucumber"}
	fmt.Println("My favorite sea creature is:", names[len(names)])
}
