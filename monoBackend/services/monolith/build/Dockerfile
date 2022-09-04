FROM golang:alpine AS build-stage

WORKDIR /build

#RUN apk update && apk add --update make

COPY ../../.. .

RUN go mod download
RUN go mod verify

#RUN make build

RUN go build -o snd-backend ./services/monolith/cmd/app/main.go


FROM golang:alpine AS production-stage

WORKDIR /usr/app

COPY --from=build-stage /build/snd-backend .
COPY --from=build-stage /build/configs/* ./configs/

EXPOSE $APP_PORT

#CMD [ "./snd-backend", "--prefork" ]
CMD [ "./snd-backend"]
