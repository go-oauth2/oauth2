package manage

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/oauth2.v2/generates"
	"gopkg.in/oauth2.v2/models"
	"gopkg.in/oauth2.v2/store/client"
	"gopkg.in/oauth2.v2/store/token"
)

func TestManager(t *testing.T) {
	Convey("Manager Test", t, func() {
		manager := NewManager()

		manager.MapClientModel(models.NewClient())
		manager.MapTokenModel(models.NewToken())
		manager.MapAuthorizeGenerate(generates.NewAuthorizeGenerate())
		manager.MapAccessGenerate(generates.NewAccessGenerate())
		manager.MapClientStorage(client.NewTempStore())
		manager.MustTokenStorage(token.NewRedisStore(
			&token.RedisConfig{Addr: "192.168.33.70:6379"},
		))

	})
}
