# Code Gen Makefile
user-rpc-gen:
	@make -f apps/user/rpc/gen.mk user-rpc-gen
social-rpc-gen:
	@make -f apps/social/rpc/gen.mk social-rpc-gen
game-rpc-gen:
	@make -f apps/game/rpc/gen.mk game-rpc-gen

# Development Makefile
user-rpc-release-test:
	@make -f deploy/make/user-rpc.mk release-test
game-rpc-release-test:
	@make -f deploy/make/game-rpc.mk release-test

# local test
user-rpc-local-test:
	@make -f deploy/make/user-rpc.mk local-test
social-rpc-local-test:
	@make -f deploy/make/social-rpc.mk local-test
game-rpc-local-test:
	@make -f deploy/make/game-rpc.mk local-test
