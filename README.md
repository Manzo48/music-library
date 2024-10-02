# Музыкальная библиотека

## Оглавление
1. [Описание задачи](#Описание-задачи)
2. [Методы API](#Методы-API)
3. [Архитектура](#Архитектура)
4. [База данных](#База-данных)
5. [Интеграция с Genius API](#Интеграция-с-Genius-API)
6. [Запуск](#Запуск)
8. [Документация](#Документация)
9. [Примеры запросов](#Примеры-запросов)


## Описание задачи

Cервис для хранения и подачи песен

## Методы API

### Получение списка песен

- **Пагинация**: на одной странице должно присутствовать 10 песен.
- **Фильтрация**: по группе, альбому, и названию песни.
- **Поля в ответе**: название песни, группа, альбом, текст песни, дата выпуска.

### Получение текста песни

- **Параметры**: ID песни, номер куплета.
- **Поля в ответе**: текст запрашиваемого куплета песни.

### Создание новой песни

- **Параметры**: название песни, группа, альбом, текст песни, дата выпуска.
- **Возвращаемое значение**: ID созданной песни или ошибка.

### Удаление песни

- **Параметры**: ID песни.
- **Возвращаемое значение**: успешное удаление или ошибка.

### Обновление информации о песне

- **Параметры**: ID песни, новые данные (группа, альбом, текст песни, дата выпуска).
- **Возвращаемое значение**: успешное обновление или ошибка.

## Архитектура

Проект построен на основе **чистой архитектуры** с разделением на слои:

- **Handler**: контроллеры для обработки запросов.
- **Service**: бизнес-логика приложения.
- **Repository**: взаимодействие с базой данных.
- **Model**: структуры данных.

## База данных

Используется **PostgreSQL**. Данные хранятся в двух таблицах:

- **songs**: основная информация о песнях (ID, название, группа, альбом, дата выпуска).
- **song_details**: детали о песнях (ID, текст песни, ссылка, дата создания, артист, альбом, приглашенные артисты).

### Пример структуры таблицы `songs`

| Поле       | Тип        | Описание               |
|------------|------------|------------------------|
| id         | SERIAL     | Уникальный идентификатор|
| name       | VARCHAR    | Название песни         |
| group_name | VARCHAR    | Группа                 |
| album      | VARCHAR    | Альбом                 |
| release_date | DATE     | Дата выпуска           |

### Пример структуры таблицы `song_details`

| Поле       | Тип        | Описание               |
|------------|------------|------------------------|
| id         | SERIAL     | Уникальный идентификатор|
| song_id    | INTEGER    | Ссылка на песню         |
| text       | TEXT       | Текст песни            |
| artist     | VARCHAR    | Артист                 |
| album      | VARCHAR    | Альбом                 |
| featured_artists | VARCHAR | Приглашенные артисты|

## Интеграция с Genius API

Приложение интегрировано с внешним API **Genius** для получения дополнительной информации о песнях:

- **Поля**: текст песни, ссылка на текст, артисты, альбом, дата выпуска.
- Используется API-клиент для обращения к Genius и обогащения данных.


## Запуск

### Использование Docker

1. Убедитесь, что **Docker** установлен и запущен.
2. В корне проекта находится файл `docker-compose.yml`, который настраивает контейнеры для приложения и базы данных.
3. Чтобы запустить проект, выполните команду:

```bash
docker-compose up --build
```
### Если проект запускается первый раз нужно сделать миграцию

```bash
migrate -path ./migrations -database 'postgres://postgres:1234@host.docker.internal:5436/music-library?sslmode=disable' up
```

### Если запускается без докера 

```bash
migrate -path ./migrations -database 'postgres://postgres:1234@localhost:5436/music-library?sslmode=disable' up
```

### Использование локальной машины

1. Убедитесь, что база данных PostgreSQL запущена на вашей локальной машине.
2. Настройте файл `configs/local_config.yml` с параметрами подключения к базе данных.
3. Запустите приложение командой:

```bash
go run cmd/main.go ./configs/local_config.yml
```

## Примеры запросов

### Пример запроса для получения списка песен

```bash
GET /api/v1/songs?page=1&group=The%20Beatles
```

### Пример запроса для получения текста песни

```bash
GET /api/v1/songs/{song_id}/lyrics?verse=1
```

### Пример создания новой песни

```bash
POST /api/v1/songs
Content-Type: application/json

{
  "name": "Yesterday",
  "group": "The Beatles",
}
```

### Пример обновления песни

```bash
PUT /api/v1/songs/{song_id}
Content-Type: application/json

{
  "name": "Hey Jude",
  "group": "The Beatles",
  "album": "1",
  "text": "Hey Jude, don't make it bad...",
  "release_date": "1968-08-26"
}
```

### Пример удаления песни

```bash
DELETE /api/v1/songs/{song_id}
```
