package scheduler

import (
	"final-project/services"
	"log"
	"time"
)

// StartCreditPaymentScheduler запускает шедулер, который каждые 12 часов обрабатывает списание платежей по кредитам.
func StartCreditPaymentScheduler(loanService *services.LoanService) {
	ticker := time.NewTicker(12 * time.Hour)
	go func() {
		for {
			<-ticker.C
			log.Println("Запуск обработки автоматических списаний кредитных платежей...")
			err := loanService.ProcessCreditPayments()
			if err != nil {
				log.Println("Ошибка при списании платежей по кредитам:", err)
			}
		}
	}()
}
