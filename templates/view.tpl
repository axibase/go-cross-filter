<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <script type="text/javascript" src="//apps.axibase.com/chartlab/portal/JavaScript/d3.min.js"></script>
    <script type="text/javascript" src="//apps.axibase.com/chartlab/portal/JavaScript/portal_init.js"></script>
    <script type="text/javascript" src="//apps.axibase.com/chartlab/portal/JavaScript/charts.min.js"></script>
    <script type="text/javascript" src="//apps.axibase.com/chartlab/portal/JavaScript/initialize.js"></script>
    <link rel="stylesheet" type="text/css" href="//apps.axibase.com/chartlab/portal/CSS/charts.min.css"/>
    <link rel="shortcut icon" href="//axibase.com/favicon.ico" />
    <style>
        body {
            padding: 30px;
        }
        h1 {
            font-size: 18px;
            font-weight: bold;
        }
        .widget-chart {
            width: 90%;

        }
        .axi-tooltip-container {
            z-index: -1
        }

    </style>
</head>
<body onload="onload()">
</body>
<script>
    function onload() {
        onBodyLoad();
        loadWidgets('chartsConfig?{{.Query}}', function (widgetConfigs) {
            widgetConfigs.forEach(function(el,i) {
                var keys = d3.keys(el.series[0].tags);
                var text = "";
                if (el.series[0].tags != undefined) {
                    text = ": ";
                    for (var j = 0; j < keys.length; ++j) {
                        text += keys[j] + "=" + el.series[0].tags[keys[j]];
                        if (j != keys.length-1) {
                            text += ","
                        }
                    }
                }
                var title = "";
                if (el.title != undefined) {
                    title = el.title
                } else {
                    el.series.forEach(function(ser, i) {
                        if (i != 0) {
                            title += ", "
                        }
                        title += ser.metric
                    })
                }
                d3.select("body").append("h1").text(title + text);
                d3.select("body").append("div").attr("id", function() {return ""+i});
                el.path = "api/v1/series?{{.Query}}";
                updateWidget(el,""+i);
            });
        });
    }
</script>
</html>