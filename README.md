[![codecov](https://codecov.io/gh/graaphscom/monogo/branch/main/graph/badge.svg)](https://codecov.io/gh/graaphscom/monogo)

# monogo

A set of tools for building another tools and apps.

| module name | purpose                                                                                                             |
|-------------|---------------------------------------------------------------------------------------------------------------------|
| asciiui     | making possible plain-text based user interfaces                                                                    |
| crawler     | providing reusable utilities for web crawling                                                                       |
| compoas     | building, composing and serving [OpenAPI Specification](https://github.com/OAI/OpenAPI-Specification) (aka Swagger) |
| dbmigrat    | maintaining database migrations in several locations (repos)                                                        |
| fa2ts       | extracting svg paths from [Font Awesome files](https://fontawesome.com/download) to TypeScript variables            |
| rspns       | building and documenting consistent HTTP responses                                                                  |

# Contributing

## Getting started

1. `brew install go golangci-lint`
2. Due to the golangci problem, diffutils are needed: `brew install diffutils`
3. Check the [Makefile](./Makefile) for available rules

## Testing unpublished code
Use the `replace` directive to have not-yet-released changes available for your local project.

See: [Requiring module code in a local directory](https://go.dev/doc/modules/managing-dependencies#local_directory)

An example excerpt of the `go.mod` file:
```
replace github.com/graaphscom/monogo/crawler => ../../monogo/crawler
```

## Releasing

## Adding a new module

1. Add an GitHub action file `./.github/workflows/<mod_name>.yml` with the following contents
   (after replacing `<mod_name>`):

```yaml
name: <mod_name>
on:
  push:
    paths:
      - "<mod_name>/**"
jobs:
  checks:
    uses: ./.github/workflows/go-mod-checks.yml
    with:
      go_module_name: <mod_name>
# Set the below inputs to make a database available for your tests.
# You can connect to the database by using its url value stored in
# the environment variable which name is defined by "postgres_url_env_var_name" input.
#     postgres_password: <some_password>
#     postgres_user: <some_username>
#     postgres_url_env_var_name: <SOME_ENV_NAME e.g. DBMIGRAT_TEST_URL_DB>
```

2. Add module name to the [Makefile](./Makefile), for example:
```diff
- modules := asciiui compoas crawler dbmigrat fa2ts rspns
+ modules := asciiui compoas crawler dbmigrat fa2ts rspns <mod_name>
```