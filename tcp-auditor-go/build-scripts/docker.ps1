#!/bin/pwsh
param(
    [Parameter(Mandatory=$true)]
    [string]$VersionNumber,
    [string]$PublishContainer="False",
    [string]$ContainerRegistryUri="",
    [string]$RegistryUser="",
    [string]$RemoveContainerAfterPush="True"
)

Write-Host "##teamcity[message text='LOG 1 - tcp-auditor-go - Docker - Starting script']"
Write-Host "##teamcity[message text='LOG 1 - VN: $VersionNumber  PC: $PublishContainer CRU: $ContainerRegistryUri RU: $RegistryUser RCAP: $RemoveContainerAfterPush']"

if("$PublishContainer" -eq $true) {
    Write-Host "##teamcity[message text='LOG 2 - tcp-auditor-go - Docker - Publishing container']"

    if([string]::IsNullOrWhiteSpace($RegistryUser)) {
        Write-Host "##teamcity[message text='LOG 3 - tcp-auditor-go - Docker - No registry user']"

        Write-Host "No Container RegistryUser Provided."
        Write-Host "Failed to publish container."
        exit 1
    } else {
        Write-Host "##teamcity[message text='LOG 4 - tcp-auditor-go - Docker - User exists']"

        Write-Host "[NOTICE] Assuming docker is logged into registry!"

        # Build private registry string
        $containerRegistryStr = "$RegistryUser/tcp-auditor-go"
        if(-not [string]::IsNullOrWhiteSpace($ContainerRegistryUri)) {
            $containerRegistryStr = "$ContainerRegistryUri/$containerRegistryStr"
        }

        Write-Host "##teamcity[message text='LOG 5 - tcp-auditor-go - Docker - Registry string']"
        Write-Host "##teamcity[message text='LOG 5 - CRS: $containerRegistryStr']"

        # Build Container
        Write-Host "Building container"
        $tags = @("$($containerRegistryStr):$VersionNumber")
        if($VersionNumber.EndsWith("-dev")) {
            $tags += "$($containerRegistryStr):latest-dev"
        } else {
            $tags += "$($containerRegistryStr):latest"
        }
        $tagStr = "-t $($tags -join ' -t ')"

        Write-Host "##teamcity[message text='LOG 6 - tcp-auditor-go - Docker - Building container']"
        Write-Host "##teamcity[message text='LOG 6 - TS: $tagStr']"

        Write-Host "Checking Docker OS"
        $osType = docker info --format='{{.OSType}}'
        $osChanged = $false
        if ($osType -eq 'windows') {
            Write-Host "Docker in $osType mode. Switching to linux."
            & 'C:\Program Files\Docker\Docker\DockerCli.exe' -SwitchDaemon
            $osChanged = $true
        }

        Write-Host "##teamcity[message text='LOG 6.5 - tcp-auditor-go - Docker - Done switching']"

        Invoke-Expression "docker build . -t tcp-auditor-go $tagStr"

        Write-Host "##teamcity[message text='LOG 7 - tcp-auditor-go - Docker - Container built']"

        Write-Host "##teamcity[message text='LOG 8 - tcp-auditor-go - Docker - Pushing containers']"

        # Push Container tags
        ForEach($tag in $tags) {
            Write-Host "##teamcity[message text='LOG 9 - tcp-auditor-go - Docker - Pushing']"
            Write-Host "##teamcity[message text='LOG 9 - T: $tag']"

            Write-Host "Pushing Container `"$tag`""
            Invoke-Expression "docker push $tag"
        }

        Write-Host "##teamcity[message text='LOG 10 - tcp-auditor-go - Docker - Done pushing']"
        Write-Host "##teamcity[message text='LOG 11 - tcp-auditor-go - Docker - Removing local containers']"

        # Remove Local Containers
        if("$RemoveContainerAfterPush" -eq $true) {
            # Removing Local Tags
            ForEach($tag in $tags) {
                Write-Host "Removing Local Container `"$tag`""
                Invoke-Expression "docker image rm $tag"
            }
        }

        Write-Host "##teamcity[message text='LOG 11.5 - tcp-auditor-go - Docker - Switching back']"

        Write-Host "Reverting Docker OS"
        if ($osChanged) {
            Write-Host "Switching back to windows mode."
            & 'C:\Program Files\Docker\Docker\DockerCli.exe' -SwitchDaemon
        }

        Write-Host "##teamcity[message text='LOG 12 - tcp-auditor-go - Docker - Complete']"
    }
}

Write-Host "##teamcity[message text='LOG 13 - tcp-auditor-go - Docker - Script finished']"