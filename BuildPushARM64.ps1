<#
    Builds testWeb-arm64 from main.go and pushes to uboo:/usr/local/docker/testWeb/
    See deploy.sh, cmd.txt and env.txt in uboo:/usr/local/docker/testWeb/
#>

Function SetConfig {
    go env -w CGO_ENABLED=0
    go env -w GOOS=linux
    go env -w GOARCH=arm64
    go env -w GOARM=6
    $script:n = "./testWeb-arm64"
}

Function BuildFile {
    $cmd1 = "go build -a -installsuffix cgo -o $($n) ."
    Invoke-Expression -Command $cmd1
}

Function SCPFile {
    $cmd2 = "scp $($n) uboo:/usr/local/docker/testWeb/$($n)"
    Invoke-Expression -Command $cmd2
}
Function ResetConfig {
    go env -w CGO_ENABLED=0
    go env -w GOOS=windows
    go env -w GOARCH=amd64
}


SetConfig
BuildFile
SCPFile
ResetConfig