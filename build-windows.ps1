# This script compiles module to WASM but uses under Windows NT
# My test release named Frank Iero.
function Invoke-GoWasmBuild {
    param(
        [string]$OutputPath = "build/frank_iero/main.wasm"
    )
    $oldGOOS = $env:GOOS
    $oldGOARCH = $env:GOARCH

    $env:GOOS = "js"
    $env:GOARCH = "wasm"
    
    go build -o $OutputPath
    
    ### Restore previous platform settings 
    if ($oldGOOS) { 
        $env:GOOS = $oldGOOS 
    } else { 
        Remove-Item Env:GOOS -ErrorAction SilentlyContinue 
    }

    if ($oldGOARCH) { 
        $env:GOARCH = $oldGOARCH
    } else { 
        Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
    }
}

Invoke-GoWasmBuild