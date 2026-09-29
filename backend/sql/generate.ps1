New-Item -ItemType Directory -Path "src\model\db" -Force
goctlwt model mysql ddl --src sql\user.sql --dir src\model\db

New-Item -ItemType Directory -Path "src\model\db" -Force
goctlwt model mysql ddl --src sql\system.sql --dir src\model\db