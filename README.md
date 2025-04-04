![CI status](https://github.com/khatibomar/secfix_task/actions/workflows/main.yml/badge.svg)

# Prerequisite:
- osquery
- dbmate
- docker
- make
- sqlc

# Usage 

I am running this on macos, so first install [queryos](https://osquery.io/downloads)

```sh
make deamon-run
```

if all good then you will see the following output

```
✅ Ready! Use this socket path in your Go app:
/tmp/osquery.omarelkhatib.29085.em
```
> the tmp socket name differs each time

after that, to take snapshot of `os version`, `os query version` and `installed apps` run this.

```sh
make app
```

time to expose this data as an API to run the api

```sh
make api
```

we can test api using `curl`


```sh
make curl
```
or with small UI app I created 

```sh
make ui
```

last we can shutdown the deamon with 
```sh 
make deamon-stop
```
