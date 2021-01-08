.PHONY: all
all: clean install test vet lint check-gofmt build

.PHONY: build
build:
	$(shell go env GOPATH)/bin/mutelight build

.PHONY: check-gofmt
check-gofmt:
	scripts/check_gofmt.sh

.PHONY: clean
clean:
	mkdir -p public/
	rm -f -r public/*

.PHONY: compile
compile: install

# Long TTL (in seconds) to set on an object in S3. This is suitable for items
# that we expect to only have to invalidate very rarely like images. Although
# we set it for all assets, those that are expected to change more frequently
# like script or stylesheet files are versioned by a path that can be set at
# build time.
LONG_TTL := 86400

# Short TTL (in seconds) to set on an object in S3. This is suitable for items
# that are expected to change more frequently like any HTML file.
SHORT_TTL := 3600

.PHONY: deploy
deploy: check-target-dir
# Note that AWS_ACCESS_KEY_ID will only be set for builds on the master branch
# because it's stored in GitHub as a secret variable. Secret variables are not
# made available to non-master branches because of the risk of being leaked
# through a script in a rogue pull request.
ifdef AWS_ACCESS_KEY_ID
	aws --version

	@echo "\n=== Syncing HTML files\n"

	# Force text/html for HTML because we're not using an extension.
	#
	# Note that we don't delete because it could result in a race condition in
	# that files that are uploaded with special directives below could be
	# removed even while the S3 bucket is actively in-use.
	aws s3 sync $(TARGET_DIR) s3://$(S3_BUCKET)/ --acl public-read --cache-control max-age=$(SHORT_TTL) --content-type text/html --exclude 'assets*' $(AWS_CLI_FLAGS)

	@echo "\n=== Syncing media assets\n"

	# Then move on to assets and allow S3 to detect content type.
	#
	# Note use of `--size-only` because mtimes may vary as they're not
	# preserved by Git. Any updates to a static asset are likely to change its
	# size though.
	aws s3 sync $(TARGET_DIR)/assets/ s3://$(S3_BUCKET)/assets/ --acl public-read --cache-control max-age=$(LONG_TTL) --follow-symlinks --size-only $(AWS_CLI_FLAGS)

	@echo "\n=== Syncing Atom feeds\n"

	# Upload Atom feed files with their proper content type.
	find $(TARGET_DIR) -name '*.atom' | sed "s|^\$(TARGET_DIR)/||" | xargs -I{} -n1 aws s3 cp $(TARGET_DIR)/{} s3://$(S3_BUCKET)/{} --acl public-read --cache-control max-age=$(SHORT_TTL) --content-type application/xml

	@echo "\n=== Syncing index HTML files\n"

	# This one is a bit tricker to explain, but what we're doing here is
	# uploading directory indexes as files at their directory name. So for
	# example, 'articles/index.html` gets uploaded as `articles`.
	#
	# We do all this work because CloudFront/S3 has trouble with index files.
	# An S3 static site can have index.html set to indexes, but CloudFront only
	# has the notion of a "root object" which is an index at the top level.
	#
	# We do this during deploy instead of during build for two reasons:
	#
	# 1. Some directories need to have an index *and* other files. We must name
	#    index files with `index.html` locally though because a file and
	#    directory cannot share a name.
	# 2. The `index.html` files are useful for emulating a live server locally:
	#    Golang's http.FileServer will respect them as indexes.
	find $(TARGET_DIR) -name index.html | egrep -v '$(TARGET_DIR)/index.html' | sed "s|^$(TARGET_DIR)/||" | xargs -I{} -n1 dirname {} | xargs -I{} -n1 aws s3 cp $(TARGET_DIR)/{}/index.html s3://$(S3_BUCKET)/{} --acl public-read --cache-control max-age=$(SHORT_TTL) --content-type text/html

	@echo "\n=== Syncing redirects.txt (tiny URLs)\n"

	cat public/redirects.txt | xargs -L1 bash -c 'echo "$$0 --> $$1" && echo "Should redirect to: $$1" | aws s3 cp - s3://$(S3_BUCKET)$$0 --acl public-read --metadata "Website-Redirect-Location=$$1"'

	@echo "\n=== Fixing robots.txt content type\n"

	# Give robots.txt (if it exists) a Content-Type of text/plain. Twitter is
	# rabid about this.
	[ -f $(TARGET_DIR)/robots.txt ] && aws s3 cp $(TARGET_DIR)/robots.txt s3://$(S3_BUCKET)/ --acl public-read --cache-control max-age=$(SHORT_TTL) --content-type text/plain $(AWS_CLI_FLAGS) || echo "no robots.txt"

else
	# No AWS access key. Skipping deploy.
endif

.PHONY: install
install:
	go install .

# Usage:
#     make PATHS="/ /archive" invalidate
.PHONY: invalidate
invalidate: check-aws-keys check-cloudfront-id
ifndef PATHS
	$(error PATHS is required)
endif
	aws cloudfront create-invalidation --distribution-id $(CLOUDFRONT_ID) --paths ${PATHS}

# Invalidates CloudFront's entire cache.
.PHONY: invalidate-all
invalidate-all: check-aws-keys check-cloudfront-id
	aws cloudfront create-invalidation --distribution-id $(CLOUDFRONT_ID) --paths /

# Invalidates CloudFront's cached assets.
.PHONY: invalidate-assets
invalidate-assets: check-aws-keys check-cloudfront-id
	aws cloudfront create-invalidation --distribution-id $(CLOUDFRONT_ID) --paths /assets

.PHONY: killall
killall:
	killall mutelight

.PHONY: lint
lint:
	$(shell go env GOPATH)/bin/golint -set_exit_status ./...

.PHONY: loop
loop:
	$(shell go env GOPATH)/bin/mutelight loop

.PHONY: sigusr2
sigusr2:
	killall -SIGUSR2 mutelight

# sigusr2 aliases
.PHONY: reboot
reboot: sigusr2
.PHONY: restart
restart: sigusr2

.PHONY: test
test:
	go test ./...

.PHONY: test-nocache
test-nocache:
	go test -count=1 ./...

.PHONY: vet
vet:
	go vet ./...

#
# Helpers
#

# Requires that variables necessary to make an AWS API call are in the
# environment.
.PHONY: check-aws-keys
check-aws-keys:
ifndef AWS_ACCESS_KEY_ID
	$(error AWS_ACCESS_KEY_ID is required)
endif
ifndef AWS_SECRET_ACCESS_KEY
	$(error AWS_SECRET_ACCESS_KEY is required)
endif

# Requires that variables necessary to update a CloudFront distribution are in
# the environment.
.PHONY: check-cloudfront-id
check-cloudfront-id:
ifndef CLOUDFRONT_ID
	$(error CLOUDFRONT_ID is required)
endif

.PHONY: check-target-dir
check-target-dir:
ifndef TARGET_DIR
	$(error TARGET_DIR is required)
endif
