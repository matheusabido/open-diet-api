FROM golang:1.25-trixie AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server \
 && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/migrate ./cmd/migrate


FROM debian:trixie-slim AS whisper

ARG WHISPER_VERSION=v1.9.4
ARG WHISPER_MODEL=small-q8_0

RUN apt-get update \
 && apt-get install -y --no-install-recommends git build-essential cmake curl ca-certificates \
 && rm -rf /var/lib/apt/lists/*

RUN curl -fL -o /out/model.bin --create-dirs \
    https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-${WHISPER_MODEL}.bin

RUN git clone --depth 1 --branch ${WHISPER_VERSION} https://github.com/ggml-org/whisper.cpp /whisper \
 && cmake -S /whisper -B /whisper/build -DCMAKE_BUILD_TYPE=Release -DBUILD_SHARED_LIBS=OFF \
 && cmake --build /whisper/build -j --target whisper-cli \
 && cp /whisper/build/bin/whisper-cli /out/whisper-cli


FROM debian:trixie-slim

RUN apt-get update \
 && apt-get install -y --no-install-recommends ffmpeg libgomp1 ca-certificates \
 && rm -rf /var/lib/apt/lists/*

COPY --from=whisper /out/whisper-cli /usr/local/bin/whisper-cli
COPY --from=whisper /out/model.bin /opt/whisper/model.bin
COPY --from=build /out/server /out/migrate /app/

ENV WHISPER_CLI_PATH=/usr/local/bin/whisper-cli \
    WHISPER_MODEL_PATH=/opt/whisper/model.bin

EXPOSE 3000

CMD ["/app/server"]
