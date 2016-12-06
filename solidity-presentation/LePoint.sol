contract LePoint {

    struct Campaign {
        address recipient;
        uint tippingPointWei;
        uint totalWei;
        address[] contributors;
        mapping (address => uint[]) contributions;
        bool tipped;
    }

    bytes32[] public campaignNames;
    mapping (bytes32 => Campaign) public campaigns;

    function createCampaign(bytes32 name, uint tippingPoint) {
        if (campaigns[name].recipient != 0) throw;

        campaignNames.push(name);
        
        address[] memory c;
        campaigns[name] = Campaign({
            recipient: msg.sender,
            totalWei: 0,
            tippingPointWei: tippingPoint,
            contributors: c,
            tipped: false
        });
    }
    
    function contributeTo(bytes32 name) {
        Campaign campaign = campaigns[name];
        campaign.totalWei += msg.value;
        campaign.contributions[msg.sender].push(msg.value);
        campaign.contributors.push(msg.sender);
        
        if (campaign.tipped)
            campaign.recipient.send(msg.value);
    }

    function checkCampaign(bytes32 name) {
        Campaign campaign = campaigns[name];
        if (campaign.tipped) return;
        
        if (campaign.totalWei >= campaign.tippingPointWei) {
            campaign.tipped = true;
            campaign.recipient.send(campaign.totalWei);
        }
    }

    function getCampaignNames() constant returns (bytes32[]) {
        return campaignNames;
    }

    function getCampaignInfo(bytes32 name) constant returns (address, bool, uint, uint) {
        Campaign c = campaigns[name];
        return (c.recipient, c.tipped, c.tippingPointWei, c.totalWei);
    }
}
