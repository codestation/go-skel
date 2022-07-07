# Go skeleton base project

## Directory summary

The following directories follow somewhat the clean architecture design.

* `api`: external layer, only processes the request data and passes it to the business layer on `app`
* `app`: business layer, doesn't put code here related to request or databases. Uses the `store`
if it needs to operate on the data.
* `store`: data repository layer. Only this layer connects to the database. Reads and writes data to structs from
  `model`.
* `model`: Data definitions. There is no business logic in this layer.

The following directories handles the application procedures unrelated to the business logic.

* `cmd`: Entrypoint for the application on the command line.
* `config`: Application configuration read from the environment, config file or command line flags.
* `db`: migration files for the database.
* `locales`: project translations.
* `testdata`: data to seed database for the testcases.
* `version`: application version info (commit and build date).
* `web`: static files served from the api.
