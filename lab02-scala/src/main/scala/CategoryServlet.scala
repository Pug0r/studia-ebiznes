package com.example

import org.scalatra._
import scala.collection.mutable

case class Category(id: Int, name: String, description: String)

class CategoryServlet extends BaseServlet {

  private val categories = mutable.Map[Int, Category]()
  private var nextId = 1

  categories(1) = Category(1, "Electronics", "Electronic devices and gadgets")
  categories(2) = Category(2, "Books", "Books and publications")
  nextId = 3

  get("/") {
    categories.values.toList
  }

  get("/:id") {
    val id = params("id").toInt
    categories.get(id) match {
      case Some(category) => category
      case None => NotFound(Map("error" -> "Category not found"))
    }
  }

  post("/") {
    val categoryData = parsedBody.extract[Map[String, Any]]
    val name = categoryData("name").asInstanceOf[String]
    val description = categoryData("description").asInstanceOf[String]
    val category = Category(nextId, name, description)
    categories(nextId) = category
    nextId += 1
    category
  }

  put("/:id") {
    val id = params("id").toInt
    categories.get(id) match {
      case Some(existingCategory) =>
        val categoryData = parsedBody.extract[Map[String, Any]]
        val name = categoryData.getOrElse("name", existingCategory.name).asInstanceOf[String]
        val description = categoryData.getOrElse("description", existingCategory.description).asInstanceOf[String]
        val updatedCategory = existingCategory.copy(name = name, description = description)
        categories(id) = updatedCategory
        updatedCategory
      case None => NotFound(Map("error" -> "Category not found"))
    }
  }

  delete("/:id") {
    val id = params("id").toInt
    categories.remove(id) match {
      case Some(category) => Ok(Map("message" -> s"Category ${category.name} deleted"))
      case None => NotFound(Map("error" -> "Category not found"))
    }
  }
}