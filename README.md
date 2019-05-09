# Go Mixpanel client

A Mixpanel client written in Go

## Warning
This project is under heavy development and not officially released yet.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

This library doesn't use 3rd party library for its core. But 3rd party library is being used for unit test. If you really want to run unit tests in the project.

```
go get -u github.com/spf13/viper
go get -u github.com/stretchr/testify/assert
```

### Installing

Use `go get` to get this library.

```
go get -u github.com/antonyho/mixpanel
```

## Running the tests

```
go test -v ./...
```

## Usage

`import github.com/antonyho/mixpanel`

```
token := "<Mixpanel token from Mixpanel setting page>"
mp := mixpanel.NewClient(token)
...
result, err := mp.Track()
```

## Contributing

Please read [CONTRIBUTING.md](https://gist.github.com/PurpleBooth/b24679402957c63ec426) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/your/project/tags). 

## Authors

* **Antony Ho** - *Initial work* - [antonyho](https://github.com/antonyho)

See also the list of [contributors](https://github.com/your/project/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* Hat tip to anyone whose code was used
* Inspiration
* etc