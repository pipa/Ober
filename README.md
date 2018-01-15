# Ober
> GoLang framework

<!-- [![NPM Version][npm-image]][npm-url]
[![Build Status][travis-image]][travis-url]
[![Downloads Stats][npm-downloads]][npm-url] -->

My attempt on creating a Go web server by creating the bits and parts needed.

<!-- ![](header.png) -->

## Installation

```sh
go get github.com/pipa/ober
```

## Usage example

```Go
package main

import (
	"fmt"
	"net/http"

	"github.com/pipa/ober"
)

func main() {
	e := ober.New()
	e.Middleware.Add(test)
  e.Middleware.Add(test) // twice...
	e.Add("/", index)
	e.Logger.Fatal(e.Start(":9999", "docs/certs/selfsigned.crt", "docs/certs/selfsigned.key")) // Only works with TLS
}

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Println("I'll get printed")
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Luis")
}
```

## Release History

* 0.0.1
    * Work in progress

## Meta

Luis Matute – [@luis_matute](https://twitter.com/luis_matute)

Distributed under the MIT license. See ``LICENSE`` for more information.

[https://github.com/pipa](https://github.com/pipa/)

## Contributing

1. Fork it (<https://github.com/yourname/yourproject/fork>)
2. Create your feature branch (`git checkout -b feature/fooBar`)
3. Commit your changes (`git commit -am 'Add some fooBar'`)
4. Push to the branch (`git push origin feature/fooBar`)
5. Create a new Pull Request

<!-- Markdown link & img dfn's -->
[npm-image]: https://img.shields.io/npm/v/datadog-metrics.svg?style=flat-square
[npm-url]: https://npmjs.org/package/datadog-metrics
[npm-downloads]: https://img.shields.io/npm/dm/datadog-metrics.svg?style=flat-square
[travis-image]: https://img.shields.io/travis/dbader/node-datadog-metrics/master.svg?style=flat-square
[travis-url]: https://travis-ci.org/dbader/node-datadog-metrics
[wiki]: https://github.com/yourname/yourproject/wiki
