# urllinter

urllinter is used to verify and validate URLs or links.This linter gathers all URLs(based on configuration), check them by calling http requests.

## The problem

Each and every URL or link in this repo must be in working state.It is practically difficult to check each of the URLs manually and ensure they are working.

## Solution

urllinter pasrses all files (path provided in configuration) and stores them.Iteratively it checks whether the URL is working or not.All failure cases are marked as Fail and success cases are marked as Pass.

## Configuration

urllinter supports few configurations.The below is the default configuration file.

```yaml
---
includeExts:
- ".yaml"
- ".sh"
- ".yml"
- ".md"
excludeLinks:
- "http://localhost"
- "https://localhost"
- "http://127.0.0.1"
- "https://127.0.0.1"
- "http://0.0.0.0"
- "https://0.0.0.0"
- "https://vault.example.com:8200"
- "$"
- "<"
- ">"
- "@"
excludeFiles:
- "cli/"
- .git/
- docs/
- ".gitignore"
- ".github/"
acceptStatusCodes:
- 200
- 201
- 302
- 401
- 403
```

## How to run urllinter

- ```go run main.go --path <any valid path or it takes only current working directory> -- config <provide config path or default path will be taken> --summary=true --details=fail```

- To manually run urllinter, user must download(clone) the source code.

- Navigate to the root directory and run the following make command ```make urllint```
