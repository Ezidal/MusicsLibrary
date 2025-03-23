# Songs Library API

This project is a Songs Library API built with Go. It allows you to manage a library of songs, including adding, updating, deleting, and retrieving songs. The API also supports pagination for song texts.

## Features

- Add, update, delete, and retrieve songs
- Pagination support for song texts
- Database migrations
- Swagger documentation

## Getting Started

### Prerequisites

- Go 1.23.5 or later
- Docker

### Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/yourusername/songs-library-api.git
    cd songs-library-api
    ```

2. Copy the example environment file and update it with your configuration:

    ```sh
    cp .env.example .env
    ```
    !!! change the address of the external api !!!

3. Start the PostgreSQL database using Docker Compose:

    ```sh
    docker-compose up -d
    ```

4. Install Go dependencies:

    ```sh
    go mod tidy
    ```

5. Run the database migrations:

    ```sh
    go run cmd/onlineLibrary/main.go
    ```

### Running the Server

Start the server:

```sh
go run cmd/onlineLibrary/main.go
```

The server will start on the port specified in the `.env` file (default is `8080`).

### API Documentation

Swagger documentation is available at:

```
http://localhost:8080/swagger/index.html
```

### API Endpoints

- `GET /songs` - Retrieve all songs
- `GET /songs/{id}/text` - Retrieve the text of a song with pagination
- `DELETE /songs/{id}` - Delete a song by ID
- `PUT /songs/{id}` - Update a song by ID
- `POST /songs` - Add a new song

### Example Requests

#### Add a Song

```sh
curl -X POST http://localhost:8080/songs \
-H "Content-Type: application/json" \
-d '{
  "group_name": "The Beatles",
  "song_name": "Hey Jude",
  "release_date": "1968-08-26",
  "text": "Hey Jude, don't make it bad...",
  "link": "https://example.com/hey-jude"
}'
```

#### Get All Songs

```sh
curl http://localhost:8080/songs
```

#### Get Song Text with Pagination

```sh
curl http://localhost:8080/songs/1/text?page=1&limit=2
```

#### Update a Song

```sh
curl -X PUT http://localhost:8080/songs/1 \
-H "Content-Type: application/json" \
-d '{
  "group_name": "The Beatles",
  "song_name": "Hey Jude",
  "release_date": "1968-08-26",
  "text": "Hey Jude, don't make it bad...",
  "link": "https://example.com/hey-jude"
}'
```

#### Delete a Song

```sh
curl -X DELETE http://localhost:8080/songs/1
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License.