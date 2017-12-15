$(function () {

    /**
     * Check a selector exist or not
     */
    $.fn.exists = function () {
        return this.length !== 0;
    }

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
    var isLoading = false;
    timer.init("#timer");

    function axiosDone(selector) {
        selector.removeClass("disabled");
        isLoading = false;
        timer.stop();
    }

    /* Fake Data / Seed Table / Drop Table: Action */
    $(".btn-init-data").click(function() {
        if (isLoading) {
            return false;
        }
        $(this).addClass("disabled");
        $(".msg").text(null);
        isLoading = true;
        selector = $(this);
        url = selector.data("path");
        timer.start();
        axios.post(url, null)
            .then(function (response) {
                $(".msg").text(response.data.msg);
                axiosDone(selector);
            })
            .catch(function (error) {
                $(".msg").text(error);
                axiosDone(selector);
            });
    })

    Highcharts.setOptions({
        lang: {
            thousandsSep: ','
        }
    })

    if ($('#user-container').exists()) {
        var userChart = Highcharts.chart('user-container', {
            chart: {
                type: 'column'
            },
            credits: {
                enabled: false
            },
            title: {
                text: 'New users by month'
            },
            xAxis: {
                title: {
                    text: null
                },
                categories: newUserByMonth.xCategories
            },
            yAxis: {
                title: {
                    text: 'Users'
                }
            },
            series: [{
                showInLegend: false,
                name: 'Users',
                data: newUserByMonth.series
            }]
        });
    }

    if ($('#access-container').exists()) {
        var accessChart = Highcharts.chart('access-container', {
            chart: {
                plotBackgroundColor: null,
                plotBorderWidth: null,
                plotShadow: false,
                type: 'pie'
            },
            credits: {
                enabled: false
            },
            title: {
                text: 'Total access by country'
            },
            tooltip: {
                pointFormat: '{series.name}: <b>{point.percentage:.1f}%</b>'
            },
            plotOptions: {
                pie: {
                    allowPointSelect: true,
                    cursor: 'pointer',
                    dataLabels: {
                        enabled: true,
                        format: '<b>{point.name}</b>: {point.percentage:.1f} %',
                        style: {
                            color: (Highcharts.theme && Highcharts.theme.contrastTextColor) || 'black'
                        }
                    }
                }
            },
            series: [{
                name: 'Access',
                colorByPoint: true,
                data: accessByCountry.series
            }]
        });
    }

    if ($('#retention-container').exists()) {
        var accessChart = Highcharts.chart('retention-container', {
            credits: {
                enabled: false
            },
            title: {
                text: 'Retention rate'
            },
            xAxis: {
                categories: retention.xCategories
            },
            yAxis: {
                title: {
                    text: 'Number service access users'
                }
            },
            plotOptions: {
                line: {
                    dataLabels: {
                        enabled: true
                    },
                    enableMouseTracking: false
                }
            },
            series: [{
                type: 'column',
                name: 'Past 30 days',
                data: retention.series.past30days
            }, {
                type: 'column',
                name: 'Comeback on this day',
                data: retention.series.days
            },
            {
                type: 'spline',
                name: 'Retention rate',
                data: retention.series.rate,
                tooltip: {
                    pointFormat: "<span style=\"color:{series.color}\">\u25CF</span> Retention rate: {point.y:,.1f}%"
                },
                marker: {
                    lineWidth: 2,
                    lineColor: Highcharts.getOptions().colors[3],
                    fillColor: 'white'
                }
            }]
        });
    }

});
