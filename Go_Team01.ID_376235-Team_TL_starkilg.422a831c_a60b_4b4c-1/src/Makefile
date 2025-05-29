
.PHONY: node1 node2 node3 client clean

HOST = 127.0.0.1
RF ?= 2

# Порты для узлов
NODE1_PORT = 50051
NODE2_PORT = 50052
NODE3_PORT = 50053

# Первый узел кластера (без bootstrap)
node1:
	@echo "Запуск первого узла на порту $(NODE1_PORT)"
	@go run server_main/server.go -host $(HOST) -port $(NODE1_PORT) -rf $(RF)

# Второй узел (подключается к первому)
node2:
	@echo "Запуск второго узла на порту $(NODE2_PORT), подключается к $(HOST):$(NODE1_PORT)"
	@go run server_main/server.go -host $(HOST) -port $(NODE2_PORT) -rf $(RF) -bh $(HOST) -bp $(NODE1_PORT)

# Третий узел (подключается к первому)
node3:
	@echo "Запуск третьего узла на порту $(NODE3_PORT), подключается к $(HOST):$(NODE1_PORT)"
	@go run server_main/server.go -host $(HOST) -port $(NODE3_PORT) -rf $(RF) -bh $(HOST) -bp $(NODE1_PORT)

# Клиент (подключается к первому узлу)
client:
	@echo "Запуск клиента, подключается к $(HOST):$(NODE1_PORT)"
	@go run client_main/client.go -H $(HOST) -P $(NODE1_PORT)

