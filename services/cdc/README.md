# Lambda CDC

This repository contains the code for the lambda dedicated to the CDC from the DynamoDB Streams that will reverse the output to the dedicated topic.

![Architecture](./architecture.png "Overview")

## Project Structure

The project follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout). Test files are written in every directory (es. internal has the both the application go files both the test ones)

- **./internal/** 
    business logic specic to this application
- **./cmd/**
    main applications for this project
- **Taskfile.yaml**
    task file for recurring commands 
    
## Taskfile

In the repository there's the taskfile with commands alias that help us during development. Here's the list of them:

- **build**
         build go binary
- **zip-linux**
        zip go binary on linux systems
- **zip-win**
        zip go binary on windwos systems
- **update-deploy**
        deploy zip on lambda

## Configuration

As we can see, the Lambda receveis from DynamoDB Stream and writes to SNS Topic. The stream is configured when lambda is created, so there's no need to reference the table in the code, while for the SNS Topic we need an enviroment variables referencing to it.

| Key      | Value |
| ----------- | ----------- |
| SNS_TOPIC_ARN      | "SNS Product Topic Name"       |


------------------