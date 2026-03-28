Write-Host "Building Docker image..."
docker build -t scalatra-crud-app .

Write-Host "Running the application in Docker..."
docker run -d -p 8080:8080 --name scalatra-app scalatra-crud-app

Start-Sleep 5

Write-Host "Starting ngrok tunnel..."
ngrok http 8080

# if it fails you need to run on host machine (after registering to ngrok)
# ngrok config add-authtoken YOUR_TOKEN

Write-Host "Stopping the application..."
docker stop scalatra-app
docker rm scalatra-app