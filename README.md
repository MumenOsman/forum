# The Daily Inkwell

This is a project that was developed as part of the Kood/Sisu learning path. The project is a website, a digital platform designed for avid readers o discuss books, share insights, and engage in meaningful conversations. 
## Features

- **User Authentication**: Secure registration and login system using bcrypt for password hashing and session-based authentication.
- **Dynamic Forum**: Create posts, contribute comments, and engage with others through a sleek, responsive interface.
- **Private Messaging (DMs)**: A fully functional internal inbox system allowing registered users to send private text messages asynchronously.
- **Interactive Sidebar**: Real-time conversation management sidebar to quickly switch between active DM threads.
- **Profile Customization**: Users can personalize their accounts with "About Me" sections and unique profile pictures.
- **Engagement System**: Like and dislike functionality for both posts and comments, with dynamic count updates.
- **Smart Filtering**: Filter discussions by category (e.g., Fantasy, Sci-Fi, Mystery), your own authored posts, or posts you've liked.
- **Global Search**: Search through the entire book archive by title or content.
- **Themed UI**: Modern theme (Nord) with premium typography and interactive hover effects.

## Technologies Used

### Backend
- **Go (Golang)**: Core application logic and HTTP server.
- **SQLite3**: Lightweight, file-based database with CGO enabled for performance.
- **Gorilla Sessions**: Management of user sessions and state.
- **Bcrypt**: Password encryption.

### Frontend
- **HTML5 Templates**: Server-side rendering for security and speed.
- **Vanilla JavaScript**: Dynamic UI interactions, including profile picture uploads and real-time sidebar updates.
- **Tailwind CSS (via CDN)**: Modern utility-first styling for a premium look and feel.
- **Google Fonts & Material Icons**: Professional typography (Inter, Serif) and consistent iconography.

### Infrastructure
- **Docker**: Multi-stage containerization for seamless deployment.
- **Alpine Linux**: Minimalist base image for the final production container.
- **Git**: Version-controlled development with mirroring between Gitea and GitHub.

## 🚀 How to Clone and Run

### 1. Prerequisites
- **Docker**: The application is fully containerized. No local Go or SQLite installation is required.

### 2. Clone the Repository
```bash
git clone https://gitea.kood.tech/mumenosman/forum.git
cd forum
```

### 3. Build and Start (Docker)
To build the Docker image and start the container on port 8080:
```bash
docker build -t forum-app .
docker run -p 8080:8080 forum-app
```
The application will be accessible at: **[http://localhost:8080](http://localhost:8080)**

## 👤 Demo Personas

The application is pre-seeded with 5 unique personas to showcase multi-user interactions. All accounts use the default password: `123456`.

| Persona | Email | Specialty |
| :--- | :--- | :--- |
| **Elara** | `test@test.com` | Fantasy & Sci-Fi (The Worldbuilder) |
| **Marcus** | `marcus@inkwell.com` | Historical Fiction & Biography (The Historian) |
| **Chloe** | `chloe@inkwell.com` | Romance & Mystery (The Escapist) |
| **Julian** | `julian@inkwell.com` | Non-Fiction & General (The Pragmatist) |
| **Sarah** | `sarah@inkwell.com` | Fiction & Mystery (The Aspiring Author) |

---
*Developed by Mumen Osman for the Kood/Sisu curriculum.*
