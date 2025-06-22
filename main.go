package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"time"
)

// Конфигурация программы
var (
	walletAuthToken = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3NTA1MDcxMDIsInVzZXIiOiIxNjI0MjYyNCIsInNoYXJkX2tleSI6IjUiLCJjbGllbnRfaWQiOiJ3YiIsInNlc3Npb25faWQiOiI5MzI5OTcyNmVjMDg0ZWQ2OTgwOTk3NjFkOWJlOTVlYiIsInZhbGlkYXRpb25fa2V5IjoiOTAyY2FjM2I3MzM0MGRlZDE0ZjJjYjhmOWVhOWViMjcyY2I5M2UwZjg4ODk1MzRjZDU1MzQyMTA2M2Y2NTdjNyIsInBob25lIjoiTVhDNGVwdXVNMjVGeXY5RllaNXdqQT09IiwidXNlcl9yZWdpc3RyYXRpb25fZHQiOjE2NzU0MTM4NTYsInZlcnNpb24iOjJ9.Sxs47unFGDC1gF4l7hUdN8ekkulA9CwO5wB3XnKXPycnTSr5gvisM6InGvrdkhad4n0kR9E960XMmzohz9HahMEcdYXlFM1annlpkUd-aEMYODrrxW3nVUTs1BCw3j2LLvQjIO31uh3I_1ou0sf4ue6HxZINnqQc1SJ-oumkyRrEWeGcWttou1Y50vnoeGhsYWnNFlW9dnPotyA_TT3K4GAEZ3zPjF5yKdeA7Iz86-vovHtWEZJ1XyyRb5F2byuS6iGCbq800o-bNkt0eq2MqqLzGb4FJkZPfxoxLWU47AOt1vEBsDuYR3UujulMpr9-QuZNnruZ1QpnO02TjH2ZQw"
	sellerAuthToken = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3NDgxMTE4MjEsInVzZXIiOiIxNjI0MjYyNCIsInNoYXJkX2tleSI6IjUiLCJjbGllbnRfaWQiOiJzZWxsZXItcG9ydGFsIiwic2Vzc2lvbl9pZCI6IjdiMmJlMjM0YTM1ODQ0YzlhMjYwMDRmMWJjZGYyNmU4IiwidmFsaWRhdGlvbl9rZXkiOiI5MDJjYWMzYjczMzQwZGVkMTRmMmNiOGY5ZWE5ZWIyNzJjYjkzZTBmODg4OTUzNGNkNTUzNDIxMDYzZjY1N2M3IiwidXNlcl9yZWdpc3RyYXRpb25fZHQiOjE2NzU0MTM4NTYsInZlcnNpb24iOjJ9.KSymBS2fwIyI4Q3MFaJYPKQnAvkeh2ghCHHWkemi4iNyfrR6ucBEgha5M4gly51e4OyjTkHUu4eeeJqBtQFvJ9xjp3UcNMUKoIKmGCAbwt1vi4474GFjQEEFiobi9DoEGTQGh0qniXC9V-vEDxyeEpKpsmN167qimwmWLUzJW2J7OufG7Bm4lSOikbDAO7zJmLw0jcO9WiEdED24SoayshOlxLruTdAmKvNQiG_weFmrbp_WeVNP4imSKNBXUieVSd4QFwkfGDn4GKoNzSBbZmM-6TWDbSdcoL7omVPREaRkGNdKB3Z9ptXcSwi07VTf0WnRQiSM7mGyZG29ttbYrA"
	sellerWbxKey    = "a4aeb81b-3aea-4bca-a373-b724a179a919"

	// Список товаров для отслеживания [nmID: минимальная_цена]
	products = map[int]float64{
		439740235: 598.00,
		363561833: 2184.00,
		355039724: 2756.00,
		420175308: 995.00,
		445719497: 351.00,
		444947468: 637.00,
	}

	// Настройки программы
	checkInterval = 5 * time.Minute // Интервал проверки цен
)

// Структуры для разбора ответов API
type WalletResponse struct {
	Payload struct {
		Payments []struct {
			UpridDiscount float64 `json:"uprid_discount"`
		} `json:"payments"`
	} `json:"payload"`
}

