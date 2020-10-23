package shield

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

type Core struct {
	URL                string `yaml:"url"`
	Session            string `yaml:"session"`
	InsecureSkipVerify bool   `yaml:"skip_verify"`
	CACertificate      string `yaml:"cacert"`
}
type Config struct {
	Path    string
	Current *Core
	SHIELDs map[string]*Core
}

func ReadConfig(path string) (*Config, error) {
	cfg := &Config{
		Path:    path,
		SHIELDs: map[string]*Core{},
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	if err := yaml.Unmarshal(b, &cfg.SHIELDs); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Write() error {
	b, err := yaml.Marshal(c.SHIELDs)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(c.Path, b, 0666)
}

func (c *Config) Select(alias string) error {
	if core, ok := c.SHIELDs[alias]; ok {
		c.Current = core
		return nil
	}
	return fmt.Errorf("unknown SHIELD core '%s'", alias)
}

func (c *Config) Add(alias string, core Core) {
	c.SHIELDs[alias] = &core
}

func (c *Config) Client(core string) (*Client, error) {
	if err := c.Select(core); err != nil {
		return nil, err
	}

	return &Client{
		URL:                c.Current.URL,
		Session:            c.Current.Session,
		InsecureSkipVerify: c.Current.InsecureSkipVerify,
		CACertificate:      c.Current.CACertificate,
		TrustSystemCAs:     true,
	}, nil
}

func EnvConfig() (*Client, error, bool) {
	url := os.Getenv("SHIELD_URL")
	username := os.Getenv("SHIELD_USERNAME")
	password := os.Getenv("SHIELD_PASSWORD")

	if url == "" || username == "" || password == "" {
		return nil, nil, false
	}

	ca := os.Getenv("SHIELD_CA")
	skip := os.Getenv("SHIELD_SKIP_VERIFY") == "yes"
	trust := os.Getenv("SHIELD_TRUST_SYSTEM_CAS") != "no"
	debug := os.Getenv("SHIELD_DEBUG") == "yes"
	trace := os.Getenv("SHIELD_TRACE") == "yes"

	var timeout int
	if s := os.Getenv("SHIELD_TIMEOUT"); s != "" {
		n, err := strconv.ParseInt(s, 10, 0)
		if err != nil {
			return nil, err, true
		}
		timeout = int(n)
	}

	c := &Client{
		URL:                url,
		InsecureSkipVerify: skip,
		CACertificate:      ca,
		TrustSystemCAs:     trust,
		Debug:              debug,
		Trace:              trace,
		Timeout:            timeout,
	}

	err := c.Authenticate(&LocalAuth{
		Username: username,
		Password: password,
	})
	return c, err, true
}
