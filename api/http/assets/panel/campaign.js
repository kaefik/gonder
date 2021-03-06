// --- Campaign table ---
$('#campaign').w2grid({
    name: 'campaign',
    header: w2utils.lang('Campaign'),
    keyboard : false,
    show: {
        header: true,
        toolbar: true,
        footer: true,
        toolbarDelete: false,
        toolbarAdd: true,
        toolbarSave: true,
        toolbarSearch: true
    },
    columns: [
        { field: 'recid', caption: w2utils.lang('Id'), size: '50px', sortable: true, attr: "align=right" },
        { field: 'name', caption: w2utils.lang('Name'), size: '100%', sortable: true, editable: { type: 'text' } }
    ],
    multiSelect: false,
    sortData: [{ field: 'recid', direction: 'DESC' }],
    url: '/api/campaigns',
    method: 'GET',
    onAdd: function (event) {
        var id, name;
        $.ajax({
            type: "GET",
            //async: false,
            dataType: 'json',
            data: {"request": JSON.stringify({"cmd": "add", "id": parseInt(w2ui['group'].getSelection()[0])})},
            url: '/api/campaigns'
        }).done(function(data) {
            if (data['status'] == 'error') {
                w2alert(w2utils.lang(data["message"]), w2utils.lang('Error'));
            } else {
                id = data["recid"];
                name = data["name"];
                w2ui.campaign.add({recid: id, name: name}, true);
                w2ui.campaign.editField(id, 1);
            }
        });

    },
    onSelect: function (event) {
        w2ui.layout.lock('main', w2utils.lang('Loading...'), true);
        var record = this.get(event.recid);
        // --- Get campaign data ---
        $.ajax({
            type: "GET",
            url: '/api/campaign',
            dataType: 'json',
            data: {"request": JSON.stringify({"cmd": "get", "id": parseInt(event.recid)})}
        }).done(function(data) {
            zone = new Date().getTimezoneOffset() * 60;
            refreshProfilesList(data["profileId"]);
            refreshSenderList(data["senderId"]);
            $('#campaignId').val(record.recid);
            $('#campaignName').val(record.name);
            $("#campaignSubject").val(data["subject"]);
            $("#campaignStartDate").val(w2utils.formatDate((new Date((data["startDate"] + zone)* 1000 )), w2utils.settings.dateFormat));
            $("#campaignStartTime").val(w2utils.formatTime((new Date((data["startDate"] + zone) * 1000)), w2utils.settings.timeFormat));
            $("#campaignEndDate").val(w2utils.formatDate((new Date((data["endDate"] + zone) * 1000)), w2utils.settings.dateFormat));
            $("#campaignEndTime").val(w2utils.formatTime((new Date((data["endDate"] + zone) * 1000)), w2utils.settings.timeFormat));
            $("#campaignSendUnsubscribe").prop("checked", data["sendUnsubscribe"]);
            $("#campaignTemplate").val(data["template"]);
            $('#campaignAcceptSend').prop('checked', data["accepted"]);

            editor.setData(data["template"]);

            w2ui['recipient'].postData["campaign"] = parseInt(event.recid);
            w2ui.layout.unlock('main');
            w2ui['sidebar'].click('parameter');
        });
        // --- /Get campaign data ---
    },
    onSave: function(event) {
        console.log(event);
    }
});
// --- /Campaign table ---

function refreshProfilesList(selectedProfile){
    var profile;
    $.ajax({
        type: "GET",
        async: false,
        url: '/api/profilelist',
        data: {"request": JSON.stringify({"cmd": "get"})},
        dataType: "json"
    }).done(function(data) {
        profile = data;
    });
    w2ui['parameter'].set('campaignProfileId', { options: { items: profile } });
    w2ui['parameter'].record['campaignProfileId'] = selectedProfile;
    w2ui['parameter'].refresh();
}

function refreshSenderList(selectedSender){
    var sender;
    $.ajax({
        type: "GET",
        async: false,
        url: '/api/senderlist',
        dataType: "json",
        data: {"request": JSON.stringify({"cmd": "get", "id": parseInt(w2ui['group'].getSelection()[0])})},
    }).done(function(data) {
        sender = data;
    });
    w2ui['parameter'].set('campaignSenderId', { options: { items: sender } });
    w2ui['parameter'].record['campaignSenderId'] = selectedSender;
    w2ui['parameter'].refresh();
}
