# Detailed Debug GetLabels API
Write-Host "=== DETAILED DEBUG GETLABELS ===" -ForegroundColor Cyan

# 1. Login as admin
Write-Host "`n1. Logging in as admin..." -ForegroundColor Yellow
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

# 2. Test GetLabels with full response logging
Write-Host "`n2. Testing GetLabels with full response..." -ForegroundColor Yellow
try {
    $headers = @{
        "Authorization" = "Bearer $token"
        "Content-Type" = "application/json"
    }
    
    Write-Host "Making request to: http://localhost:8080/api/v1/labels" -ForegroundColor Cyan
    $labelsResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/labels" -Method GET -Headers $headers
    
    Write-Host "`nFULL RESPONSE:" -ForegroundColor Green
    $labelsResponse | ConvertTo-Json -Depth 10 | Write-Host
    
    Write-Host "`nRESPONSE SUMMARY:" -ForegroundColor Green
    Write-Host "  Count: $($labelsResponse.count)" -ForegroundColor White
    Write-Host "  Limit: $($labelsResponse.limit)" -ForegroundColor White
    Write-Host "  Offset: $($labelsResponse.offset)" -ForegroundColor White
    Write-Host "  User ID: $($labelsResponse.user_id)" -ForegroundColor White
    Write-Host "  User Role: $($labelsResponse.user_role)" -ForegroundColor White
    
    if ($labelsResponse.labels) {
        Write-Host "`nLABELS DETAILS:" -ForegroundColor Green
        Write-Host "  Number of labels: $($labelsResponse.labels.Count)" -ForegroundColor White
        
        for ($i = 0; $i -lt $labelsResponse.labels.Count; $i++) {
            $label = $labelsResponse.labels[$i]
            Write-Host "`n  Label $($i + 1):" -ForegroundColor Cyan
            $label.PSObject.Properties | ForEach-Object {
                Write-Host "    $($_.Name): $($_.Value)" -ForegroundColor White
            }
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

# 3. Test with different limit values
Write-Host "`n3. Testing with different limit values..." -ForegroundColor Yellow
try {
    Write-Host "Testing with limit=1..." -ForegroundColor Cyan
    $response1 = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/labels?limit=1" -Method GET -Headers $headers
    Write-Host "  Labels with limit=1: $($response1.count)" -ForegroundColor Green
    
    Write-Host "Testing with limit=5..." -ForegroundColor Cyan
    $response5 = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/labels?limit=5" -Method GET -Headers $headers
    Write-Host "  Labels with limit=5: $($response5.count)" -ForegroundColor Green
    
    Write-Host "Testing with limit=10..." -ForegroundColor Cyan
    $response10 = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/labels?limit=10" -Method GET -Headers $headers
    Write-Host "  Labels with limit=10: $($response10.count)" -ForegroundColor Green
} catch {
    Write-Host "Limit test failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n=== DETAILED DEBUG COMPLETED ===" -ForegroundColor Cyan 