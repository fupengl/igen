CURRENTPATH = $(CURDIR)

# govendor update native package
gu:	
	govendor update igen/lib/...
	govendor update igen/probuf/...

up:
	docker-compose up

ginit:
	rm vendor/vendor.json
	govendor init
	govendor add +external
	govendor add igen/lib/logger/logf
	govendor add github.com/silenceper/gowatch
	govendor add github.com/howeyc/fsnotify
	govendor add github.com/nestgo/log
	govendor add gopkg.in/yaml.v1
	govendor add github.com/itsjamie/gin-cors
	govendor add google.golang.org/grpc/reflection
	govendor add github.com/golang/protobuf/protoc-gen-go/descriptor
	govendor add google.golang.org/grpc/reflection/grpc_reflection_v1alpha
	
dbuild:
	#docker rmi -f igen_msdemo
	docker build --rm -t igen_msdemo:latest .

drun:
	docker run --rm --name=igen_msdemo -p 8081:8081 -p 8181:8181 -v ${CURRENTPATH}:/go/src/igen/msdemo igen_msdemo

dstop:
	docker stop igen_msdemo
	docker rm igen_msdemo

drmnone:
	docker rmi $(docker images | grep "^<none>" | awk "{print $3}")
