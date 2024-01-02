package tripartite

import (
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

var DefaultClient *Client

const (
	AUTHORIZATION = "9505C098192D4BCECE6C22F77E63BFB2"
	COOKIE        = "Cookie"
	HSOT          = "host"

	AUTH_HEAD_NAME = "Authorization"
)

type Address struct {
	IP   string
	Port string
}

type Client struct {
	sync.Mutex
	*resty.Client
	targetName  string
	timeout     time.Duration
	serviceName string
	hostPort    string
	headers     map[string]string
	protocol    string
}

type Option func(*Client)

// default value is 500ms.
func Timeout(t time.Duration) Option {
	return func(c *Client) {
		c.Client.SetTimeout(t)
	}
}

func Url(url string) Option {
	return func(c *Client) {
		c.Client.SetBaseURL(c.protocol + url)
	}
}

func Protocol(proto string) Option {
	return func(c *Client) {
		c.protocol = proto + "://"
	}
}

func Header(key, value string) Option {
	return func(c *Client) {
		c.SetHeader(key, value)
	}
}

// Custom http headers
//func Headers(headers map[string]string) Option {
//	return func(c *Client) {
//		c.SetHeaders(headers)
//	}
//}

func (c *Client) SetHeader(key string, value string) {
	if c.headers == nil {
		c.headers = make(map[string]string)
	}
	c.headers[key] = value
}

func (c *Client) SetHeaders() *resty.Request {
	if c.headers == nil {
		c.headers = make(map[string]string)
	}
	c.headers["Content-Type"] = "application/x-www-form-urlencoded"
	c.headers["Cookie"] = "PHPSESSID=sh8qn1hps2vj754l990r9f4bo1; phpdisk_zcore_v2_info=c95f4F%2B1gBJoodI1%2FXC5wpgkGxfaRkI11iae9Uq6t7NtlEW0AyoiJ2wbt78CwCqR3ySG1En9wseVvuG0BiXg01GxwykX%2FZW30hDQMv1Aj62QWZxSTgAbrPUFbug2CpU7KxRGHkE2Zw; view_stat=1"
	c.headers["Host"] = "www.xunniuwp.com"
	c.headers["Connection"] = "keep-alive"
	c.headers["Referer"] = "http://www.xunniuwp.com/file-4564078.html"
	c.headers["Origin"] = "http://www.xunniuwp.com"
	c.headers["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	c.headers["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"
	c.headers["Accept-Encoding"] = "gzip, deflate"
	c.headers["X-Requested-With"] = "XMLHttpRequest"
	return c.Client.R().SetHeaders(c.headers)
}

func (c *Client) proto() string {
	return "http://"
}

//
//func (t *Client) Call(ctx context.Context, method string, args, result interface{}) error {
//	// Fetch address
//	var url string
//	var currentAddr *Address
//	if t.hostPort != "" {
//		url = "http://" + t.hostPort
//	} else {
//		var err error
//		currentAddr, err = t.discovery.GetAddress()
//		if err != nil {
//			return fmt.Errorf("GetAddress failed for %s. Error:%s\n", t.targetName, err)
//		}
//		url = "http://" + currentAddr.String()
//	}
//
//	return err
//}

// New create a new tzone client with specified service and options.
func New(opts ...Option) *Client {
	c := &Client{
		timeout: 5000 * time.Millisecond,
		Client:  resty.New(),
	}
	c.protocol = c.proto()
	c.Client.SetTimeout(c.timeout)
	for _, opt := range opts {
		opt(c)
	}

	return c
}
