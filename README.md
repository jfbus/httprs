# httprs

A ReadSeeker for
http.Response.Body [![CircleCI](https://dl.circleci.com/status-badge/img/gh/jfbus/httprs/tree/master.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/jfbus/httprs/tree/master)

## Usage

```
import "github.com/jfbus/httprs"

resp, err := http.Get(url)

rs := httprs.NewHttpReadSeeker(resp)
defer rs.Close()
io.ReadFull(rs, buf) // reads the first bytes from the response
rs.Seek(1024, 0) // moves the position
io.ReadFull(rs, buf) // does an additional range request and reads the first bytes from the second response
```

if you use a specific http.Client :

```
rs := httprs.NewHttpReadSeeker(resp, client)
```

## Doc

See https://pkg.go.dev/github.com/jfbus/httprs

## LICENSE

MIT - See LICENSE