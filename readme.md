# Avito pvz
Решение тестового задания для Avito на позицию Intern Backend разработчик.

---

## Инструкция по запуску

Запускаем докер контейнер с проектом и бд
```bash
docker-compose up --build
```

После успешного запуска API будет доступно по адресу: `http://localhost:8080/api`.

---

## API
Реализованы методы для управления ПВЗ, приёмками товаров и товарами.

### Аутентификация и получение JWT-токена
**Эндпоинт:** `POST /api/dummyLogin`

Позволяет получить тестовый JWT-токен для роли `employee` или `moderator`.

#### Пример запроса:

```bash
curl --request POST \
  --url http://localhost:8080/api/dummyLogin \
  --header "Content-Type: application/json" \
  --data '{
    "role": "moderator"
  }'
```

#### Пример успешного ответа:

```json
  "eyJhbGciOiJIUzI1..."
```

Для всех защищённых эндпоинтов необходимо передавать токен в заголовке:

```bash
Authorization: Bearer <TOKEN>
```

---

### Управление ПВЗ

#### Создание нового ПВЗ

**Эндпоинт:** `POST /api/pvz`

Создаёт новый ПВЗ. Доступно только для модераторов.

#### Пример запроса:

```bash
curl --request POST \
  --url http://localhost:8080/api/pvz \
  --header "Authorization: Bearer <TOKEN>" \
  --header "Content-Type: application/json" \
  --data '{
    "city": "Москва"
  }'
```

#### Cписок доступных для создания городов.

| Город           | 
|-----------------|
| Москва          | 
| Санкт-Петербург | 
| Казань          | 

#### Пример успешного ответа:

```json
{
  "id": "b1a7c8e2-3c4d-4f5e-8a7b-9c6d8e2f3a4b",
  "city": "Москва",
  "registrationDate": "2025-04-18T12:00:00Z"
}
```

---

### Управление приёмками товаров

#### Создание новой приёмки

**Эндпоинт:** `POST /api/receptions`

Создаёт новую приёмку товаров для указанного ПВЗ. Доступно только для сотрудников ПВЗ.

#### Пример запроса:

```bash
curl --request POST \
  --url http://localhost:8080/api/receptions \
  --header "Authorization: Bearer <TOKEN>" \
  --header "Content-Type: application/json" \
  --data '{
    "pvzId": "b1a7c8e2-3c4d-4f5e-8a7b-9c6d8e2f3a4b"
  }'
```

#### Пример успешного ответа:

```json
{
  "id": "d2b7c8e2-3c4d-4f5e-8a7b-9c6d8e2f3a4b",
  "pvzId": "b1a7c8e2-3c4d-4f5e-8a7b-9c6d8e2f3a4b",
  "status": "in_progress",
  "dateTime": "2025-04-18T12:30:00Z"
}
```

---

#### Закрытие последней приёмки

**Эндпоинт:** `POST /api/pvz/{pvzId}/close_last_reception`

Закрывает последнюю активную приёмку товаров для указанного ПВЗ. Доступно только для сотрудников ПВЗ.

#### Пример запроса:

```bash
curl --request POST \
  --url http://localhost:8080/api/pvz/b1a7c8e2-3c4d-4f5e-8a7b-9c6d8e2f3a4b/close_last_reception \
  --header "Authorization: Bearer <TOKEN>"
```

#### Пример успешного ответа:

```json
{
  "id": "d2b7c8e2-3c4d-4f5e-8a7b-9c6d8e2f3a4b",
  "pvzId": "b1a7c8e2-3c4d-4f5e-8a7b-9c6d8e2f3a4b",
  "status": "closed",
  "dateTime": "2025-04-18T12:30:00Z"
}
```

---

### Управление товарами

#### Добавление товара в приёмку

**Эндпоинт:** `POST /api/products`

Добавляет товар в текущую активную приёмку. Доступно только для сотрудников ПВЗ.

#### Пример запроса:

```bash
curl --request POST \
  --url http://localhost:8080/api/products \
  --header "Authorization: Bearer <TOKEN>" \
  --header "Content-Type: application/json" \
  --data '{
    "pvzId": "b1a7c8e2-3c4d-4f5e-8a7b-9c6d8e2f3a4b",
    "type": "электроника"
  }'
```

#### Пример успешного ответа:

```json
{
  "id": "e3c7c8e2-3c4d-4f5e-8a7b-9c6d8e2f3a4b",
  "receptionId": "d2b7c8e2-3c4d-4f5e-8a7b-9c6d8e2f3a4b",
  "type": "electronics",
  "dateTime": "2025-04-18T12:45:00Z"
}
```

---

#### Удаление последнего добавленного товара

**Эндпоинт:** `POST /api/pvz/{pvzId}/delete_last_product`

Удаляет последний добавленный товар из текущей активной приёмки. Доступно только для сотрудников ПВЗ.

#### Пример запроса:

```bash
curl --request POST \
  --url http://localhost:8080/api/pvz/b1a7c8e2-3c4d-4f5e-8a7b-9c6d8e2f3a4b/delete_last_product \
  --header "Authorization: Bearer <TOKEN>"
```

#### Пример успешного ответа:

```json
{
  "message": "OK"
}
```

---

## Тестирование

Реализованы интеграционные тесты для проверки основных сценариев работы API:

- Создание ПВЗ.
- Создание приёмки товаров.
- Добавление товаров в приёмку.
- Закрытие приёмки товаров.

Для запуска интеграционных тестов необходимо запустить прокт и убрать префикс _ из названия файла тесто tests\_integration_test.go
---

## Структура проекта

- `internal/handler` — обработчики HTTP-запросов.
- `internal/service` — бизнес-логика.
- `internal/repository` — взаимодействие с базой данных.
- `tests` — интеграционные тесты.