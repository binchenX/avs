all: build install unittest cmdlineTests
	
build:
	@go build ./cmd/avs
	@go build ./cmd/avi

install:
	@go install ./cmd/avs
	@go install ./cmd/avi
	
unittest:
	@go test ./specconv
	@# cleanup the tmp dir used by the unitest
	@rm .tmpavstest -rf

cmdlineTests: cmdlineTest1

cmdlineTest1:
	export ANDROID_BUILD_TOP="/home/binchen/d/aspen/aspen-o"
	@rm linaro/poplar -rf
	@avs i  --vendor linaro --device poplar
	@mkdir -p linaro/poplar/etc
	@mkdir -p linaro/poplar/audio
	@cp testFixtures/media_codecs.xml linaro/poplar/etc
	@cp testFixtures/audio_policy.conf linaro/poplar/audio/audio_policy.conf
	@avs v --dir linaro/poplar
	@avs u --dir linaro/poplar
