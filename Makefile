.PHONY: all
all: build

.PHONY: build
build:
	mkdir -p build && \
		gomobile bind -ldflags '-s -w' -trimpath -target=android/arm64 -tags android -v \
			-o build/libv2ray.aar github.com/2dust/AndroidLibXrayLite

.PHONY: clean
clean:
	rm -rf build
