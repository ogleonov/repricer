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
	walletAuthToken  = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3NTQ4MjYxMzcsInVzZXIiOiIzNDE2NDgyMiIsInNoYXJkX2tleSI6IjkiLCJjbGllbnRfaWQiOiJ3YiIsInNlc3Npb25faWQiOiIzODM5NWQ5NDA5MDE0YTM0ODM0MDczMGE1ZmE4NTQ0YiIsInZhbGlkYXRpb25fa2V5IjoiMGQ4OWQzMTEyNjFlODQxMTI3ZTlhOThlNjE3ZjhhODFhMTkwNDQ2MTVmY2I1ZTllN2EwMjRmNmU1ZjM3NjFkZCIsInBob25lIjoiSEF1U1B5amdPZ3JGcEFFWG1CbFJ0Zz09IiwidXNlcl9yZWdpc3RyYXRpb25fZHQiOjE2NzYwNDUyODgsInZlcnNpb24iOjJ9.gI6B5QlOy_hW8iAmCLa7dUsEbsLf9k34Ds8P7tLO2pq7OMPDcXZNVEPmnQXo5fd-VrUFYhZVNMWDHiy2erFPL5LUZZh6yeTCmtNKMvzXt9ID0bMU_A3PCY5yrIEaQhRO7DLw_JhJ7xH2lJwa1pStYpj8T6vAQu4rB20zqgbFDwz5uYiuVbbCVcLFcclY0dkepY0NCZj8O8JEdIHJZblBiiUJ7_wd7GSs53esIYukECpy705ntPbWeWY_iGVDyEHdh4pzLtGVlIhLMlX5cSDyU8hyQ6M-OzeyeXd3fhDAonL1brwskYrFmewcZ4UMl_N8gWcWERuy5phFWIbiC4BPqw"
	seller1AuthToken = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3NTQxMzkwMjEsInVzZXIiOiIxNjI0MjYyNCIsInNoYXJkX2tleSI6IjUiLCJjbGllbnRfaWQiOiJzZWxsZXItcG9ydGFsIiwic2Vzc2lvbl9pZCI6ImY2Mzc2MWEyYTA0MjRlMjBiYTRiNzg5YTBlNThkYzJjIiwidmFsaWRhdGlvbl9rZXkiOiI5MDJjYWMzYjczMzQwZGVkMTRmMmNiOGY5ZWE5ZWIyNzJjYjkzZTBmODg4OTUzNGNkNTUzNDIxMDYzZjY1N2M3IiwidXNlcl9yZWdpc3RyYXRpb25fZHQiOjE2NzU0MTM4NTYsInZlcnNpb24iOjJ9.fayL01jqU1ZJeGJ9sLSAJHeDMyntLcHOKAO-SHIExtee4EpwLcQBdGDCNuQTG5AD2qcpWYZZBBrF89hUfxJ5bd7gdEL8J9-iW6Y_ealPQqLkoI4xyrH5wfIWKWgXMaZki9a8xQxIrQ7OikMiiVdNHXCA4AfxVHWeBa6yXWIgCbRWkdQtgaizXfshwxQdoLJeU5LdeupSJYl2TEKN-VS6kCd0ilQKTGbChyOtXS05hTAMtMnwhnBO_ZQPf-f57kuGlRa-1wjI7RqFje5RKe3iWpA0-UVB3LjyascFcEde8e3oHODfE7LrkGC9KPXAibP67EeiMCTRyoNVms4VdPX4Vg"
	seller2AuthToken = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3NTQ4MzA1NjYsInVzZXIiOiIxNTg0NzMxMSIsInNoYXJkX2tleSI6IjMiLCJjbGllbnRfaWQiOiJzZWxsZXItcG9ydGFsIiwic2Vzc2lvbl9pZCI6ImRlOTFlNmRjODM5YTQ2ZDU4NWNiNDMwZTAyN2NmZWE2IiwidmFsaWRhdGlvbl9rZXkiOiI1Yjg2ZDMyMzZlOGMxYmNmOGU4NzU2Yzg1YzUwODQ2MTlmOTdiZDMxODRhNGE1NzM0MjE0ZmZjN2RlZTcxNDk0IiwidXNlcl9yZWdpc3RyYXRpb25fZHQiOjE2ODI0NTY0MjAsInZlcnNpb24iOjJ9.MAn32dd-SuqhpgySZjCaeNocQGdErkDNhtngPTbIGdH5uoUAC0TX-2ZBLyNfzbmbwlxUkEHJRB6OJEunyaBBsZz8srPw1gm4gF5xF0UG2Y5NkZj0A1hj0bYg6y41vLdo1P69EkWCbtDoRuh_UxPVSG4cqmaMZqOTWcBMaRkxCbzvB-dLN-VZ1Rp9If7_mjcy0wj4AmgZ40JoZRDJeYuGE0jNh0sJjvI-YBqNSuXmuYMQ4g7OLUmzxe6FJoIcCXhJC389VorFvFKtlroyQbfQxU5Bd4J66kIC9vZECjySBvd3Xa3CZAGrV_DcEZdGgn8DGmEsAbsGwNbhij6kWtanVQ"

	// Список товаров для первого продавца [nmID: минимальная_цена]
	productsSeller1 = map[int]float64{
		439740235: 598.00,
		363561833: 2431.00,
		355039724: 2756.00,
		420175308: 995.00,
		445719497: 351.00,
		444947468: 637.00,
		450517748: 1190.00,
		451852395: 1157.00,
		447703683: 2431.00,
		413320662: 1950.00,
		452613966: 1950.00,
		455308681: 826.00,
		455874194: 936.00,
		465007169: 533.0,
		466364173: 1612.00,
		472724832: 1352.00,
		458176275: 1664.00,
		485867509: 676.0,
		486685652: 1235.0,
		449727119: 1112.0,
		492714507: 878.00,
		492947914: 578.00,
	}

	// Список товаров для второго продавца [nmID: минимальная_цена]
	productsSeller2 = map[int]float64{
		486062217: 800.00,
		483028809: 640.00,
		473997083: 550.00,
		478334856: 615.00,
		485654591: 750.00,
		471430353: 755.00,
		472057995: 2100.00,
		473066411: 810.00,
		476823280: 750.00,
		475454890: 520.00,
		475499053: 520.00,
		480109053: 520.00,
		471832484: 750.00,
		470975205: 1200.00,
		493044219: 780.00,
		493490629: 570.00,
		495179694: 767.00,
		496076265: 670.00,
		496941899: 871.00,
		496570292: 735.00,
	}

	// Настройки программы
	checkInterval = 5 * time.Minute // Интервал проверки цен
)