type ProductInfoResponse struct {
	Data struct {
		Info struct {
			Price          float64 `json:"price"`
			Discount       float64 `json:"discount"`
			DiscountOnSite float64 `json:"discountOnSite"`
		} `json:"info"`
	} `json:"data"`
	Error bool `json:"error"`
}

func main() {
	// Настройка логирования
	logFile, err := os.OpenFile("repricer.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Ошибка открытия файла логов: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	log.Println("Запуск репрайсера Wildberries")
	log.Printf("Отслеживается %d товаров\n", len(products))
	log.Printf("Интервал проверки: %v\n", checkInterval)

	for {
		log.Println("========================================")
		log.Println("Начало нового цикла проверки...")

		// Получение скидки кошелька (целое число)
		walletDiscount, err := getWalletDiscount()
		if err != nil {
			log.Printf("Ошибка получения скидки кошелька: %v", err)
			time.Sleep(checkInterval)
			continue
		}
		log.Printf("Текущая скидка кошелька: %d%%", walletDiscount)

		// Проверка каждого товара
		for nmId, minPrice := range products {
			price, sellerDiscount, wbDiscount, err := getProductInfo(nmId)
			if err != nil {
				log.Printf("Ошибка получения информации о товаре %d: %v", nmId, err)
				continue
			}

			// Расчет финальной цены
			finalPrice := calculateFinalPrice(price, sellerDiscount, wbDiscount, walletDiscount)
			log.Printf("Товар %d - Цена: %.2f, Скидка продавца: %d%%, Скидка WB: %d%%, Финальная цена: %.2f, Минимальная цена: %.2f",
				nmId, price, sellerDiscount, wbDiscount, finalPrice, minPrice)

			// Корректировка при необходимости
			if finalPrice < minPrice || finalPrice > minPrice+1 {
				log.Printf("ТРЕБУЕТСЯ КОРРЕКТИРОВКА: Финальная цена %.2f вне диапазона [%.2f, %.2f]",
					finalPrice, minPrice, minPrice+1)

				newPrice, newDiscount := calculateNewPrice(price, sellerDiscount, wbDiscount, walletDiscount, minPrice, finalPrice)
				err = updateProductPrice(nmId, newPrice, newDiscount)
				if err != nil {
					log.Printf("Ошибка обновления цены для товара %d: %v", nmId, err)
				} else {
					log.Printf("Цена обновлена: Новая цена = %.2f, Новая скидка = %d%%", newPrice, newDiscount)
				}
			}
		}

		log.Printf("Цикл завершен. Следующая проверка через %v\n", checkInterval)
		time.Sleep(checkInterval)
	}
}

// Запрос скидки кошелька (возвращает целое число)
func getWalletDiscount() (int, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://user-grade.wildberries.ru/api/v5/grade?curr=RUB", nil)
	if err != nil {
		return 0, err
	}

	// Устанавливаем заголовки
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,zh-CN;q=0.8,zh;q=0.7,en-US;q=0.6,en;q=0.5")
	req.Header.Set("Authorization", "Bearer "+walletAuthToken)
	req.Header.Set("Dnt", "1")
	req.Header.Set("Origin", "https://www.wildberries.ru")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", "https://www.wildberries.ru/")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"137\", \"Chromium\";v=\"137\", \"Not/A)Brand\";v=\"24\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("статус: %d, тело: %s", resp.StatusCode, string(body))
	}

	var response WalletResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, err
	}

	if len(response.Payload.Payments) == 0 {
		return 0, fmt.Errorf("скидки кошелька не найдены")
	}

	// Округляем скидку кошелька до целого числа
	return int(math.Round(response.Payload.Payments[0].UpridDiscount)), nil
}

