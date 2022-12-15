[![codecov](https://codecov.io/gh/graaphscom/monogo/branch/master/graph/badge.svg)](https://codecov.io/gh/graaphscom/monogo)

# monogo

Monorepo containing Golang modules.

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
```