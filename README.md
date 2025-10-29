# 🚀 Микросервисная система управления

[![Go](https://img.shields.io/badge/Go-1.19+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![React](https://img.shields.io/badge/React-18.2.0-61DAFB?style=for-the-badge&logo=react)](https://reactjs.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-13-336791?style=for-the-badge&logo=postgresql)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-🐳-2496ED?style=for-the-badge&logo=docker)](https://www.docker.com/)
[![Redis](https://img.shields.io/badge/Redis-DC382D?style=for-the-badge&logo=redis)](https://redis.io/)

Современная микросервисная система для управления задачами и пользователями с веб-интерфейсом на React и бэкендом на Go.

## ✨ Особенности

### 🏗️ Архитектура
- **Микросервисная архитектура** - независимые сервисы
- **API Gateway** - единая точка входа
- **Контейнеризация** с Docker
- **Кэширование** с Redis
- **Реляционная БД** PostgreSQL

### 👥 Управление пользователями
- Регистрация и авторизация
- Ролевая модель (Админ, Менеджер, Пользователь)
- Полный CRUD функционал
- Красивый интерфейс с эмодзи

### ✅ Управление задачами
- Создание, редактирование, удаление задач
- Назначение задач пользователям
- Статусы и приоритеты
- Фильтрация и поиск

### 🎨 Интерфейс
- Современный React UI
- Адаптивный дизайн
- Красивые карточки и анимации
- Поддержка русского языка
- Красивый фон с листьями

## 🏗️ Структура проекта
microservices-project/
- 📁 api-gateway/ # API Gateway (Go + Gin)
- 📁 user-service/ # Сервис пользователей (Go + Redis)
- 📁 task-service/ # Сервис задач (Go)
- 📁 notification-service/ # Сервис уведомлений (Go)
- 📁 analytics-service/ # Сервис аналитики (Go)
- 📁 web-app/ # Веб-интерфейс (React)
- 📄 docker-compose.yml # Оркестрация контейнеров
- 📄 init.sql # Инициализация БД
- 📄 README.md # Документация

## 🛠️ Технологический стек

### Backend
- **Go 1.19+** - высокопроизводительный язык
- **Gin** - веб-фреймворк
- **GORM** - ORM для работы с БД
- **Redis** - кэширование и сессии

### Frontend
- **React 18** - пользовательский интерфейс
- **Axios** - HTTP клиент
- **CSS3** - стили и анимации

### Базы данных
- **PostgreSQL 13** - основная БД
- **Redis 7** - кэширование

### Инфраструктура
- **Docker** - контейнеризация
- **Docker Compose** - оркестрация

# 🚀 Быстрый старт

### Предварительные требования
- Docker
- Docker Compose

### Запуск проекта

### 1. Клонируйте репозиторий
git clone https://github.com/Tim948/microservices-project.git
cd microservices-project

###  2.Запустите все сервисы
docker-compose up -d --build

### 3. Откройте приложение
- **Веб-приложение:** http://localhost:3000
- **API Gateway:**    http://localhost:8080
- **PostgreSQL:**     localhost:5432
- **Redis:**          localhost:6379

###  Тестовые пользователи
- **Администратор:** admin / password123
- **Менеджер: manager** / password123
- **Пользователь:** user1 / password123

# 📡 API Endpoints
- **GET    /health**          # Статус сервиса
- **GET    /users**           # Список пользователей
- **POST   /users**           # Создать пользователя
- **PUT    /users/:id**       # Обновить пользователя
- **DELETE /users/:id**       # Удалить пользователя

### Task Service (:8082)
- **GET    /health**          # Статус сервиса  
- **GET    /tasks**           # Список задач
- **POST   /tasks**           # Создать задачу
- **PUT    /tasks/:id**       # Обновить задачу
- **DELETE /tasks/:id**       # Удалить задачу

###  API Gateway (:8080)
- **GET    /health**          # Статус всех сервисов
- **GET    /users/**        # Прокси к User Service
- **GET    /tasks/**         # Прокси к Task Service

# 🗃️ База данных

### Основные таблицы
- users - пользователи системы
- tasks - задачи и проекты
- projects - проекты
- notifications - уведомления
- user_activities - активность пользователей
### Инициализация
База данных автоматически инициализируется при первом запуске с тестовыми данными.

# 🔧 Разработка

### Запуск для разработки
- docker-compose up -d postgres redis
- docker-compose up -d --build user-service task-service
- docker-compose up -d --build api-gateway
- docker-compose up -d --build web-app

### Просмотр логов
### Все логи
docker-compose logs -f

### Конкретный сервис
docker-compose logs -f user-service
docker-compose logs -f web-app

### Остановка всех сервисов
docker-compose down

### Остановка с удалением volumes
docker-compose down -v

# 🎯 Функциональность
Для администратора
✅ Полное управление пользователями

✅ Создание и назначение задач

✅ Просмотр всей статистики

✅ Управление ролями

Для менеджера
✅ Управление задачами

✅ Назначение задач пользователям

✅ Просмотр прогресса проектов

Для пользователя
✅ Просмотр своих задач

✅ Изменение статуса задач

✅ Просмотр профиля

### 📊 Мониторинг

Health checks
curl http://localhost:8080/health

Проверка БД
docker exec -it microservices-project-postgres-1 psql -U micro_user -d microservices -c "SELECT * FROM users;"

Проверка Redis
docker exec -it microservices-project-redis-1 redis-cli ping

🤝 Вклад в проект

1.Форкните репозиторий
2.Создайте ветку для функции (git checkout -b feature/amazing-feature)
3.Закоммитьте изменения (git commit -m 'Add amazing feature')
4.Запушьте в ветку (git push origin feature/amazing-feature)
5.Откройте Pull Request

👨‍💻 Автор
Ваше Имя

GitHub: @Tim948

Проект: Microservices Management System

🙏 Благодарности

Go community за отличный язык
React team за прекрасный фреймворк
Docker за удобные контейнеры
Freepik за красивый фон

⭐ Не забудьте поставить звезду, если проект вам понравился!