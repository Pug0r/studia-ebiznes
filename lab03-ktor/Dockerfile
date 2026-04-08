FROM gradle:8.7-jdk17 AS build
COPY --chown=gradle:gradle . /home/gradle/src
WORKDIR /home/gradle/src

RUN gradle installDist --no-daemon

FROM eclipse-temurin:17-jre
RUN mkdir /app
COPY --from=build /home/gradle/src/build/install/ktor-discord-client /app/

WORKDIR /app

CMD ["./bin/ktor-discord-client"]