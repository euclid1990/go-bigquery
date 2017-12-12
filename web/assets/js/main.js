$(function () {
    // var myChart = Highcharts.chart('container', {
    //     chart: {
    //         type: 'bar'
    //     },
    //     title: {
    //         text: 'Fruit Consumption'
    //     },
    //     xAxis: {
    //         categories: ['Apples', 'Bananas', 'Oranges']
    //     },
    //     yAxis: {
    //         title: {
    //             text: 'Fruit eaten'
    //         }
    //     },
    //     series: [{
    //         name: 'Jane',
    //         data: [1, 0, 4]
    //     }, {
    //         name: 'John',
    //         data: [5, 7, 3]
    //     }]
    // });

    function accuracy() {
        return {
            init: function(timerElementId) {
                this.timerElementId = timerElementId;
                this.startTime = null;
                this.interval = null;
            },
            start: function() {
                $(this.timerElementId).text((0).toFixed(3));
                this.startTime = Date.now();
                this.interval = setInterval(function() {
                    var elapsedTime = Date.now() - this.startTime;
                    $(this.timerElementId).text((elapsedTime / 1000).toFixed(3));
                }.bind(this), 90);
            },
            stop: function() {
                clearInterval(this.interval);
            }
        };
    };

    var timer = accuracy();
    timer.init("#timer");
    timer.start();
});
