package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	// Настройка опций браузера
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
		chromedp.Flag("headless", true), // true для скрытого режима
	)

	// Создаем контекст браузера
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Создаем контекст таймаута (30 секунд)
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// URL товара Wildberries
	url := "https://www.wildberries.ru/catalog/355039724/detail.aspx"

	// Селектор для цены (может потребоваться обновить)
	selector := `span.price-block__final-price`

	// Ждем загрузки цены и извлекаем текст
	var priceText string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(selector, chromedp.ByQuery),
		chromedp.Text(selector, &priceText, chromedp.ByQuery),
	)
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}

	// Очищаем цену от лишних символов
	cleanPrice := strings.NewReplacer(
		"₽", "",
		" ", "",
		"\n", "",
		"\u00a0", "", // Удаляем неразрывные пробелы
	).Replace(priceText)

	fmt.Printf("Цена товара: %s руб.\n", cleanPrice)
}
