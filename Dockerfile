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

WORKDIR /go/src/github.com/tcarreira/superhero

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o /superhero


#    ______ _             _ 
#   |  ____(_)           | |
#   | |__   _ _ __   __ _| |
#   |  __| | | '_ \ / _` | |
#   | |    | | | | | (_| | |
#   |_|    |_|_| |_|\__,_|_|
#                           
#                           
FROM alpine:3.11 as final

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