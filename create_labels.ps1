# Create and test labels
Write-Host "Creating and testing labels..." -ForegroundColor Green

# 1. Login as tech user
Write-Host "`n1. Logging in as tech user..." -ForegroundColor Yellow
$loginBody = @{
    email = "tech@gmail.com"
    password = "Tech@123"
} | ConvertTo-Json

$loginResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/login" -Method POST -ContentType "application/json" -Body $loginBody
$token = $loginResponse.token
Write-Host "Login successful!" -ForegroundColor Green

# 2. Create fresh labels data
Write-Host "`n2. Creating fresh labels..." -ForegroundColor Yellow
$freshLabels = @{
    labels = @(
        @{
            LOCATION = $null
            BUNDLE_NOS = 1
            PQD = "101520002123005001"
            UNIT = "SAIL-BSP"
            TIME1 = "10:00"
            LENGTH = "STD"
            HEAT_NO = "C075500"
            PRODUCT_HEADING = "TMT BAR"
            ISI_BOTTOM = "CML 187244"
            ISI_TOP = "IS 1786:2008"
            CHARGE_DTM = "202403041000"
            MILL = "MM"
            GRADE = "IS 1786 FE550D"
            URL_APIKEY = "4700b47719c54d979af7f7a2d1dc67db"
            ID = "FRESH001"
            WEIGHT = $null
            SECTION = "TMT BAR 25"
            DATE1 = "04-MAR-24"
        },
        @{
            LOCATION = $null
            BUNDLE_NOS = 2
            PQD = "101520002123005002"
            UNIT = "SAIL-BSP"
            TIME1 = "10:05"
            LENGTH = "STD"
            HEAT_NO = "C075501"
            PRODUCT_HEADING = "TMT BAR"
            ISI_BOTTOM = "CML 187244"
            ISI_TOP = "IS 1786:2008"
            CHARGE_DTM = "202403041005"
            MILL = "MM"
            GRADE = "IS 1786 FE550D"
            URL_APIKEY = "4700b47719c54d979af7f7a2d1dc67db"
            ID = "FRESH002"
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
}

# 4. Test labels endpoint
Write-Host "`n4. Testing labels endpoint..." -ForegroundColor Yellow
try {
    $labelsResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/labels" -Method GET -Headers $headers
    Write-Host "Labels count: $($labelsResponse.count)" -ForegroundColor Green
    Write-Host "Labels: $($labelsResponse.labels | ConvertTo-Json -Depth 1)" -ForegroundColor Green
} catch {
    Write-Host "Labels test failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`nLabels creation and test completed!" -ForegroundColor Green 