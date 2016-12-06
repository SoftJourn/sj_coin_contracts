var web3 = initializeWeb3();
var instance = constructInstance(web3);

listCampaigns();
setupCreationForm();

function initializeWeb3() {
  var provider = new Web3.providers.HttpProvider("http://localhost:8545");
  var web3 = new Web3(provider);
  web3.eth.getCoinbase(function(err, coinbase) {
    if (err) return setNotice(err);
    web3.eth.defaultAccount = coinbase;
  });
  return web3;
}

function constructInstance() {
  // copied abi+address from the output of `node deploy.js`
  var abi = [{"constant":false,"inputs":[{"name":"name","type":"bytes32"}],"name":"checkCampaign","outputs":[],"type":"function"},{"constant":false,"inputs":[{"name":"name","type":"bytes32"}],"name":"contributeTo","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"","type":"uint256"}],"name":"campaignNames","outputs":[{"name":"","type":"bytes32"}],"type":"function"},{"constant":true,"inputs":[],"name":"getCampaignNames","outputs":[{"name":"","type":"bytes32[]"}],"type":"function"},{"constant":true,"inputs":[{"name":"","type":"bytes32"}],"name":"campaigns","outputs":[{"name":"recipient","type":"address"},{"name":"tippingPointWei","type":"uint256"},{"name":"totalWei","type":"uint256"},{"name":"tipped","type":"bool"}],"type":"function"},{"constant":false,"inputs":[{"name":"name","type":"bytes32"},{"name":"tippingPoint","type":"uint256"}],"name":"createCampaign","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"name","type":"bytes32"}],"name":"getCampaignInfo","outputs":[{"name":"","type":"address"},{"name":"","type":"bool"},{"name":"","type":"uint256"},{"name":"","type":"uint256"}],"type":"function"}];
  // morden
  var address = "0x403d965afb3a0f2b1e8b71cdfece31ceb6808ba7";

  return web3.eth.contract(abi).at(address);
}

function listCampaigns() {
  var campaignList = document.getElementById("campaigns");
  campaignList.innerHTML = "";
  instance.getCampaignNames(function(err, names) {
    if (err) return setNotice(err);
    names.forEach(function(rawName) {
      var name = stripNulls(web3.toAscii(rawName));
      instance.getCampaignInfo(name, function(err, rawInfo) {
        if (err) return setNotice(err);
        var contributionForm = "<form style='display: inline' onsubmit='return contribute(this, \""+name+"\")'><input name='wei' /><input type='submit' value='Send Wei' /></form>";
        var lePoint = rawInfo[2];
        var totalWei = rawInfo[3];
        var infoSpan = "<span>total wei contributed: "+totalWei+", total needed: "+lePoint+"</span>";
        var li = document.createElement("li");
        li.innerHTML = name+" "+contributionForm+" "+infoSpan;
        campaignList.appendChild(li);
      });
    });
  });
}

function setupCreationForm() {
  var form = document.getElementById("new-campaign");

  form.onsubmit = function() {
    try {
      setNotice("Campaign data sent. Waiting for confirmation...");
      var campaignName = form.elements.name.value;
      var tippingPoint = form.elements["tipping-point"].value;
      form.reset();
      instance.createCampaign(campaignName, tippingPoint, {gas: 2000000}, function(err, txHash) {
        if (err) return setNotice(err);
        handleTx(campaignName+" saved to blockchain.", txHash);
      });
    } catch(e) {
      setNotice(e);
    }
    return false;
  };
}

function contribute(form, campaignName) {
  try {
    var wei = form.elements.wei.value;
    setNotice("Sending "+wei+" wei to the "+campaignName+" campaign. Waiting for confirmation...");
    form.reset();
    instance.contributeTo(campaignName, {value: wei, gas: 2000000}, function(err, txHash) {
      if (err) return setNotice(err);
      handleTx("Contribution to "+campaignName+" saved to blockchain.", txHash);
    });
  } catch (e) {
    setNotice(e);
  }
  return false;
}

function handleTx(successMessage, txHash) {
  var interval = setInterval(function() {
    web3.eth.getTransactionReceipt(txHash, function(err, receipt) {
      if (err != null) {
        clearInterval(interval);
        setNotice(err);
      }
      if (receipt != null) {
        clearInterval(interval);
        listCampaigns();
        setNotice(successMessage);
        setTimeout(function() { setNotice(""); }, 5000);
      }
    });
  }, 1000);
}

function setNotice(message) {
  document.getElementById("notice").innerText = message;
}

function stripNulls(s) {
  return s.replace(/\0/g, '');
}