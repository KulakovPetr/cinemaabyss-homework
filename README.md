# CinemaAbyss: Monolith to Microservices

Учебный проект миграции видеосервиса с монолита на микросервисную архитектуру по паттерну Strangler Fig.

## Текущий статус

- Реализовано и поддерживается в рабочем контуре:
  - монолит `src/monolith` (users, movies, payments, subscriptions);
  - выделенный `movies-service` (`src/microservices/movies`);
  - инфраструктура запуска через Docker Compose;
  - API-тесты в Postman/Newman для доступных сервисов.
- Архитектурные артефакты для следующих этапов (`events/proxy`, Kubernetes/Helm, Kafka) сохранены в репозитории как roadmap миграции.

## Стек

- Go
- PostgreSQL
- Docker / Docker Compose
- GitHub Actions
- Postman + Newman

## Быстрый старт локально

```bash
docker-compose up -d
```

Сервисы:

- Monolith: `http://localhost:8080`
- Movies service: `http://localhost:8081`

## Тестирование

```bash
cd tests/postman
npm install
npm run test:docker -- --folder "Monolith Service"
npm run test:docker -- --folder "Movies Microservice"
```

## Архитектурная ценность проекта

- декомпозиция монолита на отдельный предметный сервис;
- формирование migration roadmap по Strangler Fig;
- проектирование API-контрактов и тестового контура;
- автоматизация CI/CD-проверок для поддерживаемой части системы.
