package proxy

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ListenAddr     string   `envconfig:"LISTEN_ADDR" default:":14000"`
	AllowedOrigins []string `envconfig:"ALLOWED_ORIGINS" default:""`
	UpstreamAddr   string   `envconfig:"UPSTREAM_ADDR" default:""`
}

func Load() (*Config, error) {

	c := &Config{}
	err := envconfig.Process("", c)
	return c, err
}

type Proxy struct {
	cfg    *Config
	server *http.Server
}

func New(cfg *Config) *Proxy {
	return &Proxy{
		cfg: cfg,
	}
}

func (p *Proxy) Serve(ctx context.Context) error {
	u, err := url.Parse(p.cfg.UpstreamAddr)
	if err != nil {
		return fmt.Errorf("failed to parse upstream address: %w", err)
	}

	rp := httputil.NewSingleHostReverseProxy(u)

	r := mux.NewRouter()
	r.Use(handlers.CompressHandler)
	if len(p.cfg.AllowedOrigins) > 0 {
		r.Use(p.cors())

		rp.ModifyResponse = func(resp *http.Response) error {
			resp.Header.Del(corsAllowOriginHeader)
			return nil
		}
	}

	r.PathPrefix("/").Handler(rp)

	p.server = &http.Server{
		Addr:    p.cfg.ListenAddr,
		Handler: r,
	}

	go func() {
		<-ctx.Done()
		p.server.Close()
	}()

	fmt.Printf("Listening on:   %s\n", p.cfg.ListenAddr)
	fmt.Printf("Proxying to:    %s\n", u)

	return p.server.ListenAndServe()
}

const (
	corsAllowOriginHeader = "Access-Control-Allow-Origin"
)

func (p *Proxy) cors() func(http.Handler) http.Handler {
	fmt.Println("setting up CORS allowed origins", p.cfg.AllowedOrigins)
	origins := strings.Join(p.cfg.AllowedOrigins, ",")
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			w.Header().Set(corsAllowOriginHeader, origins)

			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

			w.Header().Set("Access-Control-Expose-Headers", "Authorization")
			w.Header().Set("Access-Control-Request-Headers", "Authorization")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
