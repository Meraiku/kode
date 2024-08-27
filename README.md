# Тестовое задание KODE

## Реализовано

### Создание юзера по почте

#### /api/users

Поддерживает методы GET POST

- GET 

Отправляет список пользователей в формате JSON

<details>
  <summary>Комнды для тестирования</summary>

```bash
curl http://localhost:9000/api/users
```

</details>

- POST 

Позволяет создать нового пользователя

Пример тела запроса:

```json
{
    "email": "test@example.com"
}
```

- Реализовано:
  - проверка email по regex. 
  - проверка уникальности почт пользователей

Ответ в формате JSON с присвоенным UUID

<details>
  <summary>Комнды для тестирования</summary>

```bash
curl -d '{"email":"test@gmail.com"}' http://localhost:9000/api/users
curl -d '{"email":"1"}' http://localhost:9000/api/users
curl -d '{}' http://localhost:9000/api/users
```

</details>

### Создание пары access, refresh токенов

При аутентификации токена отправляются в теле ответа, но также записываются в cookie

Авторизация происходит через cookie

#### /api/tokens

Поддерживает метод POST

- POST

Пример тела запроса: 

```json
{
    "id": "6361ce03-9ea3-4d71-9028-21c20506164e"
}
```

Тело ответа:

```json
{
  "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiIxMjcuMC4wLjEiLCJzdWIiOiI0ZjFjYTI2NC0zOGIyLTQ0YTEtODFkMC03MWIxZmM0YTQzYjEiLCJleHAiOjE3MjQ3NTg5NTIsImlhdCI6MTcyNDc1ODA1Mn0.RG598gl9HMNmjlllwTHKDgO7tvMZnV4UH0bB-1y2HQi39Nk4F99ynWC6jHFLuL9mNptoYGs0S9TtMipk7XUDLw",
  "refresh_token": "aAR+JMDhCa5FJJn+ts0uEP3CUt0605VK3/WchZ7FVtI="
}
```

- Реализовано:
  - проверка нахождения пользователя в базе данных
  - выдача пары access и refresh токенов
  - хранение токенов в кэше
  - запись refresh токена в базу данных в виде bcrypt хэша

Формат access токена JWT, алгоритм SHA512.
<details>
  <summary>payload токена</summary>
  
  ```
{
  "iss": "127.0.0.1",                               - ip адрес вызова
  "sub": "6361ce03-9ea3-4d71-9028-21c20506164e",    - uuid пользователя
  "exp": 1724676917,                                - время недействительности токена(15 минут после создания)
  "iat": 1724676017                                 - время выдачи токена
}
```
  
</details>

Формат refresh токена base64.

Ответ в формате JSON с парой access refresh токенов.

<details>
  <summary>Комнды для тестирования</summary>
  
```bash
curl -d '{"id":"${UUID}"}' http://localhost:9000/api/tokens
curl -d '{"id":""}' http://localhost:9000/api/tokens
curl -d '{}' http://localhost:9000/api/tokens
```
  
</details>


### Refresh операция на токены

#### /api/tokens/refresh

Поддерживает метод POST

- POST

Пример тела запроса: 

```json
{
    "id": "6361ce03-9ea3-4d71-9028-21c20506164e",
    "refresh_token": "eHU0HXVZImadRMUyVIXFsuywhGB/FuUPCt/27ckI2Ok="
}
```

Тело ответа:

```json
{
  "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiIxMjcuMC4wLjEiLCJzdWIiOiI0ZjFjYTI2NC0zOGIyLTQ0YTEtODFkMC03MWIxZmM0YTQzYjEiLCJleHAiOjE3MjQ3NTg5NTIsImlhdCI6MTcyNDc1ODA1Mn0.RG598gl9HMNmjlllwTHKDgO7tvMZnV4UH0bB-1y2HQi39Nk4F99ynWC6jHFLuL9mNptoYGs0S9TtMipk7XUDLw",
  "refresh_token": "aAR+JMDhCa5FJJn+ts0uEP3CUt0605VK3/WchZ7FVtI="
}
```

- Реализовано:
  - проверка ip адреса вызова
  - отправка предупреждения через SMTP протокол
    на почту пользователя при вызове с другого ip адреса

- Операция защищена от:
  - ip адрес не совпадает
  - повторного использования токена
  - использования токена другого пользователя

    
<details>
  <summary>Комнды для тестирования</summary>
  
