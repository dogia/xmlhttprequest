package xmlhttprequest

const debug = true
const unsent = 0
const opened = 1
const headersRecived = 2
const loading = 3
const done = 4

var forbiddenRequestHeaders = []string{
	"accept-charset",
	"accept-encoding",
	"access-control-request-headers",
	"access-control-request-method",
	"connection",
	"content-length",
	"content-transfer-encoding",
	"cookie",
	"cookie2",
	"date",
	"expect",
	"host",
	"keep-alive",
	"origin",
	"referer",
	"te",
	"trailer",
	"transfer-encoding",
	"upgrade",
	"via",
}

var forbiddenRequestMethods = []string{
	"trace",
	"track",
	"connect",
}
