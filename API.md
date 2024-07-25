# API

### Check service status
Returns 200
#### Endpoint
`GET /healthcheck`
#### Response
`200 Ok`

### Send message
Adds a message to the db and queues it for processing
#### Endpoint
`POST /message`
#### Request
```
{
    "content": "[message content (max length is 1024 bytes)]",
}
```
#### Response
`201 Accepted`
```
{
    "Id": "[message uuid]",
    "Created": "[creation timestamp]",
    "Content": "[message content]",
    "IsProcessed": false,
}
```

### Query message
Returns message info and processed status by id
#### Endpoint
`POST /query/message`
#### Request
```
{
    "id": "[message uuid]",
}
```
#### Response
`200 Ok`
```
{
    "Id": "[message uuid]",
    "Created": "[creation timestamp]",
    "Content": "[message content]",
    "IsProcessed": [true | false],
}
```

### Query metrics
Returns global info about the messages storage
#### Endpoint
`GET /query/metrics`
#### Response
`200 Ok`
```
{
    "Messages":
    {
        "Total": [total number of messages, processed or not],
        "LastDay": [same as total, for the last 24 hours],
        "LastHour": [same as total, for the last hour],
        "LastMinute": [same as total, for the last minute],
    }
    "ProcessedTotal": [total number of processed messages],
    "ProcessedRatio": [Messages.Total / ProcessedTotal],
    "OrphanMessages": [number of messages that are >1 minute old and are not processed],
}
```
