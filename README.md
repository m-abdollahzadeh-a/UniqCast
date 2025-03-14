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

## Node.js File Processing API

