<!DOCTYPE html>
<html lang="en">
<head>
    <meta http-equiv="content-type"  content="text/html; charset=UTF8">
    <title>Cross-filter</title>
    <script src="//code.jquery.com/jquery-2.1.4.min.js"></script>

    <script src="//cdnjs.cloudflare.com/ajax/libs/d3/3.5.6/d3.min.js"></script>

    <script src="//cdnjs.cloudflare.com/ajax/libs/crossfilter/1.3.12/crossfilter.min.js"></script>
    <script src="//maxcdn.bootstrapcdn.com/bootstrap/3.3.5/js/bootstrap.min.js"></script>
    <script src="//cdn.datatables.net/1.10.9/js/jquery.dataTables.min.js"></script>
    <script src="//cdn.datatables.net/1.10.9/js/dataTables.bootstrap.min.js"></script>
    <script src="//cdn.datatables.net/scroller/1.3.0/js/dataTables.scroller.min.js"></script>
    <script src="//apps.axibase.com/chartlab/portal/JavaScript/portal_init.js"></script>
    <script src="//apps.axibase.com/chartlab/portal/JavaScript/charts.min.js"></script>
    <script src="//apps.axibase.com/chartlab/portal/JavaScript/initialize.js"></script>
    <link rel="shortcut icon" href="//axibase.com/favicon.ico" />
    <link rel="stylesheet" href="//maxcdn.bootstrapcdn.com/bootstrap/3.3.5/css/bootstrap.min.css" type="text/css">
    <link href="//cdn.datatables.net/1.10.9/css/dataTables.bootstrap.min.css" rel="stylesheet" type="text/css">
    <link rel="stylesheet" type="text/css" href="//apps.axibase.com/chartlab/portal/CSS/charts.min.css">
    <link rel="stylesheet" href="//cdn.datatables.net/scroller/1.3.0/css/scroller.dataTables.min.css" type="text/css">
    <style>
        .page-header__title button {
            margin-left: 30px
        }
        .page-header > div {
            display: inline-block;
            vertical-align: middle;
            float: none;
        }
        .menu .scrollable-menu {
            height: auto;
            max-height: 400px;
            overflow-x: hidden;
        }
        .menu > * {
            display: inline-block;
            float: right;
            margin: 5px;
        }


        #data-table_wrapper {
            margin-right: -15px;
            margin-left: -15px;
        }
        #data-table_wrapper div div {
            min-height: 0;
        }

        table.dataTable {
            margin-top: 0 !important;
            margin-bottom: 0 !important;
            border: 0;
            font-size: 13px;
        }
        #data-table th {
            border-bottom: 0;
        }
        #data-table-container .panel-body {
            padding-top: 0;
            padding-bottom: 0;
        }
        #data-table td.td-body-right {
            text-align: right;
        }
        #data-table td.td-body-center {
            text-align: center;
        }
        #pie-row > * {
            text-align: center;
        }
        #hist-row > * {
            text-align: center;
        }
        .table-bordered>thead>tr>th {
            border-bottom: none;
        }
        .dataTables_scrollBody {
            border-top: 2px solid #ddd;
        }

        .panel-block_hidden {
            display: none
        }
        .reset-btn_hidden {
            display: none
        }

        .reset_hidden {
            display: none
        }
        .widget-pie {
            padding: 0;
            width: 280px;
            height: 300px;
        }
        .widget-histogram {
            padding: 0;
        }
    </style>
