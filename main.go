package main

import (
"fmt"
"log"
"net/http"
"os"
"github.com/bootupAbdullah/do_1_server/internal"
)

func main() {
port := os.Getenv("PORT")
if port == "" {
port = "8080"
}

mux := http.NewServeMux()
mux.HandleFunc("/health", internal.HealthHandler)

addr := fmt.Sprintf(":%s", port)
log.Printf("server starting on %s", addr)

if err := http.ListenAndServe(addr, internal.LoggerMiddleware(mux)); err != nil {
log.Fatalf("server failed: %v", err)
}
}
