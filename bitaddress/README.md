# ledger nano S

## 依赖

    npm install

## 导出kuangchi BTC 地址

    node kuangchi_BTC_new_address.js add -s=1 -n=100

    -s: startIndex 开始索引
    -n: number  导出数量

## 导出kuangchi ETH 地址

    node kuangchi_ETH_new_address.js add -s=1 -n=100

    -s: startIndex 开始索引
    -n: number  导出数量

## 导出kuangchi USDT 地址

    启动omni节点 ./omnicored -txindex=1 -datadir=//Users/ma1/study/omni -server=1 -rpcuser=shao -rpcpassword=yong
    生成地址 ./kuangchi_USDT_new_address.py getnewaddress -n 100
    导出地址 ./kuangchi_USDT_new_address.py getaddressesbyaccount

    -n: number  导出数量

## 将地址导入数据库
address.txt中包含地址
使用python脚本

    pip install mysql-python
    python input_addr.py
