# WeatherHub API

A Go API built with **Clean Architecture** (Onion Architecture), featuring real-time weather data fetching, PostgreSQL persistence, and high-performance Redis caching.

---

## 🏗️ Project Architecture

The project decouples business logic from external dependencies (frameworks, DBs, APIs), allowing for high testability and maintainability.

### Key Design Patterns
* **Proxy Pattern (Caching):** A `CachedWeatherRepo` wraps the database repository. It intercepts read calls to check **Redis** for a "hit" before falling back to **Postgres**. This keeps caching logic out of the business layer.
* **Strategy Pattern:** External weather providers are abstracted via interfaces. Swapping **OpenWeatherMap** for another provider requires zero changes to the core logic.
* **Centralized Error Handling:** A unified `RespondWithError` helper maps domain errors and `go-playground` validation errors to standardized JSON responses.



---

## 🛠️ Tech Stack
* **Language:** Go (Golang)
* **Framework:** Gin Gonic (HTTP)
* **Database:** PostgreSQL
* **Caching:** Redis
* **Documentation:** Swagger (swaggo)
* **Containerization:** Docker & Docker Compose

---

## 🚀 Getting Started

### 1. Prerequisites
* [Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/install/)
* [Make](https://www.gnu.org/software/make/)

### 2. Setup & Execution
1.  **Clone the repository:**
    ```bash
    git clone [https://github.com/xoltawn/weatherhub](https://github.com/xoltawn/weatherhub)
    cd weatherhub
    ```
2.  **Initialize Configuration:**
    ```bash
    make config
    ```
3.  **Configure API Key:**
    Open the generated `.env` file and paste your `OPEN_WEATHER_MAP_API_KEY`.

4.  **Run the Application:**
    ```bash
    make up
    ```
    The API will be available at `http://localhost:8080`.

---

## 📖 API Documentation

The API is fully documented using Swagger annotations. You can interact with the endpoints directly through the UI.

* **Swagger UI:** `http://localhost:8080/swagger/index.html`

### Primary Endpoints
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/weather` | Fetch from OpenWeather & Store in DB |
| `GET` | `/weather/:id` | Get specific record by UUID |
| `GET` | `/weather/latest/:city` | Get the most recent fetch for a city |
| `GET` | `/weather` | List all stored records |
| `PUT` | `/weather/:id` | Update an existing record |
| `DELETE` | `/weather/:id` | Remove a record and invalidate cache |
| `GET` | `/api/v1/swagger/index.html` | Swagger |


---

## 📝 Project Roadmap & Enhancements

- [x] **Multi-Unit Support:** Integrated localized measurement systems (Metric/Imperial).
- [ ] **Auth (JWT):** Secure endpoints with JSON Web Tokens and Middleware.
- [ ] **Entity Normalization (Cities):** Dedicated `cities` schema to mitigate name collisions and improve search.
- [ ] **Optimized Cache Serialization:** Transition from JSON to **Protobuf** or **MessagePack** to reduce Redis memory footprint.
- [ ] **Advanced Caching:** Implement a cleaner "Cache-Aside" or "Write-Through" strategy to optimize Redis storage.

---

## 🧹 Maintenance
* **Stop Services:** `make down`
* **Check Logs:** `docker-compose logs -f app`
* **Generate Mocks:** `go generate ./...`