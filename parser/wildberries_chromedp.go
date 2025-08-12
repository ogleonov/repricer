package parser

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

// GetPriceChromedp - получает цену товара Wildberries с использованием chromedp
func GetPriceChromedp(productID string) (string, error) {
	// Настройка опций Chrome
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
	)

	// Создание контекста браузера
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Основной контекст с таймаутом
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 45*time.Second)
	defer cancel()

	// Формирование URL
	url := fmt.Sprintf("https://www.wildberries.ru/catalog/%s/detail.aspx", productID)

	// Актуальные селекторы Wildberries (июнь 2025)
	priceSelector := `span.price-block__wallet-price red-price price-block__wallet-price--pointer`
	waitSelector := `div.product-page`

	// Запуск браузера и получение цены
	var priceText string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),

		// Ждем загрузки основного контейнера товара
		chromedp.WaitVisible(waitSelector, chromedp.ByQuery),

		// Дополнительное ожидание для стабилизации страницы
		chromedp.Sleep(2*time.Second),

		// Получаем текст цены
		chromedp.Text(priceSelector, &priceText, chromedp.ByQuery),
	)
	if err != nil {
		return "", fmt.Errorf("chromedp execution failed: %w", err)
	}

	// Очистка и форматирование цены
	cleanPrice := strings.NewReplacer(
		"₽", "",
		" ", "",
		"\n", "",
		"\u00a0", "", // Удаляем неразрывные пробелы
	).Replace(priceText)

	return cleanPrice, nil
}

// Инициализация логгера для пакета
func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
