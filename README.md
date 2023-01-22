# poc-aws-dynamodb

Some sample Golang code for adding items to DynamoDB table, querying it, and removing items from it. 


## Description

Some sample code from an ongoing investigation into using DynamoDB to store some data for a new service we've been building at work. While I've used many different kinds of databases before, and have been involved in development of database technologies earlier in my career, I've never used DynamoDB on AWS before. These Golang code snippets are the result of my investigation and learning.

A few details about the project I'm involved with, as it might help you to more easily determine the suitability of this PoC for your situation, and why I tried the PoC examples I did:

* For the application we're building, the client only needs to write to and read from database located in the same AWS region, and strongly consistent read is good enough since the use cases for retrieving persisted data aren't mission-critical by any mean. In fact, eventually consistent should be good enough based on the use cases we know so far.

* Records (items) once written to the table never have to be updated. Records will only be retrieved or eventually deleted for money saving.

* Each request to the application might result in 1 or more records being written. Each request will be associated with an unique, GUID-like identifier, and all the records will be associated with the same identifier. Majority of the requests will result in only 1 record being written while a very small percentage of requests will likely result in a large number of records to be written.

* For retrieval of past requests, all records associated with each request (identified by its unique identifier) will have to be retrieved at once. But querying of past requests should be relatively rare compared to writing of new request records. So it might be beneficial to keep all records associated with the same identifier together, in the same shard, etc.

* Some customers might use the application a lot and generate a lot of requests and records, compared to other customers, so sharding strategy should not be based on customer ID.

I've included a Terraform script that I used to configure the DynamoDB database table, since I decided it wouldn't be necessary or practical to configure the database using at API for development and production situations at work anyway.


## Getting Started

### Dependencies

The PoC code was developed on macOS but should work on all operating systems with the following properly installed and configured:

* [Golang](https://go.dev/doc/install)
* [Terraform](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli)

### Downloading Code

```
gh repo clone chchench/poc-aws-dynamodb
```

### Executing program

Run terraform script to configure the DynamoDB with the table required for the Golang program, and then the Go code can be compiled and executed.

*Step for configuring DynamoDB table*

```
terraform init; terraform apply
```

*Step for running the PoC code*

```
go run main.go
```

## Version History

* 0.1
    * Initial Release

## License

This project is licensed under [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

See the [LICENSE](https://github.com/chchench/poc-aws-dynamodb/blob/main/LICENSE) file for details
