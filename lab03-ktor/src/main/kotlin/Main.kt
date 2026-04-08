package discordbot

import dev.kord.common.entity.Snowflake
import dev.kord.core.Kord
import dev.kord.core.behavior.interaction.response.respond
import dev.kord.core.event.interaction.ChatInputCommandInteractionCreateEvent
import dev.kord.core.on
import dev.kord.gateway.Intent
import dev.kord.gateway.PrivilegedIntent
import dev.kord.rest.builder.interaction.string
import io.ktor.client.*
import io.ktor.client.call.*
import io.ktor.client.engine.cio.*
import io.ktor.client.plugins.contentnegotiation.ContentNegotiation as ClientNegotiation
import io.ktor.client.request.*
import io.ktor.serialization.kotlinx.json.*
import io.ktor.server.application.*
import io.ktor.server.engine.*
import io.ktor.server.netty.*
import io.ktor.server.plugins.contentnegotiation.ContentNegotiation as ServerNegotiation
import io.ktor.server.response.*
import io.ktor.server.routing.*
import kotlinx.coroutines.launch

val mockData = mapOf(
    "electronics" to listOf("Laptop", "Smartphone", "Headphones"),
    "books" to listOf("Kotlin in Action", "Clean Code", "The Witcher"),
    "games" to listOf("Elden Ring", "Cyberpunk 2077", "Minecraft")
)

@OptIn(PrivilegedIntent::class)
suspend fun main() {
    val kord = Kord(System.getenv("DISCORD_BOT_TOKEN")!!)
    val channelId = Snowflake(System.getenv("DISCORD_CHANNEL_ID")!!.trim())

    kord.launch {
        embeddedServer(Netty, port = 8080) {
            install(ServerNegotiation) { json() }
            routing {
                get("/api/say/{text}") {
                    val text = call.parameters["text"]!!
                    kord.rest.channel.createMessage(channelId) { content = "Message from the Web: **$text**" }
                    call.respondText("Success! Sent '$text' to Discord.")
                }
                get("/api/categories") { call.respond(mockData.keys.toList()) }
                get("/api/category/{name}") {
                    mockData[call.parameters["name"]?.lowercase()]?.let { call.respond(it) } ?: call.respondText("Not Found")
                }
            }
        }.start(wait = true)
    }

    val client = HttpClient(CIO) { install(ClientNegotiation) { json() } }

    kord.createGlobalChatInputCommand("categories", "Lists categories from API")
    kord.createGlobalChatInputCommand("category", "Gets products from API") { string("name", "Category name") { required = true } }

    kord.on<ChatInputCommandInteractionCreateEvent> {
        val response = interaction.deferPublicResponse()
        val cmd = interaction.command
        when (cmd.rootName) {
            "categories" -> response.respond { content = "API says: ${client.get("http://localhost:8080/api/categories").body<List<String>>().joinToString()}" }
            "category" -> response.respond { content = "Products in ${cmd.strings["name"]}: ${client.get("http://localhost:8080/api/category/${cmd.strings["name"]}").body<List<String>>().joinToString()}" }
        }
    }

    kord.login { intents += Intent.MessageContent }
}