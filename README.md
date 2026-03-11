# Literary Lions Forum

A digital forum for the "Literary Lions" book club to discuss books, share insights, and organize categories.

## Features (Planned)
- User registration and authentication.
- Post creation with categories.
- Commenting system.
- Like and dislike functionality.
- Post filtering by category and user activity.
- Dockerized deployment.

## Tech Stack
- **Go** (Backend)
- **SQLite3** (Database)
- **Docker** (Containerization)

## How to Build and Run

### Dependencies
- **Docker**: The application is fully containerized. You only need Docker (or Docker Desktop) installed on your system to build and run the forum. No local Go or SQLite installation is required.

### Build the Application
To build the Docker image, run the following command in the root directory of the project:
```bash
docker build -t forum-app .
```

### Run the Application
Once the image is built, start the container and map it to port 8080 using:
```bash
docker run -p 8080:8080 forum-app
```
The application will be accessible in your web browser at `http://localhost:8080`.
