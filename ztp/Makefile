.PHONY: ci-job test-policygen checkSourceCrs test-policygen-kustomize test-siteconfig test-siteconfig-kustomize

ci-job: test-policygen checkSourceCrs test-policygen-kustomize test-siteconfig test-siteconfig-kustomize

test-policygen:
	@echo "ZTP: Build policy generator and run test"
	$(MAKE) -C ./policygenerator all

checkSourceCrs:
	$(MAKE) -C ./source-crs-builder check

test-policygen-kustomize:
	@echo "ZTP: Build policy generator kustomize plugin and run test"
	$(MAKE) -C ./policygenerator-kustomize-plugin test

test-siteconfig:
	@echo "ZTP: Build siteconfig generator and run test"
	$(MAKE) -C ./siteconfig-generator test

test-siteconfig-kustomize:
	@echo "ZTP: Build siteconfig generator kustomize plugin and run test"
	$(MAKE) -C ./siteconfig-generator-kustomize-plugin test
