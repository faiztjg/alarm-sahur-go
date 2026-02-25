package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
)

var (
	pumpStatus bool
	mutex      sync.Mutex
)

func main() {
	http.HandleFunc("/", serveHTML)
	http.HandleFunc("/api/on", turnOnPump)
	http.HandleFunc("/api/off", turnOffPump)
	http.HandleFunc("/api/status", getStatus)

	// Mengambil Port otomatis dari server Cloud (Render), atau pakai 8080 untuk lokal
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server berjalan di port " + port + "...")
	http.ListenAndServe(":"+port, nil)
}

func serveHTML(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
	<html>
	<head>
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Alarm Sahur Azel</title>
		<style>
			body { font-family: Arial; text-align: center; background-color: #222; color: white; margin-top: 80px; }
			.btn { background-color: #4CAF50; color: white; border: none; padding: 40px; font-size: 30px; border-radius: 20px; cursor: pointer; width: 80%; max-width: 350px; box-shadow: 0 9px #111; font-weight: bold;}
			.btn:active { background-color: #3e8e41; box-shadow: 0 5px #111; transform: translateY(4px); }
			p { font-size: 20px; margin-top: 30px; color: #00bcd4; }
		</style>
	</head>
	<body>
		<h2>Bangunin Azel Sahur!</h2>
		<button class="btn" onclick="nyalakan()">NYALAKAN AIR 💦</button>
		<p id="pesan"></p>

		<script>
			function nyalakan() {
				fetch('/api/on').then(() => {
					document.getElementById('pesan').innerText = "Air sedang menyemprot di Jakarta! 🚀";
					setTimeout(() => { document.getElementById('pesan').innerText = ""; }, 10000);
				});
			}
		</script>
	</body>
	</html>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func turnOnPump(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	pumpStatus = true
	mutex.Unlock()
	w.Write([]byte("ON"))
}

func turnOffPump(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	pumpStatus = false
	mutex.Unlock()
	w.Write([]byte("OFF"))
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	status := pumpStatus
	mutex.Unlock()
	if status {
		w.Write([]byte("ON"))
	} else {
		w.Write([]byte("OFF"))
	}
}