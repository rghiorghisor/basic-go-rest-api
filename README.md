# Basic Go REST Api
[![Build Status](https://github.com/rghiorghisor/basic-go-rest-api/workflows/build/badge.svg)](https://github.com/rghiorghisor/basic-go-rest-api/actions?query=workflow%3Abuild)
[![Go Report](https://goreportcard.com/badge/github.com/rghiorghisor/basic-go-rest-api)](https://goreportcard.com/report/github.com/rghiorghisor/basic-go-rest-api)
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://lbesson.mit-license.org/)

This application (boilerplate) is a take on developing a simple Go REST API, backed up by a MongoDb/Embedded Bolt database, using principles of Hexagonal Architecture and SOLID principles.

**Note** that this is not work done, in any way. For an application to be production ready additional features must be implemented. Some will be added to this repository in the future, other will probably remain the developer's future duties.

## Requirements

- [Git](https://git-scm.com/)
- [Go](https://golang.org/dl/)
- [mongoDB](https://www.mongodb.com/) (*Optional*)

## Details & Features

- RESTful endpoints (following [RFC 7231](https://tools.ietf.org/html/rfc7231))
- Standard CRUD operations;
- Error handling (including response JSON generation);
- Clean Architecture code organization (use case centric);
- 3tier application with:
  - RESTful API as presentation layer;
  - mongoDB or embedded BoltDB as data layer;
- Switch between local (embedded BoltDB) or remote (mongoDB) storages;
- Configurable through YAML files.

### Implementation details
Some of the implementation details one can analyze or take note from this application:
- How to load configuration YAMLs using [Viper](https://github.com/spf13/viper);
- How to specify the ENV variable name in configuration files (loaded with [Viper](https://github.com/spf13/viper));
- How to setup a configurable logger using [Logrus](https://github.com/sirupsen/logrus);
- How to setup [x-cray/logrus-prefixed-formatter](https://github.com/x-cray/logrus-prefixed-formatter);
- How to organize a use case (feature) in 3 layers;
- How to create a HTTP server using [gin-gonic/gin](https://github.com/gin-gonic/gin);
- How to implement standard RESTful CRUD operations (RFC 7231);
- How to use mongoDB GO driver ([mongodb/mongo-go-driver](https://github.com/mongodb/mongo-go-driver));
- How to setup and use [Storm](https://github.com/asdine/storm).

## Getting Started

Download modules:
```console
go mod download
```

Run application locally:
```console
go run ./cmd/api/main.go
```

Run tests:
```console
go test ./...
```

## Project layout

```
.
├── appserver              Contains the main application controls implementation;
├── cmd                    Main applications of the project;
│   └── api                The server application API (the entry point);
├── config                 Configuration logic and configuration files;
├── container              Contains the DI container implementation;
├── errors                 Application errors and error logic;
├── logger                 Application logger and logic;
├── model                  Model (entities) definitions and logic;
├── property               The entire property use case and dependencies;
│   ├── gateway            The gateways implementations;
|   |   ├── http           The HTTP gateways (Controllers);
|   |   └── storage        The storage gateway implementations;
|   |       ├── bolt      The bolt embedded database gateway (Repository); 
|   |       └── mongo      The mongoDB gateway (Repository);
│   └── service            The property business logic;
├── server                 The server application logic and dependencies;
│   ├── http               The HTTP server implementations;
|   └── storage            The server's storage overall implementation;
├── tests                  Contains additional files for testing purposes;
|   └── config             Configuration files used by config loading tests;
└── util                   Application overall utilities.
```

## Configuration
The application configuration is achieved by the `./config/config.yml` file loaded using the [Viper](https://github.com/spf13/viper) module, with a few customizations.

### Properties

| Name | Description |
| --- | --- |
| `logger.main.level` | The logger level. Accepted values are (*case insensitive*): `panic`, `fatal`, `error`, `warn`, `warning`, `info`, `debug`, `trace`. If none is present the default `info` is considered. Read more about [Logrus Levels](https://github.com/sirupsen/logrus#level-logging). |
| `logger.main.format` | The logger message format. Accepted values are (*case insensitive*): `text`, `json`. If none is present the default `info` is considered.|
| `logger.main.dir` | The directory where all log files are placed. If it does not exist its creation will be attempted. Default value is `./logs`.|
| `logger.main.file-name` | The name of the log file name. Default value is `basic-go-rest-api`.|
| `logger.main.prefix` | The default prefix associated with this logger; will be added along with all messages. Default value is `main`. |
| `logger.main.with-console` | Boolean value that if `true` will also print the log messages to console; otherwise the messages can be found only in the log files. Default value is `false`. |
| `logger.access.level` | The logger level. Accepted values are (*case insensitive*): `panic`, `fatal`, `error`, `warn`, `warning`, `info`, `debug`, `trace`. If none is present the default `info` is considered. Read more about [Logrus Levels](https://github.com/sirupsen/logrus#level-logging). |
| `logger.access.format` | The logger message format. Accepted values are (*case insensitive*): `text`, `json`. If none is present the default `info` is considered.|
| `logger.access.dir` | The directory where all log files are placed. If it does not exist its creation will be attempted. Default value is `./logs`.|
| `logger.access.file-name` | The name of the log file name. Default value is `access`.|
| `logger.access.prefix` | The default prefix associated with this logger; will be added along with all messages. Default value is `access`. |
| `logger.access.with-console` | Boolean value that if `true` will also print the log messages to console; otherwise the messages can be found only in the log files. Default value is `false`. |
| `server.http.port` | The port that the server listens on. Default value is `8080`. |
| `server.http.read-timeout` | The server read timeout (in seconds). Default value is `10`.|
| `server.http.write-timeout` | The server write timeout (in seconds). Default value is `10`.|
| `storage.type` | The storage type that must be used. Accepted values are (case insensitive): `local`, `mongo`. Default value is `local`. |
| `storage.local.name` | The location where the local storage must be created and used from. Default value is `local-storage/boltdb`. |
| `storage.mongo.uri` | The mongoDB URI. *No default value is provided*. |
| `storage.mongo.name` | The database name. *No default value is provided*. |

**Please see `config/config.default.yml` for a full sample and depiction of configuration settings.**

### Example
```json
logger: 
  application-log-console: true

server:
  http:    
    port: 8081    

storage:
  mongo:
    uri: "mongodb://localhost:27017"
    name: "testdb"    
```
### ENV variables
The configuration loading does not use the Viper ENV variables handling but a new approach. Instead of using the name of the property to determine the ENV variable key to load, the property value is used.

**Approach in Viper**
Let's consider the following simple YAML configuration file:
```json
server:  
  http:
    port: 8081
```
By enabling automatic ENV handling (`viper.AutomaticEnv()`), the name of the used variable is `$SERVER.HTTP.PORT`.

**Local approach**
Even if the above depicted approach works as expected and it can be customized according to different needs (e.g. adding name prefixes, binding properties or key transformation modules), the exact location from which a property is loaded is not very transparent. Still, for smaller environments and configuration files this might not prove itself to be an issue. But for more complex configuration it can become very difficult to determine which of the properties are loaded from the environment and which from the files.
To address this possible maintainability issue, **Basic Go REST Api** needs the desired ENV variables keys to be specified as the property string value, as it follows:
```json
server:  
  http:
    port: "$APP_SERVER_PORT"
```
In this way, by following the configuration file, one can have a better view over what properties are loaded from where, as the actual loaded variable can be seen in the configuration file.

## Usecase (Feature) Components

* **Controller** - exposes the endpoints and that must be its only job. This means that each endpoint must perform an action with 3 steps:
  - unmarshall any eventual input (e.g. query parameters, request body);
  - call i.e. delegate the action to the business component (in this case the service);
  - marshall the response, either error either success.

   :warning:_No business logic must be performed at this level, including but not limited to: input business logic validation, business logic, request or compose additional data for the response marshall action. Of course this is not a functional requirement but for code maintainability these rules must be respected._

* **Service** - implements (or delegates) the business logic with regard to certain models. This is where all the business logic must start from.

* **Repository** - implements the gateway towards a certain storage (e.g. mongoDB connector). A repository follows the Adapter Design Patter and it must implement basic operation without taking to much of the responsibility. For example, a repository must not force certain uniqueness on an entity properties if that must be imposed at service level.

## Development

### Add a new feature

In order to implement a new feature (usecase) usually the following steps must be achieved:

1. Implement a new service (e.g `property/service/property_service.go`);
2. Register the service and service creation (e.g `cmd/api/main.go`);
3. Implement a new repository (e.g `property/gateway/storage/mongo/mongo_repository.go`);
4. Register the repo and its creation (e.g `cmd/api/main.go`);
5. Implement controller (e.g `property/gateway/http/controller.go`);
6. Register the service and service creation (e.g `cmd/api/main.go`).

Even if the project provides a template for feature folder layout, the developer can decide what is the best setup for a particular case. Of course, if the decision is that no services or repositories are required to be implemented, only a Controller must be retrieved.

## Contributing
Pull requests are very welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.