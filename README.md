# Curly Succotash: Virtual Board Game Generator

**Curly Succotash** is an open-source web application that generates unique Dungeons & Dragons (D&D)-style virtual board games. Users can create custom games with specific themes, card counts, and styles, powered by Google Gemini AI for dynamic storytelling and card generation. The application features a Vue.js frontend with multi-language support, card filtering, and PDF export, backed by a Go-based API with SQLite/MySQL storage.

## Features

- **Game Generation**:
  - Create D&D-style board games with customizable themes (e.g., "Fantasy Adventure"), card counts (10-100), and styles ("D&D", "Simple", "Strategy").
  - Optional story description or AI-generated narratives via Google Gemini API.
- **Card Management**:
  - Generate cards with types (`role`, `event`, `item`), names, descriptions, and effects (e.g., "D20 + Strength ≥ 15").
  - Filter cards by type in the frontend UI.
- **Multi-Language Support**:
  - English and Chinese (Traditional) interfaces using `vue-i18n`.
  - Seamless language switching in the frontend.
- **PDF Export**:
  - Export games as structured PDFs with game theme, story, and cards in a 2x2 grid layout.
  - Uses HTML-to-PDF conversion with `wkhtmltopdf` for visually appealing card designs.
- **API-Driven Backend**:
  - RESTful API for game creation, listing, retrieval, and PDF generation.
  - SQLite or MySQL database with GORM for data persistence.
- **D&D Mechanics**:
  - Supports D20-based checks (e.g., attribute + D20 ≥ DC) and D6-based combat (e.g., damage rolls).

## Project Structure

```
curly-succotash/
├── backend/
│   ├── cmd/main/                   # Entry point for the Go server
│   ├── etc/config.yaml            # Configuration file (database, API keys)
│   ├── files/                     # Storage for generated PDFs
│   ├── global/                    # Global variables and initialization
│   │   ├── db.go                  # Database connection setup
│   │   ├── prompt.go             # AI prompt templates for Gemini API
│   │   ├── setting.go             # Configuration loading
│   ├── interfaces/                # Interface definitions
│   │   └── storageengine.go       # Storage engine interface
│   ├── internal/
│   │   ├── ai/                    # Gemini AI integration
│   │   ├── dao/                   # Data Access Objects for database operations
│   │   ├── model/                # Database models (Game, Card)
│   │   ├── service/              # Business logic services
│   ├── migrations/                # Database migrations
│   │   ├── 20250503120000_create_tables.go
│   │   ├── 20250520200000_add_game_info.go
│   │   ├── register.go
│   ├── pkg/                       # Utility packages
│   │   ├── errcode/               # Error codes
│   │   ├── limiter/               # Rate limiting
│   │   ├── logger/                # Logging utilities
│   │   ├── setting/               # Configuration utilities
│   ├── routers/                   # API routing
│   │   ├── api/v1/               # API endpoints
│   │   └── routers.go             # Router initialization
│   ├── running/                   # AI generation logic
│   │   ├── card/
│   │   ├── run_card_response.go
│   │   ├── run_gemini_api.go
│   ├── storage/logs/             # Log storage
│   ├── testing/                   # Test files
├── frontend/
│   ├── index.html                 # HTML entry point
│   ├── src/
│   │   ├── App.vue                # Root Vue component (navigation, layout)
│   │   ├── assets/                # Static assets
│   │   ├── components/
│   │   │   ├── GameForm.vue       # Main game generation and display component
│   │   │   ├── About.vue          # About page
│   │   ├── i18n.js                # Multi-language configuration
│   │   ├── index.css              # Tailwind CSS imports
│   │   ├── main.js                # Vue app initialization
│   │   ├── router/                # Vue Router configuration
│   │   ├── style.css              # Additional styles
│   ├── tailwind.config.js         # Tailwind CSS configuration
│   ├── vite.config.js             # Vite configuration
```

## Tech Stack

- **Backend**:
  - Go 1.21+
  - Gin Web Framework
  - GORM (SQLite/MySQL)
  - Google Gemini AI API
  - wkhtmltopdf (PDF generation)
- **Frontend**:
  - Vue.js 3
  - Vite (build tool)
  - Tailwind CSS
  - vue-i18n (multi-language)
  - vue-router (navigation)
- **Database**:
  - SQLite (default) or MySQL
- **Dependencies**:
  - Backend: `github.com/SebastiaanKlippert/go-wkhtmltopdf`, others in `go.mod`
  - Frontend: `npm` packages in `frontend/package.json`

## Prerequisites

- **Go**: 1.21 or higher (`go version`)
- **Node.js**: 16+ (`node --version`)
- **wkhtmltopdf**: For PDF generation
  ```bash
  sudo apt-get install wkhtmltopdf  # Ubuntu
  brew install wkhtmltopdf         # macOS
  ```
