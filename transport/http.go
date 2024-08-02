package transport

import "github.com/gin-gonic/gin"

const SupportPackageIsVersion1 = true

type Crypto struct {
	Enable         bool   `yaml:"enable"`
	AesKey         string `yaml:"aes_key" mapstructure:"aes_key"`
	AesIv          string `yaml:"aes_iv" mapstructure:"aes_iv"`
	PlainHeaderKey string `yaml:"plain_header_key" mapstructure:"plain_header_key"`
	PlainHeaderVal string `yaml:"plain_header_val" mapstructure:"plain_header_val"`
}

type Server struct {
	engine      *gin.Engine
	crypto      Crypto
	groupRoutes map[string]gin.IRoutes
	addr        []string
}

func NewServer(engine *gin.Engine, crypto Crypto, addr []string) *Server {
	if addr == nil || len(addr) == 0 {
		panic("addr is required")
	}
	return &Server{
		addr:        addr,
		engine:      engine,
		crypto:      crypto,
		groupRoutes: make(map[string]gin.IRoutes),
	}
}

func (s Server) SetGroupRoutes(middlewares map[string][]gin.HandlerFunc) {
	for operation, vals := range middlewares {
		s.groupRoutes[operation] = s.engine.Group("").Use(vals...)
	}
}

func (s Server) AddMethod(httpMethod, relativePath, operation string, handlers ...gin.HandlerFunc) {
	if _, ok := s.groupRoutes[operation]; ok {
		s.groupRoutes[operation].Handle(httpMethod, relativePath, handlers...)
	} else {
		s.groupRoutes[operation] = s.engine.Group("").Handle(httpMethod, relativePath, handlers...)
	}
}

func (s Server) Run() error {
	return s.engine.Run(s.addr...)
}
