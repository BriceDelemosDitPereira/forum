# Forum

## Overview

The "Forum" project is a web-based forum application developed in Go. It provides users with a platform for communication, allowing them to create posts, comment on posts, associate categories with posts, like and dislike posts/comments, and filter posts based on various criteria.

## Features

- **Communication:** Users can create posts and comments to communicate with each other.
- **Categories:** Posts can be associated with one or more categories, allowing for better organization and navigation.
- **Likes and Dislikes:** Registered users can express their preference for posts and comments by liking or disliking them.
- **Filtering:** Users can filter posts based on categories, posts they created, and posts they liked.
- **Authentication:** User authentication is implemented to ensure secure access to the forum. Users can register with unique email addresses and passwords, and login to create posts and comments.
- **Session Management:** Each user session is managed using cookies, allowing users to remain logged in for a specified period.
- **Database:** Data such as users, posts, comments, etc., is stored using the SQLite database library. Queries are used for operations such as creating, retrieving, and inserting data.
- **Dockerization:** The project is containerized using Docker, providing a consistent environment for development and deployment.

## Usage

To run the forum project:

1. Ensure you have Docker installed on your system.
2. Clone the project repository.
3. Navigate to the project directory.
4. Build the Docker image using the provided Dockerfile.
5. Run the Docker container.

## Authentication and Registration

- Users must register with a unique email address, username, and password.
- Passwords are encrypted before storage.
- Authentication is implemented using cookies, with each session having an expiration date.

## Communication

- Registered users can create posts and comments.
- Posts can be associated with one or more categories.
- Posts and comments are visible to all users, registered or not.

## Likes and Dislikes

- Registered users can like or dislike posts and comments.
- The number of likes and dislikes is visible to all users.

## Filter

- Users can filter posts by categories, posts they created, and posts they liked.
- Filtering options are available only for registered users and are specific to the logged-in user.

## Docker

- Docker is used for containerization of the forum project.
- Docker provides a consistent environment for development and deployment, ensuring seamless execution across different systems.

## Resources

- [SQLite Documentation](https://www.sqlite.org/docs.html)
- [Docker Documentation](https://docs.docker.com/)

## Co-developers

- [Delestre Thomas](https://github.com/Thomas-Delestre)
- [Lovergne Raphael](https://github.com/Ne0Jiku)
- [Marchais Mickael](https://github.com/Jeancrock)