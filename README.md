# MedSearch

A comprehensive medication search platform that crawls pharmaceutical websites and provides a searchable API for drug information. Built with Go backend, React frontend, MongoDB for data storage, and Elasticsearch for search capabilities.

## 🏗️ Architecture

MedSearch consists of several components:

- **Backend API** (Go/Gin): RESTful API for drug search functionality
- **Web Crawler** (Go/Colly): Scrapes drug information from pharmaceutical websites
- **Frontend** (React/Vite): Modern web interface for searching medications
- **Database** (MongoDB): Stores drug information and metadata
- **Search Engine** (Elasticsearch): Provides fast and accurate search capabilities

## 🚀 Features

- **Drug Information Crawling**: Automated web scraping of pharmaceutical websites
- **Advanced Search**: Full-text search across drug names, descriptions, and ingredients
- **RESTful API**: Clean API endpoints for integration
- **Modern Web Interface**: Responsive React frontend with real-time search
- **Docker Support**: Easy deployment with Docker Compose

## 📋 Prerequisites

- Go 1.24.3 or higher
- Node.js 18+ (for frontend development)
- Docker and Docker Compose (for containerized deployment)
- MongoDB 7.0+
- Elasticsearch 8.13.2+

## 🛠️ Installation

### Option 1: Docker Deployment (Recommended)

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd MedSearch
   ```

2. **Start the services**
   ```bash
   docker-compose up -d
   ```

   This will start:
   - MongoDB on port 27017
   - Elasticsearch on port 9200
   - Backend API on port 8080

3. **Run the crawler to populate data**
   ```bash
   docker-compose exec backend ./medsearch crawl
   ```

4. **Start the frontend**
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

### Option 2: Local Development

1. **Set up the database**
   ```bash
   # Start MongoDB and Elasticsearch using Docker
   docker-compose up -d mongo elasticsearch
   ```

2. **Configure environment variables**
   ```bash
   export MONGO_URI="mongodb://localhost:27017"
   export DB_NAME="medsearch"
   export ELASTICSEARCH_URL="http://localhost:9200"
   ```

3. **Install Go dependencies**
   ```bash
   go mod download
   ```

4. **Run the crawler**
   ```bash
   go run main.go crawl
   ```

5. **Start the API server**
   ```bash
   go run main.go api
   ```

6. **Start the frontend**
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

## 📖 Usage

### API Endpoints

#### Search Drugs
```
GET /v1/api/search?text={search_term}
```

**Parameters:**
- `text` (required): Search term for drug name, description, or ingredients

**Response:**
```json
[
  {
    "ID": "507f1f77bcf86cd799439011",
    "Name": "Paracetamol 500mg",
    "Description": "Pain reliever and fever reducer",
    "Category": "Analgesic",
    "Ingredients": ["Paracetamol"],
    "Manufacturer": "ABC Pharmaceuticals",
    "DosageForm": "Tablet",
    "Packaging": "10 tablets per blister",
    "Price": "50,000 VND"
  }
]
```

### Command Line Interface

The application supports two main commands:

1. **API Server**
   ```bash
   go run main.go api
   ```
   Starts the REST API server on port 8080

2. **Web Crawler**
   ```bash
   go run main.go crawl
   ```
   Starts the web crawler to scrape drug information

```

## 🔧 Configuration

### Environment Variables

- `MONGO_URI`: MongoDB connection string (default: `mongodb://localhost:27017`)
- `DB_NAME`: Database name (default: `medsearch`)
- `ELASTICSEARCH_URL`: Elasticsearch URL (default: `http://localhost:9200`)

### Crawler Configuration

The crawler is configured to scrape from `nhathuoclongchau.com.vn` and can be customized in `app/config/crawler_config.go`.

## 🧪 Development

### Project Structure

```
MedSearch/
├── app/
│   ├── api/           # API handlers
│   ├── config/        # Configuration files
│   ├── crawler/       # Web crawler implementation
│   ├── database/      # Database connection and repositories
│   ├── handlers/      # HTTP handlers
│   ├── logger/        # Logging utilities
│   ├── models/        # Data models
│   └── routes/        # Route definitions
├── frontend/          # React frontend application
├── docker-compose.yml # Docker services configuration
├── Dockerfile         # Backend container configuration
├── go.mod            # Go module dependencies
└── main.go           # Application entry point
```

### Running Tests

```bash
go test ./...
```

### Building

```bash
# Build backend
go build -o medsearch main.go

# Build frontend
cd frontend
npm run build
```

## 🐳 Docker

### Building Images

```bash
# Build backend image
docker build -t medsearch-backend .

# Build frontend image (if needed)
cd frontend
docker build -t medsearch-frontend .
```

### Running with Docker Compose

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

## 📝 Logging

The application logs to `app.log` with different log levels:
- INFO: General application information
- ERROR: Error messages and exceptions
- DEBUG: Detailed debugging information

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.
