include rscli.mk

build-ui: .build-web-api .build-ui

.build-web-api:
	cd pkg/web/@vervstack/matreshka && bun && bun build

.build-ui:
	cd pkg/web/@vervstack/matreshka && yarn && yarn build
	cd pkg/web/Matreshka-UI && yarn && yarn build
	cp -r pkg/web/Matreshka-UI/dist internal/transport/web
