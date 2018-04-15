var AppComponent = React.createClass({
	getInitialState: function() {
		return {location: null};
	},
	render: function() {
		return React.createElement(
			"div",
			{"id": "app"},
			React.createElement(MapComponent, {setLocation: this.setLocation}),
			React.createElement(DetailsComponent, {location: this.state.location})
		);
	},
	setLocation: function(loc) {
		this.setState({location: loc});
	}
});

var MapComponent = React.createClass({
	getInitialState: function() {
		return {loading: true};
	},
	componentDidMount: function() {
		var component = this;
		$.getJSON("/locations", function(locations) {
			component.state.map = new google.maps.Map(ReactDOM.findDOMNode(component), {
			  zoom: 2,
			  center: {lat: 0, lng: 0}
			});
			$.each(locations, function(key, loc) {
				var marker = new google.maps.Marker({
				  position: {lat: loc.lat, lng: loc.long},
				  map: component.state.map
				});
				marker.addListener("click", function() {
					component.props.setLocation(loc);
				});
			});
			component.setState({loading: false});
		});
	},
	render: function() {
		var props = {id: "map"}
		if (this.state.loading) {
			props.className = "map-loading";
		}
		return React.createElement(
			"div",
			props,
			this.state ? "Loading..." : null
		);
	}
});

var DetailsComponent = React.createClass({
	getInitialState: function() {
		return {loading: true, logs: null, unit: TEMPUNIT.CELCIUS};
	},
	componentDidUpdate: function(prevProps, prevState, snapshot) {
		if (this.props.location != prevProps.location) {
			this.getLogs();
			window.scrollTo(0, ReactDOM.findDOMNode(this).offsetTop);
		}
	},
	getLogs: function() {
		var component = this;
		$.getJSON("/logs/" + this.props.location.id, function(logs) {
			component.setState({logs: logs});
		});
	},
 	reload: function() {
		this.getLogs();
		this.forceUpdate();
	},
	render: function() {
		if (this.props.location == null) {
			return React.createElement(
				"div",
				{"id": "details"}
			);
		}
		return React.createElement(
			"div",
			{"id": "details"},
			React.createElement(
				"h1",
				{className: "details-header"},
				this.props.location.name
			),
			React.createElement(TemperatureSelectComponent, {setUnit: this.setUnit}),
			React.createElement(
				LatestDataComponent,
				{logs: this.state.logs, unit: this.state.unit}
			),
			React.createElement(
				AddLogComponent,
				{location: this.props.location, unit: this.state.unit, reload: this.reload}
			),
			React.createElement(
				LogTableComponent,
				{logs: this.state.logs, unit: this.state.unit}
			)
		);
	},
	setUnit: function(unit) {
		this.setState({unit: unit});
	}
});

var TemperatureSelectComponent = React.createClass({
	render: function() {
		var component = this;
		return React.createElement(
			"div",
			{id: "details-tempsetting", onChange: function(event){component.props.setUnit(event.target.value)}},
			React.createElement(
				"select",
				null,
				React.createElement(
					"option",
					{value: TEMPUNIT.CELCIUS},
					"Celcius"
				),
				React.createElement(
					"option",
					{value: TEMPUNIT.KELVIN},
					"Kelvin"
				),
				React.createElement(
					"option",
					{value: TEMPUNIT.FAHRENHEIT},
					"Fahrenheit"
				)
			)
		);
	}
});

var LatestDataComponent = React.createClass({
	render: function() {
		var latest = null;
		var maximum = null;
		var minimum = null;
		if (this.props.logs != null) {
			// logs from the latest 24h
			var timeLimit = (Date.now() / 1000) - 3600 * 24;
			var filtered = _.filter(this.props.logs, function(log){
				return log.time > timeLimit;
			});
			var sorted = _.sortBy(filtered, function(log) {
				return log.time;
			});
			sorted.reverse();
			if (filtered.length > 0) {
				maximum = _.max(filtered, function(log) { return log.temperature; });
				minimum = _.min(filtered, function(log) { return log.temperature; });
				latest = sorted[0];
			}
		}
		return React.createElement(
			"div",
			{id: "details-latest"},
			React.createElement(
				"h3",
				null,
				"Latest"
			),
			latest == null ? "No logs available" : React.createElement(
				LogComponent,
				{log: latest, unit: this.props.unit}
			),
			React.createElement(
				"h3",
				null,
				"24h maximum"
			),
			maximum == null ? "No logs available" : React.createElement(
				LogComponent,
				{log: maximum, unit: this.props.unit}
			),
			React.createElement(
				"h3",
				null,
				"24h minimum"
			),
			minimum == null ? "No logs available" : React.createElement(
				LogComponent,
				{log: minimum, unit: this.props.unit}
			)
		);
	}
});

