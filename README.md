# Sortlynk URL Shortener

Sortlynk is a powerful and efficient URL shortener built with Go, Gin, and GORM. It provides a robust set of features for shortening URLs, managing user accounts, and tracking usage statistics. The project is designed for high performance, using Redis for caching and a PostgreSQL database for persistent storage.

## Features

*   **URL Shortening:** Quickly and easily shorten long URLs into a more manageable format.
*   **User Authentication:** Secure user registration and login system using JWT for authentication.
*   **Rate Limiting:** Protects the service from abuse with configurable rate limits for both authenticated and unauthenticated users.
*   **Usage Statistics:** Track the number of clicks for each shortened URL.
*   **Caching:** Utilizes Redis to cache frequently accessed URLs, ensuring fast redirect times.
*   **Scalable Architecture:** The project is structured with a clear separation of concerns, making it easy to extend and maintain.

## API Endpoints

The following is a list of all available API endpoints:

| Method  | Path                  | Description                                      |
| :------ | :-------------------- | :----------------------------------------------- |
| `POST`  | `/api/v1/auth/register` | Register a new user.                             |
| `POST`  | `/api/v1/auth/login`    | Log in an existing user and receive a JWT token. |
| `POST`  | `/api/v1/urls/shorten`  | Shorten a new URL.                               |
| `GET`   | `/api/v1/urls/my`       | Get a list of all URLs created by the user.      |
| `GET`   | `/api/v1/urls/:code/stats` | Get usage statistics for a specific short URL.   |
| `GET`   | `/:code`              | Redirect to the original URL.                    |
| `GET`   | `/health`             | Health check endpoint.                           |

## Configuration

The application is configured using environment variables. The following is a list of all available options:

| Variable               | Description                                       | Default Value                               |
| :--------------------- | :------------------------------------------------ | :------------------------------------------ |
| `DB_HOST`              | The hostname of the PostgreSQL database.          | `localhost`                                 |
| `DB_PORT`              | The port of the PostgreSQL database.              | `5432`                                      |
| `DB_USER`              | The username for the PostgreSQL database.         | `postgres`                                  |
| `DB_PASSWORD`          | The password for the PostgreSQL database.         | `password`                                  |
| `DB_NAME`              | The name of the PostgreSQL database.              | `urlshortener`                              |
| `DB_SSL`               | The SSL mode for the PostgreSQL database.         | `disable`                                   |
| `REDIS_HOST`           | The hostname of the Redis server.                 | `localhost`                                 |
| `REDIS_PORT`           | The port of the Redis server.                     | `6379`                                      |
| `REDIS_PASSWORD`       | The password for the Redis server.                |                                             |
| `REDIS_DB`             | The Redis database to use.                        | `0`                                         |
| `JWT_SECRET`           | The secret key for signing JWT tokens.            | `your-secret-key-here-change-in-production` |
| `SERVER_PORT`          | The port for the application server.              | `8080`                                      |
| `UNAUTHENTICATED_LIMIT`| The number of requests allowed per minute for unauthenticated users. | `2`                                         |
| `AUTHENTICATED_LIMIT`  | The number of requests allowed per minute for authenticated users. | `1000`                                      |

## Setup and Running

To set up and run the project locally, follow these steps:

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/your-username/sortlynk.git
    cd sortlynk
    ```

2.  **Install dependencies:**

    ```bash
    go mod tidy
    ```

3.  **Set up the environment:**

    Create a `.env` file in the root of the project and add the following environment variables:

    ```
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=postgres
    DB_PASSWORD=password
    DB_NAME=urlshortener
    DB_SSL=disable

    REDIS_HOST=localhost
    REDIS_PORT=6379
    REDIS_PASSWORD=
    REDIS_DB=0

    JWT_SECRET=your-secret-key-here-change-in-production

    SERVER_PORT=8080

    UNAUTHENTICATED_LIMIT=2
    AUTHENTICATED_LIMIT=1000
    ```

4.  **Run the application:**

    ```bash
    go run main.go
    ```

The application will now be running on `http://localhost:8080`.
