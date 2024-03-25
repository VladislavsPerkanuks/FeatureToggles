# Backend
FROM golang:latest as backend
WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/cmd ./cmd
COPY backend/pkg ./pkg

RUN ls ./cmd -lah
RUN go build -o main ./cmd

EXPOSE 8081
ENV DEMO=true

CMD ["./main"]

# Frontend
FROM node:latest as frontend
WORKDIR /app

COPY frontend/package.json ./
COPY frontend/package-lock.json ./
COPY frontend/angular.json ./
COPY frontend/server.ts ./
COPY frontend/tsconfig.json ./
COPY frontend/tsconfig.app.json ./
COPY frontend/tsconfig.spec.json ./

COPY frontend/src ./src

RUN npm install -g @angular/cli@latest
RUN npm install

EXPOSE 4200

CMD ["ng", "serve", "--host","0.0.0.0"]
