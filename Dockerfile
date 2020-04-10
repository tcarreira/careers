# Docker multi-stage build

#    ____        _ _     _           
#   |  _ \      (_) |   | |          
#   | |_) |_   _ _| | __| | ___ _ __ 
#   |  _ <| | | | | |/ _` |/ _ \ '__|
#   | |_) | |_| | | | (_| |  __/ |   
#   |____/ \__,_|_|_|\__,_|\___|_|   
#                                    
#                                    
FROM golang:1.14 AS builder

# install swag command
RUN wget "https://github.com/swaggo/swag/releases/download/v1.6.5/swag_1.6.5_$(uname -s)_$(uname -m).tar.gz" \
    && wget "https://github.com/swaggo/swag/releases/download/v1.6.5/checksums.txt" \
    && sha256sum --check checksums.txt --ignore-missing \
    && tar -zxf "swag_1.6.5_$(uname -s)_$(uname -m).tar.gz" -C /usr/local/bin/ swag \
    && echo "swag installed @ /usr/local/bin/swag :" \
    && ls -alnh /usr/local/bin/swag

WORKDIR /go/src/github.com/tcarreira/superhero

# dependencies layer
COPY go.mod .
COPY go.sum .
RUN go mod download

# generate REST API documentation and build app
COPY . .
RUN swag init -g api.go \
    && CGO_ENABLED=0 go build -o /superhero

#    ______ _             _ 
#   |  ____(_)           | |
#   | |__   _ _ __   __ _| |
#   |  __| | | '_ \ / _` | |
#   | |    | | | | | (_| | |
#   |_|    |_|_| |_|\__,_|_|
#                           
#                           
FROM scratch as final

ENV PORT=8080
ENV GIN_MODE=release

# ENV DB_HOST=db
# ENV DB_PORT=
# ENV DB_NAME=
# ENV DB_USER=
# ENV DB_PASS=

COPY --from=builder /superhero /superhero

ENTRYPOINT [ "/superhero" ]
EXPOSE $PORT