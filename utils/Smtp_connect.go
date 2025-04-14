package utils

import (
	"crypto/tls"
	"fmt"
	"github.com/go-mail/mail/v2" // Импортируем библиотеку gomail
	"log"
)

// --- Конфигурация SMTP ---
const (
	smtpHost = "smtp.google.com"
	smtpPort = 587
	smtpUser = "tkachenko2223@google.com"
	smtpPass = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
)

// createMessage формирует email сообщение.
// Принимает адрес получателя, тему и тело письма (в формате HTML).
// Возвращает указатель на объект mail.Message.
func createMessage(to string, subject string, body string) *mail.Message {
	m := mail.NewMessage()

	// Устанавливаем заголовки письма
	m.SetHeader("From", smtpUser)   // От кого письмо (обычно совпадает с smtpUser)
	m.SetHeader("To", to)           // Кому отправляем
	m.SetHeader("Subject", subject) // Тема письма

	// Устанавливаем тело письма в формате HTML
	// text/html указывает почтовому клиенту, что содержимое нужно рендерить как HTML
	m.SetBody("text/html", body)

	// Опционально: можно добавить альтернативное тело в виде простого текста
	// для почтовых клиентов, которые не поддерживают HTML.
	// m.AddAlternative("text/plain", "Plain text version of the email body")

	return m
}

// createDialer настраивает и возвращает объект Dialer для подключения к SMTP-серверу.
func createDialer() *mail.Dialer {
	// Создаем Dialer с указанием хоста, порта, имени пользователя и пароля
	// smtpPort 587 обычно использует STARTTLS
	// smtpPort 465 обычно использует SSL/TLS с самого начала
	d := mail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	// Настройка TLS. Большинство современных SMTP-серверов требуют TLS.
	// InsecureSkipVerify: false - ВАЖНО! Не отключайте проверку сертификата в продакшене.
	// Установите true ТОЛЬКО для локального тестирования с самоподписанными сертификатами,
	// если вы точно знаете, что делаете.
	d.TLSConfig = &tls.Config{
		InsecureSkipVerify: false,    // Оставляем false для безопасности
		ServerName:         smtpHost, // Указываем имя хоста для проверки сертификата
	}

	return d
}

// sendEmail отправляет подготовленное сообщение с использованием заданного Dialer.
// Принимает Dialer и сообщение. Возвращает ошибку, если отправка не удалась.
func sendEmail(d *mail.Dialer, m *mail.Message) error {
	// Попытка установить соединение с SMTP-сервером и отправить письмо
	if err := d.DialAndSend(m); err != nil {
		// Логируем ошибку для отладки
		log.Printf("Ошибка при отправке email: %v", err)
		// Возвращаем ошибку, оборачивая исходную для сохранения контекста
		return fmt.Errorf("не удалось отправить email: %w", err)
	}
	// Логируем успешную отправку
	log.Println("Email успешно отправлен!")
	return nil
}

// --- Пример использования (имитация отправки уведомления о платеже) ---

func SendPaymentEmail(userEmail string, amount float64) error {
	// 1. Создание контента (тела) письма в формате HTML
	content := fmt.Sprintf(`
        <h1>Спасибо за оплату!</h1>
        <p>Ваш платеж на сумму <strong>%.2f RUB</strong> успешно обработан.</p>
        <p>Детали операции будут доступны в вашем личном кабинете.</p>
        <hr>
        <small>Это письмо сгенерировано автоматически, пожалуйста, не отвечайте на него.</small>
    `, amount)

	// 2. Подготовка сообщения с помощью нашей функции createMessage
	log.Printf("Подготовка email для %s...", userEmail)
	m := createMessage(userEmail, "Платеж успешно проведен", content)

	// 3. Настройка подключения (Dialer) с помощью нашей функции createDialer
	log.Println("Настройка SMTP подключения...")
	d := createDialer()

	// 4. Отправка письма с помощью нашей функции sendEmail
	log.Println("Отправка email...")
	if err := sendEmail(d, m); err != nil {
		// Если sendEmail вернула ошибку, пробрасываем её выше
		return err
	}

	// Если ошибок не было, логируем успешное завершение операции
	log.Printf("Email уведомление успешно отправлено на %s", userEmail)
	return nil
}
