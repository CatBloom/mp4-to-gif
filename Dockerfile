FROM golang:1.23.1

WORKDIR /app

RUN go mod init mp4togif
COPY go.mod go.sum ./
RUN go mod tidy

COPY . ./

RUN apt-get update && apt-get install -y \
    ffmpeg \
    && rm -rf /var/lib/apt/lists/*

CMD ["go", "run", "."]