---
created: 2026-03-04 01:13
tree:
tags:
---
# literary-lions-forum

You'll create a web forum that allows users to communicate, associate categories with posts, like/dislike posts & comments, and filter posts.

## The situation

The "Literary Lions" book club is bursting with passionate readers, but their paper-based discussions are becoming a tangled mess. Sticky notes overflow pages, comments get lost in scribbles, and finding past insights feels like navigating a literary labyrinth. New members feel intimidated by the chaotic paper trails, hindering their engagement and knowledge sharing.

The challenge is to design a digital forum that tames the paper pandemonium and fosters a thriving online book club.

This forum should be a haven for bookworms to:

- Spark lively discussions: Imagine a platform where members can dissect chapters, share interpretations, and debate literary merits without the clutter of paper.
- Dive deeper into themes: Enrich the reading experience through dedicated channels(categories) for specific books, genres, in-depth analyses, character studies and author interviews.`
- Preserve insights and recommendations: Past discussions and book reviews can be easily accessed and archived, creating a valuable resource for future members and revisiting past favorites.

By creating a streamlined digital space, we can empower the Literary Lions to embrace the power of online discussions, fostering deeper connections, enriching their reading journeys, and ensuring the book club roars with literary passion for years to come.

## Functional requirements

### Utilize SQLite for Data Management:

- Use [SQLite](https://sqlite.org/index.html), a lightweight and efficient embedded database, to store forum data (users, posts, comments, etc.).
- Use the [go-sqlite3](https://pkg.go.dev/github.com/mattn/go-sqlite3) driver to interface with SQLite.
- Database design: Craft a well-structured database schema using an entity relationship diagram ([ERD](https://smartdraw.com/entity-relationship-diagram/)) to model relationships between entities.
- Query execution: Interact with the database using SQL queries, ensuring the use of at least one SELECT, CREATE, and INSERT query.

### Implement user authentication

During registration, users must specify:

- A unique email address
- A unique username
- A password

To log in, users need to enter their email address and password only. The username is required for display purposes, or tagging users.

- Create login sessions using cookies with expiration dates.

### Enable communication features

- Allow registered users to create posts and comments.
- Allow registered users to associate categories with posts.
- Display posts and comments to all users, even those not registered.

### Implement like/dislike functionality

- Allow registered users to like or dislike posts and comments.
- Display the number of likes and dislikes to all users.

### Add post filtering

- Allow users to filter posts by category, created posts (for registered users), and liked posts (for registered users).

### Dockerize the application

Imagine moving your forum to a new house (server) every time it gets popular. Exhausting, right? Using Docker is like packing everything neatly into boxes (containers), which makes moving a breeze. It sets up your forum on any server instantly, keeping it smooth and buzzing with users, no matter where it lives.

Your mission now is to containerize the forum application utilizing [Docker](https://docs.docker.com/) to guarantee efficient management and deployment.

- Craft a Dockerfile to define the application's environment and dependencies within a container.
- Utilize Docker to construct an image encapsulating the application and its required components.
- Spin up a container from the created image, effectively running the application.
- Enhance organization and management by applying metadata to Docker objects such as images and containers.
- Maintain a clean environment by addressing unused objects to optimize resource usage.

## Extra requirements

### Password encryption

Store passwords securely using encryption.

### Session management

Use UUIDs for session management.

## Bonus functionality

You're welcome to implement other bonuses as you see fit. But anything you implement must not change the default functional behavior of your project.

You may use additional feature flags, command line arguments or separate builds to switch your bonus functionality on.

### Suggestions

- A search feature.
- Allow users to upload images or other files.
- User profile page.

## Resources

None

## Useful links

- [SQLite](https://sqlite.org/index.html)
- [Docker](https://docs.docker.com/)
- ([ERD](https://smartdraw.com/entity-relationship-diagram/))

## What you'll learn

During this project, you'll learn about:

- The basics of web development: HTML, HTTP, sessions, and cookies.
- Using and setting up Docker.
- SQL language and database manipulation.
- The basics of encryption.
- Docker essentials: Installation, configuration, image creation, container management, and best practices.

## Deliverables and Review Requirements

- All source code and configuration files
- A README file with:
    - Project overview
    - Setup and installation instructions
    - Usage guide
    - Any additional features or bonus functionality implemented

During the review, be prepared to:

- Demonstrate your application's functionality
- Explain your code and design choices
- Discuss any challenges you faced and how you overcame them