# Hermes

Starting chat application CODENAME HERMES

## Setup

1. clone project 
    ```sh
    git clone https://github.com/CodeStompNJ/Hermes.git
    ```
2. install Dependencies: 
    Postgres + PGAdmin
3. Using PGAdmin or otherwise, setup database:
    A. create user `seshat` with admin rights.
        ![seshat user setup](https://github.com/CodeStompNJ/Hermes/blob/master/images/setup_database-user_seshat.png?raw=true )
    B. create table `hermes` with user `seshat`.
        ![hermes database setup](https://github.com/CodeStompNJ/Hermes/blob/master/images/setup_database_hermes.png?raw=true)

## Folder structure

TODO

## Installation

1. get all dependencies : 
    ```sh
    go get -d ./...
    ```

## Running server (Backend + Frontend)

Run ```go run main.go``` to start the server and create a listener on port 8000. Current webapp is served on http://localhost:8000/.
