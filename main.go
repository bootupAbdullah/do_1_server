package main

import (
"fmt"
"log"
"net/http"
"os"
)

func main() {
port := os.Getenv("PORT")
if port == "" {
port = "8080"
}

mux := http.NewServeMux()
mux.HandleFunc("/health", healthHandler)

addr := fmt.Sprintf(":%s", port)
log.Printf("server starting on %s", addr)

if err := http.ListenAndServe(addr, loggerMiddleware(mux)); err != nil {
log.Fatalf("server failed: %v", err)
}
}
