Based on the provided code context, I'll create a comprehensive README.md file in Russian that explains how to use the financial service, supported commands, and testing instructions.

```markdown
# Финансовый сервис

Это сервис для управления банковскими счетами, картами и кредитами с REST API интерфейсом.

## Требования

- Go 1.20+
- PostgreSQL 12+
- Установленные переменные окружения (см. Настройка)

## Установка и запуск

1. Клонируйте репозиторий:
```bash
git clone <repository-url>
cd final-project
```

2. Установите зависимости:
```bash
go mod download
```

3. Настройте базу данных:
```sql
CREATE DATABASE finance;
```

4. Запустите сервис:
```bash
go run Main.go
```

Сервис будет доступен по адресу `http://localhost:8080`

## Использование API

### Аутентификация

- **Регистрация**  
  `POST /register`  
  Тело запроса:
  ```json
  {
    "username": "имя_пользователя",
    "email": "email@example.com",
    "password": "пароль"
  }
  ```

- **Вход**  
  `POST /login`  
  Тело запроса:
  ```json
  {
    "email": "email@example.com",
    "password": "пароль"
  }
  ```
  Возвращает JWT токен для авторизованных запросов.

### Управление счетами

- **Создать счет**  
  `POST /accounts`  
  Тело запроса:
  ```json
  {
    "currency": "RUB"
  }
  ```

- **Пополнить счет**  
  `POST /accounts/{id}/deposit`  
  Тело запроса:
  ```json
  {
    "amount": 1000.00
  }
  ```

- **Снять средства**  
  `POST /accounts/{id}/withdraw`  
  Тело запроса:
  ```json
  {
    "amount": 500.00
  }
  ```

- **Перевод между счетами**  
  `POST /transfer`  
  Тело запроса:
  ```json
  {
    "fromAccountID": 1,
    "toAccountID": 2,
    "amount": 300.00
  }
  ```

### Управление картами

- **Создать виртуальную карту**  
  `POST /cards`  
  Тело запроса:
  ```json
  {
    "account_id": 1,
    "cvv": "123"
  }
  ```

- **Получить информацию о карте**  
  `GET /cards/{id}`

### Кредиты

- **Оформить кредит**  
  `POST /credits`  
  Тело запроса:
  ```json
  {
    "account_id": 1,
    "amount": 10000.00,
    "term": 12
  }
  ```

- **Получить информацию о кредите**  
  `GET /credits/{id}`

- **Получить график платежей**  
  `GET /credits/{creditId}/schedule`


## Автоматические платежи

Сервис автоматически обрабатывает платежи по кредитам каждый день в полночь. Для ручного запуска:
```bash
POST /credits/process-payments
```

## Логирование

Логи сохраняются в файл `app.log` и выводятся в консоль. Уровень логирования можно настроить в `utils/logger.go`.
```
