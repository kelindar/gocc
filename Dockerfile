FROM golang:1.20 AS builder

ENV CGO_ENABLED=0

# GOPATH => /go
RUN mkdir -p /go/src
COPY . src
RUN cd src/cmd/webserver && go build . && test -x webserver

# Use crossbuild image with preinstalled LLVM toolchains for different architectures
FROM kelindar/crossbuild:latest
LABEL maintainer=roman.atachiants@gmail.com

# Use archived sources (if bullseye is archived)
#RUN echo "deb http://security.debian.org/ bullseye/updates main contrib non-free" > /etc/apt/sources.list
#RUN echo "deb http://archive.debian.org/debian bullseye main contrib non-free" > /etc/apt/sources.list
RUN apt-get update && apt-get install nano

# Install LLVM toolchain
RUN wget -O - https://apt.llvm.org/llvm-snapshot.gpg.key | apt-key add -
RUN wget https://apt.llvm.org/llvm.sh && chmod +x llvm.sh && ./llvm.sh 15

# Copy gocc binary from builder image
COPY --from=builder /go/src/cmd/webserver/webserver .
COPY --from=builder /go/src/example ./example
RUN chmod +x webserver
EXPOSE 8080
CMD ["./webserver"]