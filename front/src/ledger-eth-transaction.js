import request from 'superagent';
import Transport from "@ledgerhq/hw-transport-u2f"; // for browser
import Eth from "@ledgerhq/hw-app-eth";
import Web3 from "web3";
import Transcation from "ethereumjs-tx";

const web3 = new Web3();
const utils = web3.utils;

const contractABI = [{"constant":false,"inputs":[{"name":"dst","type":"address"},{"name":"wad","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"payable":false,"type":"function"}]

function getAddress(path) {
    return new Promise((resolve, reject) => {
        Transport.create().then(res => {
            const eth = new Eth(res);
            eth.getAddress(path).then(res => {
                resolve(res.address);
            }).catch(err => reject(err))
        }).catch(err => reject(err))
    })
}

function getTransactionCount(address) {
    return new Promise((resolve, reject) => {
        request
        .get('https://api.etherscan.io/api')
        .query({module: "proxy", 
                action: "eth_getTransactionCount", 
                address: address, 
                tag: "latest", 
                apikey: "YourApiKeyToken"
        }).then(res => {
            if ("error" in res.body) {
                var message;
                if(typeof res.body.error === "string" || res.body.error instanceof String) {
                  message = res.body.error;
                } else {
                  message = res.body.error.message;
                }
                reject(message);
                return
            }
            resolve(res.body.result);
        }).catch(err => reject(err))
    })
}

function signTx(path, tx) {
    return new Promise((resolve, reject) => {
        Transport.create().then(res => {
            const eth = new Eth(res);
            eth.signTransaction(path, tx.serialize().toString('hex')).then(res => {
                tx.v = '0x' + res.v;
                tx.r = '0x' + res.r;
                tx.s = '0x' + res.s;
                var sigtx = new Transcation(tx)
                var sigtxHex = sigtx.serialize().toString('hex');
                resolve(sigtxHex);
            }).catch(err => reject(err))
        }).catch(err => reject(err))
    });
};

//sendETH("44'/60'/0'/1", "0x99B677A1b4b326d2eed60a8DB726b4E1869C5F17", "0.006556", "40", "21000")
function sendETH (path, target, amount, gasPrice, gasLimit) {
    return new Promise((resolve, reject) => {
        getAddress(path).then(address => {
            getTransactionCount(address).then(nonce => {
                var rawTx = {
                    "nonce": nonce,
                    "gasPrice": utils.toHex(utils.toWei(gasPrice, "gwei")),
                    "gasLimit": utils.toHex(gasLimit),
                    "to": target,
                    "value": utils.toHex(utils.toWei(amount)),
                    "data": "",
                    "v": "0x01",
                    "chainId":1,
                }
                var tx = new Transcation(rawTx);
                signTx(path, tx).then(res => resolve(res)).catch(err => reject(err))
            }).catch(err => reject(err))
        }).catch(err => reject(err))
    })
};

function _fromDecimal (value, decimal) {
    var amountToSendinDecimal = value * (10 ** decimal);
    return utils.numberToHex(amountToSendinDecimal);
};

function sendToken (path, target, amount, gasPrice, gasLimit, contractAddress, decimals) {
    return new Promise((resolve, reject) => {
        getAddress(path).then(address => {
            getTransactionCount(address).then(nonce => {
                var contract = new web3.eth.Contract(contractABI, contractAddress, address);
                var rawTx = {
                    "nonce": nonce,
                    "gasPrice": utils.toHex(utils.toWei(gasPrice, "gwei")),
                    "gasLimit": utils.toHex(gasLimit),
                    "to": contractAddress,
                    "value": utils.toHex(0),
                    "data": contract.methods.transfer(target, _fromDecimal(amount, decimals)).encodeABI(),
                    "v": "0x01",
                    "chainId":1,
                }
                var tx = new Transcation(rawTx);
                signTx(path, tx).then(res => resolve(res)).catch(err => reject(err))
            }).catch(err => reject(err))
        }).catch(err => reject(err))
    })
};

export default {
    sendTransaction: (path, target, amount, gasPrice, gasLimit, contractAddress, decimals) => {
        if (contractAddress === "") {
            return sendETH(path, target, amount, gasPrice, gasLimit);
        } else {
            return sendToken(path, target, amount, gasPrice, gasLimit, contractAddress, decimals);
        }
    },
    getBalance: (address) => {
        console.log(utils);
        return new Promise((resolve, reject) => {
            request
            .get('https://api.etherscan.io/api')
            .query({module: "account", 
                    action: "balance", 
                    address: address, 
                    tag: "latest", 
                    apikey: "YourApiKeyToken"
            }).then(res => {
                if ("error" in res.body) {
                    var message;
                    if(typeof res.body.error === "string" || res.body.error instanceof String) {
                      message = res.body.error;
                    } else {
                      message = res.body.error.message;
                    }
                    reject(message);
                    return
                }
                resolve(res.body.result);
            }).catch(err => reject(err))
        })
    },
    fromWei: utils.fromWei,
    allValue: (valueHex, gasPrice, gasLimit) => {
        var value = utils.fromWei(valueHex, "wei") - utils.toWei(gasPrice, "gwei") * parseInt(gasLimit);
        return utils.fromWei(utils.numberToHex(value));
    }
}

