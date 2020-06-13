FROM golang:1.13

WORKDIR /src

COPY . .

RUN go build -o svcindex .

RUN apt-get update && apt-get install -y curl && \
    curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s \
      https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl && \
    chmod +x ./kubectl && mv ./kubectl /bin/kubectl

ENTRYPOINT [ "./svcindex" ]