// Структура для хранения данных продавца
type Seller struct {
	Name     string
	Token    string
	Cookie   string
	Products map[int]float64
}

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

	// Инициализация продавцов
	sellers := []Seller{
		{
			Name:     "Первый продавец",
			Token:    seller1AuthToken,
			Cookie:   "_wbauid=1389459591727788866; ___wbu=2afbdd4a-2b25-44d0-a488-d2343f188ea6.1727788866; wbx-validation-key=a4aeb81b-3aea-4bca-a373-b724a179a919; _ym_uid=1726235106355624516; _ym_d=1740053345; external-locale=ru; x-supplier-id-external=be41cd8a-9260-412d-9445-cc8cf1d3aad0",
			Products: productsSeller1,
		},
		{
			Name:     "Второй продавец",
			Token:    seller2AuthToken,
			Cookie:   "_wbauid=5666810631754830452; wbx-validation-key=4f7e529c-4c72-4cdc-8b17-3075eac2bce4; x-supplier-id-external=df62fdc4-c58a-41dc-9aed-caf62c76df5f",
			Products: productsSeller2,
		},
	}

	log.Println("Запуск репрайсера Wildberries")
	log.Printf("Количество продавцов: %d", len(sellers))
	for i, seller := range sellers {
		log.Printf("Продавец %d: %s, товаров: %d", i+1, seller.Name, len(seller.Products))
	}
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

		// Обработка товаров для каждого продавца
		for _, seller := range sellers {
			processSellerProducts(seller, walletDiscount)
		}

		log.Printf("Цикл завершен. Следующая проверка через %v\n", checkInterval)
		time.Sleep(checkInterval)
	}
}

