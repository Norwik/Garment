<p align="center">
    <img src="https://raw.githubusercontent.com/Clevenio/Garment/main/assets/logo.png?v=1.0.3" width="200" />
    <h3 align="center">Garment</h3>
    <p align="center">A Thread Safe Connection Pooling.</p>
    <p align="center">
        <a href="https://github.com/Clevenio/Garment/actions/workflows/build.yml">
            <img src="https://github.com/Clevenio/Garment/actions/workflows/build.yml/badge.svg">
        </a>
        <a href="https://github.com/clevenio/garment/releases">
            <img src="https://img.shields.io/badge/Version-1.0.3-cyan.svg">
        </a>
        <a href="https://goreportcard.com/report/github.com/clevenio/garment">
            <img src="https://goreportcard.com/badge/github.com/clevenio/garment?v=1.0.3">
        </a>
        <a href="https://godoc.org/github.com/clevenio/garment">
            <img src="https://godoc.org/github.com/clevenio/garment?status.svg">
        </a>
        <a href="https://github.com/clevenio/garment/blob/master/LICENSE">
            <img src="https://img.shields.io/badge/LICENSE-MIT-orange.svg">
        </a>
    </p>
</p>

Garment retains a single connection pool for different database types used inside you application (MySQL, Redis, Etcd ... etc). Please note that garment won't be needed and not recommended if you already preserve the same connection pool across your application sub packages.


## Documentation

#### Usage

Install the package with:

```zsh
$ go get github.com/clevenio/garment
```

Here is an example:

```golang
package main

import (
    "errors"
    "fmt"

    "github.com/clevenio/garment"
)

type Database struct {
    State string
}

func (d *Database) GetState() string {
    return d.State
}

func (d *Database) Terminate() {
    d.State = "disconnected"
}

func (d *Database) Close() {
    d.State = "disconnected"
}

func (d *Database) Reconnect() {
    d.State = "connected"
}

func (d *Database) Ping() bool {
    if d.State == "connected" {
        return true
    }

    return false
}

func main() {
    pool := garment.NewPool()

    ping := func(con interface{}) error {
        if con.(*Database).Ping() {
            return nil
        }

        return errors.New("DB connection is lost")
    }

    close := func(con interface{}) error {
        con.(*Database).Close()

        return nil
    }

    reconnect := func(con interface{}) error {
        con.(*Database).Reconnect()

        return nil
    }

    pool.Set("db", &Database{State: "connected"}, ping, close, reconnect)

    fmt.Println(pool.Count()) // 1

    fmt.Println(pool.Has("db"))  // true
    fmt.Println(pool.Ping("db")) // <nil>

    fmt.Println(pool.Get("db").(*Database).GetState()) // connected

    pool.Close("db")

    fmt.Println(pool.Get("db").(*Database).GetState()) // disconnected

    pool.Reconnect("db")

    fmt.Println(pool.Get("db").(*Database).GetState()) // connected

    pool.Get("db").(*Database).Terminate()

    fmt.Println(pool.Get("db").(*Database).GetState()) // disconnected

    pool.Remove("db")

    fmt.Println(pool.Count()) // 0
}
```

## Versioning

For transparency into our release cycle and in striving to maintain backward compatibility, Garment is maintained under the [Semantic Versioning guidelines](https://semver.org/) and release process is predictable and business-friendly.

See the [Releases section of our GitHub project](https://github.com/clevenio/garment/releases) for changelogs for each release version of Garment. It contains summaries of the most noteworthy changes made in each release.


## Bug tracker

If you have any suggestions, bug reports, or annoyances please report them to our issue tracker at https://github.com/clevenio/garment/issues


## Security Issues

If you discover a security vulnerability within Garment, please send an email to [hello@clivern.com](mailto:hello@clivern.com)


## Contributing

We are an open source, community-driven project so please feel free to join us. see the [contributing guidelines](CONTRIBUTING.md) for more details.


## License

© 2021, Cleven. Released under [MIT License](https://opensource.org/licenses/mit-license.php).

**Garment** is authored and maintained by [@Cleven](http://github.com/clevenio).
