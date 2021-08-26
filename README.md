<p align="center">
    <img src="https://raw.githubusercontent.com/Spacewalkio/Garment/main/assets/logo.png?v=1.0.0" width="200" />
    <h3 align="center">Garment</h3>
    <p align="center">A Thread Safe Connection Pooling.</p>
    <p align="center">
        <a href="https://github.com/Spacewalkio/Garment/actions/workflows/build.yml">
            <img src="https://github.com/Spacewalkio/Garment/actions/workflows/build.yml/badge.svg">
        </a>
        <a href="https://github.com/spacewalkio/garment/releases">
            <img src="https://img.shields.io/badge/Version-1.0.0-cyan.svg">
        </a>
        <a href="https://goreportcard.com/report/github.com/spacewalkio/garment">
            <img src="https://goreportcard.com/badge/github.com/spacewalkio/garment?v=1.0.0">
        </a>
        <a href="https://godoc.org/github.com/spacewalkio/garment">
            <img src="https://godoc.org/github.com/spacewalkio/garment?status.svg">
        </a>
        <a href="https://github.com/spacewalkio/garment/blob/master/LICENSE">
            <img src="https://img.shields.io/badge/LICENSE-MIT-orange.svg">
        </a>
    </p>
</p>


Garment is a thread safe connection pooling. It can be used to manage and reuse connections whether a database connection or other types. Most of golang packages supports connection pooling and closes idle connections if they exceed a certain number but a few don't which cause errors. Also you might want to preserve the same connection across sub packages without passing it as a parameter. In these cases you can use garment.


## Documentation

#### Usage

Install the package with:

```zsh
$ go get github.com/spacewalkio/garment
```

Here is an example:

```golang
package main

import (
    "errors"
    "fmt"

    "github.com/spacewalkio/garment"
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

See the [Releases section of our GitHub project](https://github.com/spacewalkio/garment/releases) for changelogs for each release version of Garment. It contains summaries of the most noteworthy changes made in each release.


## Bug tracker

If you have any suggestions, bug reports, or annoyances please report them to our issue tracker at https://github.com/spacewalkio/garment/issues


## Security Issues

If you discover a security vulnerability within Garment, please send an email to [hello@clivern.com](mailto:hello@clivern.com)


## Contributing

We are an open source, community-driven project so please feel free to join us. see the [contributing guidelines](CONTRIBUTING.md) for more details.


## License

© 2021, SpaceWalk. Released under [MIT License](https://opensource.org/licenses/mit-license.php).

**Garment** is authored and maintained by [@SpaceWalk](http://github.com/spacewalkio).
