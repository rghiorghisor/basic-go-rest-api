# Basic Go REST Api

This application (boilerplate) is a take on developing a simple Go REST API, backed up by a MongoDb database, using principles of Hexagonal Architecture and SOLID principles.

**Note** that this is not work done, in any way. For an application to be production ready additional features must be implemented. Some will be added to this repository in the future, other will probably remain the developer's future duties.

## Requirements

- [Git](https://git-scm.com/)
- [Go](https://golang.org/dl/)
- [mongoDB](https://www.mongodb.com/)

## Details & Features

- RESTful endpoints (following [RFC 7231](https://tools.ietf.org/html/rfc7231))
- Standard CRUD operations;
- Error handling (including response JSON generation);
- Clean Architecture code organization (use case centric);
- 3tier application with:
  - RESTful API as presentation layer;
  - mongoDB as data layer;
- Configurable through YAML files.

### Implementation details
Some of the implementation details one can analyze or take note from this application:
- How to load configuration YAMLs using VIPER;
- How to specify the ENV variable name in configuration files (loaded with Viper);
- How to organize a use case (feature) in 3 layers;
- How to implement standard RESTful CRUD operations (RFC 7231);
- How to use mongoDB GO driver ([mongodb/mongo-go-driver](https://github.com/mongodb/mongo-go-driver));

## Getting Started

Download modules:
```
go mod download
```

Run application locally:
```
go run ./cmd/api/main.go
```

Run tests:
```
go test ./...
```

## Project layout

```
.
├── cmd                    Main applications of the project;
│   └── api                The server application API (the entry point);
├── config                 Configuration logic and configuration files;
├── errors                 Application errors and error logic;
├── model                  Model (entities) definitions and logic;
├── property               The entire property use case and dependencies;
│   ├── gateway            The gateways implementations;
|   |   ├── http           The HTTP gateways (Controllers);
|   |   └── storage        The storage gateway implementations;
|   |       └── mongo      The mongoDB gateway (Repository);
│   └── service            The property business logic;
├── server                 The server application logic and dependencies;
│   ├── http               The HTTP server implementations;
|   └── storage            The server's storage overall implementation;
└── tests                  Contains additional files for testing purposes;
    └── config             Configuration files used by config loading tests.
```

## Configuration
The application configuration is achieved by the `./config/config.yml` file loaded using the [Viper](https://github.com/spf13/viper) module, with a few customizations.

### Example
```
server:  
  http:
    port: 8081    
    read-timeout: 10    
    write-timeout: 10

storage:
  mongo:
    uri: "mongodb://localhost:27017"
    name: "testdb"
    properties-collection: "properties_collection"
```
### ENV variables
The configuration loading does not use the Viper ENV variables handling but a new approach. Instead of using the name of the property to determine the ENV variable key to load, the property value is used.

**Approach in Viper**
Let's consider the following simple YAML configuration file:
```
server:  
  http:
    port: 8081
```
By enabling automatic ENV handling (`viper.AutomaticEnv()`), the name of the used variable is `$SERVER.HTTP.PORT`.

**Local approach**
Even if the above depicted approach works as expected and it can be customized according to different needs (e.g. adding name prefixes, binding properties or key transformation modules), the exact location from which a property is loaded is not very transparent. Still, for smaller environments and configuration files this might not prove itself to be an issue. But for more complex configuration it can become very difficult to determine which of the properties are loaded from the environment and which from the files.
To address this possible maintainability issue, **Basic Go REST Api** needs the desired ENV variables keys to be specified as the property string value, as it follows:
```
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

## Contributing
Pull requests are very welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.