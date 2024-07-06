# Library Management System

This project is a Library Management System built using the MVC architecture and MySQL. The system has separate client and admin portals, secure login functionalities, and comprehensive management features.

## Features

### Separate Client and Admin Portals

- **Authentication & Authorization**: Secure login functionalities for both clients and admins with role-based access control.

### Admin Features

- **Manage Book Catalog**: Admins can list, update, add, and remove books.
- **Approve/Deny Requests**:
  - Checkout and check-in requests from clients.
  - Requests from users seeking admin privileges.

### Client Features

- **View Books**: Clients can view the list of available books.
- **Request Management**:
  - Request checkout and check-in of books from the admin.
  - View their borrowing history.

### Security Features

- **Secure Password Storage**: Passwords are hashed and salted before being stored in the database.
- **JWT-based Session Management**: Implemented custom session management using JWT tokens.

## Getting Started

### Prerequisites

- Go 1.16+
- MySQL
- golang-migrate

### Installation

1. **Clone the Repository**:
   ```sh
   git clone <repository-url>
   cd <repository-directory>

2. **Build and Run the Application:**:
   ```sh
   make run

### Makefile Commands

- **help**: Print help message.
   ```sh
   make help

- **tidy**: Format code and tidy modfile.
   ```sh
   make tidy

- **audit**: Run quality control checks.
   ```sh
   make audit

- **test**: Run all tests.
   ```sh
   make test

- **test/cover**: Run all tests and display coverage.
   ```sh
   make test/cover

- **build**: Build the application.
   ```sh
   make build

- **run**: Run the application.
   ```sh
   make run

- **run/live**: Run the application with reloading on file changes.
   ```sh
   make run/live

- **migrate-up**: Run database migrations up.
   ```sh
   make migrate-up

- **migrate-down**: Run database migrations down.
   ```sh
   make migrate-down

- **migrate-new**: Create a new database migration.
   ```sh
   make migrate-new

- **push**: Push changes to the remote Git repository.
   ```sh
   make push

### Usage

1. **Admin Portal**:
    - Access via `/admin/`.
    - Manage book catalog and user requests.

1. **Client Portal**:
    - Access via `/client/`.
    - View available books, request checkouts/check-ins, and view borrowing history.
