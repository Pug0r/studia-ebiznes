plugins {
    kotlin("jvm") version "1.9.23"
    kotlin("plugin.serialization") version "1.9.23"
    application
}

group = "discordbot"
version = "1.0.0"

repositories {
    mavenCentral()
}

val ktorVersion = "2.3.9"

dependencies {
    implementation("dev.kord:kord-core:0.13.1")
    
    implementation("io.ktor:ktor-server-core:2.3.10")
    implementation("io.ktor:ktor-server-netty:2.3.10")
    implementation("io.ktor:ktor-server-content-negotiation:2.3.10")
    implementation("io.ktor:ktor-serialization-kotlinx-json:2.3.10")

    implementation("io.ktor:ktor-client-core:2.3.10")
    implementation("io.ktor:ktor-client-cio:2.3.10")
    implementation("io.ktor:ktor-client-content-negotiation:2.3.10")
}

application {
    mainClass.set("discordbot.MainKt")
}