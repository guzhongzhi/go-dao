module github.com/guzhongzhi/gmicro/test

        go 1.15

        require (
        github.com/gogo/protobuf v1.2.1
        github.com/google/wire v0.5.0
        github.com/grpc-ecosystem/grpc-gateway/v2 v2.4.0
        github.com/guzhongzhi/gmicro v0.0.0-20210508092301-a205532d8ed5
        github.com/srikrsna/protoc-gen-gotag v0.5.0
        github.com/urfave/cli/v2 v2.3.0
        google.golang.org/genproto v0.0.0-20210426193834-eac7f76ac494
        google.golang.org/grpc v1.37.0
        google.golang.org/protobuf v1.26.0
        )

        replace github.com/guzhongzhi/gmicro => ../../