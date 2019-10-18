# Go OneRoster CLI tool

Used for managing credentials for the [go-oneroster](https://github.com/fffnite/go-oneroster) api project.

## Arguments

`--new` Creates a new set of credentials with a description
`--list` lists all credentials held  
`--remove` removes a credential from the store based off client ID
`--mongo-uri` connection to a mongodb instance
`--help` display info on parameters

## Build

```
make build
```

### Usage

```
goors-cli -h
goors-cli
goors-cli -n "svc-microsoft-teams"
goors-cli -l "svc-microsoft-teams"
goors-cli -r "ae4a6ce2-f076-4f01-88a0-c16af53b6a72"
```
