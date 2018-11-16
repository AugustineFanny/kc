# coding:utf-8
from __future__ import unicode_literals

import os
import MySQLdb

HOST = 'localhost'
USERNAME = 'gcexserver'
PASSWORD = 'gcexserver'
DATABASE = 'gcexserver'

db = MySQLdb.connect(HOST, USERNAME, PASSWORD, DATABASE)
cursor = db.cursor()

def other():
    sql = ("INSERT INTO address_pool(address, currency, address_index, flag) SELECT '%s', '%s', '%s', '%s' "
           "FROM DUAL WHERE NOT EXISTS(SELECT address FROM address_pool WHERE address = '%s')")

    with open('address.txt', 'r') as f:
        for line in f.readlines():
            l = line.strip('\n').split(',')
            if len(l) == 3:
                addr, currency, index = l
                try:
                    flag = 2 if currency == "BTC" and int(index) % 10 == 0 else 0
                    cursor.execute(sql % (addr, currency, index, flag, addr))
                    db.commit()
                except Exception, e:
                    print str(e)
                    db.rollback()
                    return False
    return True

if __name__ == '__main__':
    if os.path.exists('address.txt'):
        if other():
            print '录入地址成功'
        else:
            print '录入地址失败'
    else:
        print '录入地址失败: 缺少address.txt'

    db.close()
