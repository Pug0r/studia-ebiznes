import java.sql.DriverManager

fun main() {
    try {
        Class.forName("org.sqlite.JDBC")
        
        val connection = DriverManager.getConnection("jdbc:sqlite:test.db")
        
        println("works")
        connection.close()
    } catch (e: Exception) {
        println("doesnt work")
    }
}