# Requirements:
- docker
- make
# Running locally:
```
make up
```
- default port is `8085`
# Testing:
```
curl -X POST localhost:8085/api/messages -d '{"email":"jan.kowalski@example.com","title":"Interview","content":"simple text","mailing_id":1, "insertion_time": "2020-04-24T05:42:38.725412916Z"}'
```
```
curl -X POST localhost:8085/api/messages/send -d '{"mailing_id":1}'
```
see http://localhost:8086 to view mailcatcher
# TODO
- finish api docs in /docs/api.yaml
- add deletion endpoint
- add scheduled removals
# Sending emails:
- it's done asynchronously
- background worker consumes jobs:
  - mailing list is locked for the time of sending emails from it, so no new customers can be added at that time
  - sendout is done in concurrent batches