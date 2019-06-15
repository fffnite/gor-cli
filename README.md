# Go OneRoster CLI tool

Used for managing credentials for the [go-oneroster](https://github.com/fffnite/go-oneroster) api project.

## Arguments

`-new` Creates a new set of credentials  
`-tag` Used with `-new`, attaches human readable tag to the credential  
`-list` lists all credentials held  
`-remove` removes a credential from the store based off client ID

## Build

```
make build
```

### Usage

```
gorcli -new -tag "John Smith"
```
