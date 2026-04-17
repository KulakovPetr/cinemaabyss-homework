# CinemaAbyss: Monolith to Microservices

Учебный проект миграции видеосервиса с монолита на микросервисную архитектуру по паттерну Strangler Fig.

## Что внутри

- Монолит `src/monolith` с базовыми доменами (users, movies, payments, subscriptions).
- Выделенный `movies-service` в `src/microservices/movies`.
- `events-service` для событийной интеграции через Kafka.
- `proxy-service` как API Gateway и точка поэтапного переключения трафика.
- Инфраструктура для Docker Compose, Kubernetes и Helm.
- API-тесты в Postman/Newman (`tests/postman`).

## Стек

- Go
- PostgreSQL
- Kafka + ZooKeeper
- Docker / Docker Compose
- Kubernetes / Helm
- GitHub Actions

## Быстрый старт локально

```bash
docker-compose up -d
```

Сервисы:

- API Gateway: `http://localhost:8000`
- Monolith: `http://localhost:8080`
- Movies service: `http://localhost:8081`
- Events service: `http://localhost:8082`
- Kafka UI: `http://localhost:8090`

## Тестирование

```bash
cd tests/postman
npm install
npm run test:docker
```

## Архитектурная ценность проекта

- демонстрация декомпозиции монолита;
- проектирование и эксплуатация event-driven взаимодействия;
- подготовка deployment-артефактов под Kubernetes;
- автоматизация CI/CD pipeline.
