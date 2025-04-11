# Build container 
FROM registry.fedoraproject.org/fedora-minimal:41 AS builder

RUN dnf -y update && \
    dnf -y install wget createrepo_c-devel gcc tar gzip && \
    dnf clean all

RUN wget https://go.dev/dl/go1.24.1.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.24.1.linux-amd64.tar.gz && \
    rm -f go1.24.1.linux-amd64.tar.gz

WORKDIR /app

ENV CGO_ENABLED=1 \
    GOOS=linux

COPY go.mod go.sum ./
RUN /usr/local/go/bin/go mod download

COPY . .
RUN /usr/local/go/bin/go build -o metasource-cli

# Runtime container
FROM registry.fedoraproject.org/fedora-minimal:41

RUN dnf -y update && \
    dnf -y install createrepo_c-devel && \
    dnf clean all

WORKDIR /app

COPY --from=builder /app/metasource-cli /app/metasource-cli

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]
