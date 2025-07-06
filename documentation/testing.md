# Testing

- Run all tests:
  ```sh
  go test ./... -v
  ```
- Integration tests use a real HTTP server and test the full flow using `curl`.