var TEMPUNIT = {KELVIN: 1, CELCIUS: 2, FAHRENHEIT: 3};
var UNITEXT = {};
UNITEXT[TEMPUNIT.KELVIN] = "K";
UNITEXT[TEMPUNIT.CELCIUS] = "°C";
UNITEXT[TEMPUNIT.FAHRENHEIT] = "°F";

var LogComponent = React.createClass({
	render: function() {
		return React.createElement(
			"div",
			null,
			this.props.log != null ? React.createElement(
				TemperatureComponent, 
				{temperature: this.props.log.temperature, unit: this.props.unit}
			) : null,
			this.props.log != null ? " at " + new Date(this.props.log.time * 1000).toLocaleString(): null
		);
	}
});

var TemperatureComponent = React.createClass({
	render: function() {
		var converted = Math.round(convertFromKelvin(this.props.temperature, this.props.unit) * 100) / 100;
		return React.createElement(
			"span", 
			null, 
			this.props.temperature != null ? converted.toString() + UNITEXT[this.props.unit] : ""
		);
	}
});

function convertFromKelvin(val, unit) {
	if (unit == TEMPUNIT.CELCIUS) {
		return val - 273.15;
	} else if (unit == TEMPUNIT.FAHRENHEIT) {
		return 1.8*val-459.67;
	} else if (unit == TEMPUNIT.KELVIN) {
		return val;
	}
}

function convertToKelvin(val, unit) {
	if (unit == TEMPUNIT.CELCIUS) {
		return val + 273.15;
	} else if (unit == TEMPUNIT.FAHRENHEIT) {
		return (val+459.67)/1.8;
	} else if (unit == TEMPUNIT.KELVIN) {
		return val;
	}
}

var AddLogComponent = React.createClass({
	getInitialState: function() {
		return {message: "", temperature: null, time: null};
	},
	componentDidMount: function() {
		var component = this;
		$("#details-addlog-time").datetimepicker({
			onSelect: function() {
				component.setState(
					{time: $(this).datetimepicker("getDate").getTime()}
				);
			},
			controlType: 'select', 
			showSecond:true, 
			dateFormat: $.datepicker.ISO_8601,
			separator: 'T',
			timeFormat: 'HH:mm:ssz',
			oneLine:true});
	},
	addLog: function() {
		var component = this;
		var temp = convertToKelvin(this.state.temperature, this.props.unit);
		if (temp < 0 || temp > 373.15) {
			this.setState({message: "Invalid temperature"});
			return;
		}
		var timestamp = this.state.time;
		if (timestamp < 0 || timestamp > Date.now()) {
			this.setState({message: "Invalid time"});
			return;
		}
		$.post(
			"/logs/add", 
			JSON.stringify({"locationId": component.props.location.id, "time": timestamp / 1000, "temperature": temp}), 
			function(data) {
				component.setState({message: data.message});
				component.props.reload();
			},
			"json"
		).fail(function(data) {
			component.setState({message: data.responseJSON.error});
		});
	},
	render: function() {
		var component = this;
		return React.createElement(
			"div",
			{id: "details-addlog"},
			React.createElement(
				"h3",
				null,
				"Add log:"
			),
			React.createElement(
				"form",
				{onSubmit: function(event){component.addLog(); event.preventDefault()}},
				"Temperature:",
				React.createElement("br"),
				React.createElement(
					"input", 
					{
						onChange: function(event){
							component.setState(
								{temperature: parseFloat(event.target.value)}
							);
						},
						type: "number", step: "any", size: "5", required: true
					}
				),
				React.createElement("br"),
				"Time:",
				React.createElement("br"),
				React.createElement(
					"input",
					{
						type: "text", id: "details-addlog-time", size: "30", required: true
					}
				),
				React.createElement("br"),
				React.createElement(
					"input",
					{type: "submit", value: "Add log"}
				),
				React.createElement("br"),
				React.createElement(
					"div",
					{className: "details-message"},
					this.state.message
				)
			)
		);
	}
});

var LogTableComponent = React.createClass({
	render: function() {
		var rows = [];
		var component = this;
		_.each(
			_.sortBy(this.props.logs, function(log) {
				return -log.time;
			}),
			function(log, key) {
				rows.push(React.createElement(
					"tr",
					{key: key},
					React.createElement(
						"td",
						null,
						React.createElement(
							TemperatureComponent,
							{temperature: log.temperature, unit: component.props.unit}
						)
					),
					React.createElement(
						"td",
						null,
						new Date(log.time * 1000).toLocaleString()
					)
				));
			}
		);

		return React.createElement(
			"table",
			{id: "details-logtable"},
			React.createElement(
				"tbody",
				null,
				React.createElement(
					"tr",
					null,
					React.createElement(
						"th",
						null,
						"Temperature"
					),
					React.createElement(
						"th",
						null,
						"Time"
					)
				),
				rows
			)
		)
	}
});

function initMap() {
	ReactDOM.render(
		React.createElement(AppComponent),
		$("#app")[0]
	)
}