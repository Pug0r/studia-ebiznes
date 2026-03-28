name := "scalatra-product-crud"

version := "0.1.0"

scalaVersion := "2.13.10"

libraryDependencies ++= Seq(
  "org.scalatra" %% "scalatra" % "2.8.4",
  "org.scalatra" %% "scalatra-json" % "2.8.4",
  "org.json4s" %% "json4s-jackson" % "4.0.6",
  "ch.qos.logback" % "logback-classic" % "1.4.7" % "runtime",
  "org.eclipse.jetty" % "jetty-server" % "9.4.51.v20230217",
  "org.eclipse.jetty" % "jetty-servlet" % "9.4.51.v20230217",
  "org.eclipse.jetty" % "jetty-webapp" % "9.4.51.v20230217",
  "javax.servlet" % "javax.servlet-api" % "3.1.0" % "provided"
)

enablePlugins(AssemblyPlugin)
