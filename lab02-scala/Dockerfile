FROM eclipse-temurin:11

# Install sbt
RUN apt-get update && \
    apt-get install -y curl && \
    curl -L https://github.com/sbt/sbt/releases/download/v1.9.7/sbt-1.9.7.tgz | tar xz && \
    mv sbt /usr/local && \
    ln -s /usr/local/sbt/bin/sbt /usr/local/bin/sbt && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY . /app

EXPOSE 8080

CMD ["sbt", "run"]