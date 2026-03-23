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

## Demo Credentials & Personas
The application is pre-seeded with 5 unique personas, each with their own interests and writing styles. All accounts use the password: `123456`.

| Persona | Email | Specialty |
| :--- | :--- | :--- |
| **Elara** | `test@test.com` | Fantasy & Sci-Fi (The Worldbuilder) |
| **Marcus** | `marcus@inkwell.com` | Historical Fiction & Biography (The Historian) |
| **Chloe** | `chloe@inkwell.com` | Romance & Mystery (The Escapist) |
| **Julian** | `julian@inkwell.com` | Non-Fiction & General (The Pragmatist) |
| **Sarah** | `sarah@inkwell.com` | Fiction & Mystery (The Aspiring Author) |

You can use any of these to test multi-user interactions and see how the different writing styles populate the "Daily Inkwell."
