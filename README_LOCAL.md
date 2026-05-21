# Lumina Hotel & Resorts - Go & Nuxt.js Implementation

> [!CAUTION]
> **WINDOWS USERS:** Ensure your folder path **DOES NOT** contain the `&` character or spaces. 
> If your folder is named `lumina-hotel-&-resorts`, you **MUST** rename it to `lumina-hotel-resorts` before running `npm install`.
> This project is optimized for **Node.js v22+** and **Go 1.21+**.

This project follows a clean, Laravel-inspired architecture.

## Architecture Highlights
- **Layered Backend**: Separates concerns into Controllers (HTTP), Services (Business), Repositories (Database), and Models (Data).
- **Nuxt.js Frontend**: Uses the latest Vue 3 patterns with Pinia for state management and Tailwind CSS for styling.
- **SQLite Database**: Self-contained database for easy local development.

## Setup Instructions

### 1. Backend (Go)
1. Navigate to `/backend`
2. Install dependencies: `go mod download`
3. Run the server: `go run main.go`
   - The server will run on `http://localhost:8080`
   - It will automatically create `hotel.db` and migrate schemas on start.

### 2. Frontend (Nuxt.js)
1. Navigate to `/frontend`
2. Install dependencies: `npm install`
3. Start development server: `npm run dev`
   - The frontend will run on `http://localhost:3000`

## Environment Variables
- Create a `.env` in the backend folder for JWT secrets if needed.
- Nuxt config points to `http://localhost:8080/api` by default.
