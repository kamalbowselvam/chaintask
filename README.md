# Chain Task



## Introduction 

This project aims to ease construction project by managing each tasks and providing end user with useful infos.

## Technical Documentation

### Entity Relationship Diagram

```mermaid
erDiagram
    PROJECT ||--o{ TASK : contains
    USER ||--o{ PROJECT : creates
    USER ||--|| ROLE : has
    TASK ||--o{ DOCUMENT : has
    TASK ||--o{ IMAGE: has
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
        int ID
        string Filepath
        string CreatedBy 
        date CreatedAt 
    }
```

### UML Diagram 

![UML Diagram](www.plantuml.com/plantuml/png/HST1gi8m403Gg_n_WKQGL114segsYYXKkd0tJOC9cKnACh7N5zpuTk_Vai8ywHpnpp3FQIj4XALuMJPvp4b75OWrSQ625muyu1YMfF4DNYW3bXYI4rDGlnKp0d7skEVWDErERLyzLjLNqzPoKnw7wxhhPDF9-t2eRsNvR7jvELSV)

(Url generated thanks to [this SO](https://stackoverflow.com/questions/32203610/how-to-integrate-uml-diagrams-into-gitlab-or-github/32771815#32771815))