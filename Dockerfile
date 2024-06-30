# Build stage
FROM golang:1.22.4-alpine3.20 AS build-stage

# 필요한 패키지 설치
RUN apk add --no-cache gcc musl-dev

# 작업 디렉토리 설정
WORKDIR /app

# go.mod 및 go.sum 파일 복사
COPY go.mod go.sum ./

# 종속성 다운로드
RUN go mod download

# 소스 코드 복사
COPY . .

# 프로젝트 빌드
RUN go build -o webp-example ./cmd/main

# Final stage
FROM alpine:3.17

# 필요한 패키지 설치
RUN apk add --no-cache ca-certificates

# /app 디렉토리 생성
WORKDIR /app

# 빌드한 바이너리 파일을 복사
COPY --from=build-stage /app/webp-example .

# 컨테이너가 실행할 기본 명령어 설정
CMD ["./webp-example"]
