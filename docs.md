# Реализация сервиса для блокировки запросов


---
# Разворачивание сервиса

This app was developed for the OTUS golang developer course as a final project.
Application takes a link to a source image, resizes it
and returns resized image to a user.


## Установка и запуск
#### Запускаем тесты
```shell
make test
```

#### Запускаем линтеры
```shell
make lint
```


1. Копируем `.ENV.DEV` и изменяем, если нужно
```shell
cp .env.dev .env
```

2. Собираем образ приложения
```shell
make build
```

3. Запускаем приложение
```shell 
make run
```
W
4. Останавливаем приложение
```shell
make stop
```


---

# Документация по API

## Общие сведения

Все методы, кроме `/health` и `/metrics`, требуют авторизации через заголовок `Authorization: Bearer <token>`.

---

## Ограничения

- Для большинства методов действует Rate Limiter (ограничение частоты запросов).
- Для доступа требуется Bearer-токен.

---

## Служебные методы

### GET `/health`

Проверка состояния сервиса.

**Ответ:**
- 200 OK: сервис работает.

---

### GET `/metrics`

Метрики Prometheus.

---

## Посты

### POST `/createPost`

Создать новый пост.

**Тело запроса:**
```json
{
  "text": "string",   // обязательное поле
  "data": "string",   // необязательное поле
  "ip": "string"      // необязательное поле
}
```

**Ответ:**
- 200 OK: `{"message": "Post created successfully"}`
- 400 Bad Request: некорректные данные
- 401 Unauthorized: нет авторизации
- 429 Too Many Requests: превышен лимит запросов
- 500 Internal Server Error: ошибка сервера
---

## CRUD для ботов

### POST `/crud/bot`

Создать нового бота.

**Тело запроса:**
```json
{
  "id": "string",
  "name": "string",
  "rate": 0,
  "period": "string|null"
}
```

**Ответ:**
- 200 OK: бот создан
- 400/401: ошибка

---

### GET `/crud/bot/{id}`

Получить информацию о боте по ID.

**Ответ:**
- 200 OK: структура Bot
- 404 Not Found: не найден

---

### PUT `/crud/bot/{id}`

Обновить данные бота.

**Тело запроса:**  
```json
{
  "name": "string",
  "rate": 0,
  "period": "string|null"
}
```

**Ответ:**
- 200 OK: обновлено (структура Bot)
- 400/404: ошибка

---

### DELETE `/crud/bot/{id}`

Удалить бота.

**Ответ:**
- 200 OK: удалено
- 404 Not Found: не найден

---

## CRUD для токенов

### POST `/crud/token`

Создать токен.

**Тело запроса:**
```json
{
  "id": "string",
  "bot_id": "string",
  "rate": 0,
  "period": "string|null"
}
```

**Ответ:**
- 200 OK: токен создан (структура Token)

---

### GET `/crud/token/{id}`

Получить токен по ID.

**Ответ:**
- 200 OK: структура Token
- 404 Not Found: не найден

---

### PUT `/crud/token/{id}`

Обновить токен.

**Тело запроса:**  
```json
{
  "bot_id": "string",
  "rate": 0,
  "period": "string|null"
}
```

**Ответ:**
- 200 OK: обновлено (структура Token)
- 404 Not Found: не найден

---

### DELETE `/crud/token/{id}`

Удалить токен.

**Ответ:**
- 200 OK: удалено
- 404 Not Found: не найден

---

## CRUD для хуков

### POST `/crud/hook`

Создать хук.

**Тело запроса:**
```json
{
  "id": "string",
  "name": "string",
  "rate": 0,
  "period": "string|null"
}
```

**Ответ:**
- 200 OK: хук создан (структура Hook)

---

### GET `/crud/hook/{id}`

Получить хук по ID.

**Ответ:**
- 200 OK: структура Hook
- 404 Not Found: не найден

---

### PUT `/crud/hook/{id}`

Обновить хук.

**Тело запроса:**  
```json
{
  "name": "string",
  "rate": 0,
  "period": "string|null"
}
```

**Ответ:**
- 200 OK: обновлено (структура Hook)
- 404 Not Found: не найден

---

### DELETE `/crud/hook/{id}`

Удалить хук.

**Ответ:**
- 200 OK: удалено
- 404 Not Found: не найден

---

