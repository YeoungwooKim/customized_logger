package main

import (
	"encoding/json"
	"log_test/colorLog"
	"net/http"
	"time"
)

func main() {
	// TRACE, DEBUG, INFO, WARN, ERROR, FATAL.
	colorLog.SetLogLevel(colorLog.INFO)

	http.Handle("/healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"time":    string(time.Now().Format("2006-01-02 15:04:05.000")),
			"content": "hello world",
		}

		jData, err := json.Marshal(data)
		if err != nil {
			colorLog.Error("err %v", err)
		}
		colorLog.Trace("%v", data)
		colorLog.Debug("%v", data)
		colorLog.Info("%v", data)
		colorLog.Warn("%v", data)
		colorLog.Error("%v", data)
		colorLog.Fatal("%v", data)

		w.Header().Set("Content-Type", "application/json")
		w.Write(jData)
	}))
	http.ListenAndServe(":8080", nil)
}
