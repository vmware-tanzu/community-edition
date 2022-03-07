clean:
	rm -rf public resources

prepare:
	$(CURDIR)/load-docs.sh

serve:
	hugo server \
		--buildDrafts \
		--buildFuture \
		--disableFastRender

production-build: prepare
	hugo \


preview-build: prepare
	hugo \
		--baseURL $(DEPLOY_PRIME_URL) \
		--buildDrafts \
		--buildFuture \


build:
	hugo \
	--buildDrafts \
	--buildFuture \


