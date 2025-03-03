package transport

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/trumanwong/cryptogo"
	"github.com/trumanwong/cryptogo/paddings"
	"github.com/trumanwong/gin-transport/v2/transport/errors"
)

func (s *Server) ResultError(ctx *gin.Context, err error) {
	se := errors.FromError(err)
	s.response(ctx, int(se.Code), se)
}

func (s *Server) Result(ctx *gin.Context, code int, data interface{}) {
	s.response(ctx, code, data)
}

func (s *Server) response(ctx *gin.Context, code int, data interface{}) {
	if s.crypto == nil || ctx.GetHeader(s.crypto.PlainHeaderKey) == s.crypto.PlainHeaderVal {
		ctx.JSON(code, data)
		return
	}
	plaintext, _ := json.Marshal(data)
	cipher, _ := cryptogo.AesCBCEncrypt(plaintext, []byte(s.crypto.AesKey), []byte(s.crypto.AesIv), paddings.PKCS7)
	ctx.JSON(code, gin.H{
		"Data": base64.StdEncoding.EncodeToString(cipher),
	})
}
