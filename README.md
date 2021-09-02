# GoKeeper

GoKeeper is a simple password manager.

<a href="https://ibb.co/HY6PZzL"><img src="https://i.ibb.co/Xp15Hsf/example.png" alt="GoKeeper" width="100%" /></a>


## TL;DR
The easiest way to run the app is with `docker-compose`. Navigate to `support/docker` and run
`docker-compose up`. Beware that when running the application this way, data will not persist 
between new docker instances.


## About
GoKeeper is a simple password manager, primarily made to try out some Go libraries/frameworks like
`gqlgen`, `testify`, `testcontainers-go` and `upper`. 
The backend is written in Go (version 1.16), and its whole API is GraphQL compliant. 
The backend currently works only with a postgres database for data storage.

GoKeeper uses an `argon2id` implementation to hash user's master passwords. Stored user passwords are
encrypted with an `AES-256` encryption which combines a secret salt with the hashed user's master password.

Test code coverage for backend code is 100% (excluding `main.go` and utility functions).
The tests with coverage can be run with `go test ./app/... -coverprofile coverage.out -p 1 | grep -v "no test files"`
from the project's root directory. The current implementation of docker containers for integration tests won't work
in parallel since the test database port isn't dynamically assigned, hence the `-p 1`. The `| grep -v "no test files"`
part is optional and is there to not clutter the output with results from packages that don't require tests. 
(such as generated third party packages, packages containing only data models, test utilities package, etc...)

The application currently has a frontend written in `React.js` and `Apollo GraphQL`. It is very minimal and poorly implemented,
since the focus was primarily on the backend, but I wanted some basic ui nonetheless alongside the application. 
It uses no state management framework at the moment and there are no tests for it. Also, the design is barely tolerable :)


## Improvements
This is a list of possible improvements that I think the app needs:

- [ ] Improve authentication by saving session data on backend
- [ ] Expand current JWT implementation to allow for token refresh and real sign out
- [ ] Create more options for password hashing/encryption and allow users to customize those options
- [ ] Allow users to change their master password
- [ ] Create a real sign in mechanism (e-mail confirmation, password sanity check, etc...)
- [ ] Expand the data model (allow to store additional data alongside a password, etc...)
- [ ] Use `Gqlgen` built-in features(custom data types, input validations, etc...) to reduce concerns on custom API code
- [ ] Refactor database code, possibly using some other framework or ORM
- [ ] Implement a database reconnection mechanism
- [ ] Separate go interfaces and service implementations
- [ ] Separate unit and integration tests
- [ ] Achieve 100% code coverage only with unit tests (not really feasible elegantly for database code using `Upper`)
- [ ] Allow integration tests to run in parallel by dynamically assigning ports to test docker containers
- [ ] Make a production ready build that is containerized and doesn't use a development server for the frontend
- [ ] Rewrite the whole frontend with better design, and a state management framework

## How to run it
Running the whole application locally on your machine is easy, but you'll need `go` and `npm` installed alongside `docker` and `docker-compose`.
Besides, you can run it as mentioned in the `TL;DR` section.

### Backend
You can run the backend code by just building and running the `main.go` file.
In the root project directory there is a bash script provided that builds and outputs a binary 
in the `build` directory and then runs it. You can run it with `./run.sh` from the project directory.
By running the backend application you'll get access to GraphQL playground in case you don't want to bother with the frontend.

### Database
There is a docker compose file for a configured postgres database with a persistence volume, located at
`support/docker/docker-compose-postgres.yml` from the project's root directory. 
A simple `docker-compose -f docker-compose-postgres.yml up` should set it up and run it.
You could also ignore it and configure a GoKeeper database locally with the same configuration as in the docker compose file.
Files for data migrations/tables creation are at `support/database/postgres/migration`.

### Frontend
The frontend code is located in the `ui` directory of the project root directory. From there you can run
`npm install` followed by a `npm start` to run the frontend on a separated development server.



> :warning: If it is not already obvious, this application was written for fun, learning and to try things out, 
and as such is not meant to be a substitute for a real password manager, since it lacks many features and proper infrastructure!

