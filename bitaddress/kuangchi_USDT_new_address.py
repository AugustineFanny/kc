#!/usr/bin/env python  
# _*_ coding:utf-8 _*_ 

import argparse
import subprocess
import json


def getnewaddress():
    command = './usdt_tools/omnicore-cli --rpcuser=shao --rpcpassword=yong getnewaddress "kuangchi"'
    p = subprocess.Popen(command, shell=True, stdout=subprocess.PIPE)
    out, err = p.communicate()
    return out.strip()


def getaddressesbyaccount():
    command = './usdt_tools/omnicore-cli --rpcuser=shao --rpcpassword=yong getaddressesbyaccount "kuangchi"'
    p = subprocess.Popen(command, shell=True, stdout=subprocess.PIPE)
    out, err = p.communicate()
    return out.strip()


def formataddresses(string):
    js = json.loads(string)
    with open("./address_usdt.sql", "a") as f:
        f.write("USE gcexserver;\n")
        for j in js:
            line = 'INSERT INTO kc_address_pool(address, currency, address_index, flag) SELECT "%s", "USDT", "%s", 0 FROM ' \
                   'DUAL WHERE NOT EXISTS(SELECT address FROM kc_address_pool WHERE address = "%s");\n' % (j, j, j)
            f.write(line)
            print(j)


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('command', help="getnewaddress: 生成新地址 ex: ./kuangchi_USDT_new_address.py getnewaddress -n 100\n"
                                        "getaddressesbyaccount:生成sql文件 ex: ./kuangchi_USDT_new_address.py getaddressesbyaccount")
    parser.add_argument('-n', '--num', type=int, default=1, help="command 为 getnewaddress 时 生成新地址数量 默认1")
    args = parser.parse_args()
    if args.command not in ['getnewaddress', 'getaddressesbyaccount']:
        print('ERROR: 可选command: getnewaddress, getaddressesbyaccount')
        return
    if args.command == "getnewaddress":
        for i in range(0, args.num):
            print(getnewaddress())
        return
    if args.command == "getaddressesbyaccount":
        addresses_str = getaddressesbyaccount()
        formataddresses(addresses_str)
        return


if __name__ == '__main__':
    main()
