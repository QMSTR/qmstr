# qmstr

Quartermaster is a suite of command line tools and build system plugins that instruments software builds to create
FLOSS compliance documentation and support compliance decision making. It executes as part of a software build process
to generate reports about the analysed product.

# Compile and Install

Install protobuf from your package manager

Install protobuf go compiler plugin
`go get -u github.com/golang/protobuf/protoc-gen-go`

Install the master server
`(go get github.com/QMSTR/qmstr/cmd/qmstr-master || go generate github.com/QMSTR/qmstr/cmd/master; go get github.com/QMSTR/qmstr/cmd/master)`

Install the wrapper
`go get github.com/QMSTR/qmstr/cmd/qmstr-wrapper`

Optional: install the cli
`go get github.com/QMSTR/qmstr/cmd/qmstr-client`

Or if you dare `wget -O - http://github.com/QMSTR/qmstr/raw/master/install.sh | bash`