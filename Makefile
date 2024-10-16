MOCK_OUT=pkg/mock/pkg_gen.go
MOCK_SRC=./pkg/interfaces
MOCK_INTERFACES=CloudStorageClient BigQueryClient PubSubClient PolicyClient CacheService

all: mock

mock: $(MOCK_OUT)

$(MOCK_OUT): $(MOCK_SRC)/*
	go run github.com/matryer/moq@latest -pkg mock -out $(MOCK_OUT) $(MOCK_SRC) $(MOCK_INTERFACES)
