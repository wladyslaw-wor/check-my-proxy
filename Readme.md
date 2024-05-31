# Proxy Checker

## Описание (Russian)

Это сервис на Go, который каждые 5 минут проверяет работоспособность прокси-серверов и отправляет уведомления через телеграм-бота. Сервис использует базу данных Postgres для хранения информации о прокси-серверах и пользователях.

### Требования

- Go 1.16+
- PostgreSQL
- Telegram Bot API

### Установка

1. Клонируйте репозиторий:
    ```sh
    git clone https://github.com/wladyslaw-wor/check-my-proxy.git
    cd check-my-proxy
    ```

2. Установите зависимости:
    ```sh
    go mod tidy
    ```

3. Настройте глобальные переменные среды:
    ```sh
    export POSTGRESQL_HOST=localhost
    export POSTGRESQL_PORT=5432
    export POSTGRESQL_USER=your_user
    export POSTGRESQL_PASSWORD=your_password
    export POSTGRESQL_DBNAME=your_dbname
    export BotToken=your_telegram_bot_token
    ```

4. Создайте таблицы в базе данных:
    ```sql
    create table proxy
    (
        id       serial primary key,
        ip       varchar(15) not null,
        port     varchar(15) not null,
        username varchar(50) not null,
        pass     varchar(50) not null
    );

    create table users
    (
        id      serial primary key,
        chat_id bigint unique
    );
    ```

5. Добавьте свои прокси в таблицу `proxy`.

### Запуск

1. Запустите сервис:
    ```sh
    go run main.go
    ```

### Контакты

Если у вас есть вопросы или предложения, пожалуйста, свяжитесь со мной по электронной почте.

---

## Description (English)

### Project Description

This is a Go service that checks the availability of proxy servers every 5 minutes and sends notifications via a Telegram bot. The service uses a PostgreSQL database to store information about proxy servers and users.

### Requirements

- Go 1.16+
- PostgreSQL
- Telegram Bot API

### Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/wladyslaw-wor/check-my-proxy.git
    cd check-my-proxy
    ```

2. Install dependencies:
    ```sh
    go mod tidy
    ```

3. Set up environment variables:
    ```sh
    export POSTGRESQL_HOST=localhost
    export POSTGRESQL_PORT=5432
    export POSTGRESQL_USER=your_user
    export POSTGRESQL_PASSWORD=your_password
    export POSTGRESQL_DBNAME=your_dbname
    export BotToken=your_telegram_bot_token
    ```

4. Create tables in the database:
    ```sql
    create table proxy
    (
        id       serial primary key,
        ip       varchar(15) not null,
        port     varchar(15) not null,
        username varchar(50) not null,
        pass     varchar(50) not null
    );

    create table users
    (
        id      serial primary key,
        chat_id bigint unique
    );
    ```

5. Add your proxies to the `proxy` table.

### Run

1. Run the program:
    ```sh
    go run main.go
    ```

### Contact

If you have any questions or suggestions, please contact me via email.
