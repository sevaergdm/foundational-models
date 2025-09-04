# Foundational Models
One of the common issues that I run into is a lack of a common vocabulary to describe an organization that is understandable at a technical and business level. I have always found the concept of domain entities to be something fairly easy to grasp, and generally agreed upon as something "we should have".

As data is often an afterthought when it comes to product development, I also think of this as way to start on a better foot with an entity first approach. 

To help facilitate this, I wanted to create a repository to store these "foundational models" as I am calling them where one can do basic CRUD and retrieval operations, as well as visualize the models and their relationships.

At this iteration, the setup is very naive; storing entities as json files within the repo itself and just performing file operations via local endpoints.

In the future, I would like to generate a visualization which could either be something static like an ERD or possibly utilize a graph db so that the models would become queryable.

## Setup

Everything is self contained in the repo, so it can be cloned and run directly, and then responds to basic api requests.

## Running the service

The server is started by a basic `go run .` in the root directory

The repo contains some of the testing files I was working with in the "entities" directory.

### Example GET 

Retrieve all entities in the "entities" directory:
```
curl http://localhost:8080/api/entities
```

Retrieve a specific in the "entities" directory:
```
curl http://localhost:8080/api/entities/{entityName}
```
`entityName` refers to the value in the `name` field and not the file name

### Example Validate

Validate an entity schema to ensure that it follows the expected canonical schema
```
curl http://localhost:8080/api/validate -d <JSONSCHEMA>
```

### Example Create

Create a new foundational model. This will first validate the schema and then create a new file using the `name` field in the schema

```
curl http://localhost:8080/api/entities -d <JSONSCHEMA>
```

### Example Update

Update an existing schema. Checks are done to see if there are no changes, and it is expected the the `version` will be bumped otherwise the update will fail. Upon success a report of the changes is provided.

```
curl -X PUT -H "Content-Type: application/json" -d <JSONSCHEMA> http://localhost:8080/api/entities/{entityName}
```