- **SQLite** or **MySQL**: For data storage
- **Google Gemini API Key**: Obtain from Google Cloud Console

## Installation

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/your-org/curly-succotash.git
   cd curly-succotash
   ```

2. **Backend Setup**:
   - Install Go dependencies:
     ```bash
     cd backend
     go mod tidy
     ```
   - Configure `etc/config.yaml`:
     ```yaml
     server:
       port: 8080
     database:
       type: sqlite
       path: games.db
     gemini:
       api_key: "your-gemini-api-key"
     ```
   - Run migrations:
     ```bash
     cd cmd/main
     go run main.go -config ../../etc/config.yaml
     ```

3. **Frontend Setup**:
   - Install Node.js dependencies:
     ```bash
     cd frontend
     npm install
     ```
   - Configure `vite.config.js` for API proxy:
     ```js
     server: {
       port: 5173,
       proxy: {
         '/api': {
           target: 'http://localhost:8080',
           changeOrigin: true,
         },
       },
     }
     ```

4. **Create Output Directory**:
   ```bash
   mkdir -p backend/files
   chmod 755 backend/files
   ```

## Running the Application

1. **Start Backend**:
   ```bash
   cd backend/cmd/main
   go run main.go -config ../../etc/config.yaml
   ```
   - API available at `http://localhost:8080`.

2. **Start Frontend**:
   ```bash
   cd frontend
   npm run dev
   ```
   - UI available at `http://localhost:5173`.

3. **Usage**:
   - Open `http://localhost:5173` in a browser.
   - Enter game details (theme, card count, style, optional story).
   - Generate a game, view saved games, filter cards, and download PDF.
   - Switch between English and Chinese (Traditional) via the language selector.

## API Endpoints

- **POST /api/v1/generate**
  - Description: Generate a new game.
  - Request:
    ```json
    {
      "theme": "Fantasy",
      "cardCount": 20,
      "style": "D&D",
      "description": "An epic quest"
    }
    ```
  - Response:
    ```json
    {
      "game_id": 1,
      "message": "Game generated successfully"
    }
    ```

- **GET /api/v1/games**
  - Description: List all games.
  - Response:
    ```json
    [
      {
        "id": 1,
        "theme": "Fantasy",
        "card_count": 20,
        "style": "D&D",
        "description": "An epic quest"
      },
      ...
    ]
    ```

- **GET /api/v1/games/:id**
  - Description: Get game details with cards.
  - Response:
    ```json
    {
      "id": 1,
      "theme": "Fantasy",
      "card_count": 20,
      "style": "D&D",
      "description": "An epic quest",
      "cards": [
        {
          "id": 1,
          "game_id": 1,
          "type": "role",
          "name": "Warrior",
          "description": "A brave fighter",
          "effect": "D20 + Strength >= 15"
        },
        ...
      ]
    }
    ```

- **GET /api/v1/generate-pdf/:id**
  - Description: Generate and download a PDF for a game.
  - Response: PDF file (`game_<id>.pdf`).

## Database Schema

- **Table: games**
  - `id`: Integer, primary key
  - `theme`: String, game theme
  - `card_count`: Integer, number of cards
  - `style`: String, game style (D&D, Simple, Strategy)
  - `description`: Text, story description
  - `created_at`: Timestamp
  - `is_del`: Integer (0 for active, 1 for deleted)

- **Table: cards**
  - `id`: Integer, primary key
  - `game_id`: Integer, foreign key to `games.id`
  - `type`: String (role, event, item)
  - `name`: String
  - `description`: Text
  - `effect`: Text
  - `is_del`: Integer (0 for active, 1 for deleted)

## Contributing

1. Fork the repository.
2. Create a feature branch (`git checkout -b feature/your-feature`).
3. Commit changes (`git commit -m 'Add your feature'`).
4. Push to the branch (`git push origin feature/your-feature`).
5. Open a Pull Request.

## Troubleshooting

- **Backend API Errors**:
  - Check `logs/app.log` for details.
  - Ensure `GEMINI_API_KEY` is set in `config.yaml`.
- **PDF Generation Fails**:
  - Verify `wkhtmltopdf` installation (`wkhtmltopdf --version`).
  - Ensure `files` directory has write permissions.
- **Frontend Issues**:
  - Clear npm cache: `npm cache clean --force`
  - Reinstall dependencies: `npm install`
- **Database Issues**:
  - Run migrations again if tables are missing.
  - Check `games.db` or MySQL connection.

## Acknowledgments

- Google Gemini AI for story and card generation.
- wkhtmltopdf for PDF rendering.
- Tailwind CSS for responsive UI design.