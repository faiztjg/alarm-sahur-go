package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
)

var (
	pumpStatus bool
	alarmTime  string = "03:30" // Default jam alarm
	duration   string = "10"    // Default durasi 10 detik
	mutex      sync.Mutex
)

func main() {
	http.HandleFunc("/", serveHTML)
	http.HandleFunc("/api/on", turnOnPump)
	http.HandleFunc("/api/off", turnOffPump)
	http.HandleFunc("/api/update", updateSettings)
	http.HandleFunc("/api/data", getData)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server berjalan di port " + port + "...")
	http.ListenAndServe(":"+port, nil)
}

func serveHTML(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	currTime := alarmTime
	currDur := duration
	mutex.Unlock()

	html := `<!DOCTYPE html><html><head><meta name='viewport' content='width=device-width, initial-scale=1.0'>
	<title>Alarm Sahur PAJA</title>
	<style>
		body { font-family: Arial; text-align: center; margin-top: 15px; background-color: #222; color: white;}
		h2 { color: #00bcd4; margin-bottom: 5px;}
		.box { background: #333; padding: 20px; border-radius: 10px; margin: 15px auto; width: 85%; max-width: 400px; box-shadow: 0 4px 8px rgba(0,0,0,0.5);}
		.button { padding: 15px; font-size: 16px; color: white; border: none; border-radius: 8px; cursor: pointer; width: 100%; margin-bottom: 10px; font-weight: bold;}
		.btn-on { background-color: #4CAF50; }
		.btn-off { background-color: #f44336; }
		.btn-save { background-color: #2196F3; margin-top: 10px; }
		input { width: 90%; padding: 10px; margin: 5px 0 15px 0; border-radius: 5px; border: none; font-size: 16px; text-align: center; background: #ddd;}
	</style></head><body>

	<h2>💧 Alarm Sahur PAJA</h2>
	<p style="margin-top:0; color:#aaa; font-size:14px;">Kontrol Jarak Jauh (Cloud)</p>
	
	<div class='box'>
		<h3 style="color: #4CAF50;">⚙️ Setting Otomatis</h3>
		<p style="margin-bottom: 5px;">Jam Bangun (WIB):</p>
		<input type="time" id="jam" value="` + currTime + `">
		<p style="margin-bottom: 5px;">Durasi Semprot (Detik):</p>
		<input type="number" id="durasi" value="` + currDur + `">
		<button class="button btn-save" onclick="simpanSetting()">💾 SIMPAN JADWAL</button>
		<p id="pesan-setting" style="color:#00bcd4; font-size: 14px; height: 15px; margin:0;"></p>
	</div>

	<div class='box'>
		<h3 style="color: #ff9800;">🕹️ Kontrol Manual</h3>
		<button class="button btn-on" onclick="nyalakan()">💦 NYALAKAN AIR</button>
		<button class="button btn-off" onclick="matikan()">🛑 MATIKAN AIR</button>
		<p id="pesan-manual" style="color:#00bcd4; font-size: 14px; height: 15px; margin:0;"></p>
	</div>

	<script>
		function nyalakan() {
			fetch('/api/on').then(() => {
				document.getElementById('pesan-manual').innerText = "POMPA MENYALA!";
				setTimeout(() => { document.getElementById('pesan-manual').innerText = ""; }, 4000);
			});
		}
		function matikan() {
			fetch('/api/off').then(() => {
				document.getElementById('pesan-manual').innerText = "POMPA DIMATIKAN!";
				setTimeout(() => { document.getElementById('pesan-manual').innerText = ""; }, 4000);
			});
		}
		function simpanSetting() {
			let j = document.getElementById('jam').value;
			let d = document.getElementById('durasi').value;
			fetch('/api/update?time=' + j + '&duration=' + d).then(() => {
				document.getElementById('pesan-setting').innerText = "Jadwal tersimpan di Server!";
				setTimeout(() => { document.getElementById('pesan-setting').innerText = ""; }, 4000);
			});
		}
	</script>
	</body></html>`

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

func updateSettings(w http.ResponseWriter, r *http.Request) {
	t := r.URL.Query().Get("time")
	d := r.URL.Query().Get("duration")
	if t != "" && d != "" {
		mutex.Lock()
		alarmTime = t
		duration = d
		mutex.Unlock()
	}
	w.Write([]byte("UPDATED"))
}

// Endpoint ini yang akan dibaca ESP32. Formatnya: MANUAL|JAM|DURASI (Contoh: 0|03:30|10)
func getData(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	status := "0"
	if pumpStatus {
		status = "1"
	}
	t := alarmTime
	d := duration
	mutex.Unlock()
	
	res := fmt.Sprintf("%s|%s|%s", status, t, d)
	w.Write([]byte(res))
}