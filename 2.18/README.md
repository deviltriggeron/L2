## Календарь с ивентами

HTTP-сервер для работы с небольшим календарем событий.
Для логирования используется middleware в stdout (метод, URL, время)


## 🐳 Запуск
- Команда для запуска
```
go build ./cmd/main.go
./main
```

## POST запросы

- userid ожидает uuid.UUID
- eventID можно увидеть после создания ивента

```
  curl -X POST "http://localhost:8080/create_event" \
    -H "Content-Type: application/json" \
    -d '{
      "userid": "550e8400-e29b-41d4-a716-446655440000",
      "date": {                                
        "year": 2025,         
        "month": 10,                      
        "day": 12
      },
      "info": "Встреча с клиентом"
    }'

  curl -X POST "http://localhost:8080/update_event/{eventID}" \
    -H "Content-Type: application/json" \
    -d '{
      "userid": "550e8400-e29b-41d4-a716-446655440000",
      "date": {
        "year": 2025,
        "month": 11,
        "day": 15
      },
      "info": "Новая встреча"
    }'

  curl -X POST "http://localhost:8080/delete_event/{eventID}" \
    -H "Content-Type: application/json" \
    -d '"550e8400-e29b-41d4-a716-446655440000"'

```

- Пример успешного создания ивента

```
  {"EventID":"d01c8f2a-1c23-4460-816c-1682a1b5c7d6","UserID":"550e8400-e29b-41d4-a716-446655440000"}
```
## GET запросы

- Есть поддержка query strings
```
  curl -X GET "http://localhost:8080/events_for_day?id=550e8400-e29b-41d4-a716-446655440000&date=2025-October-12"

  curl -X GET "http://localhost:8080/events_for_week?id=550e8400-e29b-41d4-a716-446655440000&date=2025-October-12"

  curl -X GET "http://localhost:8080/events_for_month?id=550e8400-e29b-41d4-a716-446655440000&date=2025-October-12"
```

- Пример успешного выполнения GET запроса

```
  {"Events: ":
    [
      {
      "EventID":"cdabd024-ef97-4aa3-ae1c-28c629652fa5",
      "Date":{"Year":2025,"Month":10,"Day":12},
      "EventInfo":"Встреча с клиентом",
      "CreateDate":"2025-10-04T22:25:01.547475+09:00"
      },
    ]
  }
```

## Тесты

- Команда для запуска тестов

```
  go test -v tests/calendar_test.go
```
