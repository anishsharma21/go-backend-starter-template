# go-backend-starter-template

Welcome to the repo, this purpose of which is to provide some boiler plate code for backend go projects. I'm using this as a starter point for a lot of future projects. It includes the following features:

- **PostgreSQL Database**: Integrated with a PostgreSQL database.
- **Docker**: Containerised with a Dockerfile.
- **Goose**: Uses Goose for database migration handling.
- **Docker Compose**: Uses Docker Compose for local development setup.
- **Air**: Supports hot module reloading with Air.

## Getting Started

To get started with this project, clone the repository and follow the instructions below (more instructions will be added soon).

```bash
git clone https://github.com/anishsharma21/go-backend-starter-template.git
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

Then, you want to install `air` - this will be used for Hot Module Relooading (HMR), which is when your code will be automatically recompiled and run when changes are made:

```bash
go install github.com/air-verse/air@latest
```

The configuration for `air` is already present in the `.air.toml` file so you can simply run the command `air` on its own from the root of the project, and your server will be started up with HMR.

### Local Database migrations (`goose`)

Use the following command to install `goose` locally as it will not be included in the project as a dependency:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

To run database migrations, you'll first need to set some environment variables that goose will use to connect to your database and locate the migration files. These environment variables are for your local database:

```bash
export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="host=localhost port=5432 user=admin password=secret dbname=mydb sslmode=disable"
export GOOSE_MIGRATION_DIR=migrations
```

Now, with your database running in the background from the previous `docker compose` instructions, check that `goose` is correctly connected to your database by running the following command:

```bash
goose status
```

Ensure your database is running, then, run the following command to run the migration up:

```bash
goose up
```

If the migration went well, you should see `OK` messages next to each applied sql file, and the final line should say `successfully migrated database to version: ...`. You can check the status again to confirm the migrations occurred successfully. Further migration files can be created using the following command:

```bash
goose create {name of migration} sql
```

With the database running, run the following command to run the migration down:

```bash
goose down
```

When updating templates or handlers that render them, make sure to reference the `globalSelectors.go` file where CSS selectors are present in to reduce hard coded values and duplication throughout the code.

You should also set the following environment variable locally for testing and development purposes:

```bash
export JWT_SECRET_KEY=secret
```

Tests run locally use the local postgres database. To replicate the CICD environment, you can clear your database before running the tests. Use the following command to run tests locally:

```bash
go test ./tests -v
```

## Production Deployment

I am using `Railway` to deploy both my postgres database and backend go server. There is a `Dockerfile` in the root of the project that is used for the backend. Private networking with the database is utilised by setting the `DATABASE_URL` and `ENV` variables. The deployment should also wait for CI - i.e. Github actions to complete, before redeploying. Database migrations will be run in production based on whether the environment variable `RUN_MIGRATION` is set to the string `true` - this also requires that you set the following environment variables:

```bash
DATABASE_URL={${{Postgres.DATABASE_URL}}}
GOOSE_DRIVER=postgres
GOOSE_DBSTRING={${{Postgres.DATABASE_URL}}}
GOOSE_MIGRATION_DIR=migrations
JWT_SECRET_KEY={SET_A_SECURE_KEY}
RUN_MIGRATION={TRUE_OR_ANYTHINGELSE}
```

Locally, you can also run/skip database migrations by either setting the `RUN_MIGRATION` environment variable to `true` to run them, or anything else to skip them.

## License

This project is licensed under the MIT License.
Feel free to customize the content further as needed!
