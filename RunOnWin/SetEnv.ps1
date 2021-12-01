$env:DB_NAME="mortgage_info"
$env:DB_PASS="9544EEiiiiDDDfghIJJklpumpkin"
$env:DB_TYPE="postgres"
$env:DB_USER="web_ro"
$env:DB_HOST="uboo.supskifun.net"
Start-Job -Name Mortgage { .\winx64_mortgage.exe}