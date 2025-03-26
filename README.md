# Clicks Counter

Небольшой сервис на Go для подсчёта кликов и получения статистики по баннерам.  

## Запуск

### С помощью Docker Compose

1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/yourusername/clicks-counter.git
   cd clicks-counter
   ```
2. Соберите и запустите:
   ```bash
   docker-compose up --build
   ```
   - Поднимется контейнер с PostgreSQL и контейнер с приложением на Go.  
   - Приложение будет доступно на `http://localhost:8080`.

3. (При необходимости) примените миграцию для создания таблицы:
   ```bash
   docker-compose exec db psql -U postgres -d clicks -f /migrations/001_init.sql
   ```
   Затем таблица `banner_clicks` будет готова к использованию.

### Пример взаимодействия

- **Проверка Health**:  
  ```bash
  curl http://localhost:8080/health
  # Ожидается "OK"
  ```
- **Инкремент клика** (bannerID=1):  
  ```bash
  curl http://localhost:8080/counter/1
  # HTTP 200 OK, тело: "OK"
  ```
- **Получить статистику**:
  ```bash
  curl -X POST http://localhost:8080/stats/1 \
    -H "Content-Type: application/json" \
    -d '{
      "tsFrom": "2025-03-01T00:00:00Z",
      "tsTo":   "2025-03-31T23:59:59Z"
    }'
  ```

## Нагрузочное тестирование

Для проверки производительности (например, 100–500 RPS) можно использовать утилиту [**wrk**](https://github.com/wg/wrk):

1. Установите wrk (в Linux через apt/brew, на Windows — через WSL).
2. Запустите тест:
   ```bash
   wrk -t4 -c100 -d30s http://localhost:8080/counter/1
   ```
   - `-t4` — 4 потока  
   - `-c100` — 100 одновременных подключений  
   - `-d30s` — 30 секунд теста  

В конце появится статистика по RPS, задержкам и ошибкам.  
