FROM golang:1.18-buster

WORKDIR /

COPY ./ddd-johnny /ddd-johnny

EXPOSE 80

# USER nonroot:nonroot

ENTRYPOINT ["/ddd-johnny"]
