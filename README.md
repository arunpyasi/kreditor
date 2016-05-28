<img align="right" src="https://kreditor.nl/assets/img/kreditor.svg">
# Kreditor

Kreditor is a lightweight web application and API for keeping track of debts.

### Usage

```sh
export KREDITOR_DATABASE_URI="username:password@tcp(localhost:3306)/databasename"
export KREDITOR_SECRET="InsertSecretStringHere"
export KREDITOR_LISTEN_ADDRESS="127.0.0.1:8080"
export KREDITOR_DEBUG="yes" #Optional variable for debugging. 

./kreditor
```

### Building

```sh
$ go get https://github.com/mdeheij/kreditor.git
$ cd kreditor
$ go build
```
### Development

Want to contribute? Great!

### Docker
Kreditor is very easy to install and deploy in a Docker container. Documentation is coming soon.
