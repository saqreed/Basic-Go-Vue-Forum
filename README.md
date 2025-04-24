# Basic Go-Vue Forum

Форум с базовым функционалом, построенный на Go (backend) и Vue.js (frontend).

## Технологии

### Backend
- Go 1.21
- Gorilla Mux (маршрутизация)
- JWT (аутентификация)
- PostgreSQL (база данных)
- SQLX (работа с БД)

### Frontend
- Vue 3
- Vue Router
- Pinia (управление состоянием)
- Vite (сборка)

## Функциональность

### Пользователи
- Регистрация и авторизация
- Профили пользователей
- Роли (пользователь/админ)
- Изменение пароля

### Посты
- Создание, редактирование, удаление постов
- Просмотр всех постов
- Фильтрация постов

### Комментарии
- Добавление комментариев к постам
- Редактирование и удаление комментариев
- Древовидная структура комментариев

### Админ-панель
- Управление пользователями
- Модерация постов
- Модерация комментариев
- Статистика

## Установка и запуск

### Требования
- Go 1.21+
- Node.js 18+
- PostgreSQL 14+

### Backend

```bash
cd backend
go mod download
go run cmd/main.go
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```

### База данных

1. Создайте базу данных PostgreSQL
2. Примените миграции из `backend/internal/database/init.sql`
3. Настройте подключение в `backend/config/config.go`

## API Endpoints

### Аутентификация
- POST /api/register - Регистрация
- POST /api/login - Авторизация
- POST /api/change-password - Смена пароля

### Профиль
- GET /api/profile - Получение профиля
- GET /api/profile/stats - Статистика пользователя
- PUT /api/profile/password - Изменение пароля

### Посты
- GET /api/posts - Список постов
- GET /api/posts/{id} - Получение поста
- POST /api/posts - Создание поста
- PUT /api/posts/{id} - Обновление поста
- DELETE /api/posts/{id} - Удаление поста

### Комментарии
- GET /api/posts/{post_id}/comments - Комментарии поста
- POST /api/posts/{post_id}/comments - Создание комментария
- PUT /api/comments/{id} - Обновление комментария
- DELETE /api/comments/{id} - Удаление комментария

### Админ
- GET /api/admin/users - Список пользователей
- GET /api/admin/posts - Список всех постов
- GET /api/admin/comments - Список всех комментариев

## Безопасность
- JWT для аутентификации
- Хеширование паролей (bcrypt)
- Middleware для проверки ролей
- CORS настройки
- Валидация входных данных

## Структура проекта

### Backend
```
backend/
├── cmd/
│   └── main.go
├── internal/
│   ├── auth/
│   ├── database/
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   └── server/
└── config/
```

### Frontend
```
frontend/
├── src/
│   ├── components/
│   ├── views/
│   ├── stores/
│   ├── router/
│   └── assets/
├── public/
└── index.html
``` 