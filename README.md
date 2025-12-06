# Gofra Market — уязвимый сервис для Attack-Defense CTF

[![Go Version](https://img.shields.io/badge/Go-1.25-blue?logo=go)](https://go.dev/dl/)
[![MongoDB](https://img.shields.io/badge/MongoDB-7.x-47A248?logo=mongodb&logoColor=white)](https://www.mongodb.com/)
[![Docker Compose](https://img.shields.io/badge/Docker%20Compose-ready-2496ED?logo=docker&logoColor=white)](https://docs.docker.com/compose/)
[![CTF Ready](https://img.shields.io/badge/Attack--Defense-ready-critical)](#)

> Тренировочный маркетплейс для проведения соревнований в формате Attack-Defense.

## Disclaimer

⚠️ Репозиторий предназначен **только** для учебных задач в контролируемых средах. Перед использованием в продакшене устраните уязвимости и внедрите защитные механизмы.

## Содержание
- [Особенности](#особенности)
- [Компоненты](#компоненты)
- [Структура репозитория](#структура-репозитория)
- [Уязвимости из коробки](#уязвимости-из-коробки)
- [Развёртывание](#развертывание)
  - [Быстрый старт (Docker Compose)](#быстрый-старт-docker-compose)
  - [Локальный запуск backend и frontend](#локальный-запуск-backend-и-frontend)
- [Документация и отладка](#документация-и-отладка)
- [Эксплойты](#эксплойты)
- [Чекер](#чекер)
- [Рекомендации по защите](#рекомендации-по-защите)

## Особенности
- 💣 Встроенные уязвимости: NoSQL-инъекция, SSRF и целочисленный underflow.
- 🛠️ Готовые инструменты: эксплойты и чекер для интеграции в Attack-Defense.
- 📚 Документация разработчика: Swagger UI и JSON-обзор пакетов в debug-режиме.

## Компоненты
- **Backend** — Go 1.25, Gin, MongoDB driver.
- **Frontend** — Quasar/Vue, собирается в nginx-образ.
- **База данных** — MongoDB с миграциями.
- **Документация** — модуль `internal/docs` (Swagger + debug-хэндлеры).

## Структура репозитория
- `gofra_service/backend`
  - `cmd/api` — точка входа сервиса.
  - `internal/app` — инициализация сервера и роутов.
  - `internal/docs` — Swagger и debug-эндпоинты.
  - `internal/db` — миграции, подключение.
  - `internal/domain` — модели данных.
  - `internal/repo`, `internal/service` — репозитории и бизнес-логика.
  - `internal/transport/http` — хэндлеры и middleware.
- `gofra_service/frontend` — SPA-приложение (Quasar).
- `exploits/` — PoС скрипты.
- `gofra_checker/` — чекер для платформы ctf01d.

## Уязвимости из коробки
| ID | Тип | Точка входа | Ключевые файлы | Что даёт |
|----|-----|-------------|----------------|----------|
| V-01 | NoSQL Injection | `GET /api/market?filter` | `internal/repo/repo_listing.go` | Отсутствиес санитизации ввода фильтра позволяет выполнять произвольные команды в консоли MongoDB. |
| V-02 | SSRF | `POST /api/listings/{id}/image_from_url` | `internal/service/service_image.go` | Неправильная настройка CORS и остутствие проверки Content-Type позволяет выполнять запросы в том числе к MongoDB. |
| V-03 | Integer Underflow | `POST /api/listings/{id}/bump` | `internal/service/service_listing.go` | При приведении типов из отрицательного int в uint можно получить огромный баланс. |

## Развёртывание

### Быстрый старт (Docker Compose)
```bash
cd gofra_service
docker compose up --build -d
```
После запуска компоненты доступны по адресу:
- Backend: http://localhost:8080
- Frontend: http://localhost:8081
- MongoDB: `mongodb://admin:S3cr3tP4ssw0rd@localhost:27017`

Backend автоматически применяет миграции и создаёт пользователя `system_seller` (пароль `system123`).

### Локальный запуск backend и frontend
**MongoDB** (если не используете docker-compose):
```bash
docker run -d --name gofra_mongo -p 27017:27017 \
  -e MONGO_INITDB_ROOT_USERNAME=admin \
  -e MONGO_INITDB_ROOT_PASSWORD=S3cr3tP4ssw0rd mongo:latest
```

**Backend**
```bash
cd gofra_service/backend
export MONGO_URL="mongodb://admin:S3cr3tP4ssw0rd@localhost:27017"
export DB_NAME="gofra"
export SERVER_PORT=8080
export GIN_MODE=debug  # включает Swagger и debug-эндпоинты
go run ./cmd/api
```

**Frontend**
```bash
cd gofra_service/frontend
npm install
npm run dev
```

## Документация и отладка
- Swagger UI (доступно при `GIN_MODE=debug`): `http://localhost:8080/swagger/index.html`
- Swagger артефакты: `backend/internal/docs/swagger/swagger.{json,yaml}`
- Перегенерация Swagger после правок API:
  ```bash
  cd gofra_service/backend
  swag init -g cmd/api/main.go -o internal/docs/swagger
  ```

## Эксплойты
- `exploit_nosql.py` — NoSQL-инъекция, вытягивает флаг из описаний листингов.
- `exploit_ssrf.py` — SSRF с последующей NoSQL-инъекцией по внутреннему сервису.
- `exploit_underflow.py` — форсирует переполнение баланса и возвращает креды пользователя.

Скрипты принимают базовый URL (`http://localhost:8080` по умолчанию) как аргумент.

## Чекер
`gofra_checker/gofra_checker.py` реализует функционал, совместимый с ctf01d:
1. Генерация учётных данных на основе флага.
2. Регистрация и логин пользователя.
3. Создание листинга с флагом в описании.
4. Проверка ключевых эндпоинтов (`/me`, `/my-listings`, `/market`, `/listings/{id}`, `/my-gofers`).

