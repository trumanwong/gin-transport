package transport

import "github.com/gin-gonic/gin"

const SupportPackageIsVersion1 = true

type Crypto struct {
	AesKey         string `yaml:"aes_key" mapstructure:"aes_key"`
	AesIv          string `yaml:"aes_iv" mapstructure:"aes_iv"`
	PlainHeaderKey string `yaml:"plain_header_key" mapstructure:"plain_header_key"`
	PlainHeaderVal string `yaml:"plain_header_val" mapstructure:"plain_header_val"`
}

type GroupMiddleware struct {
	Middleware gin.HandlerFunc
	Operations []string
}

type Server struct {
	engine      *gin.Engine
	crypto      *Crypto
	groupRoutes map[string]gin.IRoutes
	addr        []string
}

func NewServer(engine *gin.Engine, groupMiddlewares []*GroupMiddleware, addr []string) *Server {
	if addr == nil || len(addr) == 0 {
		panic("addr is required")
	}
	m := make(map[string][]gin.HandlerFunc)
	groupRoutes := make(map[string]gin.IRoutes)
	if groupMiddlewares != nil {
		for _, groupMiddleware := range groupMiddlewares {
			for _, operation := range groupMiddleware.Operations {
				if _, ok := m[operation]; !ok {
					m[operation] = []gin.HandlerFunc{groupMiddleware.Middleware}
				} else {
					m[operation] = append(m[operation], groupMiddleware.Middleware)
				}
			}
		}

		for operation, vals := range m {
			groupRoutes[operation] = engine.Group("").Use(vals...)
		}
	}
	return &Server{
		addr:        addr,
		engine:      engine,
		groupRoutes: groupRoutes,
	}
}

func (s *Server) SetCrypto(crypto *Crypto) {
	s.crypto = crypto
}

func (s *Server) AddMethod(httpMethod, relativePath, operation string, handlers ...gin.HandlerFunc) {
	if _, ok := s.groupRoutes[operation]; ok {
		s.groupRoutes[operation].Handle(httpMethod, relativePath, handlers...)
	} else {
		s.groupRoutes[operation] = s.engine.Group("").Handle(httpMethod, relativePath, handlers...)
	}
}

func (s *Server) Run() error {
	return s.engine.Run(s.addr...)
}
