init()

function init() {
     $.ajax({
         type : 'GET',
         url : "http://explorer.arkgallery.net/api/delegates/getActive",
         success : function (data) {
            delegates = data.delegates;
            for(var i = 0; i < delegates.length; i++) {
                if(delegates[i].username == "arkgallery"){
                    console.log(delegates[i]);
                    showChart(delegates[i].producedblocks, delegates[i].missedblocks);

                    $('#rate').text(delegates[i].rate +'위');
                    $('#totalvote').text(delegates[i].vote / 100000000 +'Ѧ');
                    $('#forged').text(delegates[i].forged / 100000000 +'Ѧ');
                }
            }
         }
     });
}

function showChart(produced, missed) {
    var ctxP = document.getElementById("pieChart").getContext('2d');
    var myPieChart = new Chart(ctxP, {
        type: 'pie',
        data: {
            labels: ["생성한 블록", "놓친 블록"],
            datasets: [
                {
                    data: [produced, missed],
                    backgroundColor: [ "#46BFBD", "#F7464A"],
                    hoverBackgroundColor: [ "#5AD3D1", "#FF5A5E"]
                }
            ]
        },
        options: {
            responsive: true
        }
    });
}

