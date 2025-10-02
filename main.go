package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/glebarez/sqlite" // Это самый надежный драйвер без CGO
	"gopkg.in/telebot.v3"
)

// Конфигурация программы
var (
	walletAuthToken  = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3NTg0NDc1MDgsInVzZXIiOiIzNDE2NDgyMiIsInNoYXJkX2tleSI6IjkiLCJjbGllbnRfaWQiOiJ3YiIsInNlc3Npb25faWQiOiIzODM5NWQ5NDA5MDE0YTM0ODM0MDczMGE1ZmE4NTQ0YiIsInZhbGlkYXRpb25fa2V5IjoiMGQ4OWQzMTEyNjFlODQxMTI3ZTlhOThlNjE3ZjhhODFhMTkwNDQ2MTVmY2I1ZTllN2EwMjRmNmU1ZjM3NjFkZCIsInBob25lIjoiSEF1U1B5amdPZ3JGcEFFWG1CbFJ0Zz09IiwidXNlcl9yZWdpc3RyYXRpb25fZHQiOjE2NzYwNDUyODgsInZlcnNpb24iOjJ9.QvkYCHqteG5940eiu5b8AX1CGJ_cdkdZx_D1vAnvTUFyZVPzzrCSxBn907jKBLEdj2MG50lg3Bmox_RyaeInZ2eKWpNT36KxxdNEc0Bws0RmXASb9-jdNnsFTrg7gic9dikftOmzbvdInSwwtFAEHjXafK_Cs3HYU_n3XoyMNaA-UHU5_62v-V7hnypRM8sd07mqu3XHgXnHSA0x9sYCXFcVttOUKylNf2L8HcRXxUFggqj3VH84lpb2GK_1QonvDN-5DHfWY-GIr4ibas8l5BE1Npv0NwaLqTCERRL4BlICvBv17k6AlwL3uQnOalh2yp-C3AK7JSe_KavHPX_cDA"
	seller1AuthToken = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3NTkzODQ0MDEsInVzZXIiOiIxNjI0MjYyNCIsInNoYXJkX2tleSI6IjUiLCJjbGllbnRfaWQiOiJzZWxsZXItcG9ydGFsIiwic2Vzc2lvbl9pZCI6ImY2Mzc2MWEyYTA0MjRlMjBiYTRiNzg5YTBlNThkYzJjIiwidmFsaWRhdGlvbl9rZXkiOiI5MDJjYWMzYjczMzQwZGVkMTRmMmNiOGY5ZWE5ZWIyNzJjYjkzZTBmODg4OTUzNGNkNTUzNDIxMDYzZjY1N2M3IiwidXNlcl9yZWdpc3RyYXRpb25fZHQiOjE2NzU0MTM4NTYsInZlcnNpb24iOjJ9.qkfDvjrXEM8X0MJj5-nEcMc5mwuI6jLcg5YS7_Qy9pcnr9OsVILl6iRrlsu-MwEnQ6Ik-PngybtewfnnuLgYDQyr3WE6tvRnfelBq0AxF_tp3XXBpGaVNBCHOUSCjEkqfgHeYBXo5bthoaro3qwlZa6lYSXMiQdkAHWcYoGzlcTkmlPsxwRdqm0C8f8A0wmnSoLJcA-kpk6wWI8OAXvS1X08bx0HY5Bf6EkcU1tsCxwqhk2UU0Nl7qGCl5vXHBZ3Wn3LIpoXgw-K7dUcoj88_B7xqMqogVFZiUIE-ga8m7BjMIMqQxpOBtUmOZJIC4ZyNuEBB31g2autEg6flPMksA"
	seller2AuthToken = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3NTc1MTQ1MDMsInVzZXIiOiIxNTg0NzMxMSIsInNoYXJkX2tleSI6IjMiLCJjbGllbnRfaWQiOiJzZWxsZXItcG9ydGFsIiwic2Vzc2lvbl9pZCI6ImRlOTFlNmRjODM5YTQ2ZDU4NWNiNDMwZTAyN2NmZWE2IiwidmFsaWRhdGlvbl9rZXkiOiI5OTdjMmQzYzgyODEwYzcwOGIyYjNkZTdlMjM5MzJlOGUzZjk4MzdmYzUwZjYyNjdlMWMzNmNhNDhhY2FkN2U0IiwidXNlcl9yZWdpc3RyYXRpb25fZHQiOjE2ODI0NTY0MjAsInZlcnNpb24iOjJ9.br-0_xHb2Z7TSEnHXX6fyclZ7kEP6QjDkyKYZ5-VZKw5ab7WT17mzuszrxtLHyvNzsSURksAX7QRI82fjtRSMRyyNADbGjI-uW9fcqY8pcg87mdLFbxRVdHR0ytJ5ScsVP7jOae_4RAm5p_qtt9O0vrBSs8OivIkrtFCzjElCVRyFwbcXikBPCd-zs0BsbvWKAxG5F1wvUtvHNBwxYec52-liQjJBrKUYlefNrNBNov4LuKgUz8DPUW6d4mQTLM4gVN6TAAu7hK_tWNr9w4bUwD86iURTBtXNI0N4HdTdkSlnaH7C4FqLZnXc5HcTt91hmDCoLfEuCPlC8fveoK1LQ"

	// Telegram IDs пользователей
	seller1TelegramID = int64(331871462) // ID первого продавца
	seller2TelegramID = int64(599835867) // ID второго продавца
	adminTelegramID   = int64(3572936)   // ID администратора

	// Telegram bot token
	telegramBotToken = "8083101312:AAHzCABhhWzbv5kEVxSQV6-rjNkF-9YuX7M" // Токен бота

	// Настройки программы
	checkInterval = 5 * time.Minute // Интервал проверки цен
)

