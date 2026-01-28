# Выполненные задания

## Задание 1: Проектирование архитектуры

- Создана контейнерная диаграмма C4 в формате PlantUML
- Файл: `docs/kinobezdna-to-be.puml`
- Диаграмма показывает разделение системы на микросервисы, API Gateway, Kafka, PostgreSQL

## Задание 2: Реализация прокси-сервиса и events-сервиса

### Proxy-сервис
- Реализован на Go в `src/microservices/proxy/`
- Паттерн Strangler Fig с feature flag `MOVIES_MIGRATION_PERCENT`
- Проксирует запросы `/api/movies` между монолитом и movies-service
- Проксирует `/api/events/*` в events-service
- Остальные запросы идут в монолит

### Events-сервис
- Реализован на Go в `src/microservices/events/`
- Использует библиотеку `github.com/IBM/sarama` для Kafka
- REST API endpoints:
  - `POST /api/events/movie` → топик `movie-events`
  - `POST /api/events/user` → топик `user-events`
  - `POST /api/events/payment` → топик `payment-events`
- Producer отправляет события в Kafka
- Consumer читает события из всех топиков и логирует обработку

## Задание 3: CI/CD и Kubernetes

### CI/CD
- Доработан `.github/workflows/docker-build-push.yml`
- Добавлена сборка и push для `proxy-service` и `events-service`
- Добавлен job `api-tests` для автоматического тестирования
- Workflow срабатывает при push в ветки `main` и `cinema`

### Kubernetes
- Созданы манифесты:
  - `src/kubernetes/proxy-service.yaml` (Deployment + Service)
  - `src/kubernetes/events-service.yaml` (Deployment + Service)
- Доработан `ingress.yaml` для маршрутизации:
  - `/` → proxy-service
  - `/api/events` → events-service
- Обновлён `configmap.yaml` с `EVENTS_SERVICE_URL`

## Задание 4: Helm-чарты

- Доработаны Helm-чарты в `src/kubernetes/helm/`
- Заполнены шаблоны:
  - `templates/services/proxy-service.yaml`
  - `templates/services/events-service.yaml`
- Обновлён `templates/configmap.yaml`
- В `values.yaml` настроены все необходимые параметры

## 📝 Документация

- Заполнен `Project_template.md` с описанием всех решений
- Создана инструкция для пользователя в `C:\ar\dz2\ИНСТРУКЦИЯ_ДЛЯ_СДАЧИ.md`

## 🔧 Что нужно сделать пользователю

1. **Создать свой репозиторий** на GitHub из шаблона
2. **Настроить пути к образам** в Kubernetes манифестах (заменить `ghcr.io/db-exp/cinemaabysstest/` на свой путь)
3. **Создать PAT токен** и настроить `dockerconfigsecret.yaml`
4. **Запустить локально** через docker-compose и сделать скриншоты
5. **Развернуть в Kubernetes** и сделать скриншоты
6. **Создать Pull Request** из ветки `cinema` в `main`

Подробные инструкции в файле `C:\ar\dz2\ИНСТРУКЦИЯ_ДЛЯ_СДАЧИ.md`
