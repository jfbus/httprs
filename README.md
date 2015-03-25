# httprs
A ReadSeeker for http.Response.Body

[![wercker status](https://app.wercker.com/status/b8ab18faefae7d1f88f9f23d642f0847/s/master "wercker status")](https://app.wercker.com/project/bykey/b8ab18faefae7d1f88f9f23d642f0847)

```
import "github.com/jfbus/httprs"

resp, err := http.Get(stream_url)
if err != nil {
	return err
}
rs := httprs.NewHttpReadSeeker(resp)
defer rs.Close()
io.ReadFull(rs, buf) // reads the first bytes from the response
rs.Seek(1024, 0) // does an additional range request
io.ReadFull(rs, buf) // reads the first bytes from the second response
```
if you use a specific http.Client
```
rs := httprs.NewHttpReadSeeker(resp, client)
```