FROM golang:alpine AS build-stage

WORKDIR /build

RUN apk update && apk add --update make

COPY .. .

#RUN ls -l

RUN go mod download
RUN go mod verify

RUN make build


FROM golang:alpine AS production-stage

#RUN apk --no-cache add ca-certificates

WORKDIR /usr/app

COPY --from=build-stage /build/snd-backend .
COPY --from=build-stage /build/configs/* ./configs/

EXPOSE 5000

CMD [ "./snd-backend", "--prefork" ]
