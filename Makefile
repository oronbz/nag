PREFIX ?= /usr/local

.PHONY: build run clean install release

LDFLAGS = -ldflags "-X main.version=$$(git describe --tags --always --dirty 2>/dev/null || echo dev)"

build:
	go build $(LDFLAGS) -o nag .

run:
	go run $(LDFLAGS) .

install: build
	install -d $(DESTDIR)$(PREFIX)/bin
	install -m 755 nag $(DESTDIR)$(PREFIX)/bin/nag

clean:
	rm -f nag

release:
	@latest=$$(git tag --sort=-v:refname | head -1 || echo "v0.0.0"); \
	if [ -n "$(VERSION)" ]; then \
		next="$(VERSION)"; \
		case "$$next" in v*) ;; *) next="v$$next" ;; esac; \
	else \
		patch=$$(echo "$$latest" | sed 's/v[0-9]*\.[0-9]*\.//'); \
		minor=$$(echo "$$latest" | sed 's/v[0-9]*\.\([0-9]*\)\..*/\1/'); \
		major=$$(echo "$$latest" | sed 's/v\([0-9]*\)\..*/\1/'); \
		default="v$$major.$$minor.$$((patch + 1))"; \
		printf "Version [$$default]: "; \
		read input; \
		next=$${input:-$$default}; \
		case "$$next" in v*) ;; *) next="v$$next" ;; esac; \
	fi; \
	echo "Releasing $$next..."; \
	git tag "$$next"; \
	git push origin "$$next"; \
	echo "Warming Go proxy cache..."; \
	curl -sf "https://proxy.golang.org/github.com/oronbz/nag/@v/$$next.info" > /dev/null; \
	echo "Released $$next"
