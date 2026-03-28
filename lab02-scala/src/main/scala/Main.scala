package com.example

import org.eclipse.jetty.server.Server
import org.eclipse.jetty.servlet.{ServletHolder, ServletContextHandler}

object Main extends App {
  val port = 8080
  val server = new Server(port)

  val context = new ServletContextHandler(ServletContextHandler.SESSIONS)
  context.setContextPath("/")
  context.addServlet(new ServletHolder(new ProductServlet), "/products/*")
  context.addServlet(new ServletHolder(new CategoryServlet), "/categories/*")
  context.addServlet(new ServletHolder(new CartServlet), "/carts/*")

  server.setHandler(context)
  server.start()
  server.join()
}