// Обработка товаров конкретного продавца
func processSellerProducts(seller Seller, walletDiscount int) {
	log.Printf("--- Обработка товаров для %s ---", seller.Name)

	// Проверка каждого товара
	for nmId, minPrice := range seller.Products {
		price, sellerDiscount, wbDiscount, err := getProductInfo(nmId, seller.Token, seller.Cookie)
		if err != nil {
			log.Printf("Ошибка получения информации о товаре %d: %v", nmId, err)
			continue
		}

		// Расчет финальной цены
		finalPrice := calculateFinalPrice(price, sellerDiscount, wbDiscount, walletDiscount)
		log.Printf("%s - Товар %d - Цена: %.2f, Скидка продавца: %d%%, Скидка WB: %d%%, Финальная цена: %.2f, Минимальная цена: %.2f",
			seller.Name, nmId, price, sellerDiscount, wbDiscount, finalPrice, minPrice)

		// Корректировка при необходимости
		if finalPrice < minPrice || finalPrice > minPrice+1 {
			log.Printf("ТРЕБУЕТСЯ КОРРЕКТИРОВКА: Финальная цена %.2f вне диапазона [%.2f, %.2f]",
				finalPrice, minPrice, minPrice+1)

			newPrice, newDiscount := findOptimalPrice(price, sellerDiscount, wbDiscount, walletDiscount, minPrice)
			err = updateProductPrice(nmId, newPrice, newDiscount, seller.Token, seller.Cookie)
			if err != nil {
				log.Printf("Ошибка обновления цены для товара %d: %v", nmId, err)
			} else {
				log.Printf("Цена обновлена: Новая цена = %.2f, Новая скидка = %d%%", newPrice, newDiscount)
			}
		}
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
func getProductInfo(nmId int, sellerToken string, cookie string) (price float64, sellerDiscount, wbDiscount int, err error) {
	url := fmt.Sprintf("https://discounts-prices.wildberries.ru/ns/dp-api/discounts-prices/suppliers/api/v1/nm/info?nmID=%d", nmId)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, 0, 0, err
	}

	// Устанавливаем заголовки
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,zh-CN;q=0.8,zh;q=0.7,en-US;q=0.6,en;q=0.5")
	req.Header.Set("Authorizev3", sellerToken)
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
	req.Header.Set("Cookie", cookie) // Используем переданные куки

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

// Расчет финальной цены (все скидки округляются вверх до рубля)
func calculateFinalPrice(price float64, sellerDiscount, wbDiscount, walletDiscount int) float64 {
	// Применяем скидку WB: вычисляем размер и округляем вверх
	wbDiscountAmount := price * float64(wbDiscount) / 100
	wbDiscountRounded := math.Ceil(wbDiscountAmount)
	currentPrice := price - wbDiscountRounded

	// Применяем скидку продавца
	sellerDiscountAmount := currentPrice * float64(sellerDiscount) / 100
	sellerDiscountRounded := math.Ceil(sellerDiscountAmount)
	currentPrice -= sellerDiscountRounded

	// Применяем скидку кошелька
	walletDiscountAmount := currentPrice * float64(walletDiscount) / 100
	walletDiscountRounded := math.Ceil(walletDiscountAmount)
	finalPrice := currentPrice - walletDiscountRounded

	// Округляем финальную цену до копейки
	return math.Round(finalPrice*100) / 100
}

// Поиск оптимальной цены и скидки для достижения целевой цены
func findOptimalPrice(currentPrice float64, currentDiscount, wbDiscount, walletDiscount int, minPrice float64) (float64, int) {
	const (
		discountRange = 10 // ±10% для скидки
	)

	// Инициализация лучших значений
	bestPrice := currentPrice
	bestDiscount := currentDiscount
	bestDiff := math.MaxFloat64

	// Перебираем возможные скидки
	for discountOffset := -discountRange; discountOffset <= discountRange; discountOffset++ {
		newDiscount := currentDiscount + discountOffset

		// Пропускаем недопустимые значения скидки
		if newDiscount < 0 || newDiscount > 100 {
			continue
		}

		// Находим оптимальную цену для текущей скидки
		price, diff := findPriceForDiscount(newDiscount, wbDiscount, walletDiscount, minPrice)

		// Если нашли лучший вариант
		if diff >= 0 && diff < bestDiff {
			bestPrice = price
			bestDiscount = newDiscount
			bestDiff = diff
		}
	}

	// Если не нашли подходящий вариант
	if bestDiff == math.MaxFloat64 {
		return currentPrice, currentDiscount
	}

	return bestPrice, bestDiscount
}

// Поиск оптимальной цены для фиксированной скидки
func findPriceForDiscount(discount, wbDiscount, walletDiscount int, minPrice float64) (float64, float64) {
	// Границы поиска цены
	low := minPrice
	high := minPrice * 2 // Верхняя граница - удвоенная минимальная цена

	bestPrice := low
	bestDiff := math.MaxFloat64

	// Бинарный поиск оптимальной цены
	for high-low > 0.01 { // Точность до копейки
		mid := (low + high) / 2
		finalPrice := calculateFinalPrice(mid, discount, wbDiscount, walletDiscount)
		diff := finalPrice - minPrice

		if diff >= 0 {
			if diff < bestDiff {
				bestDiff = diff
				bestPrice = mid
			}
			high = mid // Пробуем снизить цену
		} else {
			low = mid // Повышаем цену
		}
	}

	// Проверяем последнюю цену
	finalPrice := calculateFinalPrice(bestPrice, discount, wbDiscount, walletDiscount)
	finalDiff := finalPrice - minPrice

	return bestPrice, finalDiff
}

// Обновление цены товара
func updateProductPrice(nmId int, newPrice float64, newDiscount int, sellerToken string, cookie string) error {
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
	req.Header.Set("Authorizev3", sellerToken)
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
	req.Header.Set("Cookie", cookie) // Используем переданные куки

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