// Структура для хранения данных продавца
type Seller struct {
	ID         int
	Name       string
	Token      string
	Cookie     string
	TelegramID int64
}

// Структура для хранения товара
type Product struct {
	ID       int
	NmID     int
	Name     string // Новое поле — название товара
	Price    float64
	Enabled  bool
	SellerID int
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

// Глобальные переменные
var (
	db      *sql.DB
	bot     *telebot.Bot
	sellers []Seller
)

func initDB() error {
	var err error
	db, err = sql.Open("sqlite", "./products.db")
	if err != nil {
		return fmt.Errorf("ошибка открытия БД: %w", err)
	}

	// Проверяем подключение
	if err = db.Ping(); err != nil {
		return fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	// Создание таблиц
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS sellers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			token TEXT NOT NULL,
			cookie TEXT NOT NULL,
			telegram_id INTEGER NOT NULL
		);
		
		CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			nm_id INTEGER NOT NULL UNIQUE,
			name TEXT NOT NULL DEFAULT '',
			price REAL NOT NULL,
			enabled BOOLEAN NOT NULL DEFAULT 1,
			seller_id INTEGER NOT NULL,
			FOREIGN KEY (seller_id) REFERENCES sellers (id)
		);
	`)
	if err != nil {
		return fmt.Errorf("ошибка создания таблиц: %w", err)
	}

	return nil
}

func loadInitialData() error {
	// Инициализация продавцов
	sellers = []Seller{
		{
			ID:         1,
			Name:       "Первый продавец",
			Token:      seller1AuthToken,
			Cookie:     "___wbu=2afbdd4a-2b25-44d0-a488-d2343f188ea6.1727788866; wbx-validation-key=a4aeb81b-3aea-4bca-a373-b724a179a919; _ym_uid=1726235106355624516; _ym_d=1740053345; external-locale=ru; x-supplier-id-external=be41cd8a-9260-412d-9445-cc8cf1d3aad0; device_id_guru=1980e7a5501-15d7d532f29b3d94; client_ip_guru=10.109.218.211; _ga=GA1.1.1439091344.1758371046; _ga_TXRZMJQDFE=GS2.1.s1759124403$o7$g0$t1759124409$j54$l0$h0; _wbauid=474842461759384401; __zzatw-wb=MDA0dC0cTHtmcDhhDHEWTT17CT4VHThHKHIzd2UqQWchYU1bIzVRP0FaW1Q4NmdBEXUmCQg3LGBwVxlRExpceEdXeiweGnpzJ1N/EV5GRWllbQwtUlFRS19/Dg4/aU5ZQ11wS3E6EmBWGB5CWgtMeFtLKRZHGzJhXkZpdRUNDQ5iQ0ImdVo7aR9jfFsfdQ5RMywhGjBrWFQPPxY/dF9vG3siXyoIJGM1Xz9EaVhTMCpYQXt1J3Z+KmUzPGwiaEphI0dVU3wuHQ1pN2wXPHVlLwkxLGJ5MVIvE0tsP0caRFpbQDsyVghDQE1HFF9BWncyUlFRS2EQR0lrZU5TQixmG3EVTQgNND1aciIPWzklWAgSPwsmIBYIbitTCwthQEpxbxt/Nl0cOWMRCxl+OmNdRkc3FSR7dSYKCTU3YnAvTCB7SykWRxsyYV5GaXUVCTwPXHB1dSwmRGcjX0RdIEURSgopHRZ0JlZXOkFccUQmLF07VxlRDxZhDhYYRRcje0I3Yhk4QhgvPV8/YngiD2lIYCRMWFV9KRkXe3AoS3FPLH12X30beylOIA0lVBMhP05yGOqeEw==; cfidsw-wb=3Iw2TFgyiR5qtBpDWnq/sLMTbxmNfGv4yl36FIPSmLJ38sKOgQidWxgEfaja0s7LQuUA5Tr3Q9j+ovFb0V2zV4FgFi/YfmoT6j0eGvJNPedQBJdgj3RfKRlCxyOiMRgXJTQEpEBni2cwIuknMwRDAJI2LG3QJG0pBQAyMQ==",
			TelegramID: seller1TelegramID,
		},
		{
			ID:         2,
			Name:       "Второй продавец",
			Token:      seller2AuthToken,
			Cookie:     "_wbauid=5666810631754830452; wbx-validation-key=f205c486-d051-42b8-8a77-86bb72e60283; x-supplier-id-external=df62fdc4-c58a-41dc-9aed-caf62c76df5f",
			TelegramID: seller2TelegramID,
		},
	}

	// Сохранение продавцов в БД
	for _, seller := range sellers {
		_, err := db.Exec(`
			INSERT OR IGNORE INTO sellers (id, name, token, cookie, telegram_id) 
			VALUES (?, ?, ?, ?, ?)`,
			seller.ID, seller.Name, seller.Token, seller.Cookie, seller.TelegramID)
		if err != nil {
			return fmt.Errorf("ошибка сохранения продавца %d: %w", seller.ID, err)
		}
	}

	// Загрузка начальных товаров для первого продавца
	productsSeller1 := map[int]struct {
		Price float64
		Name  string
	}{
		439740235: {598.00, "Товар 439740235"},
		363561833: {2431.00, "Товар 363561833"},
		355039724: {2756.00, "Товар 355039724"},
		420175308: {995.00, "Товар 420175308"},
		445719497: {351.00, "Товар 445719497"},
		444947468: {637.00, "Товар 444947468"},
		450517748: {1190.00, "Товар 450517748"},
		451852395: {1157.00, "Товар 451852395"},
		447703683: {2431.00, "Товар 447703683"},
		413320662: {1950.00, "Товар 413320662"},
		452613966: {1950.00, "Товар 452613966"},
		455308681: {826.00, "Товар 455308681"},
		455874194: {936.00, "Товар 455874194"},
		465007169: {533.00, "Товар 465007169"},
		466364173: {1612.00, "Товар 466364173"},
		472724832: {1352.00, "Товар 472724832"},
		458176275: {1664.00, "Товар 458176275"},
		485867509: {676.00, "Товар 485867509"},
		486685652: {1235.00, "Товар 486685652"},
		449727119: {1112.00, "Товар 449727119"},
		492714507: {878.00, "Товар 492714507"},
		492947914: {578.00, "Товар 492947914"},
		498858071: {1047.00, "Товар 498858071"},
		500564205: {826.00, "Товар 500564205"},
	}

	// Загрузка начальных товаров для второго продавца
	productsSeller2 := map[int]struct {
		Price float64
		Name  string
	}{
		486062217: {800.00, "Товар 486062217"},
		483028809: {640.00, "Товар 483028809"},
		473997083: {550.00, "Товар 473997083"},
		478334856: {615.00, "Товар 478334856"},
		485654591: {750.00, "Товар 485654591"},
		471430353: {755.00, "Товар 471430353"},
		472057995: {2100.00, "Товар 472057995"},
		473066411: {810.00, "Товар 473066411"},
		476823280: {750.00, "Товар 476823280"},
		475454890: {520.00, "Товар 475454890"},
		475499053: {520.00, "Товар 475499053"},
		480109053: {520.00, "Товар 480109053"},
		471832484: {750.00, "Товар 471832484"},
		470975205: {1200.00, "Товар 470975205"},
		493044219: {780.00, "Товар 493044219"},
		493490629: {570.00, "Товар 493490629"},
		495179694: {767.00, "Товар 495179694"},
		496076265: {670.00, "Товар 496076265"},
		496941899: {871.00, "Товар 496941899"},
		496570292: {735.00, "Товар 496570292"},
		499065435: {2470.00, "Товар 499065435"},
		505166842: {598.00, "Товар 505166842"},
		517572489: {3300.00, "Товар 517572489"},
		524447299: {2236.00, "Товар 524447299"},
		525352390: {2470.00, "Товар 525352390"},
		525977899: {2691.00, "Товар 525977899"},
		534976766: {1378.00, "Товар 534976766"},
		528079095: {3627.00, "Товар 528079095"},
		497033245: {871.00, "Товар 497033245"},
	}

	// Сохранение товаров в БД
	for nmID, data := range productsSeller1 {
		_, err := db.Exec(`
			INSERT OR IGNORE INTO products (nm_id, name, price, enabled, seller_id) 
			VALUES (?, ?, ?, ?, ?)`,
			nmID, data.Name, data.Price, true, 1)
		if err != nil {
			return fmt.Errorf("ошибка сохранения товара %d для продавца 1: %w", nmID, err)
		}
	}

	for nmID, data := range productsSeller2 {
		_, err := db.Exec(`
			INSERT OR IGNORE INTO products (nm_id, name, price, enabled, seller_id) 
			VALUES (?, ?, ?, ?, ?)`,
			nmID, data.Name, data.Price, true, 2)
		if err != nil {
			return fmt.Errorf("ошибка сохранения товара %d для продавца 2: %w", nmID, err)
		}
	}

	return nil
}

func getSellerByTelegramID(telegramID int64) (*Seller, error) {
	for i := range sellers {
		if sellers[i].TelegramID == telegramID || telegramID == adminTelegramID {
			return &sellers[i], nil
		}
	}
	return nil, fmt.Errorf("пользователь не найден")
}

func getProductsBySellerID(sellerID int) ([]Product, error) {
	rows, err := db.Query("SELECT id, nm_id, name, price, enabled, seller_id FROM products WHERE seller_id = ?", sellerID)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к БД: %w", err)
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.NmID, &p.Name, &p.Price, &p.Enabled, &p.SellerID)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		products = append(products, p)
	}

	return products, nil
}

func getAllProducts() ([]Product, error) {
	rows, err := db.Query("SELECT id, nm_id, name, price, enabled, seller_id FROM products")
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к БД: %w", err)
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.NmID, &p.Name, &p.Price, &p.Enabled, &p.SellerID)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		products = append(products, p)
	}

	return products, nil
}

func addProduct(nmID int, price float64, sellerID int, name string) error {
	_, err := db.Exec("INSERT OR REPLACE INTO products (nm_id, name, price, enabled, seller_id) VALUES (?, ?, ?, ?, ?)",
		nmID, name, price, true, sellerID)
	if err != nil {
		return fmt.Errorf("ошибка добавления товара: %w", err)
	}
	return nil
}

func updateProductPriceByNmID(nmID int, newPrice float64, sellerID int) error {
	result, err := db.Exec("UPDATE products SET price = ? WHERE nm_id = ? AND seller_id = ?",
		newPrice, nmID, sellerID)
	if err != nil {
		return fmt.Errorf("ошибка обновления цены: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("товар с nmID=%d для продавца %d не найден", nmID, sellerID)
	}

	return nil
}

func updateProductName(nmID int, name string, sellerID int) error {
	result, err := db.Exec("UPDATE products SET name = ? WHERE nm_id = ? AND seller_id = ?", name, nmID, sellerID)
	if err != nil {
		return fmt.Errorf("ошибка обновления названия: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("товар с nmID=%d для продавца %d не найден", nmID, sellerID)
	}

	return nil
}

func toggleProductStatusByNmID(nmID int, sellerID int) error {
	_, err := db.Exec("UPDATE products SET enabled = NOT enabled WHERE nm_id = ? AND seller_id = ?",
		nmID, sellerID)
	if err != nil {
		return fmt.Errorf("ошибка изменения статуса: %w", err)
	}
	return nil
}

func deleteProductByNmID(nmID int, sellerID int) error {
	_, err := db.Exec("DELETE FROM products WHERE nm_id = ? AND seller_id = ?",
		nmID, sellerID)
	if err != nil {
		return fmt.Errorf("ошибка удаления товара: %w", err)
	}
	return nil
}

func getProductByNmID(nmID int, sellerID int) (*Product, error) {
	var p Product
	err := db.QueryRow("SELECT id, nm_id, name, price, enabled, seller_id FROM products WHERE nm_id = ? AND seller_id = ?",
		nmID, sellerID).Scan(&p.ID, &p.NmID, &p.Name, &p.Price, &p.Enabled, &p.SellerID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения товара: %w", err)
	}
	return &p, nil
}

func getProductByNmIDForAnySeller(nmID int) (*Product, error) {
	var p Product
	err := db.QueryRow("SELECT id, nm_id, name, price, enabled, seller_id FROM products WHERE nm_id = ?",
		nmID).Scan(&p.ID, &p.NmID, &p.Name, &p.Price, &p.Enabled, &p.SellerID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения товара: %w", err)
	}
	return &p, nil
}

func setupTelegramBot() error {
	var err error
	bot, err = telebot.NewBot(telebot.Settings{
		Token:  telegramBotToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return fmt.Errorf("ошибка инициализации бота: %w", err)
	}

	// Обработчики команд
	bot.Handle("/start", func(c telebot.Context) error {
		helpText := "Добро пожаловать в репрайсер Wildberries! Доступные команды:\n\n" +
			"🔸 /products — список ваших товаров (с артикулом, названием и ценой)\n" +
			"🔸 /add nmID цена [название] — добавить новый товар\n" +
			"🔸 /price nmID новая_цена — изменить целевую цену товара\n" +
			"🔸 /name nmID Название — установить или изменить название товара\n" +
			"🔸 /toggle nmID — включить/выключить отслеживание цены\n" +
			"🔸 /delete nmID — удалить товар из системы\n\n" +
			"💡 Чтобы снова увидеть этот список — просто отправьте команду /start"

		return c.Send(helpText)
	})

	bot.Handle("/products", handleProductsList)
	bot.Handle("/add", handleAddProduct)
	bot.Handle("/price", handleUpdatePrice)
	bot.Handle("/name", handleSetName)
	bot.Handle("/toggle", handleToggleProduct)
	bot.Handle("/delete", handleDeleteProduct)

	return nil
}

// Вспомогательная функция: отправляет сообщение + напоминание о /start
func sendWithHelp(c telebot.Context, text string) error {
	fullText := text + "\n\n💡 Чтобы посмотреть все команды — отправьте /start"
	return c.Send(fullText)
}

func handleProductsList(c telebot.Context) error {
	telegramID := c.Sender().ID
	seller, err := getSellerByTelegramID(telegramID)
	if err != nil {
		return sendWithHelp(c, "У вас нет доступа к этой функции.")
	}

	var products []Product
	if telegramID == adminTelegramID {
		products, err = getAllProducts()
	} else {
		products, err = getProductsBySellerID(seller.ID)
	}

	if err != nil {
		return sendWithHelp(c, "Ошибка получения списка товаров.")
	}

	if len(products) == 0 {
		return sendWithHelp(c, "У вас пока нет товаров.")
	}

	var message strings.Builder
	message.WriteString("Ваши товары:\n\n")

	for _, product := range products {
		status := "✅ Вкл"
		if !product.Enabled {
			status = "❌ Выкл"
		}

		sellerName := ""
		if telegramID == adminTelegramID {
			for _, s := range sellers {
				if s.ID == product.SellerID {
					sellerName = fmt.Sprintf(" (%s)", s.Name)
					break
				}
			}
		}

		name := product.Name
		if name == "" {
			name = "(без названия)"
		}

		message.WriteString(fmt.Sprintf("NM: %d | %s | Цена: %.2f | %s%s\n",
			product.NmID, name, product.Price, status, sellerName))
	}

	return sendWithHelp(c, message.String())
}

func handleAddProduct(c telebot.Context) error {
	telegramID := c.Sender().ID
	seller, err := getSellerByTelegramID(telegramID)
	if err != nil {
		return sendWithHelp(c, "У вас нет доступа к этой функции.")
	}

	args := strings.Fields(c.Message().Text)[1:]
	if len(args) < 2 {
		return sendWithHelp(c, "Использование: /add nmID цена [название]")
	}

	nmID, err := strconv.Atoi(args[0])
	if err != nil {
		return sendWithHelp(c, "Неверный формат nmID.")
	}

	price, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return sendWithHelp(c, "Неверный формат цены.")
	}

	name := "Товар без названия"
	if len(args) > 2 {
		name = strings.Join(args[2:], " ")
	}

	targetSellerID := seller.ID
	if telegramID == adminTelegramID && len(args) > 3 {
		targetSellerID, err = strconv.Atoi(args[3])
		if err != nil {
			return sendWithHelp(c, "Неверный ID продавца.")
		}
	}

	err = addProduct(nmID, price, targetSellerID, name)
	if err != nil {
		return sendWithHelp(c, "Ошибка добавления товара.")
	}

	return sendWithHelp(c, fmt.Sprintf("✅ Товар %d (%s) добавлен с ценой %.2f", nmID, name, price))
}

func handleUpdatePrice(c telebot.Context) error {
	telegramID := c.Sender().ID
	seller, err := getSellerByTelegramID(telegramID)
	if err != nil {
		return sendWithHelp(c, "У вас нет доступа к этой функции.")
	}

	args := strings.Fields(c.Message().Text)[1:]
	if len(args) < 2 {
		return sendWithHelp(c, "Использование: /price nmID новая_цена")
	}

	nmID, err := strconv.Atoi(args[0])
	if err != nil {
		return sendWithHelp(c, "Неверный формат nmID.")
	}

	newPrice, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return sendWithHelp(c, "Неверный формат цены.")
	}

	targetSellerID := seller.ID
	if telegramID == adminTelegramID {
		product, err := getProductByNmIDForAnySeller(nmID)
		if err != nil {
			return sendWithHelp(c, "Товар не найден.")
		}
		targetSellerID = product.SellerID
	} else {
		_, err = getProductByNmID(nmID, seller.ID)
		if err != nil {
			return sendWithHelp(c, "Товар не найден или у вас нет к нему доступа.")
		}
	}

	err = updateProductPriceByNmID(nmID, newPrice, targetSellerID)
	if err != nil {
		return sendWithHelp(c, "Ошибка обновления цены.")
	}

	msg := fmt.Sprintf("✅ Цена товара NM %d обновлена на %.2f", nmID, newPrice)
	return sendWithHelp(c, msg)
}

func handleSetName(c telebot.Context) error {
	telegramID := c.Sender().ID
	seller, err := getSellerByTelegramID(telegramID)
	if err != nil {
		return sendWithHelp(c, "У вас нет доступа к этой функции.")
	}

	args := strings.Fields(c.Message().Text)[1:]
	if len(args) < 2 {
		return sendWithHelp(c, "Использование: /name nmID Название")
	}

	nmID, err := strconv.Atoi(args[0])
	if err != nil {
		return sendWithHelp(c, "Неверный формат nmID.")
	}

	name := strings.Join(args[1:], " ")

	targetSellerID := seller.ID
	if telegramID == adminTelegramID {
		product, err := getProductByNmIDForAnySeller(nmID)
		if err != nil {
			return sendWithHelp(c, "Товар не найден.")
		}
		targetSellerID = product.SellerID
	} else {
		_, err = getProductByNmID(nmID, seller.ID)
		if err != nil {
			return sendWithHelp(c, "Товар не найден или у вас нет к нему доступа.")
		}
	}

	err = updateProductName(nmID, name, targetSellerID)
	if err != nil {
		return sendWithHelp(c, "Ошибка обновления названия.")
	}

	return sendWithHelp(c, fmt.Sprintf("✅ Название товара NM %d установлено: %s", nmID, name))
}

func handleToggleProduct(c telebot.Context) error {
	telegramID := c.Sender().ID
	seller, err := getSellerByTelegramID(telegramID)
	if err != nil {
		return sendWithHelp(c, "У вас нет доступа к этой функции.")
	}

	args := strings.Fields(c.Message().Text)[1:]
	if len(args) < 1 {
		return sendWithHelp(c, "Использование: /toggle nmID")
	}

	nmID, err := strconv.Atoi(args[0])
	if err != nil {
		return sendWithHelp(c, "Неверный формат nmID.")
	}

	targetSellerID := seller.ID
	if telegramID == adminTelegramID {
		product, err := getProductByNmIDForAnySeller(nmID)
		if err != nil {
			return sendWithHelp(c, "Товар не найден.")
		}
		targetSellerID = product.SellerID
	} else {
		_, err = getProductByNmID(nmID, seller.ID)
		if err != nil {
			return sendWithHelp(c, "Товар не найден или у вас нет к нему доступа.")
		}
	}

	err = toggleProductStatusByNmID(nmID, targetSellerID)
	if err != nil {
		return sendWithHelp(c, "Ошибка изменения статуса товара.")
	}

	updatedProduct, _ := getProductByNmID(nmID, targetSellerID)
	status := "включен"
	if !updatedProduct.Enabled {
		status = "выключен"
	}

	return sendWithHelp(c, fmt.Sprintf("✅ Товар NM %d %s", nmID, status))
}

func handleDeleteProduct(c telebot.Context) error {
	telegramID := c.Sender().ID
	seller, err := getSellerByTelegramID(telegramID)
	if err != nil {
		return sendWithHelp(c, "У вас нет доступа к этой функции.")
	}

	args := strings.Fields(c.Message().Text)[1:]
	if len(args) < 1 {
		return sendWithHelp(c, "Использование: /delete nmID")
	}

	nmID, err := strconv.Atoi(args[0])
	if err != nil {
		return sendWithHelp(c, "Неверный формат nmID.")
	}

	targetSellerID := seller.ID
	if telegramID == adminTelegramID {
		product, err := getProductByNmIDForAnySeller(nmID)
		if err != nil {
			return sendWithHelp(c, "Товар не найден.")
		}
		targetSellerID = product.SellerID
	} else {
		_, err = getProductByNmID(nmID, seller.ID)
		if err != nil {
			return sendWithHelp(c, "Товар не найден или у вас нет к нему доступа.")
		}
	}

	err = deleteProductByNmID(nmID, targetSellerID)
	if err != nil {
		return sendWithHelp(c, "Ошибка удаления товара.")
	}

	return sendWithHelp(c, fmt.Sprintf("✅ Товар NM %d удален", nmID))
}

func main() {
	// Настройка логирования
	logFile, err := os.OpenFile("repricer.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Ошибка открытия файла логов: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	log.Println("Проверка доступности SQLite драйвера...")
	// Проверяем, что драйвер зарегистрирован
	drivers := sql.Drivers()
	log.Printf("Доступные драйверы: %v", drivers)

	// Инициализация базы данных
	log.Println("Инициализация базы данных...")
	err = initDB()
	if err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}

	// Загрузка начальных данных
	log.Println("Загрузка начальных данных...")
	err = loadInitialData()
	if err != nil {
		log.Fatalf("Ошибка загрузки начальных данных: %v", err)
	}

	// Настройка Telegram бота
	log.Println("Настройка Telegram бота...")
	err = setupTelegramBot()
	if err != nil {
		log.Printf("Предупреждение: Не удалось настроить Telegram бота: %v", err)
		log.Printf("Продолжаем работу без бота...")
	} else {
		go func() {
			log.Println("Запуск Telegram бота...")
			bot.Start()
		}()
	}

	log.Println("Запуск репрайсера Wildberries")
	log.Printf("Количество продавцов: %d", len(sellers))
	for i, seller := range sellers {
		products, _ := getProductsBySellerID(seller.ID)
		log.Printf("Продавец %d: %s, товаров: %d", i+1, seller.Name, len(products))
	}
	log.Printf("Интервал проверки: %v\n", checkInterval)

	for {
		log.Println("========================================")
		log.Println("Начало нового цикла проверки...")

		walletDiscount, err := getWalletDiscount()
		if err != nil {
			log.Printf("Ошибка получения скидки кошелька: %v", err)
			time.Sleep(checkInterval)
			continue
		}
		log.Printf("Текущая скидка кошелька: %d%%", walletDiscount)

		for _, seller := range sellers {
			processSellerProducts(seller, walletDiscount)
		}

		log.Printf("Цикл завершен. Следующая проверка через %v\n", checkInterval)
		time.Sleep(checkInterval)
	}
}

func processSellerProducts(seller Seller, walletDiscount int) {
	log.Printf("--- Обработка товаров для %s ---", seller.Name)

	products, err := getProductsBySellerID(seller.ID)
	if err != nil {
		log.Printf("Ошибка получения товаров для %s: %v", seller.Name, err)
		return
	}

	for _, product := range products {
		if !product.Enabled {
			continue
		}

		price, sellerDiscount, wbDiscount, err := getProductInfo(product.NmID, seller.Token, seller.Cookie)
		if err != nil {
			log.Printf("Ошибка получения информации о товаре %d: %v", product.NmID, err)
			continue
		}

		finalPrice := calculateFinalPrice(price, sellerDiscount, wbDiscount, walletDiscount)
		log.Printf("%s - Товар %d - Цена: %.2f, Скидка продавца: %d%%, Скидка WB: %d%%, Финальная цена: %.2f, Минимальная цена: %.2f",
			seller.Name, product.NmID, price, sellerDiscount, wbDiscount, finalPrice, product.Price)

		if finalPrice < product.Price || finalPrice > product.Price+1 {
			log.Printf("ТРЕБУЕТСЯ КОРРЕКТИРОВКА: Финальная цена %.2f вне диапазона [%.2f, %.2f]",
				finalPrice, product.Price, product.Price+1)

			newPrice, newDiscount := findOptimalPrice(price, sellerDiscount, wbDiscount, walletDiscount, product.Price)
			err = updateProductPriceAPI(product.NmID, newPrice, newDiscount, seller.Token, seller.Cookie)
			if err != nil {
				log.Printf("Ошибка обновления цены для товара %d: %v", product.NmID, err)
			} else {
				log.Printf("Цена обновлена: Новая цена = %.2f, Новая скидка = %d%%", newPrice, newDiscount)
			}
		}
	}
}

func getWalletDiscount() (int, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://user-grade.wildberries.ru/api/v5/grade?curr=RUB", nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,zh-CN;q=0.8,zh;q=0.7,en-US;q=0.6,en;q=0.5")
	req.Header.Set("Authorization", "Bearer "+walletAuthToken)
	req.Header.Set("Dnt", "1")
	req.Header.Set("Origin", "  https://www.wildberries.ru  ")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", "https://www.wildberries.ru/  ")
	req.Header.Set("Sec-Ch-Ua", `"Google Chrome";v="137", "Chromium";v="137", "Not/A)Brand";v="24"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"macOS"`)
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

	return int(math.Round(response.Payload.Payments[0].UpridDiscount)), nil
}

func getProductInfo(nmId int, sellerToken string, cookie string) (price float64, sellerDiscount, wbDiscount int, err error) {
	url := fmt.Sprintf("https://discounts-prices.wildberries.ru/ns/dp-api/discounts-prices/suppliers/api/v1/nm/info?nmID=%d", nmId)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, 0, 0, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,zh-CN;q=0.8,zh;q=0.7,en-US;q=0.6,en;q=0.5")
	req.Header.Set("Authorizev3", sellerToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Dnt", "1")
	req.Header.Set("Origin", "  https://seller.wildberries.ru  ")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", "https://seller.wildberries.ru/  ")
	req.Header.Set("Sec-Ch-Ua", `"Google Chrome";v="137", "Chromium";v="137", "Not/A)Brand";v="24"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36")
	req.Header.Set("Cookie", cookie)

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

	return response.Data.Info.Price,
		int(math.Round(response.Data.Info.Discount)),
		int(math.Round(response.Data.Info.DiscountOnSite)),
		nil
}

