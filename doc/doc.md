### Инфра
* Docker
* Postgres
* Kafka
* pgmigrate

### Data flow
"/message" POST > api.MessageRequest
api.MessageRequest > model.Message
model.Message > INSERT db
model.Message > kafka.Write
kafka.Read > processMessage
message.IsRead = true
model.Message > UPDATE db

### API
* "/message" POST
    * in Json(api.MessageRequest)
    * out 201, Json(model.Message)
* "/query/metrics" GET
    * out Json(model.Metrics)
* "/query/message" GET
    * in Json(api.MessageQueryRequest)
    * out Json(model.Message)
* "/healthcheck" GET
    * out 200