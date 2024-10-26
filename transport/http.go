package transport

import "github.com/gin-gonic/gin"

const SupportPackageIsVersion1 = true

type Crypto struct {
	AesKey         string `yaml:"aes_key" mapstructure:"aes_key"`
	AesIv          string `yaml:"aes_iv" mapstructure:"aes_iv"`
	PlainHeaderKey string `yaml:"plain_header_key" mapstructure:"plain_header_key"`
	PlainHeaderVal string `yaml:"plain_header_val" mapstructure:"plain_header_val"`
}

type Middleware func(ctx *gin.Context, operation string) error

type Server struct {
	engine      *gin.Engine
	crypto      *Crypto
	groupRoutes map[string]gin.IRoutes
	addr        []string
	middlewares []Middleware
}

func NewServer(engine *gin.Engine, addr []string, middlewares []Middleware) *Server {
	if addr == nil || len(addr) == 0 {
		panic("addr is required")
	}
	return &Server{
		addr:        addr,
		engine:      engine,
		middlewares: middlewares,
	}
}

func (s *Server) SetCrypto(crypto *Crypto) {
	s.crypto = crypto
}

func (s *Server) AddMethod(httpMethod, relativePath string, handlers ...gin.HandlerFunc) {
	s.engine.Group("").Handle(httpMethod, relativePath, handlers...)
}

func (s *Server) GetMiddlewares() []Middleware {
	return s.middlewares
}

func (s *Server) Run() error {
	return s.engine.Run(s.addr...)
}