```bash
curl -d '{"id":"${UUID}", "refresh_token":"${Refresh token}"}' http://localhost:9000/api/tokens/refresh
curl -d '{"id":"${UUID}", "refresh_token":""}' http://localhost:9000/api/tokens/refresh
curl -d '{}' http://localhost:9000/api/tokens/refresh
```
  
</details>

### Создание и просмотр заметок

- Реализовано:
  - процесс авторизации через cookie
  - создание новых пар токенов при истечении жизни токенов

#### /api/notes

Поддерживает методы GET POST

- Get

Выводит список заметок пользователя

Тело ответа:

```json
[
  {
    "id": "b0e90346-1717-4b4d-8f57-44cdcb521c0e",
    "title": "Это Мой тайтл",
    "body": "Это моё тело",
    "created_at": "2024-08-27T11:18:53.33262Z",
    "updated_at": "2024-08-27T11:18:53.33262Z",
    "user_id": "4f1ca264-38b2-44a1-81d0-71b1fc4a43b1"
  },
  {
    "id": "7be26462-592d-4333-b7d0-b4d11241d6cf",
    "title": "Это Мой тайтл",
    "body": "Это моё тело",
    "created_at": "2024-08-27T11:18:52.45627Z",
    "updated_at": "2024-08-27T11:18:52.45627Z",
    "user_id": "4f1ca264-38b2-44a1-81d0-71b1fc4a43b1"
  }
]
```

    
<details>
  <summary>Комнды для тестирования</summary>
  
```bash
curl -v --cookie "access=${JWT token}" http://localhost:9000/api/notes
```
  
</details>


- POST

Пример тела запроса: 

```json
{
    "title": "ЭьТо Мой тайтль",
    "body": "Этттто моё телаа"
}
```
Тело ответа:

```json
{
  "title": "Это Мой тайтл",
  "body": "Это моё тело"
}
```

- Реализовано:
  - проверка орфографических ошибок при помощи Яндекс.Спеллер

    
<details>
  <summary>Комнды для тестирования</summary>
  
```bash
curl -v --cookie "access=${JWT token}" -d '{"title":"ЭьТо Мой тайтль", "body":"Этттто моё телаа"}' http://localhost:9000/api/notes
```
  
</details>

### Добавлены тесты

- Mock тесты создания токенов и пользователей с критичными ситуациями
- Тест корректности создания токенов
- Тест корректной функциональности Яндкс.Спеллера

```bash
go test ./...
```

## Настройка сервера

Реализован запуск сервера через контейнер

Переменные для запуска:

```
DB_URL=                                 -Адрес для подключения к PostgreSQL

RDB_URL=                                -Адрес для подключения к Redis

JWT_SECRET=                             -Ключ подписи JWT


SMTP_SERVER=                            -Адрес SMTP сервера(без порта). Если переменная не указана, 
                                         функционал отправки писем не будет работать,
                                         сервер сохранит свою роботоспособность

SMTP_NAME=                              -Email адрес отправителя
SMTP_PASS=                              -Пароль аутентификации SMTP сервера

```



### Запуск контейнера

Создан docker-compose файл с готовой средой для тестирования. (image выгружен в хаб)

Необходим файл .env и docker.env(для корректного подключения к базам данных)
<details>
  <summary>.env</summary>
  
```
DB_URL= postgres://postgres:test@localhost:5432/kode?sslmode=disable
RDB_URL= redis://:pass@localhost:6379/0
JWT_SECRET= RHTjGzsHH+J8uQvfgNi1N48cn8ZL6NQJXRgJZVlNWj8FlsyPkOMXgCuPdu3nx3aoMmc8VXay7iJnk4/e2mAIXA==

SMTP_SERVER="smtp.yandex.ru"

SMTP_NAME=
SMTP_PASS=
```
  
</details>
<details>
  <summary>docker.env</summary>
  
```
DB_URL= postgres://postgres:test@database:5432/kode?sslmode=disable
RDB_URL= redis://:pass@cache:6379/0
```
  
</details>

### Локальная разработка

Необходим сервер PostgreSQL, Redis, .env файл для корректной работы и запуска

Команда для запуска
```bash
make run
```

### Создание базы данных в PostgreSQL

```bash
docker exec -t database createdb -U postgres kode
```

### SQL миграции

Выполнены через goose

- Установка
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

- Использование (необходим .env файл с переменной "DB_URL")
```bash
make up
make down
```
