# test-result

## License

[GPLv3](https://www.gnu.org/licenses/quick-guide-gplv3.en.html), see LICENSE file

## Compiling

`go build -o /test-result`

## Building on Docker

`docker build -t test-result .`

## Deploying on Docker (compose)

```yml
services:
    test-result:
        build: .
        container_name: test-result
        restart: always
        ports:
            - 8080:8080
```

## Example

- [test-result.dev](https://test-result.dev)