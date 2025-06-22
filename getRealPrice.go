package repricer

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	// Настройка контекста и опций
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
		chromedp.Flag("headless", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Параметры товара
	productURL := "https://www.wildberries.ru/catalog/355039724/detail.aspx"

	// Селекторы Wildberries (могут меняться!)
	priceSelector := `//span[@class="price-block__final-price"]`

	// Извлекаем цену
	var priceText string
	err := chromedp.Run(ctx,
		chromedp.Navigate(productURL),
		chromedp.WaitVisible(priceSelector, chromedp.BySearch),
		chromedp.Text(priceSelector, &priceText, chromedp.BySearch),
	)
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	// Обработка результата
	cleanPrice := strings.ReplaceAll(priceText, "₽", "")
	cleanPrice = strings.ReplaceAll(cleanPrice, " ", "")
	cleanPrice = strings.TrimSpace(cleanPrice)

	fmt.Printf("Цена товара: %s руб.\n", cleanPrice)
}
