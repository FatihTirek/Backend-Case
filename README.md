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
   git clone [https://github.com/](https://github.com/)FatihTirek/football-league.git
   cd football-league

2. **Set up environment variables:**

   Copy the example environment file to create your local `.env` configuration.

   ```bash
   cp .env.example .env
   ```

3. **Set up environment variables:**

   Use Docker Compose to build the Go API and initialize the PostgreSQL database. The database schema will be automatically applied on the first boot.
   
   ```bash
   cp .env.example .env
   ```

4. **Verify the services:**

   The GoLang API will be available at http://localhost:8080. The PostgreSQL database is exposed on port 5432 for local inspection (e.g., via DBeaver or DataGrip).

5. **Stop the application:**
    ```bash
    docker-compose down
