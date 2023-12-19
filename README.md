# Pure-FTPd Authd Database Authenticator
[![status-badge](https://ci.cluster.lan.crystalnet.org/api/badges/5/status.svg)](https://ci.cluster.lan.crystalnet.org/repos/5)
![GitHub release (with filter)](https://img.shields.io/github/v/release/psych0d0g/pure-ftpd-paperless-dbauth)

This small Go program is designed to be plugged into Pure-FTPd's authd program, providing authentication against a Paperless-NGX MariaDB or PostgreSQL database.

## Prerequisites

- Go installed on your machine
- Pure-FTPd with authd support
- MariaDB or PostgreSQL database configured for Paperless-NGX

## Features

- Authenticate users against a Paperless-NGX database
- Support for both MariaDB and PostgreSQL databases
- Automatic dependency resolution using Go Modules

## Getting Started

### Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/psych0d0g/pure-ftpd-paperless-dbauth.git
    ```

2. Change into the project directory:

    ```bash
    cd pure-ftpd-paperless-dbauth
    ```

3. Run the build script to compile the binary:

    ```bash
    ./build.sh
    ```

### Configuration

1. Set up environment variables:

    ```bash
    export PAPERLESS_DBHOST="your_database_host"
    export PAPERLESS_DBPORT="your_database_port"
    export PAPERLESS_DBNAME="your_database_name"
    export DB_USER="your_database_user"
    export PAPERLESS_DBPASS="your_database_password"
    export PAPERLESS_DBENGINE="postgres"  # or "mysql" for MariaDB
    export PAPERLESS_CONSUMPTION_DIR="your_paperless_consumption_dir"
    ```

2. Set up authd configuration to use the compiled binary.

### Test

Run the compiled binary to authenticate users against the Paperless-NGX database.

```bash
AUTHD_ACCOUNT=username AUTHD_PASSWORD=password ./verify_pw
```

### authd configuration

please refer to pure-ftpd documentation on how to integrate it with authd: https://github.com/jedisct1/pure-ftpd/blob/master/README.Authentication-Modules

### Database Schema

Ensure your Paperless-NGX database has a table similar to the following:

```sql
CREATE TABLE auth_user (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    last_login TIMESTAMP,
    is_superuser BOOLEAN NOT NULL,
    first_name VARCHAR(30),
    last_name VARCHAR(30),
    email VARCHAR(255),
    is_staff BOOLEAN NOT NULL,
    is_active BOOLEAN NOT NULL,
    date_joined TIMESTAMP NOT NULL
);
```

### Built With

    - Go - The Go Programming Language
    - github.com/go-sql-driver/mysql - MySQL driver for Go
    - github.com/lib/pq - PostgreSQL driver for Go

### Contributing

i am open to any improvement suggestions via issues or pull reqests

### License

This project is licensed under the MIT License - see the LICENSE.md file for details.

### Acknowledgments

    -  ChatGPT
