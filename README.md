# Football League Simulation API

A GoLang-based backend API that simulates a football league, tracks match results, updates league standings, and calculates championship predictions based on Premier League rules. 

This project was built for the Insider Development Intern Hiring Day Task.

## Features
- Complete league simulation engine using GoLang.
- Interface-based design and struct composition.
- PostgreSQL database integration for persistent match data.
- Fully containerized using Docker and Docker Compose.

## Prerequisites
To run this project locally, you only need to have the following installed on your machine:
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Local Setup & Execution

1. **Clone the repository:**
   ```bash
   git clone https://github.com/FatihTirek/football-league.git
   cd football-league

2. **Set up environment variables:**

   Copy the example environment file to create your local `.env` configuration.

   ```bash
   cp .env.example .env
   ```

3. **Set up environment variables:**

   Use Docker Compose to build the Go API and initialize the PostgreSQL database. The database schema will be automatically applied on the first boot.
   
   ```bash
   docker-compose up --build -d
   ```

4. **Verify the services:**

   The GoLang API will be available at http://localhost:8080.
   
   The PostgreSQL database is exposed on port 5432 for local inspection (e.g., via DBeaver or DataGrip).

6. **Stop the application:**
    ```bash
    docker-compose down

## 🧪 Testing the API (Postman)

Since a frontend is not required for this task, the application is designed to be easily understandable and fully testable via API endpoints. 

To make the evaluation process completely frictionless, I have included a comprehensive Postman collection in the repository.

1. **Locate the File:** Find the `Football League API.postman_collection` file in the root directory of this repository.
2. **Import:** Open Postman, click **Import**, and select the `.json` file.
3. **Pre-configured Environment:** All endpoints are already set up to hit `http://localhost:8080` with the correct headers and realistic JSON request bodies for features like editing match results.
4. **Saved Examples:** I have saved **Example Responses** for every request in the collection. You can click on the examples within Postman to see the exact expected JSON inputs and outputs without even needing to spin up the Docker container first!
