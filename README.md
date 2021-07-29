# Go Calculator
A simple web application in Go which accepts math problems via the URL and returns the response in JSON.

---

By default, HTTP server runs on port 8080.

There are four endpoints, one for each basic math operation:
- `add`
- `subtract`
- `multiply`
- `divide`

All endpoints accept GET requests with two query string parameters: `x` and `y`.

Results are cached for 1 minute.

## How to run
`go run main.go`

## Run tests
`go test`

## Request

| Operation      | Example |
| ----------- | ----------- |
| add      | http://localhost:8080/add?x=2&y=2       |
| subtract   | http://localhost:8080/subtract?x=7&y=4 |
| multiply   | http://localhost:8080/multiply?x=6&y=5 |
| divide   | http://localhost:8080/divide?x=9&y=3 |

## Response
The HTTP server returns JSON response containing following fields:
- `action` - requested math operation 
- `x` - first number for math operation
- `y` - second number for math operation
- `answer` - result of math operation
- `cached` - is cached result returned

Example:
```json
{
    "action": "divide",
    "answer": 3,
    "x": 9,
    "y": 3,
    "cached": true
}
```

# Test with curl
Request:
```
curl "http://localhost:8080/add?x=4&y=3"
```
Response:
```json
{"action":"add","answer":7,"x":4,"y":3,"cached":true}
```

# Dockerize web application
```
docker build -t calcserver .
docker run -p 8080:8080 calcserver
```