// Запрос информации о товаре (возвращает целые числа для скидок)
func getProductInfo(nmId int) (price float64, sellerDiscount, wbDiscount int, err error) {
	url := fmt.Sprintf("https://discounts-prices.wildberries.ru/ns/dp-api/discounts-prices/suppliers/api/v1/nm/info?nmID=%d", nmId)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, 0, 0, err
	}

	// Устанавливаем заголовки
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,zh-CN;q=0.8,zh;q=0.7,en-US;q=0.6,en;q=0.5")
	req.Header.Set("Authorizev3", sellerAuthToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Dnt", "1")
	req.Header.Set("Origin", "https://seller.wildberries.ru")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", "https://seller.wildberries.ru/")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"137\", \"Chromium\";v=\"137\", \"Not/A)Brand\";v=\"24\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36")
	req.Header.Set("Cookie", "_wbauid=1389459591727788866; ___wbu=2afbdd4a-2b25-44d0-a488-d2343f188ea6.1727788866; wbx-validation-key=a4aeb81b-3aea-4bca-a373-b724a179a919; _ym_uid=1726235106355624516; _ym_d=1740053345; external-locale=ru; x-supplier-id-external=be41cd8a-9260-412d-9445-cc8cf1d3aad0")

	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, 0, 0, fmt.Errorf("статус: %d, тело: %s", resp.StatusCode, string(body))
	}

	var response ProductInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, 0, 0, err
	}

	if response.Error {
		return 0, 0, 0, fmt.Errorf("API вернуло ошибку")
	}

	// Округляем скидки до целых чисел
	return response.Data.Info.Price,
		int(math.Round(response.Data.Info.Discount)),
		int(math.Round(response.Data.Info.DiscountOnSite)),
		nil
}

// Расчет финальной цены (скидка кошелька округляется вверх до рубля)
func calculateFinalPrice(price float64, sellerDiscount, wbDiscount, walletDiscount int) float64 {
	// Применяем скидку продавца (целое число)
	priceAfterSeller := price * (1 - float64(sellerDiscount)/100)

	// Применяем скидку WB (целое число)
	priceAfterWB := priceAfterSeller * (1 - float64(wbDiscount)/100)

	// Рассчитываем размер скидки кошелька в рублях
	walletDiscountAmount := priceAfterWB * float64(walletDiscount) / 100

	// Округляем скидку кошелька ВВЕРХ до целого рубля
	walletDiscountAmount = math.Ceil(walletDiscountAmount)

	// Вычитаем скидку кошелька
	finalPrice := priceAfterWB - walletDiscountAmount

	// Округляем финальную цену до копейки
	return math.Round(finalPrice*100) / 100
}

// Расчет новых параметров цены (все скидки целые)
func calculateNewPrice(currentPrice float64, currentDiscount, wbDiscount, walletDiscount int, minPrice, currentFinal float64) (float64, int) {
	// Если цена уже в допустимом диапазоне - оставляем как есть
	if currentFinal >= minPrice && currentFinal <= minPrice+1 {
		return currentPrice, currentDiscount
	}

	// Определяем направление корректировки
	if currentFinal < minPrice {
		return adjustPriceUp(currentPrice, currentDiscount, wbDiscount, walletDiscount, minPrice)
	} else {
		return adjustPriceDown(currentPrice, currentDiscount, wbDiscount, walletDiscount, minPrice)
	}
}

// Корректировка при цене ниже минимальной
func adjustPriceUp(currentPrice float64, currentDiscount, wbDiscount, walletDiscount int, minPrice float64) (float64, int) {
	// Пытаемся уменьшить скидку продавца
	if currentDiscount > 0 {
		// Рассчитываем минимально необходимую цену после скидок
		requiredPriceAfterWB := minPrice / (1 - float64(walletDiscount)/100)
		requiredPriceAfterSeller := requiredPriceAfterWB / (1 - float64(wbDiscount)/100)

		// Рассчитываем новую скидку (целое число)
		newDiscount := int(math.Round(100 * (1 - requiredPriceAfterSeller/currentPrice)))
		if newDiscount < 0 {
			newDiscount = 0
		}

		// Проверяем результат
		newFinal := calculateFinalPrice(currentPrice, newDiscount, wbDiscount, walletDiscount)
		if newFinal >= minPrice && newFinal <= minPrice+1 {
			return currentPrice, newDiscount
		}
	}

	// Рассчитываем требуемую базовую цену
	requiredPriceAfterWB := minPrice / (1 - float64(walletDiscount)/100)
	requiredPriceAfterSeller := requiredPriceAfterWB / (1 - float64(wbDiscount)/100)
	newPrice := math.Ceil(requiredPriceAfterSeller)

	// Ищем минимальную цену, при которой финальная цена в диапазоне [minPrice, minPrice+1]
	for {
		// Используем 0 скидку (целое число)
		finalPrice := calculateFinalPrice(newPrice, 0, wbDiscount, walletDiscount)
		if finalPrice >= minPrice && finalPrice <= minPrice+1 {
			return newPrice, 0
		}
		newPrice += 1
	}
}

