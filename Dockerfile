FROM golang:1.15.6 AS build

WORKDIR /app

ENV PRODUCTION=1
ENV DATABASE_URL=""

COPY . /app

CMD ["go", "run", "fcpa"]