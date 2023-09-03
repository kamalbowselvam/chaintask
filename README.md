# Chain Task



## Introduction 

This project aims to ease construction project by managing each tasks and providing end user with useful infos.

## Roadmap

[Trello](https://trello.com/b/DyvI6iCi/chain-task)

## Technical Documentation

### Entity Relationship Diagram

```mermaid
erDiagram
    PROJECT ||--o{ TASK : contains
    USER ||--o{ PROJECT : creates
    USER ||--|| ROLE : has
    TASK ||--o{ DOCUMENT : has
    TASK ||--o{ IMAGE: has
    TASK ||--o{ FUNDING: has
    FUNDING ||--o{ SOURCE: has
    FUNDING{
        int Id
        int TaskId
        int SourceId
        float amount
    }
    SOURCE {
        int id
        string SourceName
        float sourceAmout
    }
    PROJECT {
        int Id
        string projectName
        date CreatedOn
        string CreatedBy
        positiion Location
        string Address
        string Responsible
        string Client
        float CompletionPercentage
        float Budget
    }
    TASK {
        int Id
        string Taskname
        date CreatedOn
        string CreatedBy
        date UpdatedOn
        string UpdatedBy
        boolean Done
        int TaskOrder
        int ProjectId
        string deliveryAddress
    }
    USER {
        string UserName
        string Password
        string FullName
        string Email
        date CreatedAt
        date PasswordChangedAt
    }
    ROLE {
        int Id
        string Role
    }
    DOCUMENT {
        int Id
        string Filepath
        string CreatedBy 
        date CreatedAt 
    }
    IMAGE {
        int Id
        string Filepath
        string CreatedBy 
        date CreatedAt 
    }
    POLICY {
        int Id
        string Policy
    }
```

### UML Diagram 

<!--[Click to Open Interactive Diagram](./chaintask.plantuml)-->


### Workflows

#### Authentification and Authorization
```mermaid
sequenceDiagram
    Client->>+chain_task: /projects
    chain_task->>+http_handler: taking incoming request
    chain_task->>+authentification_middleware:token
    authentification_middleware->>+authorization_middleware:role and policy verification
    authorization_middleware->>+service_layer:analyzing the incoming request and validates it
    service_layer->>+persistence_layer: saving data into a DB
    persistence_layer->>-service_layer: returning db object
    service_layer->>-http_handler: business loginc returning business object
    http_handler->>-chain_task: json representing the response 
    chain_task->>-Client: created project
```

#### Policies

|RESOURCES| ADMIN  | RESPONSIBLES  |  CLIENT |
|---|---|---|---|
| LOGIN | (x)  | (x)  | (x)  |
| USERS |  CRUD on every rows | RU on him-self  |  RU on him-self   |
| PROJECTS | CRUD on every rows | RU on assigned projects | RU on assigned project |
| TASKS  | CRUD on every rows | CRUD on assigned projects | CRUD on assigned project |
| DOCUMENTS | CRUD on every rows | CRUD on assigned projects | CRUD on assigned project |
| IMAGES | CRUD on every rows | CRUD on assigned projects | CRUD on assigned project |
| FUNDINGS | no rights | no rights | CRUD on assigned project |


## Deployment

Chaintask is deployed in a [EKS cluster](https://kamalbowselvam.awsapps.com/start/) on Amazon. 

Deployment is made on each push on the  `main` branch.

[Here](http://a14b4fc8215394893b5360715edc21b1-00313d1ee7230a45.elb.eu-west-3.amazonaws.com/) is the publicly avaible URL for the deployed API. You can use [this tutoriel](https://docs.aws.amazon.com/eks/latest/userguide/create-kubeconfig.html) to set up kubectl in the CloudShell Console.