package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"store/pkg/sdk/third/douyin"
)

func main() {
	http.HandleFunc("/douyin/commodity", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			http.Error(w, "url parameter is required", http.StatusBadRequest)
			return
		}

		fmt.Printf("Received request for URL: %s\n", url)

		client := douyin.NewClient()
		metadata, err := client.GetCommodityMetadata(r.Context(), url)
		if err != nil {
			fmt.Printf("Error processing URL %s: %v\n", url, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(metadata); err != nil {
			fmt.Printf("Error encoding response: %v\n", err)
		}
	})

	fmt.Println("Server starting on :9527...")
	if err := http.ListenAndServe(":9527", nil); err != nil {
		panic(err)
	}
}
