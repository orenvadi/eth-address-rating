Create a REST API service that will display the top 5 most active addresses on the Ethereum network for the last 100 blocks.
addresses on the Ethereum network for the last 100 blocks. The activity metric
Increases if the wallet sends or receives any ERC20 token.

Example:
Address 0xA sends 100 USDC
Address 0xA accepts 1 USDT
then, Activity 0xA = 2

Requirements:
- Endpoint /top.
- Using ethclient.Client
- Response Format:


```
[
    {"address":"0xA", "score": 10},
    {"address":"0xB", "score": 5},
    {"address":"0xC", "score": 4},
    {"address":"0xD", "score": 3},
    {"address":"0xE", "score": 2}
]
```
