# Football League Simulation API

A GoLang-based backend API that simulates a football league, tracks match results, updates league standings, and calculates championship predictions based on Premier League rules. 

This project was built for the **Insider Development Intern Hiring Day Task**.

---

## ☁️ Live Demo (Cloud Deployment)
This API is fully deployed and hosted on Render. You can interact with the live endpoints immediately without needing to set up the project locally.

* **Live Base URL:** `https://your-app-name.onrender.com/api/v1`

> ⚠️ **Note on Performance:** This project is hosted on a free Render tier. Render automatically spins down free instances after a period of inactivity. If the API hasn't been used recently, **the very first request may be delayed by 50 seconds or more** while the server wakes up. Subsequent requests will be fast and responsive.

---

## Prerequisites
To run this project locally, you only need to have the following installed on your machine:
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

---

## Local Setup & Execution

1. **Clone the repository:**
   ```bash
   git clone https://github.com/FatihTirek/Backend-Case.git
   cd football-league

2. **Set up environment variables:**

   Copy the example environment file to create your local `.env` configuration.

   ```bash
   cp .env.example .env
   ```

3. **Start the application:**

   Use Docker Compose to build the Go API and initialize the PostgreSQL database. The database schema will be automatically applied on the first boot.
   
   ```bash
   docker-compose up --build -d
   ```

4. **Verify the services:**

   The GoLang API will be available at http://localhost:8080/api/v1.
   
   The PostgreSQL database is exposed on port 5432 for local inspection (e.g., via DBeaver or DataGrip).

6. **Stop the application:**
    ```bash
    docker-compose down

---

## Database & SQL

As requested, the SQL schema and queries required to run this project are included:
* **Schema:** Located in `migrations/001_initial_schema.sql`. This file is automatically executed by Docker upon the initial database creation, seeding the initial teams and fixtures.
* **Queries:** All raw SQL queries used for the simulation logic (fetching standings, updating match results, etc.) are written directly in the Go codebase. You can find them clearly defined in the `repository` folder.

---

## 🧪 Testing the API (Postman)

Since a frontend is not required for this task, the application is designed to be easily understandable and fully testable via API endpoints. 

To make the evaluation process completely frictionless, I have included a comprehensive Postman collection in the repository.

1. **Locate the File:** Find the `Football League API.postman_collection` file in the root directory of this repository.
2. **Import:** Open Postman, click **Import**, and select the `.json` file.
3. **Pre-configured Environment:** The collection is designed to be used with the live cloud endpoint out of the box. You can easily switch the {{base_url}} variable to `http://localhost:8080/api/v1` if you are running it locally.
4. **Saved Examples:** I have saved **Example Responses** for every request in the collection. You can click on the examples within Postman to see the exact expected JSON inputs and outputs without even needing to spin up the Docker container first!

---

## 📖 API Documentation

**Base URL:** `https://your-app-name.onrender.com/api/v1` 
*(Or `http://localhost:8080/api/v1` for local development)*

---

### 🟢 Data Retrieval (GET)

#### 1. Get All Teams
Retrieves a list of all participating teams and their base statistics.
* **Endpoint:** `GET /teams`
* **Success Response:** `200 OK`

#### 2. Get League Standings
Retrieves the current, up-to-date league table, sorted by points and goal difference.
* **Endpoint:** `GET /standings`
* **Success Response:** `200 OK`

#### 3. Get All Matches
Retrieves the complete fixture list for the entire season, including played and unplayed matches.
* **Endpoint:** `GET /matches`
* **Success Response:** `200 OK`

#### 4. Get Matches by Week
Retrieves the fixtures and results for a specific simulated week.
* **Endpoint:** `GET /weeks/{week}/matches`
* **Path Parameter:** `week` (integer)
* **Success Response:** `200 OK`

#### 5. Get Championship Predictions
Calculates and returns the percentage chance of each team winning the league based on current standings and remaining fixtures.
* **Endpoint:** `GET /predictions`
* **Success Response:** `200 OK`

---

### 🟡 Simulation Controls (POST)

#### 6. Play Next Week
Automatically simulates all unplayed matches for the next upcoming chronological week and updates the standings.
* **Endpoint:** `POST /weeks/next/play`
* **Success Response:** `200 OK`

#### 7. Play All Remaining Matches
Simulates all remaining unplayed fixtures for the entire season in a single request.
* **Endpoint:** `POST /play-all`
* **Success Response:** `200 OK`

#### 8. Reset League
Wipes all match scores, resets the standings back to zero, and prepares the league for a fresh simulation.
* **Endpoint:** `POST /reset-league`
* **Success Response:** `200 OK`

---

### 🟠 Manual Overrides (PUT)

#### 9. Edit Match Result
Manually overrides the score of a specific match and automatically recalculates the league standings to reflect the new result.
* **Endpoint:** `PUT /matches/{id}/result`
* **Path Parameter:** `id` (integer) - The database ID of the match.
* **Request Body:**
  ```json
  {
    "home_score": 2,
    "away_score": 1
  }
* **Success Response:** `200 OK`
