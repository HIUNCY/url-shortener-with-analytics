# URL Shortener & Analytics Service

![Build Status](https://img.shields.io/badge/build-passing-brightgreen)
![Go Version](https://img.shields.io/badge/go-1.18%2B-blue)
![License](https://img.shields.io/badge/license-MIT-green)

A robust and feature-rich URL Shortener service built with Go, complete with click tracking and in-depth analytics. Designed with a **Clean Architecture** for scalability and maintainability.

---

## Table of Contents

- [About The Project](#about-the-project)
- [Key Features](#key-features)
- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Getting Started](#getting-started)
- [Usage](#usage)
- [API Testing with Postman](#api-testing-with-postman)
- [API Documentation](#api-documentation)
- [Contact](#contact)

---

## About The Project

This project provides a URL shortening service that is not only functional but also delivers valuable insights through analytics. Users can shorten URLs, use custom aliases, and track the performance of each link in detail, from click counts to the geographical data of visitors.

The project is designed following **Clean Architecture** principles, clearly separating business logic (domain), services, and the data access layer (repository). This approach ensures the code remains clean, testable, and ready for future development.

---

## Key Features

-   👤 **User Management**: Registration, Login (JWT), Profile Management, and API Key authentication.
-   🔗 **URL Management**: Create, view, update, and delete short URLs with customization options (alias, title, password, expiration date).
-   ➡️ **Fast Redirection**: An efficient redirection process with asynchronous click tracking.
-   📊 **In-Depth Analytics**: Track total clicks, referrers, geography (country, city), devices, browsers, and OS for each URL.
-   🔳 **QR Code Generation**: Generate and download QR codes for every short URL.
-   📚 **API Documentation**: Interactive API documentation automatically generated using Swagger.

---

## Architecture

This project adopts **Clean Architecture** principles to ensure a clear separation of concerns and a directed flow of dependencies.
-   **Handlers**: Responsible for receiving HTTP requests, validating input (DTOs), and calling the appropriate service.
-   **Services**: Contain all core business logic. This layer has no knowledge of HTTP or database details.
-   **Repository**: Acts as an abstraction for the data layer. It defines interfaces for CRUD operations and complex queries.
-   **Domain**: Represents the core business entities (`User`, `URL`, `Click`) and repository contracts.

---

## Tech Stack

-   **Language**: [Go](https://golang.org/)
-   **Web Framework**: [Gin Gonic](https://gin-gonic.com/)
-   **Database**: [PostgreSQL](https://www.postgresql.org/)
-   **ORM**: [GORM](https://gorm.io/)
-   **Configuration**: [Viper](https://github.com/spf13/viper)
-   **API Documentation**: [Swag (Swagger)](https://github.com/swaggo/swag)
-   **GeoIP**: [MaxMind GeoLite2](https://www.maxmind.com)

---

## Getting Started

To get a local copy up and running, follow these simple steps.

### Prerequisites

-   Go version 1.18 or higher.
-   A running PostgreSQL server.
-   A [MaxMind](https://www.maxmind.com/en/geolite2/signup) account to download the GeoIP database.
-   [Postman](https://www.postman.com/downloads/) (Optional, for API testing).

### Installation

1.  **Clone the repo:**
    ```sh
    git clone https://github.com/HIUNCY/url-shortener-with-analytics
    cd url-shortener-with-analytics
    ```

2.  **Install Go dependencies:**
    ```sh
    go mod tidy
    ```

3.  **Set up the Database:**
    -   Create a new database in PostgreSQL.
    -   Run the SQL script provided in `url_shortener.sql` to create all necessary tables and indexes.

4.  **Download the GeoIP Database:**
    -   Download the `GeoLite2-City.mmdb` file from your MaxMind account.
    -   Create a `geoip` directory in the project root and place the `.mmdb` file inside it.

5.  **Set up Environment Variables:**
    -   Copy the `.env.example` file to a new file named `.env`.
    -   Fill in all the variables in the `.env` file with your database credentials, JWT secret, and the GeoIP database path.

---

## Usage

After the installation is complete, you can run the server.

1.  **Run the Server:**
    ```sh
    go run cmd/main.go
    ```
    The server will be running at `http://localhost:8080` (or your configured port).

---

## API Testing with Postman

This project includes a Postman collection and environment file for easy API testing.

1.  **Open Postman** and go to `File > Import`.
2.  **Import the Collection:** Select the `URL Shortener with Analytics.postman_collection.json` file from the project root.
3.  **Import the Environment:** Import the `uswa.postman_environment.json` file.
4.  **Select the Environment:** In the top-right corner of Postman, select the "uswa" environment from the dropdown menu.
5.  You can now use the pre-configured requests in the collection to test the API. The collection is set up to automatically handle authentication tokens after you log in.

---

## API Documentation

For a more detailed and interactive API specification, you can use the built-in Swagger documentation.

-   **Access the Docs:** With the server running, open your browser and navigate to:
    `http://localhost:8080/swagger/index.html`

-   **Main Endpoints**:
    -   `/auth`: Register, Login, Refresh Token, Logout.
    -   `/profile`: User profile management.
    -   `/urls`: CRUD operations for short URLs.
    -   `/analytics`: Dashboard and detailed URL analytics.
    -   `/{short_code}`: The public endpoint for redirection.

---

## Contact

Muhamad Zainul Kamal - zainulkamal393@gmail.com

Project Link: [https://github.com/HIUNCY/url-shortener-with-analytics](https://github.com/HIUNCY/url-shortener-with-analytics)