</head>
<body onload="onLoad()">
<div class="container">
    <div class="row page-header">
        <div class="col-xs-6 page-header__title">
            <h3><span>Cross-filter</span><button class="btn btn-default btn-xs reset-btn reset-btn_hidden" >Reset All</button></h3>
        </div><div class="clearfix visible-xs-block"></div><!--
        --><div class="col-xs-12 col-sm-6 menu">
        <div class="dropdown menu__entity-group-list menu-button">
            <button class="btn btn-default dropdown-toggle" type="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="true">
                <span class="menu-button__text">Entity group</span>
                <span class="menu-button__caret"></span>
            </button>
            <ul class="dropdown-menu scrollable-menu">
            </ul>
        </div>
        <div class="dropdown menu__tables-list menu-button">
            <button class="btn btn-default dropdown-toggle" type="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="true">
                <span class="menu-button__text">Select Dataset</span>
                <span class="menu-button__caret"></span>
            </button>
            <ul class="dropdown-menu scrollable-menu">
            </ul>
        </div>
    </div>

    </div>
    <div class="row panel-block panel-block_hidden panel-block_collapsible">
        <div class="panel panel-default col-xs-12">
            <div class="panel-heading row" role="button">
                <h3 class="panel-title panel-block__title"><span class="glyphicon glyphicon-filter" style="font-size: small" aria-hidden="true"></span> Tag filters</h3>
            </div>
            <div class="panel-body row collapse in panel-block__body" id="pie-row">

            </div>
        </div>
    </div>
    <div class="row panel-block panel-block_hidden panel-block_collapsible">
        <div class="panel panel-default col-xs-12">
            <div class="panel-heading collapsible row" role="button">
                <h3 class="panel-title panel-block__title"><span class="glyphicon glyphicon-filter" style="font-size: small" aria-hidden="true"></span> Metric filters</h3>
            </div>
            <div class="panel-body row collapse in panel-block__body" id="hist-row">

            </div>
        </div>
    </div>
    <div class="row panel-block panel-block_hidden">
        <div class="panel panel-default col-xs-12" id="data-table-container">
            <div class="panel-heading row">
                <h3 class="panel-title panel-block__title">
                    Records <button class="btn btn-default btn-xs reset-btn" style="margin-left: 30px">Reset All</button>
                    <span class="table-update-time" style="float: right"></span>
                </h3>
            </div>
            <div class="panel-body row panel-block__body">
                <table class="table table-bordered table-striped table-condensed" width="100%" id="data-table">
                </table>
            </div>

        </div>
    </div>
</div>
<script>
    screener = window.screener || {};
    screener.PieChart = (function (axibaseCharts) {
        function PieChart(config) {
            var chart = Object.create(PieChart.prototype);
            axibaseCharts.EventListener(chart);

            var widgetConfig = {
                type: 'pie',
                script: 'config.url = null',
                timezone: 'UTC',
                selectormode: "highlight",
                legendposition: "none",
                onseriesclick: function(seriesConfig, widget, config) {
                    chart.fire("click", seriesConfig.label);
                    var slice = widget.getSlice(this);
                    widget.selectSlice(slice);
                },
                onseriesdoubleclick: function() {
                    chart.fire("doubleclick", arguments);
                }
            };

            chart._widget = updateWidgetConfig(config.id, widgetConfig, chart);
            chart._id = config.id;
            chart._widgetConfig = chart._widget.config;
            return chart;
        }
        PieChart.prototype.setKeys = function(keys) {
            var chart = this;
            chart._widgetConfig.series = keys.map(function (key) {
                return {
                    label: key,
                    entity: "e",
                    metric: key
                }
            });

            this._widget= updateWidgetConfig(chart._id,  chart._widgetConfig, chart);
        };
        PieChart.prototype.destroy = function() {
            var handlers = this.on();
            d3.keys(handlers).forEach(function (key) {
                delete handlers[key];
            });
            this._widget.destroy();
        };
        PieChart.prototype.setData = function(data) {
            var series = data.map(function(el){
                return {
                    label: el.key,
                    entity: "e",
                    metric: el.key,
                    data: {t:[0], v: [el.value], len: 1}
                }
            });

            this._widget.loader.defaultURL({ full: "" }).loadMethod(function(url, callback){
                setTimeout(function(){
                    callback(series);
                });
                return {};
            });
            this._widget.reload();
        };

        PieChart.prototype.getSelectionStatus = function() {
            return this._widget.getSlices().map(function(slice){ return {key: slice.getTitle(), selected: slice.isSelected};});
        };
        PieChart.prototype.setSelectionStatus = function(status) {
            var selectedSlices = this._widget.getSlices().filter(function(slice){ return status[slice.getTitle()];});
            this._widget.selectSlices(selectedSlices);
        };

        PieChart.prototype.resetSelection = function () {
            this._widget.selectSlices([]);
        };

        function updateWidgetConfig(id, config, chart) {
            var widget = axibaseCharts.updateWidget(config, id);
            widget.on("selectslices", function (slices) {
                chart.fire("selectslices", slices);
            });
            return widget;
        }

        return PieChart
    }(window));