func calculateFinalPrice(price float64, sellerDiscount, wbDiscount, walletDiscount int) float64 {
	wbDiscountAmount := price * float64(wbDiscount) / 100
	wbDiscountRounded := math.Ceil(wbDiscountAmount)
	currentPrice := price - wbDiscountRounded

	sellerDiscountAmount := currentPrice * float64(sellerDiscount) / 100
	sellerDiscountRounded := math.Ceil(sellerDiscountAmount)
	currentPrice -= sellerDiscountRounded

	walletDiscountAmount := currentPrice * float64(walletDiscount) / 100
	walletDiscountRounded := math.Ceil(walletDiscountAmount)
	finalPrice := currentPrice - walletDiscountRounded

	return math.Round(finalPrice*100) / 100
}

func findOptimalPrice(currentPrice float64, currentDiscount, wbDiscount, walletDiscount int, minPrice float64) (float64, int) {
	const discountRange = 10

	bestPrice := currentPrice
	bestDiscount := currentDiscount
	bestDiff := math.MaxFloat64

	for offset := -discountRange; offset <= discountRange; offset++ {
		newDiscount := currentDiscount + offset
		if newDiscount < 0 || newDiscount > 100 {
			continue
		}

		price, diff := findPriceForDiscount(newDiscount, wbDiscount, walletDiscount, minPrice)
		if diff >= 0 && diff < bestDiff {
			bestPrice = price
			bestDiscount = newDiscount
			bestDiff = diff
		}
	}

	if bestDiff == math.MaxFloat64 {
		return currentPrice, currentDiscount
	}

	return bestPrice, bestDiscount
}

