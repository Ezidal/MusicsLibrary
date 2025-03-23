package main

import (
	"LibMusic/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {

	// Регистрируем обработчик для пути /info
	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, что метод запроса — GET
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// Получаем параметры из query
		group := r.URL.Query().Get("group")
		song := r.URL.Query().Get("song")
		fmt.Println(group, song)

		// Проверяем, что параметры переданы
		if group == "" || song == "" {
			http.Error(w, "group and song parameters are required", http.StatusBadRequest)
			return
		}

		// Пример данных (в реальном приложении данные можно брать из базы данных или другого источника)
		songDetail := models.DetailSong{
			ReleaseDate: "16.07.2006",
			Text: `Ooh baby, don't you know I suffer?
Ooh baby, can you hear me moan?
You caught me under false pretenses
How long before you let me go?

Ooh
You set my soul alight
Ooh
You set my soul alight`,
			Link: "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
		}

		// Устанавливаем заголовок Content-Type
		w.Header().Set("Content-Type", "application/json")

		// Кодируем данные в JSON и отправляем ответ
		if err := json.NewEncoder(w).Encode(songDetail); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	})

	// Запускаем сервер на порту 8080
	println("Server is running on http://localhost:8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}
