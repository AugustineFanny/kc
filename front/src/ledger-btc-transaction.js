import request from 'superagent'
import Transport from "@ledgerhq/hw-transport-u2f"; // for browser
import AppBtc from "@ledgerhq/hw-app-btc";
import bitcore from "bitcore-lib";

const signTx = async (raws, satoshis, target) => {
    var transport;
    try {
        transport = await Transport.create();  
    }catch(e) {
        alert("error: " + e);
    }

  const btc = new AppBtc(transport);

  var inputs = [];
  var associatedKeysets = [];
  for (var i in raws) {
    inputs.push([btc.splitTransaction(raws[i][0]), raws[i][1]]);
    associatedKeysets.push(raws[i][2]);
  }

  var output = bitcore.Transaction.Output({
        satoshis: satoshis,
        script: bitcore.Script.buildPublicKeyHashOut(target).toBuffer()
    }).toBufferWriter()
  var tx = {outputs: [{amount: output.bufs[0], script: output.bufs[2]}]};

  var outputScript = btc.serializeTransactionOutputs(tx).toString('hex');

  btc.createPaymentTransactionNew(inputs, associatedKeysets, undefined, outputScript)
  .then(async (res) => {
    const txid = await providers.broadcasting(res);
    alert(txid.txid);
  }).catch(error => {
    alert("error: " + error);
  });
};

const signTxHasChange = async (raws, satoshis, target, changeSatoshis, addr) => {
    var transport;
    try {
        transport = await Transport.create();  
    }catch(e) {
        alert("error: " + e);
    }

  const btc = new AppBtc(transport);

  var inputs = [];
  var associatedKeysets = [];
  for (var i in raws) {
    inputs.push([btc.splitTransaction(raws[i][0]), raws[i][1]]);
    associatedKeysets.push(raws[i][2]);
  }

  var output1 = bitcore.Transaction.Output({
        satoshis: satoshis,
        script: bitcore.Script.buildPublicKeyHashOut(target).toBuffer()
    }).toBufferWriter()
  var output2 = bitcore.Transaction.Output({
        satoshis: changeSatoshis,
        script: bitcore.Script.buildPublicKeyHashOut(addr).toBuffer()
    }).toBufferWriter()
  var tx = {outputs: [
    {amount: output1.bufs[0], script: output1.bufs[2]},
    {amount: output2.bufs[0], script: output2.bufs[2]},
    ]};

  var outputScript = btc.serializeTransactionOutputs(tx).toString('hex');

  btc.createPaymentTransactionNew(inputs, associatedKeysets, undefined, outputScript)
  .then(async (res) => {
    const txid = await providers.broadcasting(res);
    alert(txid.txid);
  }).catch(error => {
    alert("error: " + error);
  });
};

var providers = {
    unspent: function(addr, associatedKeyset) {
        return request.get('https://blockexplorer.com/api/addr/' + addr + '/utxo?noCache=1').send().then(function (res) {
            return res.body.map(function (e) {
                return {
                    txid: e.txid,
                    vout: e.vout,
                    satoshis: e.satoshis,
                    confirmations: e.confirmations,
                    associatedKeyset: associatedKeyset
                };
            });
        });
    },
    rawTx: function(txid, vout, associatedKeyset) {
        return request.get('https://blockexplorer.com/api/rawtx/' + txid).send().then(function (res) {
            return [res.body['rawtx'], vout, associatedKeyset];
        });
    },
    // feeName: "fastest" or "halfHour" or "hour"
    fee: function (feeName) {
        return request.get('https://bitcoinfees.earn.com/api/v1/fees/recommended').send().then(function (res) {
            return res.body[feeName + "Fee"];
        });
    },
    broadcasting: function (rawTx) {
        return request.post('https://blockexplorer.com/api/tx/send')
                      .send({'rawtx': rawTx})
                      .set({'Content-Type': 'application/x-www-form-urlencoded'})
                      .then(function (res) {
                        return res.body;
                      })
    }
}

function getTransactionSize (numInputs, numOutputs) {
    return numInputs*180 + numOutputs*34 + 10 + numInputs;
}

export default {
    sendTransaction: function(addrs, target) {
        var addrIterable = [providers.fee("fastest")];
        for (var i in addrs) {
            addrIterable.push(providers.unspent(addrs[i].addr, addrs[i].associatedKeyset));
        }
        Promise.all(addrIterable).then(function(res) {
            var feePerByte = 10;
            var inputs = [];
            var ninputs = 0;
            var availableSat = 0;
            for (var i = 1; i < res.length; i++) {
                var utxos = res[i];
                for (var j = 0; j < utxos.length; j++) {
                    var utxo = utxos[j];
                    if(utxo.confirmations >= 1) {
                        inputs.push([utxo.txid, utxo.vout, utxo.associatedKeyset]);
                        availableSat += utxo.satoshis;
                        ninputs++;
                    }
                }
            }
            var fee = getTransactionSize(ninputs, 1) * feePerByte;
            alert("转到地址: " + target + 
                  "\n转账金额: " + (availableSat - fee) / 100000000 +
                  "\n手续费: " + fee / 100000000);
            if (availableSat < 10000) {
                alert("金额过低");
                return
            }
            var inputIterable = [];
            for (var i in inputs) {
                inputIterable.push(providers.rawTx(inputs[i][0], inputs[i][1], inputs[i][2]))
            }
            Promise.all(inputIterable).then(function(res) {   
                signTx(res, availableSat - fee, target);
            }).catch(error => {
                alert("error: " + error);
            })
        }).catch(error => {
            alert("error: " + error);
        })
    },
    withdraw: function(addr, associatedKeyset, target, amount) {
        var satoshis = amount * 100000000;
        Promise.all([providers.unspent(addr, associatedKeyset)]).then(function(res) {
            utxos = res[0];
            var inputs = [];
            for (var i = 0; i < utxos.length; i++) {
                var utxo = utxos[i];
                if(utxo.confirmations >= 1) {
                    inputs.push([utxo.txid, utxo.vout, utxo.associatedKeyset]);
                    availableSat += utxo.satoshis;
                    ninputs++;
                }
            }

            var fee = getTransactionSize(ninputs, 1) * feePerByte;
            if (availableSat < satoshis + fee) {
                alert("btc不足");
                return
            }
            
            alert("转到地址: " + target + 
                  "\n转账金额: " + satoshis / 100000000 +
                  "\n手续费: " + fee / 100000000);
            if (availableSat < 10000) {
                return
            }
            var inputIterable = [];
            for (var i in inputs) {
                inputIterable.push(providers.rawTx(inputs[i][0], inputs[i][1], inputs[i][2]))
            }
            Promise.all(inputIterable).then(function(res) {
                signTxHasChange(res, satoshis, target, availableSat - satoshis - fee, addr);
            }).catch(error => {
                alert("error: " + error);
            })
        }).catch(error => {
            alert("error: " + error);
        })
    }
};