</script>
<script>
    screener = window.screener || {};
    screener.HistChart = (function (axibaseCharts) {
        function HistogramChart(config) {
            var chart = Object.create(HistogramChart.prototype);
            axibaseCharts.EventListener(chart);

            var widgetConfig = {
                type: 'histogram',
                script: 'config.url = null',
                timezone: 'UTC',
                barcount: 20,
                minrangeforce: config.minrangeforce,
                maxrangeforce: config.maxrangeforce,
                rangeoffset: config.rangeoffset,
                legendposition: "none",
                series: [{
                    entity: "e",
                    metric: "m"
                }]
            };
            chart._widget= axibaseCharts.updateWidget(widgetConfig, config.id);
            chart._widget.addRangeSelector()
                    .on('rangeselectstart', function(min, max){
                        chart.fire("rangeselectstart", min, max);
                    })
                    .on('rangeselectend', function (min, max) {
                        chart.fire("rangeselectend", min, max);
                    })
                    .on("rangeselectexit", function () {
                        chart.fire("rangeselectexit");
                    });
            chart._id = config.id;
            chart._widgetConfig = widgetConfig;
            return chart
        }
        HistogramChart.prototype.destroy = function() {
            var handlers = this.on();
            d3.keys(handlers).forEach(function (key) {
                delete handlers[key];
            });

            this._widget.destroy();
        };

        HistogramChart.prototype.setData = function (data) {
            var t = [], v = [];
            data.forEach(function (el,i) {
                t[i] = el.t;
                v[i] = el.v;
            });
            var series = [{
                entity: "e",
                metric: "m",
                data: {t:t, v: v, len: t.length}
            }];
            this._widget.loader.defaultURL({ full: "" }).loadMethod(function(url, callback) {
                setTimeout(function() {
                    callback(series);
                }, 1);
                return {};
            });
            this._widget.reload();
            this._widget.loader(2);
        };
        HistogramChart.prototype.resetSelection = function () {
            this._widget.resetRangeSelector();
            this._widget.zoom.fire("reset", this._widget.axis[0].scale());
        };

        return HistogramChart;
    }(window));
</script>
<script>
    var pageHeader = (function () {

        var onEntityGroupSelect = function(){};
        var onDatasetSelect = function(){};
        return {
            setTitle: function (text) {
                d3.select(".page-header .page-header__title span").text(text);
            },
            setEntityGroupsList: function(entityGroups) {
                entityGroups = ["all"].concat(entityGroups);
                var list = d3.select(".menu__entity-group-list ul").selectAll("li").data(entityGroups, function(x) {return x}).attr();
                list.enter().append("li").append("a")
                        .text(function(d){return d})
                        .attr("role", function(d){ return "button"})
                        .on("click", function(d) {
                            d3.event.preventDefault();
                            clickEntityGroupsListHandler(d);
                        });
                list.exit().remove();
            },
            setTablesList: function(tables) {
                var list = d3.select(".menu__tables-list ul").selectAll("li").data(tables, function(x) {return x});
                list.enter().append("li").append("a")
                        .text(function(d){return d;})
                        .attr("role", function(d){ return "button"})
                        .on("click", function(d) {
                            d3.event.preventDefault();
                            clickTablesListHandler(d);
                        });
                list.exit().remove();
            },
            onEntityGroupSelect: function(handler) {
                onEntityGroupSelect = handler;
            },
            onDatasetSelect: function(handler) {
                onDatasetSelect = handler;
            }
        };
        function clickTablesListHandler(d){
            d3.select(".menu__tables-list .menu-button__text").text(d);
            var entityGroup = d3.select(".menu__entity-group-list .menu-button__text").text();
            entityGroup = (entityGroup === "Entity group" || entityGroup === "all") ? "" : entityGroup;
            onDatasetSelect(d, entityGroup);
        }
        function clickEntityGroupsListHandler(d){
            d3.select(".menu__entity-group-list .menu-button__text").text(d);
            var entityGroup = (d === "Entity group" || d === "all") ? "" : d;
            onEntityGroupSelect(d3.select(".menu__tables-list .menu-button__text").text(),entityGroup);
        }
    }());
