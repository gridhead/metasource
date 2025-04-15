# Build container 
FROM registry.fedoraproject.org/fedora-minimal:41 AS builder

RUN dnf -y update && \
    dnf -y install wget createrepo_c-devel gcc tar gzip golang && \
    dnf clean all

WORKDIR /app

ENV CGO_ENABLED=1 \
    GOOS=linux

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o meta

# Runtime container
FROM registry.fedoraproject.org/fedora-minimal:41

RUN dnf -y update && \
    dnf -y install createrepo_c-devel && \
    dnf clean all

RUN mkdir -p /app/metasource_db

VOLUME ["/app/metasource_db"]

WORKDIR /app

COPY --from=builder /app/meta /app/meta

EXPOSE 8080
