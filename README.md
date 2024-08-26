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

#### /api/tokens

Поддерживает метод POST

- POST

Пример тела запроса: 

```json
{
    "id": "6361ce03-9ea3-4d71-9028-21c20506164e"
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
curl -d '{"id":"6361ce03-9ea3-4d71-9028-21c20506164e"}' http://localhost:9000/api/tokens
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
curl -d '{"id":"6361ce03-9ea3-4d71-9028-21c20506164e", "refresh_token":"eHU0HXVZImadRMUyVIXFsuywhGB/FuUPCt/27ckI2Ok="}' http://localhost:9000/api/tokens/refresh
curl -d '{"id":"6361ce03-9ea3-4d71-9028-21c20506164e", "refresh_token":""}' http://localhost:9000/api/tokens/refresh
curl -d '{}' http://localhost:9000/api/tokens/refresh
```
  
</details>

### Добавлены тесты

- Mock тесты эндпоинтов с критичными ситуациями
- Создание токенов

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
DB_URL= postgres://postgres:test@localhost:5432/medods?sslmode=disable
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
DB_URL= postgres://postgres:test@database:5432/medods?sslmode=disable
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
