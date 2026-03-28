package com.example

import org.scalatra._
import scala.collection.mutable

case class Cart(id: Int, userId: Int, productIds: List[Int])

class CartServlet extends BaseServlet {

  private val carts = mutable.Map[Int, Cart]()
  private var nextId = 1

  carts(1) = Cart(1, 1, List(1, 2))  
  carts(2) = Cart(2, 2, List(1)) 
  nextId = 3

  get("/") {
    carts.values.toList
  }

  get("/:id") {
    val id = params("id").toInt
    carts.get(id) match {
      case Some(cart) => cart
      case None => NotFound(Map("error" -> "Cart not found"))
    }
  }

  post("/") {
    val cartData = parsedBody.extract[Map[String, Any]]
    val userId = cartData("userId").asInstanceOf[BigInt].toInt
    val productIds = cartData("productIds").asInstanceOf[List[BigInt]].map(_.toInt)
    val cart = Cart(nextId, userId, productIds)
    carts(nextId) = cart
    nextId += 1
    cart
  }

  put("/:id") {
    val id = params("id").toInt
    carts.get(id) match {
      case Some(existingCart) =>
        val cartData = parsedBody.extract[Map[String, Any]]
        val userId = cartData.getOrElse("userId", existingCart.userId).asInstanceOf[BigInt].toInt
        val productIds = cartData.getOrElse("productIds", existingCart.productIds).asInstanceOf[List[BigInt]].map(_.toInt)
        val updatedCart = existingCart.copy(userId = userId, productIds = productIds)
        carts(id) = updatedCart
        updatedCart
      case None => NotFound(Map("error" -> "Cart not found"))
    }
  }

  delete("/:id") {
    val id = params("id").toInt
    carts.remove(id) match {
      case Some(cart) => Ok(Map("message" -> s"Cart for user ${cart.userId} deleted"))
      case None => NotFound(Map("error" -> "Cart not found"))
    }
  }
}