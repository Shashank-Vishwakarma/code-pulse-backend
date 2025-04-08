FROM alpine:latest
RUN apk add --no-cache g++ libstdc++ 
WORKDIR /app
COPY . /app
RUN g++ -o main main.cpp
CMD ["./main"]