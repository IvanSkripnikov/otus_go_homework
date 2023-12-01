## Overview

This repository is a template for a Todo application in the Golang that can run on any platform as long as Docker is running.

MySQL is used as the data store.

## Endpoints

Method | Path | Description
--- | --- | ---
GET | `/tasks` | Get all tasks |
POST | `/tasks` | Create a new task |
GET | `/tasks/{id}` | Get a task |
PUT | `/tasks/{id}` | Update a task |
DELETE | `/tasks/{id}` | Delete a task |

## Usage

```shell
$ cp .env.example .env

$ docker compose up --build
```

## Tips

### Connect db container

```shell
$ docker compose exec -it db /bin/bash -c "mysql -uroot -p<PASSWORD>"
```

### Import API test file with Thunder client

1. Install extension from [here](https://marketplace.visualstudio.com/items?itemName=rangav.vscode-thunder-client)
2. Import [JSON file](https://github.com/ega4432/go-rest-api-docker/blob/main/.vscode/thunder-collection_go-rest-api-docker.json)
