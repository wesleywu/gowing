buf-generate:
	cd protobuf && buf generate

buf-push: buf-generate
	cd protobuf && buf push