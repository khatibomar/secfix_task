prerequiests:
- osquery

# Installation

I am running this on macos, so first install [queryos](https://osquery.io/downloads)

```sh
sh run_osquery_temp.sh
```

if all good then you will see the following output

```
✅ Ready! Use this socket path in your Go app:
/tmp/osquery.omarelkhatib.29085.em
```
> the tmp socket name differs each time

after that run go app with specifying temp socket

```sh
go run main.go --socket-path=/tmp/osquery.omarelkhatib.29085.em
```
