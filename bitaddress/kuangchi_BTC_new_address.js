let fs=require("fs")
const ledger = require('./src');
let args = require('process.args') ('add', {
    s: 'startIndex',
    n: 'number'
})

function argsError() {
    console.error("error: invalid command")
    console.error("example: ")
    console.error("     node gcexserver_BTC_new_address.js add -s=1000 -n=500");
    console.error("     -s:起始index,大于0");
    console.error("     -n:数量number,大于0");
}

if(!('startIndex' in args) || !('number' in args)) {
    argsError();
    return
}

let startIndex = parseInt(args.startIndex);
let number = parseInt(args.number);

if(isNaN(startIndex) || isNaN(number) || startIndex < 0 || number < 1) {
    argsError();
    return
}

const gcexserver_BTC = 0;
const gcexserver_BTC_Start = gcexserver_BTC + startIndex;
const gcexserver_BTC_End = gcexserver_BTC_Start + number;

function getAddress(btc, addressIndex) {
    var path = "44'/0'/1'/0/" + addressIndex;
    btc.getWalletPublicKey_async(path)
    .then(result => {
        sql = `INSERT INTO kc_address_pool(address, currency, address_index, flag) SELECT "${result["bitcoinAddress"]}", "BTC", "${path}", 0 \
    FROM DUAL WHERE NOT EXISTS(SELECT address FROM address_pool WHERE address = "${result["bitcoinAddress"]}");\n`
        fs.writeFile('./address.sql', sql, {'flag': 'a'});
        console.log(result["bitcoinAddress"], path);
    })
    .catch(error => {
            console.log(error);
    })
}

ledger
    .comm_node
    .create_async()
    .then(comm => {
    console.log(comm.device.getDeviceInfo());
var btc = new ledger.btc(comm);
fs.writeFile('./address.sql', "USE gcexserver;\n", {'flag': 'a'});
for(var i = gcexserver_BTC_Start; i < gcexserver_BTC_End; i ++) {
    getAddress(btc, i);
}
})
.catch(reason => {
    console.log('An error occured: ', reason);
});
