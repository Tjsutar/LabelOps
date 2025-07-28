# Check database directly
Write-Host "Checking database for labels..." -ForegroundColor Green

# Login to get token
$loginBody = @{
    email = "tech@gmail.com"
    password = "Tech@123"
} | ConvertTo-Json

$loginResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/login" -Method POST -ContentType "application/json" -Body $loginBody
$token = $loginResponse.token

# Test labels endpoint with admin token
Write-Host "`nTesting with admin login..." -ForegroundColor Yellow
try {
    $adminLoginBody = @{
        email = "admin@gmail.com"
        password = "admin123"
    } | ConvertTo-Json

    $adminResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/login" -Method POST -ContentType "application/json" -Body $adminLoginBody
    $adminToken = $adminResponse.token
    
    $headers = @{
        "Authorization" = "Bearer $adminToken"
        "Content-Type" = "application/json"
    }
    
    $adminLabelsResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/labels" -Method GET -Headers $headers
    Write-Host "Admin labels count: $($adminLabelsResponse.count)" -ForegroundColor Green
    Write-Host "Admin labels: $($adminLabelsResponse.labels | ConvertTo-Json -Depth 1)" -ForegroundColor Green
} catch {
    Write-Host "Admin test failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test with tech user
Write-Host "`nTesting with tech user..." -ForegroundColor Yellow
try {
    $headers = @{
        "Authorization" = "Bearer $token"
        "Content-Type" = "application/json"
    }
    
    $techLabelsResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/labels" -Method GET -Headers $headers
    Write-Host "Tech user labels count: $($techLabelsResponse.count)" -ForegroundColor Green
    Write-Host "Tech user labels: $($techLabelsResponse.labels | ConvertTo-Json -Depth 1)" -ForegroundColor Green
} catch {
    Write-Host "Tech user test failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`nDatabase check completed!" -ForegroundColor Green 