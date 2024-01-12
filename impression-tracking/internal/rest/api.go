package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/SleepNFire/mediakeys/impression-tracking/internal/data"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type ginService struct {
	Public    *gin.Engine
	Internal  *gin.Engine
	Technical *gin.Engine
}

func newGinService() ginService {
	gs := ginService{
		Public:    gin.New(),
		Internal:  gin.New(),
		Technical: gin.New(),
	}
	gs.Public.Use(gin.Logger())
	gs.Public.Use(gin.Recovery())

	gs.Internal.Use(gin.Logger())
	gs.Internal.Use(gin.Recovery())

	gs.Technical.Use(gin.Logger())
	gs.Technical.Use(gin.Recovery())

	return gs
}

type MicroService struct {
	Public    http.Server
	Internal  http.Server
	Technical http.Server
}

func NewServer(port string, router *gin.Engine) http.Server {
	return http.Server{
		Addr:    port,
		Handler: router,
	}
}

func NewMicroService(routers ginService) *MicroService {
	return &MicroService{
		Public:    NewServer(":8080", routers.Public),
		Internal:  NewServer(":8081", routers.Internal),
		Technical: NewServer(":8082", routers.Technical),
	}
}

func Init(lc fx.Lifecycle, mysql *data.RedisAccessor) (*MicroService, error) {
	routers := newGinService()

	mysql.RegisterEndpoints(routers.Technical)

	ms := NewMicroService(routers)

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go ms.startServer(&ms.Public, "Public", 8080)
				go ms.startServer(&ms.Internal, "Internal", 8081)
				go ms.startServer(&ms.Technical, "Technical", 8082)

				return nil
			},
			OnStop: func(ctx context.Context) error {
				return ms.Stop()
			},
		},
	)

	return ms, nil
}

func (ms *MicroService) startServer(server *http.Server, name string, port int) {
	fmt.Printf("Démarrage du serveur %s sur le port %d\n", name, port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Erreur en démarrant le serveur %s : %s\n", name, err)
	}
}

func (ms *MicroService) stopServer(server *http.Server, name string) {
	fmt.Printf("Arrêt du serveur %s...\n", name)
	if err := server.Shutdown(context.Background()); err != nil {
		fmt.Printf("Erreur lors de l'arrêt du serveur %s : %s\n", name, err)
	}
}

func (ms *MicroService) Stop() error {
	// Arrêt des serveurs
	ms.stopServer(&ms.Public, "Public")
	ms.stopServer(&ms.Internal, "Internal")
	ms.stopServer(&ms.Technical, "Technical")

	return nil
}
