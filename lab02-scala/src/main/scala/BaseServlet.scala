package com.example

import org.scalatra._
import org.scalatra.json._
import org.json4s.{DefaultFormats, Formats}

abstract class BaseServlet extends ScalatraServlet with JacksonJsonSupport {

  protected implicit lazy val jsonFormats: Formats = DefaultFormats

 // special method defined by scaltra, this runs before every request 
  before() {
    contentType = formats("json")
    response.setHeader("Access-Control-Allow-Origin", "http://localhost:3000, http://localhost:8081")
    response.setHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    response.setHeader("Access-Control-Allow-Headers", "Content-Type")
  }
}