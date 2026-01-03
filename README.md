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
    git clone https://github.com/xoltawn/weatherhub
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
# copy config file
make config

# run in container
make up

# shutdown container
make down
```

### Project Roadmap & Enhancements
[x] Multi-Unit Support: Integrated localized measurement systems (Metric/Imperial) supporting internationalized weather data standards.

[ ] Entity Normalization (Cities): Implement a dedicated cities schema and search API. This mitigates namespace collisions (ambiguous city names) and improves data integrity by using unique identifiers (e.g., OpenWeather City IDs or Geo-coordinates).

[ ] Optimized Cache Serialization: Transition from standard JSON to a binary serialization format (e.g., Protobuf or MessagePack) to reduce Redis memory footprint and improve I/O throughput.

- [ ] Reading data from cache needs cleaner approach to save Redis storage 