</script>
<script>
    var loader = (function () {
        return {
            loadTablesList: function(callback) {
                d3.json("tablesList", function(error, json) {
                    if (error) return console.warn(error);
                    callback(json);
                });
            },
            loadEntityGroupsList: function(table, callback) {
                d3.json("entityGroupsList?table="+encodeQuery(table), function(error, json) {
                    if (error) return console.warn(error);
                    callback(json);
                });
            },
            loadTable: function (table, entityGroup, callback) {
                d3.xhr("table?table="+encodeQuery(table)+"&entityGroup="+encodeQuery(entityGroup), function(xmlHttpRequest){
                    var json = JSON.parse(xmlHttpRequest.response);
                    callback(json)
                });
            }
        }
    }());
</script>
<script>
    var filterManager = (function(){
        var manager = {};
        EventListener(manager);

        var ndx;
        var globalDim;
        var dims = {};
        var filters = {};

        function setData(data) {
            if (ndx && ndx.size() != 0) {
                throw Error("crossfilter data is not empty")
            }
            ndx = crossfilter(data);
            globalDim = ndx.dimension(function (d) { return d;});

        }

        function getDimData(key) {
            dims[key].filterAll();
            var data  = dims[key].bottom(Infinity);
            dims[key].filter(filters[key]);
            return data;
        }
        function getGroupByCountData(key) {
            return dims[key].group().top(Infinity);
        }
        function getData() {
            return globalDim.bottom(Infinity);
        }
        function setDim(key) {
            if (dims[key]) {
                dims[key].dispose();
            }
            return dims[key] = ndx.dimension(function(d){return d[key]});
        }
        function filterDim(key, filter) {
            if (filters[key] !== undefined) {
                dims[key].filterAll();
                delete filters[key];
            }
            if (filter !== null) {
                dims[key].filter(filter);
                filters[key] = filter;
            }
            manager.fire("filter", key);
        }
        function reset() {
            d3.keys(dims).forEach(function (key) {
                dims[key].filterAll();
            });
            globalDim.filterAll();
            manager.fire("filter", "all");
        }

        function clear() {
            var handlers = manager.on();
            d3.keys(handlers).forEach(function (key) {
                delete handlers[key];
            });
            d3.keys(dims).forEach(function (key) {
                dims[key].dispose();
                delete dims[key];
            });
            globalDim.dispose();
            ndx.remove();
            manager.fire("clear");
        }


        manager.reset = reset;
        manager.clear = clear;
        manager.getDimData = getDimData;
        manager.setDim = setDim;
        manager.getData = getData;
        manager.filterDim = filterDim;
        manager.setData = setData;
        manager.getGroupByCountData = getGroupByCountData;

        return manager;
    }());

    var pies = [], hists = [], table;

    var isFirst = true;

    function onLoad() {
        d3.selectAll('.panel-block.panel-block_collapsible .panel-heading').on('click', function () {
            $(this).siblings(".panel-block__body").collapse("toggle");
        });
        d3.selectAll(".reset-btn").on("click", function() {
            filterManager.reset();
        });

        pageHeader.onDatasetSelect(function(tableName) {
            loader.loadEntityGroupsList(tableName, pageHeader.setEntityGroupsList);
            d3.select(".menu__entity-group-list .menu-button__text").text("all");
            loader.loadTable(tableName, "", onLoadTableHandler);
            window.history.replaceState(null, null, "?table="+encodeQuery(tableName));

        });
        pageHeader.onEntityGroupSelect(function(tableName, entityGroupName) {
            loader.loadTable(tableName, entityGroupName, onLoadTableHandler);
            if (entityGroupName) {
                window.history.replaceState(null, null, "?table="+encodeQuery(tableName)+"&entityGroup="+encodeQuery(entityGroupName));
            } else {
                window.history.replaceState(null, null, "?table="+encodeQuery(tableName));
            }
        });
        loader.loadTablesList(function(tables) {
            pageHeader.setTablesList(tables);
            var table =  getQueryVariable("table");
            var entityGroup = getQueryVariable("entityGroup");

            if (table && tables.indexOf(table) !== -1) {
                loader.loadEntityGroupsList(table, function(entityGroups) {
                    if (entityGroup && entityGroups.indexOf(entityGroup) !== -1) {
                        d3.select(".menu__tables-list .menu-button__text").text(table);
                        d3.select(".menu__entity-group-list .menu-button__text").text(entityGroup);
                        loader.loadTable(table, entityGroup, onLoadTableHandler);
                    } else {
                        window.history.replaceState(null, null, "?table="+table);
                        d3.select(".menu__tables-list .menu-button__text").text(table);
                        loader.loadTable(table, "", onLoadTableHandler);
                    }
                    pageHeader.setEntityGroupsList(entityGroups);
                });
            } else {
                window.history.replaceState(null, null, window.location.pathname);
            }
        });
    }
    function getQueryVariable(variable) {
        var query = window.location.search.substring(1);
        var vars = query.split('&');
        for (var i = 0; i < vars.length; i++) {
            var pair = vars[i].split('=');
            if (decodeURIComponent(pair[0]) == variable) {
                return decodeURIComponent(pair[1].replace(/\+/g, '%20'));
            }
        }
    }

    function onLoadTableHandler(json) {
        d3.selectAll(".panel-block.panel-block_hidden").classed("panel-block_hidden", false);
        d3.selectAll(".reset-btn.reset-btn_hidden").classed("reset-btn_hidden", false);
        d3.selectAll('.panel-block.panel-block_collapsible .panel-heading').forEach(function(d){
            $(d).siblings(".panel-block__body").collapse("show");
        });
        $(".table-update-time").text(new Date(json.updateTime).toUTCString());
        if (!isFirst) {
            clearPage();
        }
        pageHeader.setTitle(json.name);

        if (json.rows.length > 0) {
            appendPortalLinkToColumnsAndRows(json.name, json.columns, json.rows);
            var data = arrayRowToObjectRow(json.columns, json.rows);

            isFirst = false;
            filterManager.setData(data);
            var i = 0, j = 0;
            json.columns.forEach(function (el) {
                if (el.filter !== undefined) {
                    if (el.numeric) {
                        hists.push({
                            id: "hist-" + i,
                            chartHeader: el.label,
                            vKey: el.label,
                            dimKey: el.label,
                            range: el.filter.range
                        });
                        i++;
                    } else {
                        pies.push({
                            id: "pie-" + j,
                            chartHeader: el.label,
                            dimKey: el.label
                        });
                        j++;
                    }
                    filterManager.setDim(el.label);
                }
            });
            createHistCharts(hists);
            createPieCharts(pies);

            var ignoreUpdate = false;
            pies.forEach(function (pie) {
                pie.widget.on("selectslices", function (slices) {
                    if (!ignoreUpdate) {
                        var selectedKeys = slices.filter(function (slice) {
                            return slice.isSelected;
                        }).map(function (slice) {
                            return slice.getTitle();
                        });
                        filterManager.filterDim(pie.dimKey, function (d) {
                            if (selectedKeys.length == 0) {
                                return true;
                            } else {
                                return selectedKeys.indexOf(d) != -1;
                            }
                        });
                    }
                });
                filterManager.on("filter", function (key) {
                    if (key === "all") {
                        ignoreUpdate = true;
                        pie.widget.resetSelection();
                        ignoreUpdate = false;
                    }
                    if (key !== pie.dimKey) {
                        pie.widget.setData(filterManager.getGroupByCountData(pie.dimKey));
                    }
                });
            });
            hists.forEach(function (hist) {
                hist.widget
                        .on("rangeselectend", function (min, max) {
                            filterManager.filterDim(hist.dimKey, [min, max]);
                        })
                        .on("rangeselectexit", function () {
                            if (!ignoreUpdate) {
                                filterManager.filterDim(hist.dimKey, null);
                            }
                        });
                filterManager.on("filter", function (key) {
                    if (key === "all") {
                        hist.widget.resetSelection();
                    }
                    if (key !== hist.dimKey) {
                        hist.widget.setData(convertToHistData(filterManager.getDimData(hist.dimKey), hist.vKey));
                    }
                });
            });

            var columns = createDatatableColumns(json.columns);
            //to render table after widgets data reload
            setTimeout(function () {
                table = setupTable(columns, data);

                filterManager.on("filter", function () {
                    table.widget.clear();
                    table.widget.rows.add(filterManager.getData());
                    table.widget.draw();
                });
            });
        }
        function clearPage() {
            pies.forEach(function (pie) {
                pie.widget.destroy();
            });
            hists.forEach(function (hist) {
                hist.widget.destroy();
            });
            pies.length = 0;
            hists.length = 0;

            if (table) {
                table.widget.clear();
                table.widget.destroy();
                table = undefined;
                $("#data-table").empty();
            }
            filterManager.clear();
            d3.select("#hist-row").selectAll("div").remove();
            d3.select("#pie-row").selectAll("div").remove();
        }
        function appendPortalLinkToColumnsAndRows(table, columns, rows) {
            if (columns.length > 0) {
                var link = "link";
                var keys = columns.map(function (col) {
                    return col.label;
                });

                while (keys.indexOf(link) !== -1) {
                    link = "_" + link;
                }
            }
            rows.forEach(function (el) {
                var row = {};
                columns.forEach(function (col, i) {
                    if (col._type !== "link") {
                        row[columns[i].label] = el[i]
                    }
                });
                el.unshift(createPortalLink(table, row));
            });
            columns.unshift({
                label: link,
                _type: "link"
            });
        }
        function createPortalLink(table, row) {
            var keys = d3.keys(row);
            var link = "charts?table="+encodeQuery(table);
            keys.forEach(function (key) {
                link+="&col="+encodeQuery(key+"\\:"+row[key]);
            });
            return link;
        }
        function arrayRowToObjectRow(columns, array) {
            var objects = [];
            array.forEach(function (el) {
                var obj = {};
                columns.forEach(function (col, i) {
                    if (col.numeric) {
                        obj[col.label] = +el[i];
                        obj[col.label] = (isNaN(obj[col.label])) ? -1: obj[col.label];
                    } else {
                        if (el[i] == "null") {
                            obj[col.label] = "-";
                        } else {
                            obj[col.label] = el[i];
                        }

                    }

                });
                objects.push(obj);
            });
            return objects;
        }
        function setupTable(columns, data) {
            table = {};
            table.widget = $("#data-table").DataTable({
                "paging": true,
                "lengthChange": false,
                "searching": false,
                "autoWidth": true,
                "info": false,
                "scrollCollapse": true,
                "scrollY":        500,
                "deferRender":    true,
                "scroller":       true,
                "data": data,
                "columns": columns,
                "order": [[ 1, "asc" ]]
            });
            return table;
        }
        function createDatatableColumns(responseColumns) {
            return responseColumns.map(function(el){
                return {
                    data: el.label.replace(/\./g, "\\."),
                    render: function( data, type, full, meta ) {
                        if (el.formatter) {
                            if (el.formatter.round !== undefined || el.formatter.multiplier !== undefined) {
                                var mult = (el.formatter.multiplier == undefined) ? 1 : el.formatter.multiplier;
                                if (el.formatter.round == undefined) {
                                    return data * mult;
                                } else {
                                    return (data * mult).toFixed(el.formatter.round);
                                }
                            }

                        } else if (el._type == "link") {
                            return '<a href="'+data+'" target="_blank"><span class="glyphicon glyphicon-share-alt"></span></a>'
                        }else {
                            return data;
                        }
                    },
                    orderable: (el._type == "link") ? false : true,
                    visible: !el.hidden,
                    title: (el._type == "link") ? "&nbsp;" : el.label.replace(/ /g, "&nbsp;"),
                    className: (el.numeric) ? "td-body-right": (el._type == "link") ? "td-body-center": ""
                }
            });
        }

        function createHistCharts(hists) {
            var histDivs = d3.select("#hist-row").selectAll("div").data(hists);
            var div = histDivs.enter().append("div")
                    .classed("col-xs-12 col-sm-12 col-md-12 col-lg-6", true)
                    .attr("id", function (d) { return d.id })
                    .append("div");
            div.append("strong").html(function(d){return d.chartHeader+" ";});
            div.append("a")
                    .attr("role", "button")
                    .text(function(d) {return "reset"})
                    .classed("reset", true).classed("reset_hidden", true);
            histDivs.exit().remove();
            //for each div setup hist widget
            hists.forEach(function (el, i) {
                var minrangeforce, maxrangeforce;
                if (el.range) {
                    minrangeforce = el.range[0];
                    maxrangeforce = el.range[1];
                }
                el.widget = screener.HistChart({id: el.id, rangeoffset: 15, minrangeforce: minrangeforce, maxrangeforce: maxrangeforce});
                d3.select("#"+el.id+" .reset").on("click", function() { el.widget.resetSelection(); });
                el.widget
                        .on("rangeselectend", function() {d3.select("#"+el.id+" .reset").classed("reset_hidden", false)})
                        .on("rangeselectexit", function() {d3.select("#"+el.id+" .reset").classed("reset_hidden", true)});
                el.widget.setData(convertToHistData(filterManager.getDimData(el.dimKey), el.vKey));
            });
        }
        function createPieCharts(pies) {
            var pieDivs = d3.select("#pie-row").selectAll("div").data(pies);
            var div = pieDivs.enter().append("div")
                    .classed("col-xs-6 col-sm-4 col-md-3 col-lg-3", true)
                    .attr("id", function (d) { return d.id })
                    .append("div");
            div.append("strong").html(function(d){return d.chartHeader+" ";});
            div.append("a")
                    .attr("role", "button")
                    .text(function(d) {return "reset"})
                    .classed("reset", true).classed("reset_hidden", true);

            pieDivs.exit().remove();
            //for each div setup pie widget
            pies.forEach(function (el, i) {
                el.widget = screener.PieChart({id: el.id});
                d3.select("#"+el.id+" .reset").on("click", function() { el.widget.resetSelection(); });
                el.widget
                        .on("selectslices", function(slices) {
                            if (slices.filter(function (slice) {return slice.isSelected;}).length === 0) {
                                d3.select("#"+el.id+" .reset").classed("reset_hidden", true)
                            } else {
                                d3.select("#"+el.id+" .reset").classed("reset_hidden", false)
                            }
                        });
                el.widget.setKeys(filterManager.getGroupByCountData(el.dimKey).map(function (el) { return el.key; }));
                el.widget.setData(filterManager.getGroupByCountData(el.dimKey));
            });
        }
    }
    function encodeQuery(string) {
        return encodeURIComponent(string).replace(/%20/g, '+');
    }

    function convertToHistData(data, vKey) {
        return data.map(function (d, i) {
            return {t: i, v: d[vKey]}
        });
    }
</script>
</body>
</html>
