# WeatherHub API

A Go API built with **Clean Architecture**, featuring real-time weather data fetching, PostgreSQL persistence, and high-performance Redis caching.

---

## 🏗️ Project Architecture

This project is built using the **Onion Architecture** (Clean Architecture) pattern. This ensures that the business logic remains independent of frameworks, databases, and external APIs.



### Key Design Patterns
* **Proxy Pattern (Caching):** A `CachedWeatherRepo` wraps the database repository. It intercepts read calls to check **Redis** for a "hit" before falling back to **Postgres**. This keeps caching logic out of the business layer.
* **Strategy Pattern:** External weather providers are abstracted via interfaces. Swapping **OpenWeatherMap** for another provider requires zero changes to the core logic.
* **Centralized Error Handling:** A unified `RespondWithError` helper maps domain errors and `go-playground` validation errors to standardized JSON responses.

---

## 🚀 Getting Started

### 1. Prerequisites
* [Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/install/)
* [Make](https://www.gnu.org/software/make/)

### 2. Environment Setup
1.  **Clone the repository:**
    ```bash
    git clone [https://github.com/xoltawn/weatherhub.git](https://github.com/xoltawn/weatherhub.git)
    cd weatherhub
    ```
2.  **Create your `.env` file:**
    ```bash
    cp .env.example .env
    ```
3.  **Configure API Key:**
    Open `.env` and paste your `OPEN_WEATHER_MAP_API_KEY`.

### 3. Running the App
The `Makefile` automates the Docker lifecycle:
```bash
# copt config file
make config
make up

#shutdown
make down
```

### Todos
- [x] Add basic unit system (metric. imperial )
- [ ] A table and api for `cities` to provider search city functionality and mitigate same city name problem
- [ ] Reading data from cache needs cleaner approach to save Redis storage