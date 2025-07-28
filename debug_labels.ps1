# Debug labels script
Write-Host "=== DEBUG LABELS SCRIPT ===" -ForegroundColor Cyan

# 1. Login and get token
Write-Host "`n1. Logging in..." -ForegroundColor Yellow
try {
    $loginBody = @{
        email = "tech@gmail.com"
        password = "Tech@123"
    } | ConvertTo-Json

    $loginResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/login" -Method POST -ContentType "application/json" -Body $loginBody
    $token = $loginResponse.token
    $user = $loginResponse.user
    Write-Host "Login successful! User ID: $($user.id), Role: $($user.role)" -ForegroundColor Green
} catch {
    Write-Host "Login failed: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# 2. Create fresh labels with unique IDs
Write-Host "`n2. Creating fresh labels..." -ForegroundColor Yellow
$timestamp = Get-Date -Format "yyyyMMddHHmmss"
$freshLabels = @{
    labels = @(
        @{
            LOCATION = $null
            BUNDLE_NOS = 1
            PQD = "101520002123005$timestamp"
            UNIT = "SAIL-BSP"
            TIME1 = "10:00"
            LENGTH = "STD"
            HEAT_NO = "C075$timestamp"
            PRODUCT_HEADING = "TMT BAR"
            ISI_BOTTOM = "CML 187244"
            ISI_TOP = "IS 1786:2008"
            CHARGE_DTM = "202403041000"
            MILL = "MM"
            GRADE = "IS 1786 FE550D"
            URL_APIKEY = "4700b47719c54d979af7f7a2d1dc67db"
            ID = "DEBUG$timestamp"
            WEIGHT = $null
            SECTION = "TMT BAR 25"
            DATE1 = "04-MAR-24"
        },
        @{
            LOCATION = $null
            BUNDLE_NOS = 2
            PQD = "101520002123005$($timestamp)2"
            UNIT = "SAIL-BSP"
            TIME1 = "10:05"
            LENGTH = "STD"
            HEAT_NO = "C075$($timestamp)2"
            PRODUCT_HEADING = "TMT BAR"
            ISI_BOTTOM = "CML 187244"
            ISI_TOP = "IS 1786:2008"
            CHARGE_DTM = "202403041005"
            MILL = "MM"
            GRADE = "IS 1786 FE550D"
            URL_APIKEY = "4700b47719c54d979af7f7a2d1dc67db"
            ID = "DEBUG$($timestamp)2"
            WEIGHT = $null
            SECTION = "TMT BAR 20"
            DATE1 = "04-MAR-24"
        }
    )
} | ConvertTo-Json -Depth 3

# 3. Upload fresh labels
Write-Host "`n3. Uploading fresh labels..." -ForegroundColor Yellow
try {
    $headers = @{
        "Authorization" = "Bearer $token"
        "Content-Type" = "application/json"
    }
    
    $uploadResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/labels/batch" -Method POST -Headers $headers -Body $freshLabels
    Write-Host "Upload response: $($uploadResponse | ConvertTo-Json)" -ForegroundColor Green
} catch {
    Write-Host "Upload failed: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "Response: $($_.Exception.Response)" -ForegroundColor Red
}

# 4. Test labels endpoint
Write-Host "`n4. Testing labels endpoint..." -ForegroundColor Yellow
try {
    $labelsResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/labels" -Method GET -Headers $headers
    Write-Host "Labels count: $($labelsResponse.count)" -ForegroundColor Green
    if ($labelsResponse.labels) {
        Write-Host "First label: $($labelsResponse.labels[0] | ConvertTo-Json)" -ForegroundColor Green
    } else {
        Write-Host "No labels returned" -ForegroundColor Red
    }
} catch {
    Write-Host "Labels test failed: $($_.Exception.Message)" -ForegroundColor Red
}

# 5. Test with different query parameters
Write-Host "`n5. Testing with different parameters..." -ForegroundColor Yellow
try {
    $labelsResponse2 = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/labels?limit=10&offset=0" -Method GET -Headers $headers
    Write-Host "Labels with params count: $($labelsResponse2.count)" -ForegroundColor Green
} catch {
    Write-Host "Labels with params failed: $($_.Exception.Message)" -ForegroundColor Red
}

# 6. Check health endpoint
Write-Host "`n6. Checking health endpoint..." -ForegroundColor Yellow
try {
    $health = Invoke-RestMethod -Uri "http://localhost:8080/health" -Method GET
    Write-Host "Health: $($health.status)" -ForegroundColor Green
} catch {
    Write-Host "Health check failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n=== DEBUG COMPLETED ===" -ForegroundColor Cyan 