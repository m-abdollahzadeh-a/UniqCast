# UniqCast
# UniqCast

# Project Style Guide

This project follows the Uber Go Style Guide as a foundation for development. We aim to maintain consistency,
readability, and best practices in our codebase by adhering to the principles outlined in the guide.

For detailed guidelines, please refer to the official Uber Go Style Guide:
📄 [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md
)

## What it do:
This application consists of two main components: a Node.js app and an MP4Processor.
The Node.js app interacts with users via REST APIs, while the MP4Processor handles the actual processing of MP4 files.
The two components communicate via NATS messaging.the nodejs app interact with user via four rest API.

The Node.js app provides four REST APIs for user interaction and it built with Express that provides an API for managing file processing tasks.
The application integrates with NATS for messaging, PostgreSQL for data storage,and includes [Swagger documentation for API endpoints](http://localhost:3000/api-docs/.).
- **/process**: Accepts a filePath (must be a valid path on the Node.js app's filesystem) and sends this path to the MP4Processor via NATS.
- **/list/all**: Lists information about all processed files.
- **/list/detail/:id**: Retrieves details about a specific processed file.
- **/delete/:id**: Deletes information about a specific file

MP4Processor:
Extracts the initial segment of the specified MP4 file and writes it to a file on the shared filesystem.
It then sends a message back through NATS with one of two statuses:

- **Successful**: Indicates the file was processed successfully.
- **Failed**: Indicates an error occurred during processing, along with a related message.

## App Structure
The application is developed using **Node.js** with the **Yarn** package manager.

The **config section** handles database configurations, reading settings from environment variables.

The **DAO (Data Access Object)** is responsible for interacting with the database models to store or retrieve data.

The **model section** includes the **Protocol** model, which represents the message received as a response from the processor.

The **handlers** are divided into three main parts:
- **API Handler**: Manages user requests across four routes and uses postgresHandler functions to respond to users.
- **Postgres Handler**: Utilizes the DAO to read from and write to the PostgreSQL database.
- **NATS Handler**: A utility that facilitates connecting to, subscribing to, and publishing messages over NATs.

```
.
├── config
│   └── database.js
├── dao
│   └── protocol.js
├── Dockerfile
├── handlers
│   ├── apiHandler.js
│   ├── natsHandler.js
│   └── postgresHandler.js
├── index.js
├── models
│   └── Protocol.js
├── package.json
├── swagger.js
└── yarn.lock
```

## Processor Structure
As you can see in the tree bellow the processor has three main section: config, model, processor
```
├── config
│   └── config.go
├── Dockerfile
├── go.mod
├── go.sum
├── main.go
├── main_test.go
├── model
│   ├── box.go
│   └── message.go
├── processor
│   ├── fileutils.go
│   ├── fileutils_test.go
│   └── process.go
```
The configuration section includes code designed to read settings from command-line arguments or environment variables,
utilizing the Viper library for this purpose.

The model section is divided into two parts: the box model, which represents the structure of each box within an MP4 file,
and the message model, which defines the format of the response sent back via NATs to the Node.js application.

The processor component features file utilities that facilitate writing byte arrays to the filesystem.
The core functionality, however, resides in process.go,
which contains a function to extract the initial segment of data and write it to a file using these utilities.

In main.go, the processNatsMessage function leverages the capabilities of process.go to locate the initial segment
and subsequently publishes the outcome through NATs.


## HOW TO RUN
Simply execute docker-compose up --build to get this code up and running! :)

If you don’t have Docker and Docker Compose installed, you can refer to the official [Docker documentation](https://docs.docker.com/engine/install/) for installation instructions.