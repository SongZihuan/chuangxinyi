New-Item -ItemType Directory -Path "src\service\v1" -Force
goctlwt api go --api api\v1\service.api --dir src\service\v1