// Корректировка при цене выше минимальной
func adjustPriceDown(currentPrice float64, currentDiscount, wbDiscount, walletDiscount int, minPrice float64) (float64, int) {
	// Пытаемся увеличить скидку продавца (целыми шагами)
	if currentDiscount < 100 {
		// Начинаем с текущей скидки
		newDiscount := currentDiscount

		// Постепенно увеличиваем скидку (целыми значениями)
		for newDiscount <= 100 {
			newDiscount += 1
			newFinal := calculateFinalPrice(currentPrice, newDiscount, wbDiscount, walletDiscount)

			// Если достигли целевого диапазона
			if newFinal >= minPrice && newFinal <= minPrice+1 {
				return currentPrice, newDiscount
			}

			// Если перешли ниже минималки - откатываем на шаг и используем
			if newFinal < minPrice {
				newDiscount -= 1
				newFinal = calculateFinalPrice(currentPrice, newDiscount, wbDiscount, walletDiscount)
				if newFinal >= minPrice {
					return currentPrice, newDiscount
				}
				break
			}
		}
	}

	// Если не получилось решить скидкой, понижаем базовую цену
	newPrice := currentPrice
	for {
		newPrice -= 1
		if newPrice < 1 {
			return currentPrice, currentDiscount // Защита от отрицательных цен
		}

		// Рассчитываем финальную цену с текущей скидкой
		finalPrice := calculateFinalPrice(newPrice, currentDiscount, wbDiscount, walletDiscount)

		if finalPrice >= minPrice && finalPrice <= minPrice+1 {
			return newPrice, currentDiscount
		}

		// Если цена опустилась ниже минимальной - переключаемся на режим повышения
		if finalPrice < minPrice {
			return adjustPriceUp(newPrice, currentDiscount, wbDiscount, walletDiscount, minPrice)
		}
	}
}

func updateProductPrice(nmId int, newPrice float64, newDiscount int) error {
	url := "https://discounts-prices.wildberries.ru/ns/dp-api/discounts-prices/suppliers/api/v1/nm/upload/task?checkChange=true"

	// Структура запроса (все поля целочисленные)
	type PriceData struct {
		NmID     int `json:"nmID"`
		Price    int `json:"price"`
		Discount int `json:"discount"`
	}

	type RequestPayload struct {
		Data PriceData `json:"data"`
	}

	payload := RequestPayload{
		Data: PriceData{
			NmID:     nmId,
			Price:    int(math.Round(newPrice)),
			Discount: newDiscount,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Логируем отправляемые данные
	log.Printf("Отправка данных для товара %d: %s", nmId, string(jsonData))

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// Устанавливаем заголовки
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,zh-CN;q=0.8,zh;q=0.7,en-US;q=0.6,en;q=0.5")
	req.Header.Set("Authorizev3", sellerAuthToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Dnt", "1")
	req.Header.Set("Origin", "https://seller.wildberries.ru")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", "https://seller.wildberries.ru/")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"137\", \"Chromium\";v=\"137\", \"Not/A)Brand\";v=\"24\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36")
	req.Header.Set("Cookie", "_wbauid=1389459591727788866; ___wbu=2afbdd4a-2b25-44d0-a488-d2343f188ea6.1727788866; wbx-validation-key=a4aeb81b-3aea-4bca-a373-b724a179a919; _ym_uid=1726235106355624516; _ym_d=1740053345; external-locale=ru; x-supplier-id-external=be41cd8a-9260-412d-9445-cc8cf1d3aad0")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("статус: %d, тело: %s", resp.StatusCode, string(body))
	}

	// Логируем отправленные данные
	log.Printf("Отправка данных для nmID %d: цена=%d, скидка=%d", nmId, int(math.Round(newPrice)), newDiscount)

	return nil
}
