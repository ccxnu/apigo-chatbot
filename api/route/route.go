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

	// Initialize clients
	httpClient := httpclient.NewHTTPClient(paramCache)

	// Initialize services
	embeddingService := embedding.NewOpenAIEmbeddingService(paramCache, httpClient)

	// Initialize use cases
	paramUseCase := usecase.NewParameterUseCase(paramRepo, paramCache, timeout)
	docUseCase := usecase.NewDocumentUseCase(docRepo, paramCache, timeout)
	chunkUseCase := usecase.NewChunkUseCase(chunkRepo, statsRepo, paramCache, embeddingService, timeout)
	statsUseCase := usecase.NewChunkStatisticsUseCase(statsRepo, paramCache, timeout)

	// Register all routes
	// All routes are now registered via Huma which uses the ServeMux
	// JWT middleware can be added later when needed

	// Parameter routes
	NewParameterRouter(paramUseCase, mux, humaAPI)

	// Knowledge module routes
	NewDocumentRouter(docUseCase, mux, humaAPI)
	NewChunkRouter(chunkUseCase, mux, humaAPI)
	NewChunkStatisticsRouter(statsUseCase, mux, humaAPI)
}
