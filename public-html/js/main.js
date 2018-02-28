var currentLocation;

function initMap() {
	var map = new google.maps.Map($("#map")[0], {
	  zoom: 2,
	  center: {lat: 0, lng: 0}
	});
	$.getJSON("/locations", function(locations) {
		$.each(locations, function(key, loc) {
			var marker = new google.maps.Marker({
			  position: {lat: loc.lat, lng: loc.long},
			  map: map
			});
			marker.addListener("click", function() {
				displayLocationDetails(loc);
			});
		});
	});
}

function displayLocationDetails(loc) {
	currentLocation = loc;
	$.getJSON("/logs/" + loc.id, function(logs) {
		if (logs.length > 0) {
			logs = _.sortBy(logs, function(log) {
				return log.time;
			});
			logs.reverse();
			// logs from the latest 24h
			var timeLimit = (Date.now() / 1000) - 3600 * 24
			var filtered = _.filter(logs, function(log){
				return log.time > timeLimit;
			});
			var max = _.max(filtered, function(log) { return log.temperature; });
			var min = _.min(filtered, function(log) { return log.temperature; });
			if (max != Number.POSITIVE_INFINITY) {
				$("#maximum").html(createLogElement(max))
			} else {
				$("#maximum").html("No logs available")
			}
			if (min != Number.NEGATIVE_INFINITY) {
				$("#minimum").html(createLogElement(min))
			} else {
				$("#minimum").html("No logs available")
			}
			$("#latest").html(createLogElement(logs[0]));
			
			$("#history").html("<tr><th>Temperature (&deg;C)</th><th>Time</th></tr>");
			_.each(logs, function(log) {
				var row = $("<tr></tr>");
				$("<td>" + (Math.round((log.temperature - 273.15) * 100) / 100) + "</td>").appendTo(row);
				$("<td>" + new Date(log.time * 1000).toLocaleString() + "</td>").appendTo(row);
				row.appendTo($("#history"));
			});
			
		} else {
			$("#latest").html("No logs available")
			$("#maximum").html("No logs available")
			$("#minimum").html("No logs available")
			$("#history").html("No logs available")
		}
		$("#details-name").html(loc.name);
		$("#details").show();
		$('html, body').animate({
			scrollTop: $("#details").offset().top
		}, 1000);
	});
}

function createLogElement(log) {
	return $("<div></div>").html((Math.round((log.temperature - 273.15) * 100) / 100) + "&deg;C at " + (new Date(parseInt(log.time) * 1000)).toLocaleString());
}


function message(err, msg) {
	if (err) {
		$("#input-message").addClass("message-error");
	} else {
		$("#input-message").removeClass("message-error");
	}
	$("#input-message").html(msg)
}

function addLog() {
	$("#input-message").html()
	var kelvinValue;
	var temperatureUnit = $("#input-temperatureunit").val();
	var value = parseFloat($("#input-temperature").val());
	if (temperatureUnit == "celcius") {
		kelvinValue = value + 273.15;
	} else if (temperatureUnit == "fahrenheit") {
		kelvinValue = (value - 32) / 1.8 + 273.15;
	} else {
		kelvinValue = value;
	}
	if (kelvinValue < 0 || kelvinValue > 373.15) {
		message(true, "Invalid temperature");
		return
	}
	var timestamp = $("#input-time").datetimepicker("getDate").getTime();
	if (timestamp < 0 || timestamp > Date.now()) {
		message(true, "Invalid time");
		return
	}
	$.post(
		"/logs/add", 
		JSON.stringify({"locationId": currentLocation.id, "time": timestamp / 1000, "temperature": kelvinValue}), 
		function(data) {
			message(false, data.message);
			displayLocationDetails(currentLocation);
		},
		"json"
	).fail(function(data) {
		message(true, data.responseJSON.error);
	});
}

$("#input").submit(addLog);
$.extend($.datepicker,{_checkOffset:function(inst,offset,isFixed){return offset}});
$("#input-time").datetimepicker({
	controlType: 'select', 
	showSecond:true, 
	dateFormat: $.datepicker.ISO_8601,
    separator: 'T',
    timeFormat: 'HH:mm:ssz',
	oneLine:true});
	
$("#details").hide();