# Create admin user
Write-Host "=== CREATING ADMIN USER ===" -ForegroundColor Cyan

# Wait a moment for server to start
Start-Sleep -Seconds 3

# 1. Check if server is running
Write-Host "`n1. Checking server health..." -ForegroundColor Yellow
try {
    $health = Invoke-RestMethod -Uri "http://localhost:8080/health" -Method GET
    Write-Host "Server is running! Status: $($health.status)" -ForegroundColor Green
} catch {
    Write-Host "Server not running. Please start the backend first." -ForegroundColor Red
    exit 1
}

# 2. Register admin user
Write-Host "`n2. Registering admin user..." -ForegroundColor Yellow
try {
    $registerBody = @{
        first_name = "Admin"
        last_name = "User"
        email = "admin@labelops.com"
        password = "Admin@123"
        role = "admin"
    } | ConvertTo-Json

    $registerResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/register" -Method POST -ContentType "application/json" -Body $registerBody
    Write-Host "Admin user created successfully!" -ForegroundColor Green
    Write-Host "User ID: $($registerResponse.user.id)" -ForegroundColor Green
    Write-Host "Email: $($registerResponse.user.email)" -ForegroundColor Green
    Write-Host "Role: $($registerResponse.user.role)" -ForegroundColor Green
} catch {
    Write-Host "Registration failed: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.Exception.Response) {
        $response = $_.Exception.Response.GetResponseStream()
        $reader = New-Object System.IO.StreamReader($response)
        $responseBody = $reader.ReadToEnd()
        Write-Host "Response: $responseBody" -ForegroundColor Red
    }
}

# 3. Test login with admin user
Write-Host "`n3. Testing admin login..." -ForegroundColor Yellow
try {
    $loginBody = @{
        email = "admin@labelops.com"
        password = "Admin@123"
    } | ConvertTo-Json

    $loginResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/login" -Method POST -ContentType "application/json" -Body $loginBody
    $token = $loginResponse.token
    $user = $loginResponse.user
    Write-Host "Admin login successful!" -ForegroundColor Green
    Write-Host "User ID: $($user.id)" -ForegroundColor Green
    Write-Host "Role: $($user.role)" -ForegroundColor Green
    Write-Host "Token: $($token.Substring(0, 50))..." -ForegroundColor Green
} catch {
    Write-Host "Admin login failed: $($_.Exception.Message)" -ForegroundColor Red
}

# 4. Test labels endpoint with admin
Write-Host "`n4. Testing labels endpoint with admin..." -ForegroundColor Yellow
try {
    $headers = @{
        "Authorization" = "Bearer $token"
        "Content-Type" = "application/json"
    }
    
    $labelsResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/labels" -Method GET -Headers $headers
    Write-Host "Admin labels count: $($labelsResponse.count)" -ForegroundColor Green
    if ($labelsResponse.labels) {
        Write-Host "Labels found: $($labelsResponse.labels.Count)" -ForegroundColor Green
    } else {
        Write-Host "No labels found" -ForegroundColor Yellow
    }
} catch {
    Write-Host "Labels test failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n=== ADMIN USER CREATION COMPLETED ===" -ForegroundColor Cyan 