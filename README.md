# diceBlacklist
 
一个为 TRPG 骰子机器人服务的非官方云端黑名单，目前主要支持 [SealDice](https://github.com/sealdice/sealdice-core)。

## 配置

运行后会在程序目录下创建 blacklist.sqlite.db 和 appid.json，前者是黑名单数据库，后者用于存放客户端 ID。

同时，每 24 小时会进行一次备份，备份将存放于 backups 文件夹（没有则会新建）。

## 杂项

目录中的 migrate.py 用于将 SealDice 导出的黑名单（JSON 格式）写入数据库。Dockerfile 和 build_script.sh 用于编译 Windows 可执行文件。

## API

// TODO

