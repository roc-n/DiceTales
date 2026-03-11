game-rpc-gen:
	goctl rpc protoc apps/game/rpc/game.proto \
		--go_out=apps/game/rpc \
		--go-grpc_out=apps/game/rpc \
		--zrpc_out=apps/game/rpc
