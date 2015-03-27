package httprs

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

// A HttpReadSeeker reads from a http.Response.Body. It can Seek
// by doing range queries.
type HttpReadSeeker struct {
	c       *http.Client
	req     *http.Request
	res     *http.Response
	r       io.ReadCloser
	pos     int64
	canSeek bool
}

var _ io.ReadCloser = (*HttpReadSeeker)(nil)
var _ io.Seeker = (*HttpReadSeeker)(nil)

var ErrNoContentLength = errors.New("Content-Length was not set")
var ErrRangeRequestsNotSupported = errors.New("Range requests are not supported by the remote server")

// Builds a HttpReadSeeker, using the http.Response and, optionaly, the http.Client
// that was used for the query. If no http.Client is passed, http.DefaultClient will
// be used for range queries.
func NewHttpReadSeeker(res *http.Response, client ...*http.Client) *HttpReadSeeker {
	r := &HttpReadSeeker{
		req:     res.Request,
		res:     res,
		r:       res.Body,
		canSeek: (res.Header.Get("Accept-Ranges") == "bytes"),
	}
	if len(client) > 0 {
		r.c = client[0]
	} else {
		r.c = http.DefaultClient
	}
	return r
}

func (r *HttpReadSeeker) Read(p []byte) (n int, err error) {
	if r.r == nil {
		err = r.rangeRequest()
	}
	if err == nil {
		n, err = r.r.Read(p)
		r.pos += int64(n)
	}
	return
}

func (r *HttpReadSeeker) Close() error {
	if r.r != nil {
		return r.r.Close()
	}
	return nil
}

func (r *HttpReadSeeker) Seek(offset int64, whence int) (int64, error) {
	if !r.canSeek {
		return 0, ErrRangeRequestsNotSupported
	}
	var err error
	switch whence {
	case 0:
	case 1:
		offset += r.pos
	case 2:
		if r.res.ContentLength <= 0 {
			return 0, ErrNoContentLength
		}
		offset = r.res.ContentLength - offset
	}
	if r.r != nil && r.pos != offset {
		err = r.r.Close()
		r.r = nil
	}
	r.pos = offset
	return r.pos, err
}

func (r *HttpReadSeeker) rangeRequest() error {
	r.req.Header.Set("Range", fmt.Sprintf("bytes=%d-", r.pos))
	res, err := r.c.Do(r.req)
	if err != nil {
		return err
	}
	r.r = res.Body
	return nil
}
