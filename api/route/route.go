package route

import (
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"api-chatbot/api/dal"
	"api-chatbot/domain"
	"api-chatbot/internal/embedding"
	"api-chatbot/internal/httpclient"
	"api-chatbot/internal/jwttoken"
	"api-chatbot/repository"
	"api-chatbot/usecase"
)

func Setup(paramCache domain.ParameterCache, timeout time.Duration, db *pgxpool.Pool, mux *http.ServeMux, humaAPI huma.API) {
	// Use shared parameter cache from App initialization

	// Initialize DAL
	dataAccess := dal.NewDAL(db)

	// Initialize repositories
	paramRepo := repository.NewParameterRepository(dataAccess)
	docRepo := repository.NewDocumentRepository(dataAccess)
	chunkRepo := repository.NewChunkRepository(dataAccess)
	statsRepo := repository.NewChunkStatisticsRepository(dataAccess)
	sessionRepo := repository.NewWhatsAppSessionRepository(dataAccess)
	convRepo := repository.NewConversationRepository(dataAccess)
	adminRepo := repository.NewAdminRepository(dataAccess)

	// Initialize clients
	httpClient := httpclient.NewHTTPClient(paramCache)

	// Initialize services
	embeddingService := embedding.NewOpenAIEmbeddingService(paramCache, httpClient)
	tokenService := jwttoken.NewTokenService(paramCache)

	// Initialize use cases
	paramUseCase := usecase.NewParameterUseCase(paramRepo, paramCache, timeout)
	docUseCase := usecase.NewDocumentUseCase(docRepo, paramCache, timeout)
	chunkUseCase := usecase.NewChunkUseCase(chunkRepo, statsRepo, paramCache, embeddingService, timeout)
	statsUseCase := usecase.NewChunkStatisticsUseCase(statsRepo, paramCache, timeout)
	sessionUseCase := usecase.NewWhatsAppSessionUseCase(sessionRepo, paramCache, timeout)
	convUseCase := usecase.NewConversationUseCase(convRepo, paramCache, timeout)
	adminUseCase := usecase.NewAdminUseCase(adminRepo, tokenService, paramCache)

	// Register all routes
	// All routes are now registered via Huma which uses the ServeMux
	// JWT middleware can be added later when needed

	// Parameter routes
	NewParameterRouter(paramUseCase, mux, humaAPI)

	// Knowledge module routes
	NewDocumentRouter(docUseCase, mux, humaAPI)
	NewChunkRouter(chunkUseCase, mux, humaAPI)
	NewChunkStatisticsRouter(statsUseCase, mux, humaAPI)

	// WhatsApp admin routes
	NewWhatsAppAdminRouter(sessionUseCase, convUseCase, mux, humaAPI)

	// Admin authentication routes
	NewAdminAuthRouter(adminUseCase, mux, humaAPI)
}