func findPriceForDiscount(discount, wbDiscount, walletDiscount int, minPrice float64) (float64, float64) {
	low, high := minPrice, minPrice*2
	bestPrice := low
	bestDiff := math.MaxFloat64

	for high-low > 0.01 {
		mid := (low + high) / 2
		finalPrice := calculateFinalPrice(mid, discount, wbDiscount, walletDiscount)
		diff := finalPrice - minPrice

		if diff >= 0 {
			if diff < bestDiff {
				bestDiff = diff
				bestPrice = mid
			}
			high = mid
		} else {
			low = mid
		}
	}

	finalPrice := calculateFinalPrice(bestPrice, discount, wbDiscount, walletDiscount)
	return bestPrice, finalPrice - minPrice
}

func updateProductPriceAPI(nmId int, newPrice float64, newDiscount int, sellerToken string, cookie string) error {
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

	jsonData, _ := json.Marshal(payload)

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://discounts-prices.wildberries.ru/ns/dp-api/discounts-prices/suppliers/api/v1/nm/upload/task?checkChange=true", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,zh-CN;q=0.8,zh;q=0.7,en-US;q=0.6,en;q=0.5")
	req.Header.Set("Authorizev3", sellerToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Dnt", "1")
	req.Header.Set("Origin", "  https://seller.wildberries.ru  ")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", "https://seller.wildberries.ru/  ")
	req.Header.Set("Sec-Ch-Ua", `"Google Chrome";v="137", "Chromium";v="137", "Not/A)Brand";v="24"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36")
	req.Header.Set("Cookie", cookie)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("статус: %d, тело: %s", resp.StatusCode, string(body))
	}

	log.Printf("Отправка данных для nmID %d: цена=%d, скидка=%d", nmId, int(math.Round(newPrice)), newDiscount)
	return nil
}
