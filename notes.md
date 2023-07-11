 ```bash
 curl "https://api.trello.com/1/members/me/boards?filter=open&fields=id,name,lists&lists=open&list_fields=id,name&cards=open&card_fields=id,name&key=${KEY}&token=${TOKEN}" > /tmp/board-lists.json ; jq < /tmp/board-lists.json
 ```
- We get some functionality for free from the Trello API with filtering and requesting nested fields
