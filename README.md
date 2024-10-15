# 5 Most active Ethereum addresses API


## How to run 

**Install dependencies**

```bash
go mod download
```

**Set ETH_RPC_URL in .env**

i have used one in example


**Run the server**

```
go run main.go
```

## Fetching api

```bash
curl -X GET "http://localhost:8080/top"
```


### NOTE

It takes aproximitally 1 minute to get the response from api, if using Infura free api to get ethereum blocks info, because of rate limit.
