let fs=require("fs")
const ledger = require('./src');
let args = require('process.args') ('add', {
    s: 'startIndex',
    n: 'number'
})

function argsError() {
    console.error("error: invalid command")
    console.error("example: ")
    console.error("     node gcexserver_ETH_new_address.js add -s=1000 -n=500");
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

const gcexserver_ETH = 0;
const gcexserver_ETH_Start = gcexserver_ETH + startIndex;
const gcexserver_ETH_End = gcexserver_ETH_Start + number;

function getAddress(eth, addressIndex) {
    var path = "44'/60'/1'/" + addressIndex;
    eth.getAddress_async(path)
    .then(result => {
        sql = `INSERT INTO kc_address_pool(address, currency, address_index, flag) SELECT "${result["address"]}", "ETH", "${path}", 0 \
FROM DUAL WHERE NOT EXISTS(SELECT address FROM address_pool WHERE address = "${result["address"]}");\n`
        fs.writeFile('./address.sql', sql, {'flag': 'a'});
        console.log(result["address"], path);
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
let eth = new ledger.eth(comm);
fs.writeFile('./address.sql', "USE gcexserver;\n", {'flag': 'a'});
for(let i = gcexserver_ETH_Start; i < gcexserver_ETH_End; i ++) {
    getAddress(eth, i);
}
})
.catch(reason => {
    console.log('An error occured: ', reason);
});
