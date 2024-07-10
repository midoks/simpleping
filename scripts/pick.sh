#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin:/opt/homebrew/bin

# https://github.com/FiloSottile/homebrew-musl-cross
# brew install FiloSottile/musl-cross/musl-cross --without-x86_64 --with-i486 --with-aarch64 --with-arm

# brew install mingw-w64
# sudo port install mingw-w64

VERSION=1.0
curPath=`pwd`
rootPath=$(dirname "$curPath")

PACK_NAME=simpleping

# go tool dist list
mkdir -p $rootPath/tmp/build
mkdir -p $rootPath/tmp/package
source ~/.bash_profile


echo "abspath:$rootPath"
echo $LDFLAGS
build_app(){

	mkdir -p $rootPath/tmp
	cd $rootPath

	if [ -f $rootPath/tmp/build/simpleping ]; then
		rm -rf $rootPath/tmp/build/simpleping
		rm -rf $rootPath/simpleping
	fi

	echo "build_app" $1 $2
	echo "export CGO_ENABLED=1 GOOS=$1 GOARCH=$2"
	export CGO_ENABLED=1 GOOS=$1 GOARCH=$2
	echo "cd $rootPath && go build main.go"

	
	echo "arch:$2"
	
	# export CGO_ENABLED=1 GOOS=linux GOARCH=amd64
	if [ $1 == "darwin" ]; then
		echo "go build -o ${PACK_NAME} -v -ldflags '${LDFLAGS}'"
		cd $rootPath && go build -o ${PACK_NAME} -v -ldflags "${LDFLAGS}" 
		cp $rootPath/${PACK_NAME} $rootPath/tmp/build
	fi

	if [ $1 == "linux" ]; then
		export CC=x86_64-linux-musl-gcc
		if [ $2 == "amd64" ]; then
			echo "CC=x86_64-linux-musl-gcc"
			export CC=x86_64-linux-musl-gcc
		fi

		if [ $2 == "386" ]; then
			echo "CC=i486-linux-musl-gcc"
			export CC=i486-linux-musl-gcc
		fi

		if [ $2 == "arm64" ]; then
			echo "CC=aarch64-linux-musl-gcc"
			export CC=aarch64-linux-musl-gcc
		fi

		if [ $2 == "arm" ]; then
			echo "CC=arm-linux-musl-gcc"
			export CC=arm-linux-musleabi-gcc
		fi

		export CGO_LDFLAGS="-static"
		echo "go build -o ${PACK_NAME} -v -ldflags '${LDFLAGS}' main.go "
		cd $rootPath && go build -v -ldflags "${LDFLAGS}" -o ${PACK_NAME} main.go
		cp $rootPath/${PACK_NAME} $rootPath/tmp/build
	fi

	# cp -rf $rootPath/conf $rootPath/tmp/build
	cp -rf $rootPath/scripts $rootPath/tmp/build
	cp -rf $rootPath/LICENSE $rootPath/tmp/build
	cp -rf $rootPath/README.md $rootPath/tmp/build

	cd $rootPath/tmp/build && xattr -c * && rm -rf ./*/.DS_Store && rm -rf ./*/*/.DS_Store
	cd $rootPath/tmp/build && tar -zcvf $rootPath/tmp/package/${PACK_NAME}_$1_$2.tar.gz ./
}

golist=`go tool dist list`
echo $golist

# build_app linux amd64
# build_app linux arm64
build_app darwin amd64
# build_app darwin arm64
