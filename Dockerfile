FROM cgr.dev/chainguard/go:latest-dev AS builder

# Info
LABEL org.opencontainers.image.authors="r@nice.pink"
LABEL org.opencontainers.image.source="https://github.com/nice-pink/"

WORKDIR /app

# get go module ready
COPY go.mod go.sum ./
RUN go mod download

# copy module code
COPY . .

# build all
RUN ./build

# plain
# FROM cgr.dev/chainguard/glibc-dynamic:latest AS runner

# git 
# FROM cgr.dev/chainguard/git:latest-root-dev AS runner
FROM cgr.dev/chainguard/git:latest AS runner

# Info
LABEL org.opencontainers.image.authors="r@nice.pink"
LABEL org.opencontainers.image.source="https://github.com/nice-pink/"

WORKDIR /app

# copy executable
COPY --from=builder /app/bin/* /app/
ENTRYPOINT [ "/app/helm-updater" ]
