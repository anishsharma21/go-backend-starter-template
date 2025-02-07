# go-backend-starter-template

Welcome to the repo, this purpose of which is to provide some boiler plate code for backend go projects. I'm using this as a starter point for a lot of future projects. It includes the following features:

- **PostgreSQL Database**: Integrated with a PostgreSQL database.
- **Docker**: Containerised with a Dockerfile.
- **Goose**: Uses Goose for database migration handling.
- **Docker Compose**: Uses Docker Compose for local development setup.
- **Air**: Supports hot module reloading with Air.

## Work in Progress

This project is a work in progress. Commands to set up the project locally will be provided soon.

## Getting Started

To get started with this project, clone the repository and follow the instructions below (more instructions will be added soon).

```bash
git clone https://github.com/yourusername/go-backend-starter-template.git
cd go-backend-starter-template
```

Then, run the following command to install all the go dependencies:

```bash
go mod download
```

### Developing locally

To develop locally, first run an instance of the local database using `docker-compose`. If you don't have `docker` or `docker-compose`, its pretty easy to install them so go ahead and do that - you can use [this link](https://docs.docker.com/desktop/). Then, once you have both installed (which you can check by running `docker version` and `docker compose version`), you can run the following command to start a local postgres database which will have its data persisted:

```bash
docker compose up -d
```

The `-d` flag is to run it in detached mode - without it, all the logs will appear in your terminal and you will have start a new terminal session to run further commands. It's useful to learn about `docker` and `docker compose` so you understand how to build images and manage containers locally. You can leave this postgres database running, but if you ever want to stop it, you can run `docker compose down`.

## License

This project is licensed under the MIT License.
Feel free to customize the content further as needed!
