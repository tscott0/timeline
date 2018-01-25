{{define "timeline.js"}}

    // DOM element where the Timeline will be attached
    var container = document.getElementById('visualization');

    var items = new vis.DataSet({{.EventArray}});

    // Configuration for the Timeline
    var options = {{.OptionsObject}};

    // Create a Timeline
    var timeline = new vis.Timeline(container, items, options);

    // Event listeners
    timeline.on('click', function (properties) {
        console.log(properties);
    });
    
{{end}}