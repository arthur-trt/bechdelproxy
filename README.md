# Bechdel Proxy

Bechdel Proxy is a project designed to integrate the Bechdel Test score into movie profiles on Overseerr, a platform that already displays movie ratings from IMDb and Rotten Tomatoes. The Bechdel Test is a feminist evaluation based on three key criteria. This proxy aims to cache and optimize Bechdel Test data for local use, minimizing direct API requests to [bechdeltest.com](http://bechdeltest.com), a community-driven site.

## Features

- **Cache Bechdel Test scores:** Fetch and store data locally from [bechdeltest.com](http://bechdeltest.com) for efficient use.
- **Integration with Overseerr:** Enhance movie profiles with Bechdel Test scores.
- **Dual functionality:**
  - Serve as an API when run normally.
  - Update the database with fresh data using the `--update` flag.
- **Optimized performance:** Filters and deduplicates data from [bechdeltest.com](http://bechdeltest.com) to ensure high efficiency.
- **Debug mode:** Allows loading data from a local `raw_data.json` file instead of fetching from the Bechdel Test API.

## Installation

### Prerequisites

- [Go](https://go.dev/) (latest version)
- [PostgreSQL](https://www.postgresql.org/) database
- Docker (optional, for containerized deployment)

### Clone the Repository
```bash
git clone https://github.com/arthur-trt/bechdel-proxy.git
cd bechdel-proxy
```

### Environment Variables

Set the following environment variables in a `.env` file or your system:

| Variable       | Description                                                | Default         |
|----------------|------------------------------------------------------------|-----------------|
| `LOG_LEVEL`    | Logging level (`DEBUG`, `INFO`, `WARN`, `ERROR`)           | `WARN`          |
| `PGSQL_HOST`   | PostgreSQL database host                                   |                 |
| `PGSQL_USER`   | PostgreSQL username                                        |                 |
| `PGSQL_PASS`   | PostgreSQL password                                        |                 |
| `PGSQL_DB`     | PostgreSQL database name                                   |                 |
| `PGSQL_PORT`   | PostgreSQL port                                            | `5432`          |
| `PGSQL_TZ`     | PostgreSQL timezone                                        | `Europe/Paris`  |
| `PGSQL_BATCH_SIZE` | Number of records to process in a single batch during updates | `1000`          |

### Build the Application

```bash
go build -o bechdelproxy .
```

## Usage

### Running the API
To start the API server:
```bash
./bechdelproxy
```

### Updating the Database
To fetch and update Bechdel Test data from [bechdeltest.com](http://bechdeltest.com):
```bash
./bechdelproxy --update
```

### Debug Mode
If `LOG_LEVEL` is set to `DEBUG`, the program will load movie data from a local file named `raw_data.json` instead of fetching it from the Bechdel Test API.

## Docker Deployment

A Docker setup is included to simplify deployment. Use the provided `Dockerfile` and `docker-compose.yml` for containerized execution.

### Build the Docker Image
```bash
docker build -t bechdelproxy .
```

### Run with Docker Compose

```bash
docker-compose up
```

The `docker-compose.yml` includes:
- A container for the Bechdel Proxy application.
- A PostgreSQL database container.

## API Endpoints

### Get Bechdel Rating by IMDb ID

**Endpoint:** `GET /imdb/:imdbid`

**Description:** Retrieves the Bechdel Test rating for a movie by its IMDb ID.

**Response:**
```json
{
    "ID": 0,
    "CreatedAt": "string",
    "UpdatedAt": "string",
    "DeletedAt": null,
    "title": "string",
    "id": 0,
    "imdbid": "string",
    "rating": 0
}
```

## Database Schema

The application uses the following schema for the `movies` table:

```go
type Movie struct {
	gorm.Model
	Title     string `json:"title"`
	BechdelID int    `gorm:"uniqueIndex:idx_bechdel_imdb" json:"id"`
	IMDBID    string `gorm:"uniqueIndex:idx_bechdel_imdb" json:"imdbid"`
	Rating    int    `json:"rating"`
}
```

- **Title:** The title of the movie.
- **BechdelID:** Unique ID from Bechdel Test data.
- **IMDBID:** IMDb ID (prefixed with `tt` if not already present).
- **Rating:** The Bechdel Test rating (0-3).

## Logging

The application uses a custom logging package based on `logrus`. Logging levels can be adjusted using the `LOG_LEVEL` environment variable. Logs include details about:
- API operations
- Database updates
- Errors and warnings

## Contributing

Contributions are welcome! Feel free to submit issues or pull requests to improve the project.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.

## Acknowledgments

- [bechdeltest.com](http://bechdeltest.com) for providing the data.
- [Overseerr](https://overseerr.dev/) for inspiration and integration.

