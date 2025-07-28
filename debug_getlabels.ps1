# Debug GetLabels API - Comprehensive Test
Write-Host "=== DEBUG GETLABELS API ===" -ForegroundColor Cyan

# 1. Test server health
Write-Host "`n1. Testing server health..." -ForegroundColor Yellow
try {
    $healthResponse = Invoke-RestMethod -Uri "http://localhost:8080/health" -Method GET
    Write-Host "Server is running: $($healthResponse.status)" -ForegroundColor Green
} catch {
    Write-Host "Server not running: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# 2. Login as admin
Write-Host "`n2. Logging in as admin..." -ForegroundColor Yellow
try {
    $loginBody = @{
        email = "admin@labelops.com"
        password = "Admin@123"
    } | ConvertTo-Json

    $loginResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/login" -Method POST -ContentType "application/json" -Body $loginBody
    $token = $loginResponse.token
    $user = $loginResponse.user
    Write-Host "Admin login successful! User ID: $($user.id), Role: $($user.role)" -ForegroundColor Green
} catch {
    Write-Host "Admin login failed: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# 3. Test GetLabels with detailed logging
Write-Host "`n3. Testing GetLabels with admin..." -ForegroundColor Yellow
try {
    $headers = @{
        "Authorization" = "Bearer $token"
        "Content-Type" = "application/json"
    }
    
    Write-Host "Making request to: http://localhost:8080/api/v1/labels" -ForegroundColor Cyan
    $labelsResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/labels" -Method GET -Headers $headers
    
    Write-Host "SUCCESS! Response received:" -ForegroundColor Green
    Write-Host "  Count: $($labelsResponse.count)" -ForegroundColor White
    Write-Host "  Limit: $($labelsResponse.limit)" -ForegroundColor White
    Write-Host "  Offset: $($labelsResponse.offset)" -ForegroundColor White
    Write-Host "  User ID: $($labelsResponse.user_id)" -ForegroundColor White
    Write-Host "  User Role: $($labelsResponse.user_role)" -ForegroundColor White
    
    if ($labelsResponse.labels) {
        Write-Host "`nLabels found: $($labelsResponse.labels.Count)" -ForegroundColor Green
        for ($i = 0; $i -lt $labelsResponse.labels.Count; $i++) {
            $label = $labelsResponse.labels[$i]
            Write-Host "`nLabel $($i + 1):" -ForegroundColor Cyan
            Write-Host "  ID: $($label.id)" -ForegroundColor White
            Write-Host "  Label ID: $($label.label_id)" -ForegroundColor White
            Write-Host "  Status: $($label.status)" -ForegroundColor White
            Write-Host "  User ID: $($label.user_id)" -ForegroundColor White
            Write-Host "  Created At: $($label.created_at)" -ForegroundColor White
        }
    } else {
        Write-Host "No labels found in response" -ForegroundColor Yellow
    }
} catch {
    Write-Host "GetLabels failed: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.Exception.Response) {
        $stream = $_.Exception.Response.GetResponseStream()
        $reader = New-Object System.IO.StreamReader($stream)
        $body = $reader.ReadToEnd()
        Write-Host "Response body: '$body'" -ForegroundColor Red
    }
}

# 4. Test with different parameters
Write-Host "`n4. Testing with different parameters..." -ForegroundColor Yellow
try {
    Write-Host "Testing with limit=100..." -ForegroundColor Cyan
    $labelsResponse2 = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/labels?limit=100" -Method GET -Headers $headers
    Write-Host "Labels with limit=100: $($labelsResponse2.count)" -ForegroundColor Green
    
    Write-Host "Testing with no filters..." -ForegroundColor Cyan
    $labelsResponse3 = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/labels?limit=50&offset=0" -Method GET -Headers $headers
    Write-Host "Labels with no filters: $($labelsResponse3.count)" -ForegroundColor Green
} catch {
    Write-Host "Parameter test failed: $($_.Exception.Message)" -ForegroundColor Red
}

# 5. Test with tech user to compare
Write-Host "`n5. Testing with tech user..." -ForegroundColor Yellow
try {
    $techLoginBody = @{
        email = "tech@gmail.com"
        password = "Tech@123"
    } | ConvertTo-Json

    $techLoginResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/login" -Method POST -ContentType "application/json" -Body $techLoginBody
    $techToken = $techLoginResponse.token
    $techUser = $techLoginResponse.user
    Write-Host "Tech login successful! User ID: $($techUser.id), Role: $($techUser.role)" -ForegroundColor Green
    
    $techHeaders = @{
        "Authorization" = "Bearer $techToken"
        "Content-Type" = "application/json"
    }
    
    $techLabelsResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/labels" -Method GET -Headers $techHeaders
    Write-Host "Tech user labels count: $($techLabelsResponse.count)" -ForegroundColor Green
    Write-Host "Tech user role: $($techLabelsResponse.user_role)" -ForegroundColor Green
} catch {
    Write-Host "Tech user test failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n=== DEBUG COMPLETED ===" -ForegroundColor Cyan 