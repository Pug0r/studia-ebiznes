package com.example

import org.scalatra._
import scala.collection.mutable

case class Product(id: Int, name: String, price: Double, description: String)

class ProductServlet extends BaseServlet {

  private val products = mutable.Map[Int, Product]()
  private var nextId = 1

  products(1) = Product(1, "Laptop", 999.99, "High-performance laptop")
  products(2) = Product(2, "Mouse", 29.99, "Wireless mouse")
  nextId = 3

  get("/") {
    products.values.toList
  }

  get("/:id") {
    val id = params("id").toInt
    products.get(id) match {
      case Some(product) => product
      case None => NotFound(Map("error" -> "Product not found"))
    }
  }

  post("/") {
    val productData = parsedBody.extract[Map[String, Any]]
    val name = productData("name").asInstanceOf[String]
    val price = productData("price").asInstanceOf[Double]
    val description = productData("description").asInstanceOf[String]
    val product = Product(nextId, name, price, description)
    products(nextId) = product
    nextId += 1
    product
  }

  put("/:id") {
    val id = params("id").toInt
    products.get(id) match {
      case Some(existingProduct) =>
        val productData = parsedBody.extract[Map[String, Any]]
        val name = productData.getOrElse("name", existingProduct.name).asInstanceOf[String]
        val price = productData.getOrElse("price", existingProduct.price).asInstanceOf[Double]
        val description = productData.getOrElse("description", existingProduct.description).asInstanceOf[String]
        val updatedProduct = existingProduct.copy(name = name, price = price, description = description)
        products(id) = updatedProduct
        updatedProduct
      case None => NotFound(Map("error" -> "Product not found"))
    }
  }

  delete("/:id") {
    val id = params("id").toInt
    products.remove(id) match {
      case Some(product) => Ok(Map("message" -> s"Product ${product.name} deleted"))
      case None => NotFound(Map("error" -> "Product not found"))
    }
  }
}