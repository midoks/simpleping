#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin

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


echo $LDFLAGS
build_app(){

	if [ -f $rootPath/tmp/build/simpleping ]; then
		rm -rf $rootPath/tmp/build/simpleping
		rm -rf $rootPath/simpleping
	fi

	if [ -f $rootPath/tmp/build/simpleping.exe ]; then
		rm -rf $rootPath/tmp/build/simpleping.exe
		rm -rf $rootPath/simpleping.exe
	fi

	echo "build_app" $1 $2

	echo "export CGO_ENABLED=1 GOOS=$1 GOARCH=$2"
	echo "cd $rootPath && go build simpleping.go"

	# export CGO_ENABLED=1 GOOS=linux GOARCH=amd64

	if [ $1 != "darwin" ];then
		export CGO_ENABLED=1 GOOS=$1 GOARCH=$2
		export CGO_LDFLAGS="-static"
	fi


	if [ $1 == "linux" ]; then
		export CC=x86_64-linux-musl-gcc
		if [ $2 == "amd64" ]; then
			export CC=x86_64-linux-musl-gcc

		fi

		if [ $2 == "386" ]; then
			export CC=i486-linux-musl-gcc
		fi

		if [ $2 == "arm64" ]; then
			export CC=aarch64-linux-musl-gcc
		fi

		if [ $2 == "arm" ]; then
			export CC=arm-linux-musleabi-gcc
		fi

		cd $rootPath && go build -ldflags "${LDFLAGS}" simpleping.go 
	fi

	if [ $1 == "darwin" ]; then
		echo "cd $rootPath && go build -v -ldflags '${LDFLAGS}'"
		cd $rootPath && go build -v -ldflags "${LDFLAGS}"
		
		cp $rootPath/simpleping $rootPath/tmp/build
	fi
	

	cp -rf $rootPath/scripts $rootPath/tmp/build
	cp -rf $rootPath/LICENSE $rootPath/tmp/build
	cp -rf $rootPath/README.md $rootPath/tmp/build
	cp -rf $rootPath/conf $rootPath/tmp/build

	cd $rootPath/tmp/build && xattr -c * && rm -rf ./*/.DS_Store && rm -rf ./*/*/.DS_Store
	cd $rootPath/tmp/build && rm -rf ./conf/app.conf && rm -rf ./conf/locale


	# zip
	#cd $rootPath/tmp/build && zip -r -q -o ${PACK_NAME}_${VERSION}_$1_$2.zip  ./ && mv ${PACK_NAME}_${VERSION}_$1_$2.zip $rootPath/tmp/package
	# tar.gz
	cd $rootPath/tmp/build && tar -zcvf ${PACK_NAME}_${VERSION}_$1_$2.tar.gz ./ && mv ${PACK_NAME}_${VERSION}_$1_$2.tar.gz $rootPath/tmp/package
	# bz
	#cd $rootPath/tmp/build && tar -jcvf ${PACK_NAME}_${VERSION}_$1_$2.tar.bz2 ./ && mv ${PACK_NAME}_${VERSION}_$1_$2.tar.bz2 $rootPath/tmp/package

}

golist=`go tool dist list`
echo $golist

# build_app linux amd64
# build_app linux 386
# build_app linux arm64
# build_app linux arm
# build_app darwin amd64
build_app darwin arm

