# Verve Go Application

This repo contains a Go application that serves as a RESTful service for counting unique requests based on a provided ID. The application is build with Gin framework for HTTP  handling and  uses Zerolog for structured logging.

## Table of Contents
>
>- Architecture Overview
>- Implementation
>- Extensions
>   - Extension 1: Sending Data via HTTP POST
>   - Extension 2: ID Deduplication Across Instances
>   - Extension 3: Sending Data to a Streaming Service
>- Steps to run code

## Architecture Overview
The application is designed as a REST service with the following key components:

- `HTTP Endpoint:` /api/verve/accept which accepts an integer id as a mandatory query parameter and an optional endpoint for further processing.
- `Unique ID Counting:` Counts unique IDs received within a minute using a shared storage mechanism (Redis can be used in extended version).
- `Logging Mechanism:` Sends the count of unique IDs to a logging system (initially to a log file and later extended to Kafka).


## Implementation

### Core Features
`Single Endpoint:` The service has one GET endpoint /api/verve/accept.

`Unique ID Tracking:` The application tracks unique IDs using a Redis instance to ensure uniqueness across multiple requests.

`Logging:` Logs the count of unique requests every minute using Zerolog.


### Code Structure
The main components of the code are:

`HandleVerveAccept:` Handles incoming requests, checks for the unique ID, and responds accordingly.

`LogUniqueRequestsEveryMinute:` Logs the count of unique IDs every minute.

`sendRequestToEndpoint:` Fires an HTTP GET request to a specified endpoint with the count of unique IDs.

`makePostRequest:` Fires an HTTP POST request to a specified endpoint with the count of unique IDs.


## Extensions
###  Extension 1: Sending Data via HTTP POST
`Objective:` Instead of firing an HTTP GET request to the endpoint, the application fires a POST request. The data structure of the content can be freely decided.

`Implementation:`

- The makePostRequest function constructs and sends a POST request to the provided endpoint with the unique request count.

### Extension 2: ID Deduplication Across Instances
`Objective:` Ensure that ID deduplication works even when the service is behind a load balancer and multiple instances receive the same ID simultaneously.

`Implementation:`

- A Redis instance can be used to track unique IDs. The application checks Redis for the existence of an ID before counting it.
- The SetNX command ensures that the ID is only counted once, and a TTL of one minute allows for new requests in subsequent minutes.

### Extension 3: Sending Data to a Streaming Service
`Objective:`  Instead of writing the count of unique received IDs to a log file, send the count to a distributed streaming service of choice (Kafka).

`Implementation:`

- A Kafka producer can be initialized, and the unique request count is sent to a Kafka topic every minute.
- The logUniqueRequestsEveryMinute function constructs a message containing the unique request count and sends it to the Kafka topic.


## Steps to run code
1. Install required Go packages
2. Run the application: `go run main.go`