# Сервис назначения ревьюеров для Pull Request’ов

Реализация [тестового задания](https://github.com/avito-tech/tech-internship/blob/main/Tech%20Internships/Backend/Backend-trainee-assignment-autumn-2025/Backend-trainee-assignment-autumn-2025.md) для стажера Backend (осенняя волна 2025)

## Требования для запуска 

- [Docker](https://www.docker.com/get-started) установлен на вашей машине.
- [Docker Compose](https://docs.docker.com/compose/install/) установлен, так как сервис запускается с его помощью.

---

## Запуск проекта

1. Клонируйте репозиторий:

```bash
git clone https://github.com/kgugunava/avito-tech-internship
cd avito-tech-internship
```

2. Создайте файл .env в корне проекта с содержимым из .env.example 

3. Запустите проект с помощью docker compose
```bash
docker compose up
```