# EventFinder (review)

Это упрощённый MVP без Docker, без авторизации и без базы данных.

## Что работает
- backend на Go с in-memory данными
- frontend на HTML/CSS/JS
- загрузка списка событий
- добавление/удаление избранного

## Запуск backend
1. Откройте терминал в `review/backend`
2. Выполните:

```bash
cd review/backend
go run main.go
```

Backend запустится на `http://localhost:8081`.

## Запуск frontend
Откройте другой терминал в `review/frontend`.

Если у вас установлен Python:

```bash
cd review/frontend
python -m http.server 8000
```

Откройте в браузере `http://localhost:8000`.

Если Python нет, можно использовать любой локальный статический сервер.

## API
- `GET /api/events` — список событий
- `GET /api/favorites` — список избранного
- `POST /api/favorites` — добавить в избранное, JSON `{ "event_id": 1 }`
- `DELETE /api/favorites/:event_id` — убрать из избранного
