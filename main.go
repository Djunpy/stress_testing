package main

import (
	"fmt"
	vegeta "github.com/tsenart/vegeta/lib"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func generateUniqueUserData() (string, string, int32) {
	// Генерация уникальных данных
	rand.Seed(time.Now().UnixNano())
	uniqueID := strconv.Itoa(rand.Intn(1000000)) // Генерация случайного числа до 1 000 000

	username := "username" + uniqueID
	email := "user" + uniqueID + "@gmail.com"
	number := rand.Intn(10000000) + 1000000 // Генерация случайного 7-значного номера

	return username, email, int32(number)
}

func main() {
	rate := vegeta.Rate{Freq: 50, Per: time.Second} // 50 запросов в секунду
	duration := 10 * time.Minute                    // Продолжительность теста

	targeter := func(tg vegeta.Targeter) vegeta.Targeter {
		return func(tgt *vegeta.Target) error {
			username, email, number := generateUniqueUserData()
			body := fmt.Sprintf(`{
				"username": "%s",
				"email": "%s",
				"photo": "sdf/sdfljsdf/sdfsdf",
				"password1": "7798439814Ka",
				"password2": "7798439814Ka",
				"first_name": "yurii",
				"last_name": "vovnenko",
				"phone": {
					"number": %d,
					"country_code": "+373"
				}
			}`, username, email, number)

			*tgt = vegeta.Target{
				Method: "POST",
				URL:    "http://localhost:8082/api/v1/auth/public/sign-up", // Замените на ваш реальный URL эндпоинта
				Body:   []byte(body),
				Header: map[string][]string{
					"Content-Type": {"application/json"},
				},
			}
			return nil
		}
	}

	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter(vegeta.NewStaticTargeter()), rate, duration, "Big Bang!") {
		metrics.Add(res)
		fmt.Printf("Status Code: %d\n", res.Code)
		fmt.Printf("Response Body: %s\n", string(res.Body))
		fmt.Printf("Latency: %s\n", res.Latency)

	}
	metrics.Close()

	fmt.Printf("Requests: %d\n", metrics.Requests)
	fmt.Printf("Duration: %s\n", metrics.Duration)
	fmt.Printf("Latencies: %s\n", metrics.Latencies)
	fmt.Printf("Bytes In: %d\n", metrics.BytesIn.Total)
	fmt.Printf("Bytes Out: %d\n", metrics.BytesOut.Total)
	fmt.Printf("Success: %f\n", metrics.Success)
	fmt.Printf("Status Codes: %v\n", metrics.StatusCodes)
	fmt.Printf("Errors: %v\n", metrics.Errors)

	// Вывод результатов в файл для дальнейшего анализа
	reportFile, err := os.Create("report.html")
	if err != nil {
		fmt.Printf("Error creating report file: %v\n", err)
		return
	}
	defer reportFile.Close()
	vegeta.NewJSONReporter(&metrics).Report(reportFile)
}
